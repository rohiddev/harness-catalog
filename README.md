# Harness IDP — Payments Platform Catalog Setup

This repo contains the catalog entity definitions and layout YAMLs for the Payments Platform.
Follow the steps below **in order** to configure everything in the Harness IDP platform.

---

## Step 1 — Register Catalog Entities

**IDP > Software Catalog > Register Existing Component**

Register the `catalog-info.yaml` file. This creates all 6 entities in one shot:

| Entity | Kind | Type |
|---|---|---|
| payments-platform | System | — |
| payments-service | Component | service |
| payments-rest-api | API | openapi |
| payments-postgres | Resource | database |
| payments-rabbitmq | Resource | messaging-queue |
| payments-redis | Resource | cache |

After registration, verify all 6 entities appear in the catalog.

---

## Step 2 — Configure Entity Layouts

**IDP > Admin > Layout > Catalog Entities**

Paste each YAML file into the corresponding layout editor. **Save after each one.**

| Layout File | Where to Paste |
|---|---|
| `layout-system.yaml` | Layout > System > Edit |
| `layout-component-service.yaml` | Layout > Component > service > Edit |
| `layout-resource-database.yaml` | Layout > Resource > database > Edit |
| `layout-resource-messaging-queue.yaml` | Layout > Resource > messaging-queue > Edit |
| `layout-resource-cache.yaml` | Layout > Resource > cache > Edit |

After saving, open each entity in the catalog and verify the correct tabs and info cards appear.

---

## Step 3 — Configure Scorecard Checks

**IDP > Configure > Scorecards > Checks tab > Create Custom Check**

Create the following checks one by one:

### Check 1 — Has Owner Defined
- **Name:** Has owner defined
- **Mode:** Basic
- **Data Source:** Catalog
- **Data Point:** Annotation Exists
- **Operator:** Equal to
- **Value:** `github.com/project-slug`

### Check 2 — Has Runbook Link
- **Name:** Has runbook link
- **Mode:** Basic
- **Data Source:** Catalog
- **Data Point:** Annotation Exists
- **Operator:** Equal to
- **Value:** `bank.com/oncall-schedule`

### Check 3 — Has TechDocs Configured
- **Name:** Has TechDocs configured
- **Mode:** Basic
- **Data Source:** Catalog
- **Data Point:** Annotation Exists
- **Operator:** Equal to
- **Value:** `backstage.io/techdocs-ref`

### Check 4 — Has On-Call Schedule
- **Name:** Has on-call schedule
- **Mode:** Basic
- **Data Source:** Catalog
- **Data Point:** Annotation Exists
- **Operator:** Equal to
- **Value:** `bank.com/oncall-schedule`

### Check 5 — Has API Definition
- **Name:** Has API definition
- **Mode:** Advanced (JEXL)
- **Expression:**
  ```
  catalog.annotationExists."github.com/project-slug" == true
  ```

### Check 6 — Pen Test Within 12 Months
- **Name:** Pen test within 12 months
- **Mode:** Advanced (JEXL)
- **Expression:**
  ```
  <+metadata.annotations['bank.com/last-pen-test']> != null
  ```

### Check 7 — Branch Protection Enabled
- **Name:** Branch protection enabled
- **Mode:** Basic
- **Data Source:** GitHub
- **Data Point:** Is Branch Protected
- **Operator:** Equal to
- **Value:** `true`

> **Prerequisite for Check 7:** GitHub must be connected as a data source.
> Go to **IDP > Configure > Scorecards > Data Sources** and connect your GitHub organisation.

---

## Step 4 — Create the Scorecard

**IDP > Configure > Scorecards > Create New Scorecard**

| Field | Value |
|---|---|
| Name | Production Readiness |
| Description | Evaluates production readiness of payment services against bank standards |
| Kind | Component |
| Type | service |
| Lifecycle | production |

**Add all 7 checks** created in Step 3, then click **Publish Scorecard**.

---

## Step 5 — Verify

Open **payments-service** in the catalog and confirm:

- [ ] Overview tab shows the **Production Readiness** score card (summary, half-width)
- [ ] **Scorecard** tab shows all 7 checks with Pass/Fail results
- [ ] Check 7 (Branch Protection) shows FAIL with an AI-generated fix recommendation

---

## Reference

| File | Purpose |
|---|---|
| `catalog-info.yaml` | All 6 entity definitions |
| `layout-system.yaml` | System entity layout |
| `layout-component-service.yaml` | Component / service layout (includes Scorecard tab) |
| `layout-resource-database.yaml` | Resource / database layout |
| `layout-resource-messaging-queue.yaml` | Resource / messaging-queue layout |
| `layout-resource-cache.yaml` | Resource / cache layout |
| `preview.html` | Offline UI preview — open in browser to simulate the entity pages |

---

*Author: Rohid Dev · github.com/rohiddev*
