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

*Author: Rohid Dev · github.com/rohiddev*
