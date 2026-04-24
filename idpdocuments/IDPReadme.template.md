# <Service Display Name>

> **SYSID:** `SYSID-NNNNN` | **Team:** `<team-name>` | **Tier:** `Tier 1 / 2 / 3`
> **Lifecycle:** `production / experimental / deprecated`

## What This Service Does

<!-- REQUIRED: One paragraph. What does this service do? What is it responsible for?
     What does it explicitly NOT do? Who are the consumers? -->

---

## System

<!-- REQUIRED: Which platform / system does this service belong to?
     Must match the System entity registered in Harness IDP.
     Example: Payments Platform (system:account/payments_platform) -->

**Platform:** `<platform-name>`
**IDP Catalog:** [View in Harness IDP](https://app.harness.io/ng/account/<account-id>/module/idp/catalog/<service-name>)

---

## Ownership

| Role | Contact |
|---|---|
| **Team** | `<team-name>` |
| **Primary On-Call** | `<nid>@bank.com` |
| **Secondary On-Call** | `<nid>@bank.com` |
| **Cloud Architect** | `<nid>@bank.com` |
| **Risk Partner (NID)** | `<NID>` |
| **Assignment Group** | `<ServiceNow assignment group>` |
| **Support DL** | `<support-team@bank.com>` |

---

## Repository

| Item | Value |
|---|---|
| **Repo** | `github.com/<org>/<repo>` |
| **Primary Branch** | `master` |
| **Language** | `<Java / Python / Go / .NET / Node.js>` |
| **Framework** | `<Spring Boot / FastAPI / Gin / Express / Other>` |
| **Build Tool** | `<Maven / Gradle / pip / go build / Other>` |
| **Bitbucket Project Key** | `<PROJ>` |

---

## Infrastructure

| Resource | Type | Details |
|---|---|---|
| **Database** | `<PostgreSQL / MySQL / MSSQL / MongoDB>` | `<RDS / Azure SQL / On-prem>` — `<region>` |
| **Cache** | `<Redis / Memcached / None>` | `<ElastiCache / Azure Cache>` |
| **Message Bus** | `<Kafka / MQ / None>` | Topics: `<topic-names>` |
| **Object Storage** | `<S3 / Azure Blob / None>` | Bucket: `<bucket-name>` |
| **FTP** | `<Yes / No>` | `<inbound / outbound / both>` |
| **Cloud** | `<AWS / Azure>` | Division: `<division>` |
| **Kubernetes Namespace** | `<namespace / N/A>` | |
| **VPC** | `<internal-vpc>` / `<ingress-vpc>` | |

---

## Dependencies

### Services This Depends On

<!-- REQUIRED: List every upstream service this calls at runtime.
     Reference the IDP catalog entry where possible. -->

| Service | Owner Team | Type | Notes |
|---|---|---|---|
| `<upstream-service-name>` | `<team>` | REST / Kafka / FTP | `<what it uses from this service>` |
| `<upstream-service-name>` | `<team>` | REST / Kafka / FTP | `<what it uses from this service>` |

### Services That Depend On This

<!-- OPTIONAL: List known downstream consumers of this service. -->

| Service | Owner Team | Type | Notes |
|---|---|---|---|
| `<downstream-service-name>` | `<team>` | REST / Kafka / FTP | `<what they consume>` |

### Enterprise Applications

<!-- Check all that apply -->

- [ ] API Connect
- [ ] Data Power (DMZ)
- [ ] Autosys
- [ ] IIB-MQ
- [ ] OcpSharedToDB
- [ ] OcpCotsToDb
- [ ] EksSharedToDb

---

## APIs Exposed

<!-- OPTIONAL: Complete if this service exposes an API.
     Full OpenAPI / AsyncAPI definition lives in catalog-info.yaml. -->

| API | Type | Endpoint / Topic | Auth |
|---|---|---|---|
| `<api-name>` | REST / Kafka / FTP | `<url or topic>` | `<OAuth2 / mTLS / SASL>` |

---

## APIs Consumed

| API | Provider | Type | Topic / Endpoint |
|---|---|---|---|
| `<api-name>` | `<service-name>` | REST / Kafka / FTP | `<url or topic>` |

---

## Environments

| Environment | Bitbucket Project Key | Cloud Account | Region | Notes |
|---|---|---|---|---|
| Dev | `<key>` | `<account-id>` | `us-east-1` | |
| Test | `<key>` | `<account-id>` | `us-east-1` | |
| UAT | `<key>` | `<account-id>` | `us-east-1` | |
| Prod | `<key>` | `<account-id>` | `us-east-1` | |

---

## CI/CD

| Item | Value |
|---|---|
| **Harness Project** | `SYSID-NNNNN` |
| **Pipeline — Build** | `<pipeline-id>` |
| **Pipeline — Deploy** | `<pipeline-id>` |
| **Deploy Strategy** | `<Rolling / Blue-Green / Canary>` |
| **Approval Required** | `<Yes — 4-eyes / No>` |

---

## Observability

| Tool | Link |
|---|---|
| **Datadog Dashboard** | `https://app.datadoghq.com/dashboard/<id>` |
| **Log Query** | `<Datadog log filter or Splunk query>` |
| **Alerts** | `<PagerDuty service ID or oncall link>` |
| **SonarQube** | `<sonar-project-key>` |

---

## Security & Compliance

| Item | Value |
|---|---|
| **ISBIA Risk Rating** | `Critical / Major / Significant / Important / Minor / Unknown` |
| **Data Classification** | `Public / Private / Restricted / Confidential / Internal` |
| **SOX Relevant** | `Yes / No` |
| **Public Facing** | `Yes / No / Both` |
| **AD Group — Primary Approver** | `corp/XXXXXXX` |
| **AD Group — Secondary Approver** | `corp/XXXXXXX` |
| **Delegators (EmployeeId)** | `D108991, D109002` |

---

## Runbook

<!-- OPTIONAL: Inline quick-reference. Link to full runbook for detail. -->

### Common Failure Modes

| Symptom | Likely Cause | First Action |
|---|---|---|
| `<symptom>` | `<cause>` | `<action>` |
| `<symptom>` | `<cause>` | `<action>` |

### Restart Procedure

```bash
# <Add kubectl / harness CLI / restart command here>
```

### Rollback Procedure

```bash
# <Add rollback steps here>
```

**Full Runbook:** [runbooks.bank.com/<service-name>](https://runbooks.bank.com/<service-name>)

---

## Local Development

```bash
# Clone
git clone https://github.com/<org>/<repo>.git

# Install dependencies
<npm install / mvn install / pip install -r requirements.txt / go mod download>

# Run locally
<npm start / mvn spring-boot:run / python app.py / go run main.go>

# Run tests
<npm test / mvn test / pytest / go test ./...>
```

**Prerequisites:**
- `<JDK 17 / Python 3.11 / Go 1.21 / Node 20>`
- `<Docker>`
- `<Access to dev secrets — request via IDP Self-Service>`

---

## Links

| Resource | URL |
|---|---|
| Confluence Space | `https://confluence.bank.com/display/<PROJ>` |
| Jira Board | `https://jira.bank.com/projects/<PROJ>` |
| Architecture Diagram | `<link>` |
| Data Flow Diagram | `<link>` |
| Prod Readiness Scorecard | [View in Harness IDP](https://app.harness.io/ng/account/<account-id>/module/idp/scorecards) |
