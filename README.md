# Harness IDP — Payments Platform Catalog Setup

This repo contains catalog entity definitions and layout YAMLs for the Payments Platform.
Follow the steps below **in order** to configure everything in Harness IDP.

---

## Step 0 (Optional) — ServiceNow CMDB Integration

If your organisation manages services in ServiceNow CMDB, configure this integration first.
It auto-populates governance and compliance fields directly from CMDB — eliminating manual
annotation maintenance in `catalog-info.yaml`.

**IDP > Configure > Integrations > + New Integration > ServiceNow**

### Prerequisites
- Feature flag `IDP_CATALOG_CD_AUTO_DISCOVERY` must be enabled (contact Harness Support)
- Harness CD enabled on the same account as IDP
- A Kubernetes agent installed in your environment (prompted on first setup)
- A Harness ServiceNow connector (username/password or OAuth)

### Field Mappings for Payments Platform

Configure the following mappings when setting up the integration.
CMDB table: **`cmdb_ci_service`**

| CMDB Column | Catalog YAML Field | Notes |
|---|---|---|
| `sys_id` | `metadata.name` | Used as unique identifier |
| `u_app_key` | `spec.properties.appKey` | |
| `u_domain` | `spec.properties.domain` | |
| `u_tier` | `spec.properties.tier` | |
| `u_product_owner` | `spec.properties.productOwner` | |
| `u_tech_lead` | `spec.properties.techLead` | |
| `u_team` | `spec.owner` | Maps to catalog owner |
| `u_email_dl` | `spec.properties.emailDl` | |
| `u_data_classification` | `spec.properties.dataClassification` | |
| `u_sox_in_scope` | `spec.properties.soxInScope` | |
| `u_gdpr_in_scope` | `spec.properties.gdprInScope` | |
| `u_last_pen_test` | `spec.properties.lastPenTest` | |
| `u_last_dr_test` | `spec.properties.lastDrTest` | |
| `u_rto_hours` | `spec.properties.rtoHours` | |
| `u_rpo_hours` | `spec.properties.rpoHours` | |
| `lifecycle` | `spec.lifecycle` | |
| `u_sla_uptime` | `spec.properties.slaUptime` | |

**Display Name column:** `u_app_key` (or `name`)
**Unique Identifier column:** `sys_id`
**Correlation field:** `name` — matches existing catalog entities by name so CMDB data merges into them rather than creating duplicates

### Sync Schedule
Set to daily or on-demand depending on how frequently CMDB records change.

### After Sync
Discovered records appear in the **Discovered** tab. For the payments entities:
- Entities already in the catalog → select **Merge** to enrich with CMDB fields
- New entities not yet in the catalog → select **Register** to create them

> **Layout note:** CMDB-sourced fields land in `spec.properties`, not `metadata.annotations`.
> The layout YAMLs in this repo currently read from `metadata.annotations['bank.com/...']`.
> If you use the CMDB integration, update the layout `value:` fields to use
> `<+spec.properties.fieldName>` for any field sourced from CMDB.

---

## Step 1 — Import Catalog Entities from Git

**IDP > Software Catalog > Import from Git**

Import `catalog-info.yaml` from this repository. Harness IDP 2.0 automatically converts
the Backstage YAML format into the Harness Catalog Entity Model and creates all 6 entities:

| Entity | Kind | Type |
|---|---|---|
| payments-platform | System | — |
| payments-service | Component | service |
| payments-rest-api | API | openapi |
| payments-postgres | Resource | database |
| payments-rabbitmq | Resource | messaging-queue |
| payments-redis | Resource | cache |

> **Note:** In IDP 2.0 entities are **inline** — their full lifecycle (edit, delete) is managed
> through the UI or API after import. You do not need to keep the Git file in sync.

After import, verify all 6 entities appear in the Software Catalog.

---

## Step 2 — Configure Entity Layouts

**IDP > Admin > Configure > Layout > Catalog Entities**

Paste each YAML file into the corresponding layout editor and click **Save**.

| Layout File | Where to Paste |
|---|---|
| `layout-system.yaml` | Layout > System > Edit |
| `layout-component-service.yaml` | Layout > Component > service > Edit |
| `layout-resource-database.yaml` | Layout > Resource > database > Edit |
| `layout-resource-messaging-queue.yaml` | Layout > Resource > messaging-queue > Edit |
| `layout-resource-cache.yaml` | Layout > Resource > cache > Edit |

After saving, open each entity in the catalog and verify the correct tabs and cards appear.

**What each layout adds:**

| Entity | Tabs |
|---|---|
| System | Overview (Ownership & Compliance, Platform & SLA), Components, APIs, Resources |
| Component / service | Overview (About, Links, Ownership, Deployment, Scorecard card), Dependencies, API, Kubernetes, CI/CD, Scorecard, Docs |
| Resource / database | Overview (Database Details), Dependencies |
| Resource / messaging-queue | Overview (Queue Details, Consumers), Dependencies |
| Resource / cache | Overview (Cache Details), Dependencies |

---

## Step 3 — Connect GitHub as a Scorecard Data Source

> Skip this step if GitHub is already connected.

**IDP > Configure > Scorecards > Data Sources**

Connect your GitHub organisation. This is required for the **Branch Protection** check (Check 7).
All other checks use Catalog data only and work without GitHub.

---

## Step 4 — Create Scorecard Checks

**IDP > Configure > Scorecards > Checks tab > Create Custom Check**

Create the following 7 checks:

### Check 1 — Has Owner Defined
| Field | Value |
|---|---|
| Name | Has owner defined |
| Mode | Basic |
| Data Source | Catalog |
| Data Point | Annotation Exists |
| Operator | Equal to |
| Value | `github.com/project-slug` |

### Check 2 — Has Runbook Link
| Field | Value |
|---|---|
| Name | Has runbook link |
| Mode | Basic |
| Data Source | Catalog |
| Data Point | Annotation Exists |
| Operator | Equal to |
| Value | `bank.com/oncall-schedule` |

### Check 3 — Has TechDocs Configured
| Field | Value |
|---|---|
| Name | Has TechDocs configured |
| Mode | Basic |
| Data Source | Catalog |
| Data Point | Annotation Exists |
| Operator | Equal to |
| Value | `backstage.io/techdocs-ref` |

### Check 4 — Has On-Call Schedule
| Field | Value |
|---|---|
| Name | Has on-call schedule |
| Mode | Basic |
| Data Source | Catalog |
| Data Point | Annotation Exists |
| Operator | Equal to |
| Value | `bank.com/oncall-schedule` |

### Check 5 — Has API Definition
| Field | Value |
|---|---|
| Name | Has API definition |
| Mode | Advanced (JEXL) |
| Expression | `catalog.annotationExists."github.com/project-slug" == true` |

### Check 6 — Pen Test Within 12 Months
| Field | Value |
|---|---|
| Name | Pen test within 12 months |
| Mode | Advanced (JEXL) |
| Expression | `catalog.annotationExists."bank.com/last-pen-test" == true` |

### Check 7 — Branch Protection Enabled
| Field | Value |
|---|---|
| Name | Branch protection enabled |
| Mode | Basic |
| Data Source | GitHub |
| Data Point | Is Branch Protected |
| Operator | Equal to |
| Value | `true` |

---

## Step 5 — Create the Scorecard

**IDP > Configure > Scorecards > Create New Scorecard**

| Field | Value |
|---|---|
| Name | Production Readiness |
| Description | Evaluates production readiness of payment services against bank standards |
| Kind | Component |
| Type | service |
| Lifecycle | production |

Add all 7 checks created in Step 4, then click **Publish Scorecard**.

> Scorecard calculations run against the **default branch** of the linked Git repository
> (`github.com/project-slug` annotation on the Component entity).

---

## Step 6 — Verify

Open **payments-service** in the Software Catalog and confirm:

- [ ] **Overview tab** — Production Readiness scorecard summary card is visible (half-width)
- [ ] **Scorecard tab** — All 7 checks listed with Pass / Fail results
- [ ] Check 7 (Branch Protection) shows FAIL with an AI-generated fix recommendation
- [ ] Checks 1–6 all show PASS

---

## Production Readiness Governance

> **Scorecard as a production gate.** Both `catalog-info.yaml` and `now.yaml` must be
> complete before a service can be deployed to production. The scorecard enforces this
> automatically — incomplete files block the deployment pipeline.

### How It Works

```
Developer fills catalog-info.yaml + now.yaml
            ↓
    Scorecard evaluates both files
            ↓
   Score < 100% → pipeline BLOCKED
   Score = 100% → production deployment APPROVED
            ↓
  now.yaml syncs to cmdb_rel_ci via RapDev CSDM as Code
            ↓
  Dependency graph visible on catalog entity detail page
```

### Scorecard Checks — catalog-info.yaml

| Check | Validates |
|---|---|
| Has SYS ID | `bank.com/sysid` annotation present |
| Has owner defined | `spec.owner` populated |
| Has TechDocs | `backstage.io/techdocs-ref` annotation present |
| Has on-call schedule | `bank.com/oncall-schedule` annotation present |
| Has data classification | `bank.com/data-classification` annotation present |
| Has SOX scope declared | `bank.com/sox-in-scope` annotation present |
| Has email DL | `bank.com/email-dl` annotation present |
| Has GitHub slug | `github.com/project-slug` annotation present |

### Scorecard Checks — now.yaml

| Check | Validates |
|---|---|
| now.yaml exists in repo | File `now.yaml` present on main branch (GitHub data source) |
| Has dependencies declared | `spec.dependsOn` populated in catalog entity |
| Has tier declared | `bank.com/tier` annotation present |
| Has team declared | `bank.com/team` annotation present |

> **Note:** now.yaml checks require GitHub connected as a scorecard data source (Step 3).
> now.yaml schema is not yet finalised — checks will be added once RapDev confirms the spec.

### Scorecard Configuration

Create a single **Production Readiness** scorecard (Step 5) with all 12 checks.
Set the pass threshold to **100%** — partial completion does not satisfy the gate.

### Pipeline Integration — Prod Deployment Gate

Add this step before the production stage in every Harness CD pipeline:

```yaml
- step:
    name: Check Production Readiness Scorecard
    type: Http
    spec:
      url: https://app.harness.io/idp/api/scorecards/<entity>/score
      method: GET
      headers:
        - key: x-api-key
          value: <+secrets.getValue("harness_api_key")>
      assertion: <+json.select("score", httpResponseBody)> == 100
    failureStrategies:
      - onFailure:
          errors: [AssertionError]
          action:
            type: Abort
```

A failing scorecard aborts the pipeline — the developer must fix `catalog-info.yaml`
or `now.yaml`, wait for the next scorecard evaluation, then re-trigger the deployment.

### What Unlocks on the Catalog Page After Both Files Pass

| Tab | What becomes visible |
|---|---|
| **Overview** | Production Readiness scorecard card — green 100%, all checks PASS |
| **Scorecard** | Full check table — catalog-info.yaml and now.yaml checks with PASS/FAIL |
| **Dependencies** | Downstream: what this system depends on (sourced from now.yaml) |
| **Dependencies** | Upstream: who depends on this system — blast radius view |

The dependency graph only populates after now.yaml syncs to cmdb_rel_ci and the
`dependsOn` edges appear in the IDP catalog. This creates a natural incentive —
developers complete both files to unblock production and unlock their full catalog page.

---

## File Reference

| File | Purpose |
|---|---|
| `catalog-info.yaml` | All 6 entity definitions (System, Component, API, 3× Resource) |
| `layout-system.yaml` | System entity page layout |
| `layout-component-service.yaml` | Component / service page layout (includes Scorecard tab) |
| `layout-resource-database.yaml` | Resource / database page layout |
| `layout-resource-messaging-queue.yaml` | Resource / messaging-queue page layout |
| `layout-resource-cache.yaml` | Resource / cache page layout |
| `preview.html` | Offline UI preview — open in browser to simulate all entity pages |

---

---

## Service Mapping Initiative

> **Critical company-wide initiative.** This section describes how the IDP catalog
> integrates with the broader service mapping programme across 19,000 repositories.

### Overview

Every repo in the organisation will have a `now.yaml` file alongside `catalog-info.yaml`.
The `now.yaml` declares the direct dependencies of that system (1 level deep).
RapDev's CSDM as Code integration reads `now.yaml` on deploy and writes the relationships
into ServiceNow's `cmdb_rel_ci` table — building a complete company-wide dependency graph.

**Goal:** If System B has an outage, immediately identify every upstream system impacted —
1 hop, 2 hops, N hops — using a recursive SQL query across `cmdb_rel_ci`.

### now.yaml — Schema Status

> **Schema is not yet finalised.** RapDev is working with the team to define the spec.
> Do not create `now.yaml` files until the schema is confirmed.

Each `now.yaml` will declare:
- The system's CMDB CI class and key attributes (tier, team, classification)
- Direct dependencies — what systems this system calls or depends on
- Relationship type (`Depends on::Used by`, `Calls::Called by`, etc.)

### cmdb_rel_ci — Relationship Table

ServiceNow stores all CI-to-CI relationships here. Each row has:

| Column | Meaning |
|---|---|
| `parent` | The system that depends on another |
| `child` | The system being depended on |
| `type` | Relationship type e.g. `Depends on::Used by` |

The `::` in the type name splits two English readings: `[parent→child]::[child→parent]`.
Always query **both** parent and child columns (UNION) — different tools insert from different ends.

### Blast Radius — Recursive SQL

To find all systems impacted by an outage in a given system:

```sql
WITH RECURSIVE blast_radius AS (
  SELECT p.name AS impacted_system, p.sys_id AS impacted_sys_id, 1 AS hops
  FROM cmdb_rel_ci r
  JOIN cmdb_ci p ON r.parent = p.sys_id
  JOIN cmdb_ci c ON r.child  = c.sys_id
  JOIN cmdb_rel_type t ON r.type = t.sys_id
  WHERE c.name = :failed_system AND t.name LIKE 'Depends on%'

  UNION ALL

  SELECT p.name, p.sys_id, br.hops + 1
  FROM cmdb_rel_ci r
  JOIN cmdb_ci p ON r.parent = p.sys_id
  JOIN cmdb_ci c ON r.child  = c.sys_id
  JOIN cmdb_rel_type t ON r.type = t.sys_id
  JOIN blast_radius br ON c.sys_id = br.impacted_sys_id
  WHERE t.name LIKE 'Depends on%' AND br.hops < 10
)
SELECT impacted_system, MIN(hops) AS blast_radius_hops
FROM blast_radius
GROUP BY impacted_system
ORDER BY blast_radius_hops;
```

### IDP Catalog — Dependency Visibility

When `dependsOn` is populated on a catalog entity, IDP automatically creates **reverse edges**:
- System A's Dependencies tab shows → what A depends on (downstream)
- System B's Dependencies tab shows → who depends on B (blast radius / upstream)

No extra configuration needed. The reverse `dependencyOf` edge is created automatically.

### Scorecard Check — Mapping Completeness

Track initiative progress across all repos:

| Field | Value |
|---|---|
| Name | Has dependency mapping |
| Mode | Advanced (JEXL) |
| Expression | `catalog.annotationExists."bank.com/sysid" == true && spec.dependsOn != null` |

This surfaces as a dashboard showing which systems have completed dependency mapping —
the key progress metric for the initiative.

### Rollout Priority

| Phase | Target | Rationale |
|---|---|---|
| 1 | Tier-1 systems | Highest blast radius, most critical for incident response |
| 2 | Tier-2 systems | ~50% graph coverage, useful for major incidents |
| 3 | All remaining repos | Near-complete graph, reliable for all outage scenarios |

> Blast radius accuracy is directly proportional to now.yaml completion rate.
> A system with no `now.yaml` appears as a leaf node — its upstream dependents are invisible.

### Current Limitation

The Harness IDP CMDB integration reads **CI field properties** only — it does not yet read
`cmdb_rel_ci` to auto-populate `dependsOn` in the catalog. Until Harness adds this support,
`dependsOn` in `catalog-info.yaml` must be maintained manually or via the Catalog Ingestion API.
Raise with Harness Support as a feature request referencing the `cmdb_rel_ci` table.

---

*Author: Rohid Dev · github.com/rohiddev*
