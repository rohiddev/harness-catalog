# Harness IDP — Platform Setup Guide

Step-by-step instructions for the platform team to configure Harness IDP.
Run these steps once when setting up IDP for the first time.

> **For service teams** registering their own service — see `README.md`.

---

## Prerequisites

- Harness IDP 2.0 enabled on your account
- GitHub connector configured in Harness pointing to your GitHub organisation
- Access to **IDP Admin** in Harness

---

## Step 1 — Import Catalog Entities from Git

**IDP > Software Catalog > Import from Git**

| Field | Value |
|---|---|
| Repository | `https://github.com/rohiddev/harness-catalog` |
| Branch | `master` |
| File path | `catalog-info.yaml` |

Click **Import**. Harness creates all entities from the file:

| Entity | Kind | Type |
|---|---|---|
| `payments-platform` | System | — |
| `payments-service` | Component | service |
| `payments-rest-api` | API | openapi |
| `payments-ftp-outbound` | API | asyncapi |
| `payments-ftp-inbound` | API | asyncapi |
| `payments-kafka-producer` | API | asyncapi |
| `payments-kafka-consumer` | API | asyncapi |
| `payments-postgres` | Resource | database |
| `payments-redis` | Resource | cache |
| `payments-ftp` | Resource | ftp-server |
| `payments-kafka` | Resource | messaging-queue |

After import, verify all entities appear under **IDP > Software Catalog**.

---

## Step 2 — Configure Homepage & Sidenav

### 2a — Upload Bank Logo (portal-wide)

> The logo appears in the top-left of the sidenav on every page. Set it once here — it does not need to be in any YAML file.

1. Log in to Harness IDP as a **Platform Admin**
2. Go to **Admin > Appearance**
3. Upload your bank logo under the **Logo** field
   - Format: PNG or SVG
   - Recommended size: 200×40px (horizontal) or 40×40px (icon)
   - Use a transparent background (sidenav is dark)
4. Click **Save**

**Production swap:** Repeat this step on the prod Harness account. No YAML changes needed.

---

### 2b — Apply Sidenav Layout

> File: `sidenav.yaml` (in this repo)

1. Go to **Admin > Layout**
2. Select **Sidenav**
3. Click **Edit YAML**
4. Paste the full contents of `sidenav.yaml`
5. Click **Save**

**Expected result — left sidebar shows exactly 6 items:**

| # | Label | Navigates to |
|---|---|---|
| 1 | Home | IDP homepage |
| 2 | Catalog | Software catalog |
| 3 | Self-Service | Workflows (100+ available) |
| 4 | Scorecards | Service quality scores |
| 5 | TechDocs | Service documentation |
| 6 | My Starred | User's bookmarked entities |

---

### 2c — Apply Homepage Layout

> File: `homepage.yaml` (in this repo)

**Before pasting — replace these placeholders with real internal URLs:**

| Placeholder in homepage.yaml | Replace with |
|---|---|
| `https://confluence.bank.com/display/ENG/freeze-calendar` | Your Confluence freeze calendar URL |
| `https://servicenow.bank.com/change` | Your ServiceNow change board URL |
| `https://oncall.bank.com/incidents` | Your incident dashboard URL |
| `https://oncall.bank.com` | Your on-call schedule URL |
| `platform-team@bank.com` (×2) | Your platform team email address |

> `https://app.harness.io` and `https://app.datadoghq.com` are real — leave them as-is.

**Steps:**

1. Go to **Admin > Layout**
2. Select **Homepage**
3. Click **Edit YAML**
4. Paste the full contents of `homepage.yaml` (with placeholders replaced)
5. Click **Save**

**Expected result — homepage layout:**

| Section | Content |
|---|---|
| **Banner** | Green strip: "No active production freezes" (update manually during freezes) |
| **Header** | Personalised greeting: `Good morning, <first name>. Welcome to the Bank Developer Portal.` |
| **Quick links** | Catalog · Self-Service · Scorecards · TechDocs |
| **Row 1** | Platform Status (full width) — freeze calendar, CAB, incidents, on-call, pipelines, Datadog |
| **Row 2** | Recently Visited + Top Visited |
| **Row 3** | My Starred + What's New |
| **Row 4** | Getting Started checklist + IDP 101 key concepts |

**Verify with a non-admin developer account:**
- `<+greeting>` renders as `Good morning` / `Good afternoon` / `Good evening`
- `<+first_name>` renders as the developer's actual first name (not the literal text)
- All 4 rows of cards are visible
- Recently Visited populates after browsing a few catalog entities
- Starred Entities populates after starring an entity (click ⭐ on any entity page)

---

### 2d — Managing the Freeze Banner

The banner is **static text** — updated manually by the Platform Engineering team in ~2 minutes.

**To activate a production freeze:**

1. Go to **Admin > Layout > Homepage > Edit YAML**
2. Find the `banner:` section
3. Change the text to:
   ```yaml
   banner:
     text: "🔴 PROD FREEZE ACTIVE — no deployments to production until [DATE e.g. Friday 1 May 17:00 UTC]. Contact platform-team@bank.com"
   ```
4. Click **Save**

**To deactivate:**

1. Go to **Admin > Layout > Homepage > Edit YAML**
2. Restore the banner to:
   ```yaml
   banner:
     text: "✅ No active production freezes. For platform support contact platform-team@bank.com"
   ```
3. Click **Save**

> Total time: ~2 minutes. No code change, no PR, no deployment required.

---

## Step 3 — Configure Entity Layouts

**IDP > Admin > Configure > Layout > Catalog Entities**

Open each layout file from this repo, paste the full YAML into the corresponding editor and click **Save**.

| Layout file | Where to paste |
|---|---|
| `layout-system.yaml` | Layout > System > Edit |
| `layout-component-service.yaml` | Layout > Component > service > Edit |
| `layout-resource-database.yaml` | Layout > Resource > database > Edit |
| `layout-resource-messaging-queue.yaml` | Layout > Resource > messaging-queue > Edit |
| `layout-resource-cache.yaml` | Layout > Resource > cache > Edit |

**What each layout adds:**

| Entity | Tabs |
|---|---|
| System | Overview (Ownership & Compliance, Platform & SLA), Components, APIs, Resources |
| Component / service | Overview, Dependencies (graph + blast radius), API, Kubernetes, CI/CD, Scorecard, Docs |
| Resource / database | Overview (Database Details), Dependencies |
| Resource / messaging-queue | Overview (Queue Details), Dependencies |
| Resource / cache | Overview (Cache Details), Dependencies |

After saving each layout, open the corresponding entity in the catalog and verify the correct tabs and cards appear.

---

## Step 4 — Connect GitHub as Scorecard Data Source

> Skip if GitHub is already connected as a data source.

**IDP > Configure > Scorecards > Data Sources > + Add**

Connect your GitHub organisation using the existing GitHub connector.
This is required for:
- Branch protection check
- `now.yaml` file existence check (Phase 2 — add once now.yaml schema confirmed)

All other checks use Catalog data only and work without this step.

---

## Step 5 — Create Scorecard Checks

**IDP > Configure > Scorecards > Checks > + Create Check**

### catalog-info.yaml Checks (create now)

#### Check 1 — Has SYS ID
| Field | Value |
|---|---|
| Name | Has SYS ID |
| Data Source | Catalog |
| Data Point | Annotation Exists |
| Operator | Equal to |
| Value | `bank.com/sysid` |

#### Check 2 — Has Owner Defined
| Field | Value |
|---|---|
| Name | Has owner defined |
| Data Source | Catalog |
| Data Point | Annotation Exists |
| Operator | Equal to |
| Value | `github.com/project-slug` |

#### Check 3 — Has TechDocs
| Field | Value |
|---|---|
| Name | Has TechDocs configured |
| Data Source | Catalog |
| Data Point | Annotation Exists |
| Operator | Equal to |
| Value | `backstage.io/techdocs-ref` |

#### Check 4 — Has Contact Details
| Field | Value |
|---|---|
| Name | Has contact details |
| Data Source | Catalog |
| Data Point | Annotation Exists |
| Operator | Equal to |
| Value | `pagerduty.com/service-id` |

#### Check 5 — Has Datadog Dashboard
| Field | Value |
|---|---|
| Name | Has Datadog dashboard |
| Data Source | Catalog |
| Data Point | Annotation Exists |
| Operator | Equal to |
| Value | `datadoghq.com/dashboard-url` |

> **Optional check** — mark this as informational in the scorecard if not all teams have
> Datadog dashboards. Mandatory for Tier-1 and Tier-2 services.

#### Check 6 — Has Jira Project
| Field | Value |
|---|---|
| Name | Has Jira project key |
| Data Source | Catalog |
| Data Point | Annotation Exists |
| Operator | Equal to |
| Value | `jira/project-key` |

#### Check 7 — Has CI/CD Configured
| Field | Value |
|---|---|
| Name | Has CI/CD configured |
| Data Source | Catalog |
| Data Point | Annotation Exists |
| Operator | Equal to |
| Value | `harness.io/project-url` |

#### Check 8 — Branch Protection Enabled
| Field | Value |
|---|---|
| Name | Branch protection enabled |
| Data Source | GitHub |
| Data Point | Is Branch Protected |
| Operator | Equal to |
| Value | `true` |

---

### now.yaml Checks

> Schema is production-ready. Create these checks now.
> Reference: `now.yaml` (payments-service worked example) and `now.template.yaml` (golden template).

#### Check 9 — now.yaml Exists in Repo
| Field | Value |
|---|---|
| Name | now.yaml exists in repo |
| Data Source | GitHub |
| Data Point | File exists |
| File path | `now.yaml` |
| Branch | `master` |

#### Check 10 — Has Dependency Mapping
| Field | Value |
|---|---|
| Name | Has dependency mapping |
| Mode | Advanced (JEXL) |
| Expression | `catalog.annotationExists."bank.com/sysid" == true && spec.dependsOn != null` |

#### Check 11 — Has Tier Declared
| Field | Value |
|---|---|
| Name | Has tier declared |
| Data Source | Catalog |
| Data Point | Annotation Exists |
| Operator | Equal to |
| Value | `bank.com/tier` |

#### Check 12 — Has Team Declared
| Field | Value |
|---|---|
| Name | Has team declared |
| Data Source | Catalog |
| Data Point | Annotation Exists |
| Operator | Equal to |
| Value | `bank.com/team` |

---

## Step 6 — Create the Production Readiness Scorecard

**IDP > Configure > Scorecards > + New Scorecard**

| Field | Value |
|---|---|
| Name | Production Readiness |
| Description | Gates production deployment — both catalog-info.yaml and now.yaml must be complete |
| Kind | Component |
| Type | service |
| Lifecycle | production |
| Pass threshold | **100%** |

Add all checks created in Step 4. Click **Publish**.

> Add Checks 9–12 to this scorecard once now.yaml schema is confirmed.

---

## Step 7 — Add Production Gate to Deployment Pipelines

Add this step before the production stage in every Harness CD pipeline.
It aborts the deployment if the scorecard is below 100%.

```yaml
- step:
    name: Check Production Readiness Scorecard
    identifier: CheckProductionReadinessScorecard
    type: Http
    spec:
      url: https://app.harness.io/idp/api/scorecards/<+pipeline.variables.entityName>/score
      method: GET
      headers:
        - key: x-api-key
          value: <+secrets.getValue("harness_api_key")>
        - key: Content-Type
          value: application/json
      assertion: <+json.select("score", httpResponseBody)> == 100
    failureStrategies:
      - onFailure:
          errors: [AssertionError]
          action:
            type: Abort
```

Add `entityName` as a pipeline variable (value = catalog entity name, e.g. `payments-service`).

---

## Step 8 — Verify

Open **payments-service** in the catalog and confirm:

- [ ] **Overview tab** — Scorecard summary card visible showing 12/12 checks
- [ ] **Dependencies tab** — Dependency graph renders, blast radius cards visible
- [ ] **API tab** — REST API, Kafka producer/consumer, FTP inbound/outbound all listed
- [ ] **Kubernetes tab** — K8s workload visible (requires Kubernetes connector)
- [ ] **CI/CD tab** — Harness pipeline runs visible
- [ ] **Scorecard tab** — All 8 catalog-info.yaml checks listed with Pass/Fail results

---

## Step 9 — Roll Out to Service Teams

Share the repo with all teams:

> **`https://github.com/rohiddev/harness-catalog`**
>
> Each team — two files, both required for 100% scorecard:
>
> **catalog-info.yaml** (IDP catalog)
> 1. Copy `catalog-info.template.yaml` → save as `catalog-info.yaml` in repo root
> 2. Fill all `# REQUIRED` fields, remove sections that don't apply
> 3. Reference `catalog-info.yaml` (payments-service) as a worked example
>
> **now.yaml** (ServiceNow CMDB + blast radius)
> 4. Copy `now.template.yaml` → save as `now.yaml` in repo root
> 5. Fill `sysId`, `name`, `projectName`
> 6. Declare every dependency: APIs called, databases, Kafka topics, FTP paths
> 7. Reference `now.yaml` (payments-service) as a worked example
>
> **Go live**
> 8. Raise a PR — platform team imports the service into IDP once merged
> 9. Achieve 100% Production Readiness scorecard before deploying to production

---

## How Relationships Render in IDP

Once a developer adds `spec.dependsOn` to their `catalog-info.yaml` and pushes to Git,
IDP renders three things automatically on the entity page — no extra configuration needed.

### 1 — Dependency Graph

**Dependencies tab → `EntityCatalogGraphCard`**

An interactive visual graph showing the full dependency picture:

```
                    ┌──────────────────┐
  web-checkout ────▶│                  │──▶ auth-service
  mobile-app   ────▶│ payments-service │──▶ fraud-detection-api
  batch-settle ────▶│                  │──▶ payments-postgres
                    └──────────────────┘──▶ payments-kafka
```

- **Right side** — what this service depends on (declared in `spec.dependsOn`)
- **Left side** — who depends on THIS service (blast radius — auto-derived by IDP)

IDP automatically creates the reverse `dependencyOf` edge. If Service A declares
`dependsOn: [payments-service]`, payments-service's graph shows Service A on the left
with no extra work from the payments team.

### 2 — Dependency List Cards

**Dependencies tab — three list cards below the graph**

| Card | Shows |
|---|---|
| `EntityDependsOnComponentsCard` | Services this entity calls — auth-service, fraud-detection-api |
| `EntityDependsOnResourcesCard` | Infrastructure this entity uses — postgres, redis, kafka, ftp |
| `EntityHasComponentsCard` | Services that depend on THIS entity — blast radius |

### 3 — Scorecard Summary Card

**Overview tab → `EntityScoreCard`**

Shows the Production Readiness score at a glance. Once both `catalog-info.yaml` and
`now.yaml` score 100%, the card shows green and production deployment is unblocked.

### Developer responsibility

Both files live in Git — both are the developer's responsibility to keep in sync:

| File | Updated by | Drives |
|---|---|---|
| `now.yaml` | Developer | RapDev webhook → `cmdb_rel_ci` in ServiceNow (instant) |
| `catalog-info.yaml` `spec.dependsOn` | Developer | IDP Git webhook → dependency graph in catalog |

Git is the source of truth. No pipeline pulls from `cmdb_rel_ci` — the developer
updates both files in the same commit when dependencies change.

---

## Scorecard — Production Readiness

The scorecard is the production gate. **A service cannot deploy to production until it
scores 100%.** The pipeline aborts automatically if the score is below 100%.

### 12 Checks

**catalog-info.yaml checks — build these now (8 checks)**

| # | Check | Data Source | What it validates |
|---|---|---|---|
| 1 | Has SYS ID | Catalog | `bank.com/sysid` annotation present — links entity to CMDB |
| 2 | Has owner defined | Catalog | `github.com/project-slug` present — service has a team |
| 3 | Has TechDocs | Catalog | `backstage.io/techdocs-ref` present — documentation exists |
| 4 | Has contact details | Catalog | `pagerduty.com/service-id` present — on-call contact wired |
| 5 | Has Datadog dashboard | Catalog | `datadoghq.com/dashboard-url` present — APM in place (optional for Tier-3) |
| 6 | Has Jira project | Catalog | `jira/project-key` present — work tracking linked |
| 7 | Has CI/CD configured | Catalog | `harness.io/project-url` present — CI/CD tab will show pipeline runs |
| 8 | Branch protection enabled | GitHub | Master branch is protected |

**now.yaml checks — add once RapDev confirms schema (4 checks)**

| # | Check | Data Source | What it validates |
|---|---|---|---|
| 9 | now.yaml exists in repo | GitHub | File present on master branch |
| 10 | Relationships declared | GitHub | `relationships:` key present in now.yaml — cmdb_rel_ci will be populated |
| 11 | Dependencies synced to catalog | Catalog | `spec.dependsOn` not empty — IDP graph will render |
| 12 | Has tier declared | Catalog | `bank.com/tier` annotation present — CMDB tier captured |

### What the scorecard page looks like

```
Production Readiness — payments-service          100% ✅

catalog-info.yaml (8/8)
  ✅  Has SYS ID
  ✅  Has owner defined
  ✅  Has TechDocs
  ✅  Has contact details
  ✅  Has Datadog dashboard
  ✅  Has Jira project key
  ✅  Has CI/CD configured
  ✅  Branch protection enabled

now.yaml (4/4)
  ✅  now.yaml exists in repo
  ✅  Relationships declared in now.yaml
  ✅  Dependencies synced to catalog
  ✅  Has tier declared

🟢  Production deployment unblocked
```

### Tier requirements

| Tier | Requirement |
|---|---|
| Tier-1 | All 12 checks — 100% mandatory |
| Tier-2 | All 12 checks — 100% mandatory |
| Tier-3 | Checks 1–8 only — now.yaml checks optional until rollout complete |

### How the production gate works

```
Developer raises PR to deploy to production
        ↓
Harness CD pipeline runs
        ↓
Step: Check Production Readiness Scorecard
        ↓
Score < 100%  →  pipeline ABORTS — developer fixes failing checks
Score = 100%  →  deployment proceeds
```

Scorecard evaluates every 15 minutes. After fixing a failing check, the developer
waits for the next evaluation then re-triggers the pipeline.

### Pipeline gate YAML

Add this step before the production stage in every Harness CD pipeline:

```yaml
- step:
    name: Check Production Readiness Scorecard
    identifier: CheckProductionReadinessScorecard
    type: Http
    spec:
      url: https://app.harness.io/idp/api/scorecards/<+pipeline.variables.entityName>/score
      method: GET
      headers:
        - key: x-api-key
          value: <+secrets.getValue("harness_api_key")>
        - key: Content-Type
          value: application/json
      assertion: <+json.select("score", httpResponseBody)> == 100
    failureStrategies:
      - onFailure:
          errors: [AssertionError]
          action:
            type: Abort
```

Add `entityName` as a pipeline variable (e.g. `payments-service`).

---

## Optional — ServiceNow CMDB Integration

**IDP > Configure > Integrations > + New Integration > ServiceNow**

> Configure this if your organisation manages services in ServiceNow CMDB.
> It auto-populates governance and compliance fields from CMDB into the catalog.

### Prerequisites
- Feature flag `IDP_CATALOG_CD_AUTO_DISCOVERY` enabled (contact Harness Support)
- Harness CD enabled on the same account as IDP
- A Harness ServiceNow connector (username/password or OAuth)

### Field Mappings

CMDB table: `cmdb_ci_service`

| CMDB Column | Catalog Field | Notes |
|---|---|---|
| `sys_id` | `metadata.name` | Unique identifier |
| `name` | `metadata.title` | Display name |
| `u_tier` | `spec.properties.tier` | Tier 1/2/3 |
| `u_team` | `spec.owner` | Maps to catalog owner |
| `u_data_classification` | `spec.properties.dataClassification` | |
| `u_sox_in_scope` | `spec.properties.soxInScope` | |
| `u_email_dl` | `spec.properties.emailDl` | |
| `lifecycle` | `spec.lifecycle` | production / deprecated |

**Correlation field:** `name` — merges CMDB data into existing catalog entities rather than creating duplicates.

### After Sync
- Entities already in catalog → select **Merge**
- New entities not yet in catalog → select **Register**

---

## Reference

| File | Purpose |
|---|---|
| `catalog-info.template.yaml` | Golden template — teams copy this for catalog-info.yaml |
| `catalog-info.yaml` | Worked example — Payments Platform (catalog-info) |
| `now.template.yaml` | Golden template — teams copy this for now.yaml |
| `now.yaml` | Worked example — Payments Platform (now.yaml) |
| `homepage.yaml` | Homepage layout — apply via Admin > Layout > Homepage |
| `sidenav.yaml` | Sidenav layout — apply via Admin > Layout > Sidenav |
| `layout-component-service.yaml` | Component / service page layout |
| `layout-system.yaml` | System page layout |
| `layout-resource-*.yaml` | Resource page layouts |
| `preview.html` | Offline UI preview — open in browser (includes homepage view) |
| `wiki/integrations/` | ServiceNow, blast radius runbook, dependency ingestion docs |

---

*Author: Rohid Dev · github.com/rohiddev*
