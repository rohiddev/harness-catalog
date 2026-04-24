# Harness IDP Homepage Design

**Date:** 2026-04-24
**Status:** Approved

---

## Goal

A world-class developer portal homepage for a bank with multiple internal teams (Loans, Payments, Money Transfer). Developers log in and immediately find: their recently visited and starred services, platform status, self-service workflows, scorecards, and onboarding help — all via native Harness IDP features only.

---

## Constraints

- Native Harness IDP features only — no custom plugins, no iFrames, no API-driven widgets
- Bank logo set via Admin > Appearance (portal-wide) — not in homepage YAML
- Banner is static text, manually updated by Platform Engineering team during freeze events
- ~100 self-service workflows — link to workflows page, do not tile them on homepage

---

## Files

| File | Purpose | Apply via |
|---|---|---|
| `homepage.yaml` | Homepage layout | Admin > Layout > Homepage > Edit YAML |
| `sidenav.yaml` | Left navigation | Admin > Layout > Sidenav > Edit YAML |

---

## Sidenav

6 items, pinned:

| Item | URL | Icon |
|---|---|---|
| Home | `/` | home |
| Catalog | `/catalog` | catalog |
| Self-Service | `/create` | scaffolder |
| Scorecards | `/scorecards` | quality |
| TechDocs | `/docs` | docs |
| My Starred | `/catalog?filters[starred]=true` | star |

---

## Homepage Layout

### Header
- Title: `<+greeting>, <+first_name>. Welcome to the Bank Developer Portal.`
- Subtitle: `Your single source of truth for services, APIs, and self-service.`
- Quick links: Catalog · Self-Service · Scorecards · TechDocs
- `<+greeting>` = native Harness IDP time-based variable (Good morning/afternoon/evening)
- `<+first_name>` = native Harness IDP user variable

### Banner (static, admin-managed)
- **Default state:** `✅ No active production freezes. For platform support contact platform-team@bank.com`
- **Freeze state:** `🔴 PROD FREEZE ACTIVE — no deployments to production until [date]. Contact platform-team@bank.com`
- Updated manually in Admin > Layout > Homepage during freeze events
- **Why not API-driven:** Native banner does not support dynamic content. Admin update takes 2 minutes.

### Row 1 — Platform Status (12md)
Full-width Markdown card. Links to:
- Freeze Calendar (Confluence)
- Change Advisory Board (ServiceNow)
- Incident Dashboard
- On-Call Schedule
- Harness Pipelines
- Datadog

Update the URLs in `homepage.yaml` to point to internal bank systems before go-live.

### Row 2 — Personalised (6md + 6md)
- `RecentlyVisited` — native card, auto-populated per user
- `TopVisited` — native card, most visited entities across the portal

### Row 3 — My Workspace (6md + 6md)
- `StarredEntities` — native card, user's bookmarked entities
- `What's New` — Markdown card, maintained by Platform Engineering team with platform release notes

### Row 4 — Onboarding (6md + 6md)
- `Getting Started` — 5-step checklist for new joiners: register service → CMDB → TechDocs → workflows → star
- `IDP 101` — concept reference table: Catalog, catalog-info.yaml, now.yaml, Workflow, Scorecard, TechDocs, System, Starred

---

## Production Go-Live Checklist

- [ ] Bank logo uploaded via Admin > Appearance > Logo
- [ ] Update `homepage.yaml` Platform Status card URLs to real internal URLs (replace `bank.com` placeholders)
- [ ] Update `What's New` card content with real announcements
- [ ] Apply `homepage.yaml` via Admin > Layout > Homepage
- [ ] Apply `sidenav.yaml` via Admin > Layout > Sidenav
- [ ] Test with a developer account — verify `<+first_name>` and `<+greeting>` render correctly
- [ ] Brief Platform Engineering team on banner update procedure for freeze events

---

## What's Not on This Homepage (intentional)

| Item | Reason |
|---|---|
| Workflow tiles | ~100 workflows — link to page instead |
| Pending approvals | No native homepage card — developers use Harness platform directly |
| API-driven freeze status | Banner is native static only — manual update is sufficient |
| Team/domain filter | Catalog tagging not yet consistent — future enhancement |
| 4 big pillar tiles | Replaced by sidenav — tiles add no value when nav is always visible |
