# Update Catalog Entry — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Workflow + pipeline that lets a team add Resources and APIs to an existing Component's `catalog-info.yaml`, with a Go validator ensuring correctness before commit.

**Architecture:** Three independent deliverables built in order: Go validator binary → IDP workflow YAML → Harness pipeline YAML. The pipeline calls the validator before committing — if validation fails, nothing is committed and the requestor receives an email.

**Tech Stack:** Go 1.21 (gopkg.in/yaml.v3), Harness IDP 2.0 workflow YAML (harness.io/v1), Harness pipeline YAML, GitHub PAT via HashiCorp Vault.

**Spec:** `docs/superpowers/specs/2026-04-27-update-catalog-entry-design.md`

**Key confirmed patterns (memory):**
- Workflow: `SelectFieldFromApi`, `dependencies/oneOf`, `type: array` with object items
- Pipeline: Http steps preferred, retry 2×30s, `delegateSelectors: act-delegate-k8s-ephub-p2r1`, ASCII step names only
- Secrets: HashiCorp Vault only — never GitHub org secrets
- Notifications: Email API (`POST /send`) — no Slack
- Entity format: `identifier`/`name`/`type`/`owner` at top level, `owner: group:account/<team>`, `spec.system` is a list

---

## File Structure

```
harness-catalog/
├── tools/
│   ├── catalog-validator/
│   │   └── main.go              # Validates catalog-info.yaml before commit
│   └── catalog-appender/
│       └── main.go              # Reads existing catalog-info.yaml, appends new entities
├── workflows/
│   └── update-catalog-entry.yaml    # IDP workflow form
└── pipelines/
    └── update-catalog-entry-pipeline.yaml  # Harness pipeline
```

---

### Task 1: Go Catalog Validator

**Files:**
- Create: `harness-catalog/tools/catalog-validator/main.go`

The validator reads a `catalog-info.yaml` (multi-document YAML) and checks:
1. Every document has `apiVersion: harness.io/v1`
2. Every document has top-level `identifier`, `name`, `kind`, `owner`
3. `owner` matches `^group:account/.+`
4. `kind` is one of: System, Component, Resource, API
5. `spec.system` entries (if present) match `^system:account/.+`
6. `spec.dependsOn` entries (if present) match `^(resource|component|api):account/.+`
7. No duplicate `identifier` values across all documents
8. All `name` values match `^[a-z][a-z0-9-]+$`

Exit 0 = valid. Exit 1 = prints errors, one per line.

- [ ] **Step 1: Create the validator**

```go
// harness-catalog/tools/catalog-validator/main.go
package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	nameRegex  = regexp.MustCompile(`^[a-z][a-z0-9-]+$`)
	ownerRegex = regexp.MustCompile(`^group:account/.+`)
	refRegex   = regexp.MustCompile(`^(system|resource|component|api):account/.+`)
	validKinds = map[string]bool{"System": true, "Component": true, "Resource": true, "API": true}
)

type Entity struct {
	APIVersion string                 `yaml:"apiVersion"`
	Kind       string                 `yaml:"kind"`
	Identifier string                 `yaml:"identifier"`
	Name       string                 `yaml:"name"`
	Owner      string                 `yaml:"owner"`
	Spec       map[string]interface{} `yaml:"spec"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: catalog-validator <path-to-catalog-info.yaml>")
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read file: %v\n", err)
		os.Exit(1)
	}

	var errors []string
	seenIdentifiers := map[string]bool{}

	decoder := yaml.NewDecoder(strings.NewReader(string(data)))
	docIndex := 0
	for {
		var entity Entity
		err := decoder.Decode(&entity)
		if err != nil {
			break
		}
		if entity.APIVersion == "" && entity.Kind == "" {
			docIndex++
			continue
		}
		docIndex++
		prefix := fmt.Sprintf("doc[%d]", docIndex)

		if entity.APIVersion != "harness.io/v1" {
			errors = append(errors, fmt.Sprintf("%s: apiVersion must be harness.io/v1, got %q", prefix, entity.APIVersion))
		}
		if !validKinds[entity.Kind] {
			errors = append(errors, fmt.Sprintf("%s: kind must be System|Component|Resource|API, got %q", prefix, entity.Kind))
		}
		if entity.Identifier == "" {
			errors = append(errors, fmt.Sprintf("%s: missing top-level identifier", prefix))
		} else if seenIdentifiers[entity.Identifier] {
			errors = append(errors, fmt.Sprintf("%s: duplicate identifier %q", prefix, entity.Identifier))
		} else {
			seenIdentifiers[entity.Identifier] = true
		}
		if entity.Name == "" {
			errors = append(errors, fmt.Sprintf("%s: missing top-level name", prefix))
		} else if !nameRegex.MatchString(entity.Name) {
			errors = append(errors, fmt.Sprintf("%s: name %q must match ^[a-z][a-z0-9-]+$", prefix, entity.Name))
		}
		if entity.Owner == "" {
			errors = append(errors, fmt.Sprintf("%s: missing top-level owner", prefix))
		} else if !ownerRegex.MatchString(entity.Owner) {
			errors = append(errors, fmt.Sprintf("%s: owner %q must be group:account/<id>", prefix, entity.Owner))
		}
		if entity.Spec != nil {
			if sys, ok := entity.Spec["system"]; ok {
				refs, ok := sys.([]interface{})
				if !ok {
					errors = append(errors, fmt.Sprintf("%s: spec.system must be a list", prefix))
				} else {
					for _, r := range refs {
						s, _ := r.(string)
						if !refRegex.MatchString(s) {
							errors = append(errors, fmt.Sprintf("%s: spec.system entry %q must be system:account/<id>", prefix, s))
						}
					}
				}
			}
			if deps, ok := entity.Spec["dependsOn"]; ok {
				refs, ok := deps.([]interface{})
				if !ok {
					errors = append(errors, fmt.Sprintf("%s: spec.dependsOn must be a list", prefix))
				} else {
					for _, r := range refs {
						s, _ := r.(string)
						if !refRegex.MatchString(s) {
							errors = append(errors, fmt.Sprintf("%s: spec.dependsOn entry %q must be resource:account/<id>", prefix, s))
						}
					}
				}
			}
		}
	}

	if len(errors) > 0 {
		for _, e := range errors {
			fmt.Println(e)
		}
		os.Exit(1)
	}
	fmt.Println("catalog-info.yaml is valid")
}
```

- [ ] **Step 2: Initialise Go module**

```bash
cd /Users/chandinidev/Documents/rohid-code/harness-catalog/tools/catalog-validator
go mod init catalog-validator
go get gopkg.in/yaml.v3
```

- [ ] **Step 3: Smoke test — valid file**

```bash
cd /Users/chandinidev/Documents/rohid-code/harness-catalog/tools/catalog-validator
go run main.go ../../catalog-info.yaml
```

Expected output: `catalog-info.yaml is valid`
If it fails: the payments platform catalog-info.yaml may use `group:default/` — that's expected at this point, fix the template file first if needed.

- [ ] **Step 4: Smoke test — invalid file**

Create `/tmp/bad-catalog.yaml`:
```yaml
apiVersion: harness.io/v1
kind: Component
name: My Bad Service
owner: team-payments
```

Run:
```bash
go run main.go /tmp/bad-catalog.yaml
```

Expected output (errors, one per line):
```
doc[1]: missing top-level identifier
doc[1]: name "My Bad Service" must match ^[a-z][a-z0-9-]+$
doc[1]: owner "team-payments" must be group:account/<id>
```

- [ ] **Step 5: Commit**

```bash
cd /Users/chandinidev/Documents/rohid-code/harness-catalog
git add tools/catalog-validator/
git commit -m "feat: add catalog-info.yaml validator"
```

---

### Task 2: Go Catalog Appender

**Files:**
- Create: `harness-catalog/tools/catalog-appender/main.go`

Reads environment variables (set by pipeline) containing JSON arrays, reads the existing `catalog-info.yaml`, appends new entity documents, writes the file back. The pipeline sets env vars from trigger payload.

Input env vars:
- `CATALOG_FILE` — path to catalog-info.yaml
- `COMPONENT_IDENTIFIER` — e.g. `payments_service`
- `SYSID` — e.g. `SYSID-00123`
- `SYSTEM_IDENTIFIER` — e.g. `payments_platform`
- `OWNER_GROUP` — e.g. `team-payments`
- `UPDATE_SYSTEM` — `"true"` or `"false"`
- `SYSTEM_TITLE`, `SYSTEM_DESCRIPTION`, `SYSTEM_DOMAIN`, `SYSTEM_CONFLUENCE`, `SYSTEM_RUNBOOK`, `SYSTEM_TAGS`
- `DATABASES` — JSON: `[{"name":"orders-db","engine":"postgresql","description":"..."}]`
- `CACHES` — JSON: `[{"name":"orders-cache","engine":"redis","description":"..."}]`
- `MESSAGING` — JSON: `[{"name":"orders-kafka","type":"kafka","description":"..."}]`
- `REST_APIS` — JSON: `[{"name":"orders-rest-api","endpoint":"https://api.bank.com/orders","description":"..."}]`
- `ASYNC_PRODUCERS` — JSON: `[{"name":"orders-kafka-producer","protocol":"kafka","topic":"orders.events.order-created","description":"..."}]`
- `ASYNC_CONSUMERS` — JSON: `[{"name":"orders-kafka-consumer","protocol":"kafka","topic":"fraud.events.fraud-decision","consumerGroup":"orders-fraud-consumer","description":"..."}]`

- [ ] **Step 1: Create the appender**

```go
// harness-catalog/tools/catalog-appender/main.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type DBEntry struct {
	Name        string `json:"name"`
	Engine      string `json:"engine"`
	Description string `json:"description"`
}

type CacheEntry struct {
	Name        string `json:"name"`
	Engine      string `json:"engine"`
	Description string `json:"description"`
}

type MessagingEntry struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type RestAPIEntry struct {
	Name        string `json:"name"`
	Endpoint    string `json:"endpoint"`
	Description string `json:"description"`
}

type AsyncProducerEntry struct {
	Name        string `json:"name"`
	Protocol    string `json:"protocol"`
	Topic       string `json:"topic"`
	Description string `json:"description"`
}

type AsyncConsumerEntry struct {
	Name          string `json:"name"`
	Protocol      string `json:"protocol"`
	Topic         string `json:"topic"`
	ConsumerGroup string `json:"consumerGroup"`
	Description   string `json:"description"`
}

func env(key string) string { return os.Getenv(key) }

func parseJSON[T any](envKey string) []T {
	raw := env(envKey)
	if raw == "" || raw == "[]" {
		return nil
	}
	var result []T
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse %s: %v\n", envKey, err)
		os.Exit(1)
	}
	return result
}

func main() {
	catalogFile := env("CATALOG_FILE")
	if catalogFile == "" {
		catalogFile = "catalog-info.yaml"
	}

	owner := fmt.Sprintf("group:account/%s", env("OWNER_GROUP"))
	sysid := env("SYSID")
	systemIdentifier := env("SYSTEM_IDENTIFIER")
	systemRef := fmt.Sprintf("system:account/%s", systemIdentifier)

	existing, err := os.ReadFile(catalogFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read %s: %v\n", catalogFile, err)
		os.Exit(1)
	}

	var newDocs []string

	// System update
	if env("UPDATE_SYSTEM") == "true" {
		tags := []string{}
		if t := env("SYSTEM_TAGS"); t != "" {
			for _, tag := range strings.Split(t, ",") {
				tags = append(tags, fmt.Sprintf("    - %s", strings.TrimSpace(tag)))
			}
		}
		doc := fmt.Sprintf(`---
apiVersion: harness.io/v1
kind: System
identifier: %s
name: %s
type: system
owner: %s
metadata:
  title: %s
  description: >
    %s
  annotations:
    bank.com/sysid: %s
    jira/project-key: ""
  tags:
%s
  links:
    - url: %s
      title: Confluence Space
      icon: docs
    - url: %s
      title: Platform Runbook
      icon: book
spec:
  domain: %s
`,
			systemIdentifier,
			systemIdentifier,
			owner,
			env("SYSTEM_TITLE"),
			env("SYSTEM_DESCRIPTION"),
			sysid,
			strings.Join(tags, "\n"),
			env("SYSTEM_CONFLUENCE"),
			env("SYSTEM_RUNBOOK"),
			env("SYSTEM_DOMAIN"),
		)
		newDocs = append(newDocs, doc)
	}

	// Databases
	for _, db := range parseJSON[DBEntry]("DATABASES") {
		identifier := strings.ReplaceAll(db.Name, "-", "_")
		doc := fmt.Sprintf(`---
apiVersion: harness.io/v1
kind: Resource
identifier: %s
name: %s
type: database
owner: %s
metadata:
  title: %s Database
  description: >
    %s
  annotations:
    bank.com/sysid: %s
  tags:
    - %s
    - database
spec:
  system:
    - %s
`, identifier, db.Name, owner, db.Name, db.Description, sysid, db.Engine, systemRef)
		newDocs = append(newDocs, doc)
	}

	// Caches
	for _, c := range parseJSON[CacheEntry]("CACHES") {
		identifier := strings.ReplaceAll(c.Name, "-", "_")
		doc := fmt.Sprintf(`---
apiVersion: harness.io/v1
kind: Resource
identifier: %s
name: %s
type: cache
owner: %s
metadata:
  title: %s Cache
  description: >
    %s
  annotations:
    bank.com/sysid: %s
  tags:
    - %s
    - cache
spec:
  system:
    - %s
`, identifier, c.Name, owner, c.Name, c.Description, sysid, c.Engine, systemRef)
		newDocs = append(newDocs, doc)
	}

	// Messaging
	for _, m := range parseJSON[MessagingEntry]("MESSAGING") {
		identifier := strings.ReplaceAll(m.Name, "-", "_")
		doc := fmt.Sprintf(`---
apiVersion: harness.io/v1
kind: Resource
identifier: %s
name: %s
type: messaging-queue
owner: %s
metadata:
  title: %s
  description: >
    %s
  annotations:
    bank.com/sysid: %s
  tags:
    - %s
    - messaging
spec:
  system:
    - %s
`, identifier, m.Name, owner, m.Name, m.Description, sysid, m.Type, systemRef)
		newDocs = append(newDocs, doc)
	}

	// REST APIs
	for _, r := range parseJSON[RestAPIEntry]("REST_APIS") {
		identifier := strings.ReplaceAll(r.Name, "-", "_")
		doc := fmt.Sprintf(`---
apiVersion: harness.io/v1
kind: API
identifier: %s
name: %s
type: openapi
owner: %s
metadata:
  title: %s
  description: >
    %s
  annotations:
    bank.com/sysid: %s
  tags:
    - rest
    - openapi
spec:
  type: openapi
  lifecycle: production
  system:
    - %s
  definition: |
    openapi: "3.0.0"
    info:
      title: %s
      version: "1.0.0"
    servers:
      - url: %s
    paths: {}
`, identifier, r.Name, owner, r.Name, r.Description, sysid, systemRef, r.Name, r.Endpoint)
		newDocs = append(newDocs, doc)
	}

	// Async Producers
	for _, p := range parseJSON[AsyncProducerEntry]("ASYNC_PRODUCERS") {
		identifier := strings.ReplaceAll(p.Name, "-", "_")
		doc := fmt.Sprintf(`---
apiVersion: harness.io/v1
kind: API
identifier: %s
name: %s
type: asyncapi
owner: %s
metadata:
  title: %s
  description: >
    %s
  annotations:
    bank.com/sysid: %s
  tags:
    - %s
    - asyncapi
    - producer
spec:
  type: asyncapi
  lifecycle: production
  system:
    - %s
  definition: |
    asyncapi: "2.6.0"
    info:
      title: %s
      version: "1.0.0"
    channels:
      %s:
        publish:
          summary: %s
`, identifier, p.Name, owner, p.Name, p.Description, sysid, p.Protocol, systemRef, p.Name, p.Topic, p.Description)
		newDocs = append(newDocs, doc)
	}

	// Async Consumers
	for _, c := range parseJSON[AsyncConsumerEntry]("ASYNC_CONSUMERS") {
		identifier := strings.ReplaceAll(c.Name, "-", "_")
		groupLine := ""
		if c.ConsumerGroup != "" {
			groupLine = fmt.Sprintf("              groupId: %s\n", c.ConsumerGroup)
		}
		doc := fmt.Sprintf(`---
apiVersion: harness.io/v1
kind: API
identifier: %s
name: %s
type: asyncapi
owner: %s
metadata:
  title: %s
  description: >
    %s
  annotations:
    bank.com/sysid: %s
  tags:
    - %s
    - asyncapi
    - consumer
spec:
  type: asyncapi
  lifecycle: production
  system:
    - %s
  definition: |
    asyncapi: "2.6.0"
    info:
      title: %s
      version: "1.0.0"
    channels:
      %s:
        subscribe:
          summary: %s
          bindings:
            %s:
%s`, identifier, c.Name, owner, c.Name, c.Description, sysid, c.Protocol, systemRef, c.Name, c.Topic, c.Description, c.Protocol, groupLine)
		newDocs = append(newDocs, doc)
	}

	if len(newDocs) == 0 {
		fmt.Println("nothing to append")
		return
	}

	f, err := os.OpenFile(catalogFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot open %s for append: %v\n", catalogFile, err)
		os.Exit(1)
	}
	defer f.Close()

	for _, doc := range newDocs {
		if _, err := f.WriteString("\n" + doc); err != nil {
			fmt.Fprintf(os.Stderr, "write error: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("appended %d entities to %s\n", len(newDocs), catalogFile)
}
```

- [ ] **Step 2: Initialise Go module**

```bash
cd /Users/chandinidev/Documents/rohid-code/harness-catalog/tools/catalog-appender
go mod init catalog-appender
go get gopkg.in/yaml.v3
```

- [ ] **Step 3: Smoke test**

```bash
cd /Users/chandinidev/Documents/rohid-code/harness-catalog/tools/catalog-appender
cp ../../entities/component-payments-service.yaml /tmp/test-catalog.yaml

CATALOG_FILE=/tmp/test-catalog.yaml \
OWNER_GROUP=team-payments \
SYSID=SYSID-00000 \
SYSTEM_IDENTIFIER=payments_platform \
DATABASES='[{"name":"orders-db","engine":"postgresql","description":"Primary transactional store"}]' \
CACHES='[]' MESSAGING='[]' REST_APIS='[]' ASYNC_PRODUCERS='[]' ASYNC_CONSUMERS='[]' UPDATE_SYSTEM=false \
go run main.go

cat /tmp/test-catalog.yaml
```

Expected: file ends with a new `---` Resource block for `orders-db`.

- [ ] **Step 4: Commit**

```bash
cd /Users/chandinidev/Documents/rohid-code/harness-catalog
git add tools/catalog-appender/
git commit -m "feat: add catalog-info.yaml appender"
```

---

### Task 3: Workflow YAML

**Files:**
- Create: `harness-catalog/workflows/update-catalog-entry.yaml`

Key decisions:
- Component picker: `SelectFieldFromApi` pulling from Harness Catalog proxy
- System toggle: `dependencies/oneOf` on `updateSystem` enum (Yes/No)
- Arrays: `type: array` with object `items` — each item has its own fields
- SYSID + system auto-filled via `ContextViewer` + `setContextData` from Component picker
- All arrays serialized to JSON in `trigger:harness-custom-pipeline` payload using `{{ parameters.field | dump }}`

- [ ] **Step 1: Create the workflow**

```yaml
# harness-catalog/workflows/update-catalog-entry.yaml
apiVersion: harness.io/v1
kind: Workflow
name: update-catalog-entry
identifier: update_catalog_entry
type: service
owner: platform_engineering
orgIdentifier: default
projectIdentifier: cfg-IDP-project
icon: scaffolder
metadata:
  description: >
    Add Resources (database, cache, messaging) and APIs (REST, async) to an existing
    service Component, and optionally update the parent System — all committed directly
    to master with schema validation guardrails.
  tags:
    - catalog
    - platform
    - self-service

spec:
  parameters:
    - title: Select Service
      required:
        - componentIdentifier
      properties:
        componentIdentifier:
          title: Service Component
          type: string
          description: "Select the service you want to update"
          ui:field: SelectFieldFromApi
          ui:options:
            path: proxy/catalog-api/v1/entities?filter=kind=Component
            valueSelector: identifier
            labelSelector: metadata.title
          ui:autofocus: true

        sysid:
          title: SYSID
          type: string
          ui:field: ContextViewer
          readonly: true
          getContextData: "{{formContext.sysid}}"

        systemIdentifier:
          title: System
          type: string
          ui:field: ContextViewer
          readonly: true
          getContextData: "{{formContext.systemIdentifier}}"

        ownerGroup:
          title: Owner Group
          type: string
          ui:field: ContextViewer
          readonly: true
          getContextData: "{{formContext.ownerGroup}}"

        repoSlug:
          title: Repository
          type: string
          ui:field: ContextViewer
          readonly: true
          getContextData: "{{formContext.repoSlug}}"

    - title: System (optional)
      properties:
        updateSystem:
          title: Update System information?
          type: string
          enum:
            - "No"
            - "Yes"
          default: "No"

      dependencies:
        updateSystem:
          oneOf:
            - properties:
                updateSystem:
                  const: "No"
            - properties:
                updateSystem:
                  const: "Yes"
                systemTitle:
                  title: System Title
                  type: string
                  description: "Display name — e.g. Payments Platform"
                systemDescription:
                  title: System Description
                  type: string
                  ui:widget: textarea
                  ui:options:
                    rows: 3
                systemDomain:
                  title: Domain
                  type: string
                  enum:
                    - banking
                    - lending
                    - risk
                    - payments
                    - operations
                    - data
                systemConfluence:
                  title: Confluence URL
                  type: string
                  description: "https://confluence.bank.com/display/PROJ"
                systemRunbook:
                  title: Runbook URL
                  type: string
                  description: "https://runbooks.bank.com/platform-name"
                systemTags:
                  title: Tags
                  type: string
                  description: "Comma-separated — e.g. payments,tier-1"
              required:
                - systemTitle
                - systemDescription
                - systemDomain

    - title: Resources
      properties:
        databases:
          title: Databases
          type: array
          description: "Add one entry per database this service owns"
          items:
            type: object
            required:
              - name
              - engine
              - description
            properties:
              name:
                title: Resource Name
                type: string
                pattern: "^[a-z][a-z0-9-]+$"
                description: "kebab-case — e.g. orders-db"
                ui:placeholder: "Resource Name - e.g. orders-db"
              engine:
                title: Engine
                type: string
                enum:
                  - postgresql
                  - mysql
                  - oracle
                  - mssql
              description:
                title: Description
                type: string
                description: "Engine version, hosting, HA setup"
                ui:widget: textarea
                ui:options:
                  rows: 2

        caches:
          title: Caches
          type: array
          description: "Add one entry per cache this service owns"
          items:
            type: object
            required:
              - name
              - engine
              - description
            properties:
              name:
                title: Resource Name
                type: string
                pattern: "^[a-z][a-z0-9-]+$"
                description: "kebab-case — e.g. orders-cache"
                ui:placeholder: "Resource Name - e.g. orders-cache"
              engine:
                title: Engine
                type: string
                enum:
                  - redis
                  - memcached
              description:
                title: Description
                type: string
                ui:widget: textarea
                ui:options:
                  rows: 2

        messaging:
          title: Messaging
          type: array
          description: "Add one entry per messaging resource this service owns"
          items:
            type: object
            required:
              - name
              - type
              - description
            properties:
              name:
                title: Resource Name
                type: string
                pattern: "^[a-z][a-z0-9-]+$"
                description: "kebab-case — e.g. orders-kafka"
                ui:placeholder: "Resource Name - e.g. orders-kafka"
              type:
                title: Type
                type: string
                enum:
                  - kafka
                  - ibm-mq
                  - rabbitmq
                  - ftp-server
              description:
                title: Description
                type: string
                ui:widget: textarea
                ui:options:
                  rows: 2

    - title: APIs
      properties:
        restApis:
          title: REST APIs
          type: array
          description: "Add one entry per REST API this service exposes"
          items:
            type: object
            required:
              - name
              - endpoint
              - description
            properties:
              name:
                title: API Name
                type: string
                pattern: "^[a-z][a-z0-9-]+$"
                description: "kebab-case — e.g. orders-rest-api"
                ui:placeholder: "API Name - e.g. orders-rest-api"
              endpoint:
                title: Base URL
                type: string
                description: "https://api.bank.com/service-name"
                ui:placeholder: "Base URL - e.g. https://api.bank.com/orders"
              description:
                title: Description
                type: string
                ui:widget: textarea
                ui:options:
                  rows: 2

        asyncProducers:
          title: Async Producers
          type: array
          description: "Add one entry per topic or FTP path this service publishes to"
          items:
            type: object
            required:
              - name
              - protocol
              - topic
              - description
            properties:
              name:
                title: API Name
                type: string
                pattern: "^[a-z][a-z0-9-]+$"
                description: "kebab-case — e.g. orders-kafka-producer"
                ui:placeholder: "API Name - e.g. orders-kafka-producer"
              protocol:
                title: Protocol
                type: string
                enum:
                  - kafka
                  - ftp
              topic:
                title: Topic / Path
                type: string
                description: "Kafka topic or FTP path"
                ui:placeholder: "Topic / Path - e.g. orders.events.order-created"
              description:
                title: Description
                type: string
                ui:widget: textarea
                ui:options:
                  rows: 2

        asyncConsumers:
          title: Async Consumers
          type: array
          description: "Add one entry per topic or FTP path this service subscribes to"
          items:
            type: object
            required:
              - name
              - protocol
              - topic
              - description
            properties:
              name:
                title: API Name
                type: string
                pattern: "^[a-z][a-z0-9-]+$"
                description: "kebab-case — e.g. orders-kafka-consumer"
                ui:placeholder: "API Name - e.g. orders-kafka-consumer"
              protocol:
                title: Protocol
                type: string
                enum:
                  - kafka
                  - ftp
              topic:
                title: Topic / Path
                type: string
                description: "Kafka topic or FTP path this service reads"
                ui:placeholder: "Topic / Path - e.g. fraud.events.fraud-decision"
              consumerGroup:
                title: Consumer Group
                type: string
                description: "Kafka consumer group ID (leave blank for FTP)"
                ui:placeholder: "Consumer Group - e.g. orders-fraud-consumer"
              description:
                title: Description
                type: string
                ui:widget: textarea
                ui:options:
                  rows: 2

  steps:
    - id: log_inputs
      name: Log Inputs
      action: debug.log
      input:
        message: >
          component={{ parameters.componentIdentifier }},
          sysid={{ parameters.sysid }},
          system={{ parameters.systemIdentifier }},
          updateSystem={{ parameters.updateSystem }}

    - id: trigger_pipeline
      name: Trigger Update Catalog Pipeline
      action: trigger:harness-custom-pipeline
      input:
        pipeline: update_catalog_entry_pipeline
        payload:
          componentIdentifier: "{{ parameters.componentIdentifier }}"
          repoSlug: "{{ parameters.repoSlug }}"
          sysid: "{{ parameters.sysid }}"
          systemIdentifier: "{{ parameters.systemIdentifier }}"
          ownerGroup: "{{ parameters.ownerGroup }}"
          updateSystem: "{{ parameters.updateSystem }}"
          systemTitle: "{{ parameters.systemTitle or '' }}"
          systemDescription: "{{ parameters.systemDescription or '' }}"
          systemDomain: "{{ parameters.systemDomain or '' }}"
          systemConfluence: "{{ parameters.systemConfluence or '' }}"
          systemRunbook: "{{ parameters.systemRunbook or '' }}"
          systemTags: "{{ parameters.systemTags or '' }}"
          databases: "{{ parameters.databases | dump }}"
          caches: "{{ parameters.caches | dump }}"
          messaging: "{{ parameters.messaging | dump }}"
          restApis: "{{ parameters.restApis | dump }}"
          asyncProducers: "{{ parameters.asyncProducers | dump }}"
          asyncConsumers: "{{ parameters.asyncConsumers | dump }}"
          requestorName: "${{ user.entity.metadata.name }}"
          requestorEmail: "${{ user.entity.spec.profile.email or user.entity.metadata.name }}"

  output:
    text:
      title: Catalog Update Submitted
      content: >
        Your catalog update has been submitted. The pipeline will append the new entities
        to your service's catalog-info.yaml, validate the YAML, and commit to master.
        Each new entity will be registered in IDP automatically.
      summary:
        component: "{{ parameters.componentIdentifier }}"
        system: "{{ parameters.systemIdentifier }}"
        updateSystem: "{{ parameters.updateSystem }}"
    links:
      - title: View Pipeline Execution
        url: "https://app.harness.io/ng/account/ACCOUNT_ID/module/idp-admin/orgs/default/projects/cfg-IDP-project/pipelines/update_catalog_entry_pipeline/executions"
```

- [ ] **Step 2: Spec-review the workflow YAML**

Check:
- All arrays use `type: array` with `items.properties` — not `type: string`
- `dependencies/oneOf` used for System toggle — not `allOf/if/then`
- `ContextViewer` fields have `readonly: true` and `getContextData: "{{formContext.key}}"` — all three present
- Step syntax: `{{ parameters.x }}` (no `$`) in `debug.log` and payload body values quoted
- `requestorName`/`requestorEmail` use `${{ }}` syntax (user context injection)

- [ ] **Step 3: Commit**

```bash
cd /Users/chandinidev/Documents/rohid-code/harness-catalog
git add workflows/update-catalog-entry.yaml
git commit -m "feat: add update-catalog-entry workflow"
```

---

### Task 4: Pipeline YAML

**Files:**
- Create: `harness-catalog/pipelines/update-catalog-entry-pipeline.yaml`

The pipeline receives the trigger payload, clones the target repo, runs the appender, runs the validator, commits, and registers entities.

Key patterns from confirmed memory:
- `delegateSelectors: [act-delegate-k8s-ephub-p2r1]` on every ShellScript step
- Step names: ASCII only, no em dash
- Retry: 2 retries × 30s, `onRetryFailure: Abort`
- Secrets: HashiCorp Vault — `<+secrets.getValue("vault://...path")>`
- Failure notification: last stage with `when: pipelineStatus: Failure`, email API only
- `<+trigger.payload.KEY>` for all payload fields

- [ ] **Step 1: Create the pipeline**

```yaml
# harness-catalog/pipelines/update-catalog-entry-pipeline.yaml
pipeline:
  name: Update Catalog Entry Pipeline
  identifier: update_catalog_entry_pipeline
  orgIdentifier: default
  projectIdentifier: cfg-IDP-project
  tags: {}
  variables:
    - name: componentIdentifier
      type: String
      description: "Component identifier from workflow — e.g. payments_service"
      value: <+trigger.payload.componentIdentifier>
    - name: repoSlug
      type: String
      description: "GitHub repo slug — e.g. bank-org/payments-service"
      value: <+trigger.payload.repoSlug>
    - name: sysid
      type: String
      description: "ServiceNow SYSID — e.g. SYSID-00123"
      value: <+trigger.payload.sysid>
    - name: systemIdentifier
      type: String
      description: "Parent system identifier — e.g. payments_platform"
      value: <+trigger.payload.systemIdentifier>
    - name: ownerGroup
      type: String
      description: "Owner group identifier — e.g. team-payments"
      value: <+trigger.payload.ownerGroup>
    - name: updateSystem
      type: String
      description: "Yes or No — whether to update System entity"
      value: <+trigger.payload.updateSystem>
    - name: systemTitle
      type: String
      description: "System display title"
      value: <+trigger.payload.systemTitle>
    - name: systemDescription
      type: String
      description: "System description"
      value: <+trigger.payload.systemDescription>
    - name: systemDomain
      type: String
      description: "System domain — e.g. banking"
      value: <+trigger.payload.systemDomain>
    - name: systemConfluence
      type: String
      description: "Confluence URL for this system"
      value: <+trigger.payload.systemConfluence>
    - name: systemRunbook
      type: String
      description: "Runbook URL for this system"
      value: <+trigger.payload.systemRunbook>
    - name: systemTags
      type: String
      description: "Comma-separated tags"
      value: <+trigger.payload.systemTags>
    - name: databases
      type: String
      description: "JSON array of database entries"
      value: <+trigger.payload.databases>
    - name: caches
      type: String
      description: "JSON array of cache entries"
      value: <+trigger.payload.caches>
    - name: messaging
      type: String
      description: "JSON array of messaging entries"
      value: <+trigger.payload.messaging>
    - name: restApis
      type: String
      description: "JSON array of REST API entries"
      value: <+trigger.payload.restApis>
    - name: asyncProducers
      type: String
      description: "JSON array of async producer entries"
      value: <+trigger.payload.asyncProducers>
    - name: asyncConsumers
      type: String
      description: "JSON array of async consumer entries"
      value: <+trigger.payload.asyncConsumers>
    - name: requestorEmail
      type: String
      description: "Requestor email for failure notification"
      value: <+trigger.payload.requestorEmail>
    - name: requestorName
      type: String
      description: "Requestor display name"
      value: <+trigger.payload.requestorName>

  stages:
    - stage:
        name: Clone and Update Catalog
        identifier: clone_and_update_catalog
        type: Custom
        spec:
          execution:
            steps:
              - step:
                  name: Clone Target Repo
                  identifier: clone_target_repo
                  type: ShellScript
                  spec:
                    shell: Bash
                    delegateSelectors:
                      - act-delegate-k8s-ephub-p2r1
                    source:
                      type: Inline
                      spec:
                        script: |
                          set -euo pipefail
                          GITHUB_PAT=<+secrets.getValue("vault://secret/data/platform/github#pat")>
                          REPO_SLUG=<+pipeline.variables.repoSlug>
                          git config --global user.email "rohid@rohiddev.com"
                          git config --global user.name "Rohid Dev"
                          rm -rf /tmp/catalog-update-repo
                          git clone --retry 3 \
                            "https://${GITHUB_PAT}@github.com/${REPO_SLUG}.git" \
                            /tmp/catalog-update-repo
                          echo "cloned ${REPO_SLUG}"
                  failureStrategies:
                    - onFailure:
                        errors:
                          - AllErrors
                        action:
                          type: Retry
                          spec:
                            retryCount: 2
                            retryIntervals:
                              - 30s
                            onRetryFailure:
                              action:
                                type: Abort

              - step:
                  name: Append Catalog Entries
                  identifier: append_catalog_entries
                  type: ShellScript
                  spec:
                    shell: Bash
                    delegateSelectors:
                      - act-delegate-k8s-ephub-p2r1
                    source:
                      type: Inline
                      spec:
                        script: |
                          set -euo pipefail
                          CATALOG_FILE=/tmp/catalog-update-repo/catalog-info.yaml

                          # Create minimal catalog-info.yaml if not present
                          if [ ! -f "$CATALOG_FILE" ]; then
                            echo "catalog-info.yaml not found — creating skeleton"
                            cat > "$CATALOG_FILE" <<EOF
                          apiVersion: harness.io/v1
                          kind: Component
                          identifier: <+pipeline.variables.componentIdentifier>
                          name: <+pipeline.variables.componentIdentifier>
                          type: service
                          owner: group:account/<+pipeline.variables.ownerGroup>
                          metadata:
                            title: <+pipeline.variables.componentIdentifier>
                            description: Auto-generated by IDP catalog update workflow.
                            annotations:
                              bank.com/sysid: <+pipeline.variables.sysid>
                          spec:
                            lifecycle: production
                            system:
                              - system:account/<+pipeline.variables.systemIdentifier>
                          EOF
                          fi

                          cd /tmp/catalog-update-repo

                          # Run appender from harness-catalog repo (pre-cloned by CI)
                          CATALOG_FILE="$CATALOG_FILE" \
                          OWNER_GROUP="<+pipeline.variables.ownerGroup>" \
                          SYSID="<+pipeline.variables.sysid>" \
                          SYSTEM_IDENTIFIER="<+pipeline.variables.systemIdentifier>" \
                          UPDATE_SYSTEM="<+pipeline.variables.updateSystem>" \
                          SYSTEM_TITLE="<+pipeline.variables.systemTitle>" \
                          SYSTEM_DESCRIPTION="<+pipeline.variables.systemDescription>" \
                          SYSTEM_DOMAIN="<+pipeline.variables.systemDomain>" \
                          SYSTEM_CONFLUENCE="<+pipeline.variables.systemConfluence>" \
                          SYSTEM_RUNBOOK="<+pipeline.variables.systemRunbook>" \
                          SYSTEM_TAGS="<+pipeline.variables.systemTags>" \
                          DATABASES='<+pipeline.variables.databases>' \
                          CACHES='<+pipeline.variables.caches>' \
                          MESSAGING='<+pipeline.variables.messaging>' \
                          REST_APIS='<+pipeline.variables.restApis>' \
                          ASYNC_PRODUCERS='<+pipeline.variables.asyncProducers>' \
                          ASYNC_CONSUMERS='<+pipeline.variables.asyncConsumers>' \
                          go run /harness-catalog/tools/catalog-appender/main.go

                          echo "append complete"
                  failureStrategies:
                    - onFailure:
                        errors:
                          - AllErrors
                        action:
                          type: Abort

              - step:
                  name: Validate Catalog YAML
                  identifier: validate_catalog_yaml
                  type: ShellScript
                  spec:
                    shell: Bash
                    delegateSelectors:
                      - act-delegate-k8s-ephub-p2r1
                    source:
                      type: Inline
                      spec:
                        script: |
                          set -euo pipefail
                          go run /harness-catalog/tools/catalog-validator/main.go \
                            /tmp/catalog-update-repo/catalog-info.yaml
                          echo "validation passed"
                  failureStrategies:
                    - onFailure:
                        errors:
                          - AllErrors
                        action:
                          type: Abort

              - step:
                  name: Commit and Push to Master
                  identifier: commit_and_push
                  type: ShellScript
                  spec:
                    shell: Bash
                    delegateSelectors:
                      - act-delegate-k8s-ephub-p2r1
                    source:
                      type: Inline
                      spec:
                        script: |
                          set -euo pipefail
                          cd /tmp/catalog-update-repo
                          git add catalog-info.yaml
                          git diff --cached --quiet && echo "nothing to commit" && exit 0
                          git commit -m "chore: update catalog-info.yaml via IDP workflow

                          Component: <+pipeline.variables.componentIdentifier>
                          Requestor: <+pipeline.variables.requestorName>
                          SYSID: <+pipeline.variables.sysid>"
                          git push --retry 3 origin master
                          echo "pushed to master"
                  failureStrategies:
                    - onFailure:
                        errors:
                          - AllErrors
                        action:
                          type: Retry
                          spec:
                            retryCount: 2
                            retryIntervals:
                              - 30s
                            onRetryFailure:
                              action:
                                type: Abort

    - stage:
        name: Register Entities in IDP
        identifier: register_entities
        type: Custom
        spec:
          execution:
            steps:
              - step:
                  name: Register All Entities
                  identifier: register_all_entities
                  type: ShellScript
                  spec:
                    shell: Bash
                    delegateSelectors:
                      - act-delegate-k8s-ephub-p2r1
                    source:
                      type: Inline
                      spec:
                        script: |
                          set -euo pipefail
                          HARNESS_API_KEY=<+secrets.getValue("vault://secret/data/platform/harness#api_key")>
                          ACCOUNT_ID=<+account.identifier>
                          REPO_SLUG=<+pipeline.variables.repoSlug>

                          # Register each document in catalog-info.yaml individually
                          # Split multi-doc YAML into individual files and import each
                          CATALOG_FILE=/tmp/catalog-update-repo/catalog-info.yaml
                          DOC_INDEX=0

                          csplit --quiet --prefix=/tmp/catalog-doc- --suffix-format="%03d.yaml" \
                            "$CATALOG_FILE" "/^---/" "{*}" 2>/dev/null || true

                          for doc in /tmp/catalog-doc-*.yaml; do
                            [ -f "$doc" ] || continue
                            # Skip empty docs
                            grep -q "apiVersion:" "$doc" || continue

                            echo "Registering $doc"
                            HTTP_STATUS=$(curl -s -o /tmp/register-response.json -w "%{http_code}" \
                              --retry 3 --retry-delay 5 --retry-connrefused \
                              -X POST \
                              "https://app.harness.io/gateway/idp/api/v1/entities?accountIdentifier=${ACCOUNT_ID}" \
                              -H "x-api-key: ${HARNESS_API_KEY}" \
                              -H "Content-Type: application/yaml" \
                              --data-binary @"$doc")

                            if [ "$HTTP_STATUS" -ge 400 ]; then
                              echo "WARNING: registration failed for $doc (HTTP $HTTP_STATUS)"
                              cat /tmp/register-response.json
                              # Continue — partial registration is better than aborting
                            else
                              echo "Registered $doc (HTTP $HTTP_STATUS)"
                            fi
                          done

                          echo "entity registration complete"
                  failureStrategies:
                    - onFailure:
                        errors:
                          - AllErrors
                        action:
                          type: Ignore

    - stage:
        name: Notify on Failure
        identifier: notify_on_failure
        type: Custom
        when:
          pipelineStatus: Failure
        failureStrategies:
          - onFailure:
              errors:
                - AllErrors
              action:
                type: Ignore
        spec:
          execution:
            steps:
              - step:
                  name: Send Failure Email
                  identifier: send_failure_email
                  type: ShellScript
                  spec:
                    shell: Bash
                    delegateSelectors:
                      - act-delegate-k8s-ephub-p2r1
                    source:
                      type: Inline
                      spec:
                        script: |
                          set -euo pipefail
                          EMAIL_API_KEY=<+secrets.getValue("vault://secret/data/platform/email#api_key")>
                          curl -s --retry 3 --retry-delay 5 --retry-connrefused \
                            -X POST https://email.bank.com/send \
                            -H "Authorization: Bearer ${EMAIL_API_KEY}" \
                            -H "Content-Type: application/json" \
                            -d "{
                              \"to\": \"<+pipeline.variables.requestorEmail>\",
                              \"subject\": \"Catalog Update Failed - <+pipeline.variables.componentIdentifier>\",
                              \"body\": \"Your catalog update for <+pipeline.variables.componentIdentifier> failed. Please check the pipeline execution for details: https://app.harness.io/ng/account/<+account.identifier>/module/idp-admin/orgs/default/projects/cfg-IDP-project/pipelines/update_catalog_entry_pipeline/executions\"
                            }"
                          echo "failure notification sent"
```

- [ ] **Step 2: Commit**

```bash
cd /Users/chandinidev/Documents/rohid-code/harness-catalog
git add pipelines/update-catalog-entry-pipeline.yaml
git commit -m "feat: add update-catalog-entry pipeline"
```

---

### Task 5: Fix catalog-info.template.yaml

**Files:**
- Modify: `harness-catalog/idpdocuments/catalog-info.template.yaml`

The template currently uses old Backstage format — `identifier`/`name`/`type`/`owner` must be top-level per confirmed IDP 2.0 format. `group:default/` must be `group:account/`. `spec.system` must be a list.

- [ ] **Step 1: Fix the template structure**

For each entity kind, move `name`, `identifier`, `type`, `owner` to top level and fix all references:

```yaml
# System — correct format
apiVersion: harness.io/v1
kind: System
identifier: <platform-name>              # REQUIRED snake_case — e.g. payments_platform
name: <platform-name>                    # REQUIRED kebab-case — e.g. payments-platform
type: system
owner: group:account/<team-identifier>  # REQUIRED e.g. group:account/team-payments
metadata:
  title: <Platform Display Name>
  description: >
    <One paragraph describing what this platform does.>
  annotations:
    bank.com/sysid: <SYSID-NNNNN>
    jira/project-key: <PROJ>
  tags:
    - <domain>
    - <tier>
    - sysid-<NNNNN>
  links:
    - url: https://confluence.bank.com/display/<PROJ>
      title: Confluence Space
      icon: docs
    - url: https://runbooks.bank.com/<platform-name>
      title: Platform Runbook
      icon: book
spec:
  domain: <domain>
```

And for Component `spec.system` must be a list:
```yaml
spec:
  type: service
  lifecycle: production
  system:
    - system:account/<platform-name>     # LIST not string
  dependsOn:
    - resource:account/<service-name>-db
    - resource:account/<service-name>-cache
```

And all `group:default/` → `group:account/` throughout.

- [ ] **Step 2: Commit**

```bash
cd /Users/chandinidev/Documents/rohid-code/harness-catalog
git add idpdocuments/catalog-info.template.yaml
git commit -m "fix: update catalog-info.template.yaml to Harness IDP 2.0 format"
```

---

### Task 6: Push All

- [ ] **Step 1: Push develop branch**

```bash
cd /Users/chandinidev/Documents/rohid-code/harness-catalog
git push origin develop
```

Expected: all 5 commits pushed.

---

## Smoke Test Checklist (5 pilot repos)

After pipeline is live in Harness:
1. Import `workflows/update-catalog-entry.yaml` in IDP Admin > Workflows
2. Import `pipelines/update-catalog-entry-pipeline.yaml` in Harness pipeline studio
3. Run workflow against 1 test repo — add 1 database + 1 REST API
4. Check pipeline execution — validate step should pass, commit should appear on master
5. Check IDP catalog — new Resource and API entities should appear
6. Run validator locally against the committed file: `go run tools/catalog-validator/main.go /tmp/test-catalog.yaml`
