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
| Branch | `main` |
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

## Step 2 — Configure Entity Layouts

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

## Step 3 — Connect GitHub as Scorecard Data Source

> Skip if GitHub is already connected as a data source.

**IDP > Configure > Scorecards > Data Sources > + Add**

Connect your GitHub organisation using the existing GitHub connector.
This is required for:
- Branch protection check
- `now.yaml` file existence check (Phase 2 — add once now.yaml schema confirmed)

All other checks use Catalog data only and work without this step.

---

## Step 4 — Create Scorecard Checks

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

#### Check 4 — Has PagerDuty
| Field | Value |
|---|---|
| Name | Has PagerDuty service |
| Data Source | Catalog |
| Data Point | Annotation Exists |
| Operator | Equal to |
| Value | `pagerduty.com/service-id` |

#### Check 5 — Has Grafana Dashboard
| Field | Value |
|---|---|
| Name | Has Grafana dashboard |
| Data Source | Catalog |
| Data Point | Annotation Exists |
| Operator | Equal to |
| Value | `grafana/dashboard-url` |

#### Check 6 — Has Datadog Dashboard
| Field | Value |
|---|---|
| Name | Has Datadog dashboard |
| Data Source | Catalog |
| Data Point | Annotation Exists |
| Operator | Equal to |
| Value | `datadoghq.com/dashboard-url` |

#### Check 7 — Has Jira Project
| Field | Value |
|---|---|
| Name | Has Jira project key |
| Data Source | Catalog |
| Data Point | Annotation Exists |
| Operator | Equal to |
| Value | `jira/project-key` |

#### Check 8 — Branch Protection Enabled
| Field | Value |
|---|---|
| Name | Branch protection enabled |
| Data Source | GitHub |
| Data Point | Is Branch Protected |
| Operator | Equal to |
| Value | `true` |

---

### now.yaml Checks (add once RapDev confirms schema)

> Do not create these until now.yaml schema is finalised. See `wiki/integrations/now-yaml-draft-schema.md`.

#### Check 9 — now.yaml Exists in Repo
| Field | Value |
|---|---|
| Name | now.yaml exists in repo |
| Data Source | GitHub |
| Data Point | File exists |
| File path | `now.yaml` |
| Branch | `main` |

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

## Step 5 — Create the Production Readiness Scorecard

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

## Step 6 — Add Production Gate to Deployment Pipelines

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

## Step 7 — Verify

Open **payments-service** in the catalog and confirm:

- [ ] **Overview tab** — Scorecard summary card visible showing 8/12 checks (4 now.yaml checks pending)
- [ ] **Dependencies tab** — Dependency graph renders, blast radius cards visible
- [ ] **API tab** — REST API, Kafka producer/consumer, FTP inbound/outbound all listed
- [ ] **Kubernetes tab** — K8s workload visible (requires Kubernetes connector)
- [ ] **CI/CD tab** — Harness pipeline runs visible
- [ ] **Scorecard tab** — All 8 catalog-info.yaml checks listed with Pass/Fail results

---

## Step 8 — Roll Out to Service Teams

Share the repo with all teams:

> **`https://github.com/rohiddev/harness-catalog`**
>
> Each team:
> 1. Copies `catalog-info.template.yaml` → saves as `catalog-info.yaml` in their repo root
> 2. Fills all `# REQUIRED` fields, removes sections that don't apply
> 3. References `catalog-info.yaml` (payments-service) as a worked example
> 4. Raises a PR — platform team imports the service into IDP once merged
> 5. Achieves 100% Production Readiness scorecard before deploying to production

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
| `catalog-info.template.yaml` | Golden template — teams copy this |
| `catalog-info.yaml` | Worked example — Payments Platform |
| `layout-component-service.yaml` | Component / service page layout |
| `layout-system.yaml` | System page layout |
| `layout-resource-*.yaml` | Resource page layouts |
| `preview.html` | Offline UI preview — open in browser |
| `wiki/integrations/` | ServiceNow, now.yaml, Phase 2 pipeline docs |

---

*Author: Rohid Dev · github.com/rohiddev*
