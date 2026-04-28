# Update Catalog Entry — Design Spec

> Approved by user 2026-04-27

**Goal:** A single Harness IDP workflow that lets a team add Resources (database, cache, messaging) and APIs (REST, async producer, async consumer) to an existing Component, and optionally update the parent System — all committed to master with validation guardrails.

---

## Scope

Three deliverables built in order:

| Deliverable | File | Purpose |
|---|---|---|
| Go validator | `tools/catalog-validator/main.go` | Validates generated YAML before commit |
| Workflow YAML | `workflows/update-catalog-entry.yaml` | IDP form — user input |
| Pipeline YAML | `pipelines/update-catalog-entry-pipeline.yaml` | Clone → append → validate → commit → register |

---

## Workflow Form Design

**Section 1 — Select Component**
- `SelectFieldFromApi` pulling from Harness Catalog API
- Auto-fills: SYSID, repo slug, system identifier (read-only via ContextViewer)

**Section 2 — System (optional toggle)**
- `dependencies/oneOf` on toggle field
- If "Yes": title, description, domain (enum), Confluence URL, Runbook URL, tags (comma-separated string)

**Section 3 — Resources (type: array)**
- Databases — name (regex), engine enum (postgresql/mysql/oracle/mssql), description
- Caches — name (regex), engine enum (redis/memcached), description
- Messaging — name (regex), type enum (kafka/ibm-mq/rabbitmq/ftp-server), description

**Section 4 — APIs (type: array)**
- REST APIs — name (regex), endpoint URL, description
- Async Producers — name (regex), protocol enum (kafka/ftp), topic/path, description
- Async Consumers — name (regex), protocol enum (kafka/ftp), topic/path, consumer group (kafka only), description

**Guardrails:**
- Component from live catalog — no free text
- SYSID auto-filled from selected Component
- All name fields: `pattern: "^[a-z][a-z0-9-]+$"`
- Engine/type/protocol: enums — no free text
- Arrays serialized to JSON strings in pipeline payload

**Steps:**
1. `debug.log` — log all inputs
2. `trigger:harness-custom-pipeline` — fire pipeline with full JSON payload

---

## Pipeline Design

**Trigger payload fields:**
- `componentIdentifier` — e.g. `payments_service`
- `repoSlug` — e.g. `bank-org/payments-service`
- `sysid` — e.g. `SYSID-00123`
- `systemIdentifier` — e.g. `payments_platform`
- `ownerGroup` — e.g. `team-payments`
- `updateSystem` — `"true"` or `"false"`
- `systemTitle`, `systemDescription`, `systemDomain`, `systemConfluence`, `systemRunbook`, `systemTags`
- `databases` — JSON array string
- `caches` — JSON array string
- `messaging` — JSON array string
- `restApis` — JSON array string
- `asyncProducers` — JSON array string
- `asyncConsumers` — JSON array string
- `requestorName`, `requestorEmail`

**Pipeline stages:**

| Stage | Steps | Notes |
|---|---|---|
| clone_repo | ShellScript: git clone | GitHub PAT from Vault |
| append_entries | ShellScript: go run tools/catalog-appender/main.go | Reads existing catalog-info.yaml, appends new entries |
| validate | ShellScript: go run tools/catalog-validator/main.go | Exits non-zero on any error |
| commit_push | ShellScript: git commit + push to master | Signed commit, Rohid Dev author |
| register_entities | Http steps (one per entity type) | POST to Harness Catalog API per entity |
| notify_on_failure | ShellScript: email API | Only on pipeline failure |

---

## Go Validator — Checks

1. `apiVersion: harness.io/v1` present on every document
2. Top-level `identifier`, `name`, `kind`, `owner` present
3. `owner` format: `group:account/<id>`
4. `spec.system` is a list (not string), each entry: `system:account/<id>`
5. `spec.dependsOn` entries: `resource:account/<id>`
6. No duplicate `identifier` values within the file
7. All `name` fields match `^[a-z][a-z0-9-]+$`
8. `kind` is one of: `System`, `Component`, `Resource`, `API`

Exit 0 = valid. Exit 1 = prints human-readable errors, one per line.

---

## Guardrails Summary

| Risk | Guardrail |
|---|---|
| Typos in names | Regex `^[a-z][a-z0-9-]+$` on all name fields |
| Wrong owner/system format | Validator checks `group:account/` and `system:account/` prefixes |
| Duplicate identifiers | Validator checks for duplicate `identifier` fields |
| Broken YAML | Validator parses entire file before any commit |
| Wrong component | Live catalog picker — no free text |
| Missing SYSID | Auto-filled from Component — not user-entered |
| Bad commit | Validation failure → pipeline aborts, nothing committed |
