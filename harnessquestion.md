# Question for Harness Engineer

**Subject:** Catalog Ingestion API — can it patch `spec.dependsOn` on a Git-backed entity?

---

We are building an automated pipeline that reads relationship data from ServiceNow `cmdb_rel_ci`
and needs to update `spec.dependsOn` on catalog entities in Harness IDP.

Our entities are registered via **Import from Git** — they are Git-backed, sourced from
`catalog-info.yaml` in each repo.

We want to avoid committing machine-generated changes back to every repo's `catalog-info.yaml`.
We would prefer to push `spec.dependsOn` directly via the Catalog Ingestion API.

---

## Please confirm the following

### 1 — Can the Catalog Ingestion API patch `spec.dependsOn` on a Git-backed entity?

Or does it only work for entities that are fully ingested (not Git-backed)?

### 2 — If it can — does the next Git sync overwrite the API-pushed value?

Our concern: if a developer pushes a commit to their repo, Harness re-reads `catalog-info.yaml`
from Git. Does this overwrite the `spec.dependsOn` value we set via the API?

### 3 — Is there a merge strategy?

Can specific fields be marked as API-managed so they survive Git re-syncs? For example —
Git owns everything except `spec.dependsOn`, which is owned by the ingestion pipeline.

### 4 — What is the exact API call to update `spec.dependsOn`?

Please share the endpoint, method, and payload format. For example:

```
PATCH /idp/api/v1/catalog/entities/component/default/payments-service
{
  "spec": {
    "dependsOn": [
      "component:default/auth-service",
      "resource:default/payments-postgres"
    ]
  }
}
```

Is this the correct format or does it require the full entity body?

### 5 — If Git-backed entities are not supported — is this on the roadmap?

We would like to raise this as a feature request. The use case is: automated dependency sync
from an external CMDB (`cmdb_rel_ci`) into `spec.dependsOn` without requiring Git commits
to every repo.

---

## Context

- **Use case:** Service Mapping Initiative — 19,000 repos, each gets a `now.yaml` declaring
  1-level-deep dependencies. RapDev CSDM as Code syncs `now.yaml` → `cmdb_rel_ci` via webhook
  (instant). We need `cmdb_rel_ci` → `spec.dependsOn` to complete the chain so the native
  dependency graph and blast radius cards render in IDP.

- **Our fallback:** If the Catalog Ingestion API cannot do this, we will commit
  `catalog-info.yaml` back to each repo via the GitHub API (Option A). But we want to
  confirm the Ingestion API option first as it is significantly cleaner — no machine commits
  in developer repos, no GitHub write token required.

- **Critical question:** Question 2 is the most important. Even if Question 1 is yes, if
  Git sync overwrites the API-pushed value on the next commit, the approach does not work
  and we fall back to Option A.
