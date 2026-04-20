---
name: ServiceNow Service Mapping Initiative
description: Critical company-wide initiative - now.yaml in 19K repos for blast radius analysis via cmdb_rel_ci
type: project
---

Company-wide critical initiative to add `now.yaml` to all 19,000 repos, enabling full service dependency mapping and blast radius analysis in ServiceNow CMDB.

**Why:** Outage impact analysis — if System B fails, identify all upstream systems affected. 1-level-deep dependencies per repo, recursive traversal via SQL CTE for full blast radius.

**How to apply:** This is a key strategic initiative. Any catalog, CMDB, or dependency work should reference now.yaml as the source of truth. IDP catalog `dependsOn` fields should align with now.yaml declarations.

**Status (2026-04-19):** now.yaml schema is NOT yet finalised — RapDev is working with the team to define the spec. Do not create now.yaml examples until schema is confirmed.

**Key facts:**
- 19,000 repos in the company
- Uses RapDev CSDM as Code to sync now.yaml → ServiceNow cmdb_rel_ci
- Blast radius SQL: recursive CTE on cmdb_rel_ci traversing upstream dependents
- Tier-1 systems mandated first
- IDP catalog shows blast radius via reverse dependencyOf edges
- Harness IDP CMDB integration does not yet read cmdb_rel_ci (feature request needed)
- Scorecard check "Has dependency mapping" tracks initiative progress

**Wiki:** `obsidian/harnessidp/wiki/integrations/servicenow-service-mapping.md`
