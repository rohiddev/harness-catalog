# TechDocs Setup — How to enable native docs in Harness IDP

TechDocs renders your `docs/` folder natively inside the Harness IDP catalog.
The `catalog-info.yaml` annotation `backstage.io/techdocs-ref: dir:.` activates it.

## Files to add to every repo

```
your-repo/
├── catalog-info.yaml          # already has techdocs-ref: dir:.
├── mkdocs.yml                 # copy from harness-catalog/mkdocs.template.yml
└── docs/
    └── index.md               # copy from harness-catalog/IDPReadme.template.md
```

## Steps

1. Copy `mkdocs.template.yml` → `mkdocs.yml` in your repo root. Fill in `site_name`.
2. Create a `docs/` folder in your repo root.
3. Copy `IDPReadme.template.md` → `docs/index.md`. Fill in all REQUIRED fields.
4. Commit and push.
5. Import (or re-sync) your `catalog-info.yaml` in Harness IDP.
6. The **Docs** tab on your catalog entity will render `docs/index.md` natively.

## catalog-info.yaml annotation (already in template)

```yaml
metadata:
  annotations:
    # Points TechDocs at the master branch repo root
    backstage.io/techdocs-ref: url:https://github.com/<org>/<repo>/tree/master/
    backstage.io/source-location: url:https://github.com/<org>/<repo>/tree/master/
```

> `catalog-info.yaml`, `mkdocs.yml`, and `docs/` must all be on the **master** branch.

## Optional additional pages

Add more markdown files to `docs/` and reference them in `mkdocs.yml` nav:

```
docs/
├── index.md          # Overview (IDPReadme)
├── architecture.md   # Architecture diagrams + decisions
├── runbook.md        # Full incident runbook
├── apis.md           # API reference detail
└── oncall.md         # On-call guide
```
