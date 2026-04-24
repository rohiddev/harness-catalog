# Payments Platform

> **SYSID:** `SYSID-00000` | **Team:** `team-payments` | **Tier:** `Tier 1`
> **Lifecycle:** `production`

## What This System Does

The Payments Platform is the end-to-end payments processing system responsible for payment
initiation, fraud decision integration, ledger posting and lifecycle event publishing.
It covers all payment channels — REST API, bulk FTP file exchange and real-time Kafka event
streams. It does **not** own fraud scoring, ledger accounting, or downstream notification logic;
those are provided by dependent systems consumed via APIs.

---

## System Boundary

| Entity | Type | Description |
|---|---|---|
| `payments-service` | Component | Core transaction processor — initiates, validates and settles payments |
| `payments-postgres` | Resource | Primary transactional database (PostgreSQL 14, AWS RDS Multi-AZ) |
| `payments-redis` | Resource | Session cache and idempotency key store (Redis 7, ElastiCache) |
| `payments-ftp` | Resource | FTP server for secure file exchange with third-party processors |
| `payments-kafka` | Resource | Kafka cluster for payment lifecycle event streaming |
| `payments-rest-api` | API | REST API for payment initiation, status queries and refunds |
| `payments-ftp-outbound` | API | AsyncAPI — settlement and reconciliation files to processors |
| `payments-ftp-inbound` | API | AsyncAPI — payment instruction files from processors |
| `payments-kafka-producer` | API | AsyncAPI — payment lifecycle events published to Kafka |
| `payments-kafka-consumer` | API | AsyncAPI — fraud decisions and ledger confirmations consumed |

**[View in Harness IDP](https://app.harness.io/ng/account/<account-id>/module/idp/catalog/payments_platform)**

---

## Ownership

| Role | Contact |
|---|---|
| **Team** | `team-payments` |
| **Primary On-Call** | `oncall-primary@bank.com` |
| **Secondary On-Call** | `oncall-secondary@bank.com` |
| **Cloud Architect** | `architect@bank.com` |
| **Assignment Group** | `Payments Platform Support` |
| **Support DL** | `team-payments@bank.com` |

---

## Infrastructure

| Resource | Type | Details |
|---|---|---|
| **Database** | PostgreSQL 14 | AWS RDS Multi-AZ — `us-east-1` — 35-day backup retention |
| **Cache** | Redis 7 | AWS ElastiCache — in-transit encryption |
| **Message Bus** | Kafka | SASL/SCRAM-SHA-256 auth, TLS in transit — managed by platform team |
| **FTP** | Yes | Inbound + outbound — `ftp.bank.com` |
| **Cloud** | AWS | Division: payments |
| **Kubernetes Namespace** | `payments-p` | |

---

## External Dependencies

### Upstream Services Consumed

| Service | Owner Team | Type | Notes |
|---|---|---|---|
| `auth-api` | team-identity | REST | OAuth2 token validation on every inbound request |
| `fraud-detection-api` | team-fraud | Kafka | Fraud decisions consumed to approve/block/review transactions |
| `ledger-api` | team-ledger | Kafka | Ledger confirmations consumed to complete settlement |

### Downstream Consumers

| Service | Owner Team | Type | Notes |
|---|---|---|---|
| `notification-service` | team-notifications | Kafka | Consumes `payments.events.payment-completed` |
| `audit-service` | team-audit | Kafka | Consumes all payment lifecycle events |
| `reporting-service` | team-reporting | Kafka | Consumes completed events for daily reconciliation |

### Enterprise Applications

- [x] API Connect
- [ ] Data Power (DMZ)
- [ ] Autosys
- [ ] IIB-MQ
- [ ] OcpSharedToDB
- [ ] OcpCotsToDb
- [x] EksSharedToDb

---

## Environments

| Environment | Cloud Account | Region | Notes |
|---|---|---|---|
| Dev | `<dev-account-id>` | `us-east-1` | |
| Test | `<test-account-id>` | `us-east-1` | |
| UAT | `<uat-account-id>` | `us-east-1` | |
| Prod | `<prod-account-id>` | `us-east-1` | |

---

## CI/CD

| Item | Value |
|---|---|
| **Harness Project** | `SYSID-00000` |
| **Pipeline — Build** | `payments-ci` |
| **Pipeline — Deploy** | `payments-deploy` |
| **Deploy Strategy** | Blue-Green |
| **Approval Required** | Yes — 4-eyes |

---

## Observability

| Tool | Link |
|---|---|
| **Datadog Dashboard** | [app.datadoghq.com/dashboard/payments](https://app.datadoghq.com/dashboard/payments) |
| **Alerts** | [oncall.bank.com/payments](https://oncall.bank.com/payments) |
| **SonarQube** | `payments-service` |

---

## Security & Compliance

| Item | Value |
|---|---|
| **ISBIA Risk Rating** | `Critical` |
| **Data Classification** | `Restricted` |
| **SOX Relevant** | `Yes` |
| **Public Facing** | `Yes` |
| **AD Group — Primary Approver** | `corp/XXXXXXX` |
| **AD Group — Secondary Approver** | `corp/XXXXXXX` |

---

## Links

| Resource | URL |
|---|---|
| Confluence Space | [confluence.bank.com/display/PAY](https://confluence.bank.com/display/PAY) |
| Architecture Diagram | [confluence.bank.com/display/PAY/architecture](https://confluence.bank.com/display/PAY/architecture) |
| Jira Board | [jira.bank.com/projects/PAY](https://jira.bank.com/projects/PAY) |
| Platform Runbook | [runbooks.bank.com/payments-platform](https://runbooks.bank.com/payments-platform) |
| Prod Readiness Scorecard | [View in Harness IDP](https://app.harness.io/ng/account/<account-id>/module/idp/scorecards) |
