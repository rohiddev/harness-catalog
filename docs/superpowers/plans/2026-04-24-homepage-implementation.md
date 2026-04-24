# Harness IDP Homepage Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Configure the Harness IDP homepage and sidenav to give bank developers a world-class, navigation-first portal experience using only native Harness IDP features.

**Architecture:** Two YAML files applied via the Harness IDP Admin UI. `sidenav.yaml` controls the persistent left navigation. `homepage.yaml` controls what developers see on login — personalised greeting, platform status, recently visited/starred services, what's new, and onboarding help. No custom plugins or code required.

**Tech Stack:** Harness IDP native — homepage YAML, sidenav YAML, Markdown cards, native cards (RecentlyVisited, TopVisited, StarredEntities), banner.

---

### Task 1: Finalize homepage.yaml with real internal URLs

**Files:**
- Modify: `harness-catalog/homepage.yaml`

The current file has `bank.com` placeholder URLs in the Platform Status card. Replace them with real internal URLs before the file is applied in the Admin UI.

- [ ] **Step 1: Open homepage.yaml and locate the Platform Status card**

Open `harness-catalog/homepage.yaml`. Find the `## Platform Status` section (Row 1, 12md card). It contains 6 placeholder URLs:

```yaml
  - type: custom
    width: 12md
    content: |
      ## Platform Status

      | | Resource | Description |
      |---|---|---|
      | 📅 | [Freeze Calendar](https://confluence.bank.com/display/ENG/freeze-calendar) | Scheduled production freeze windows |
      | 🎫 | [Change Advisory Board](https://servicenow.bank.com/change) | Active and pending change requests |
      | 🚨 | [Incident Dashboard](https://oncall.bank.com/incidents) | Live incidents and P1/P2 status |
      | 📞 | [On-Call Schedule](https://oncall.bank.com) | Who is on-call right now |
      | 🏗 | [Harness Pipelines](https://app.harness.io) | CI/CD pipeline runs |
      | 📊 | [Datadog](https://app.datadoghq.com) | Service health and dashboards |
```

- [ ] **Step 2: Replace placeholder URLs**

Update each `bank.com` URL to the real internal URL:

| Placeholder | Replace with |
|---|---|
| `https://confluence.bank.com/display/ENG/freeze-calendar` | Real Confluence freeze calendar URL |
| `https://servicenow.bank.com/change` | Real ServiceNow change board URL |
| `https://oncall.bank.com/incidents` | Real incident dashboard URL |
| `https://oncall.bank.com` | Real on-call schedule URL |

> `https://app.harness.io` and `https://app.datadoghq.com` are real URLs — leave them as-is.

- [ ] **Step 3: Update banner email address**

In the `banner:` section, replace `platform-team@bank.com` with the real platform team email:

```yaml
banner:
  text: "✅ No active production freezes. For platform support contact <real-email@yourbank.com>"
```

- [ ] **Step 4: Update Getting Started card email**

In Row 4's `Getting Started` card, replace the email in the last line:

```yaml
      Need help? Email [platform-team@bank.com](mailto:platform-team@bank.com)
```

Replace both occurrences of `platform-team@bank.com` with the real email.

- [ ] **Step 5: Update What's New card with real announcements**

In Row 3's `What's New` card, replace the example entries with real platform announcements relevant to your developers. Keep the same markdown format:

```markdown
## What's New

**[Month Year]**
- [announcement]
- [announcement]

**[Month Year]**
- [announcement]

---
[Full platform release notes →](https://developer.harness.io/release-notes/internal-developer-portal)
```

- [ ] **Step 6: Commit the finalised homepage.yaml**

```bash
cd /Users/chandinidev/Documents/rohid-code/harness-catalog
git add homepage.yaml
git commit -m "feat: finalize homepage.yaml with real internal URLs"
```

---

### Task 2: Upload bank logo via Admin > Appearance

**Files:** None — configured in Harness IDP Admin UI only.

The logo is set portal-wide via Admin > Appearance. It appears in the top-left of the sidenav and the homepage header on every page — not in the YAML.

- [ ] **Step 1: Prepare the logo file**

The logo must be:
- Format: PNG or SVG
- Recommended size: 200×40px (horizontal lockup) or 40×40px (icon only)
- Background: transparent if possible (dark sidenav background in Harness IDP)

- [ ] **Step 2: Upload the logo**

1. Log in to Harness IDP as a Platform Admin
2. Go to **Admin > Appearance**
3. Find the **Logo** field
4. Upload your bank logo file
5. Click **Save**

- [ ] **Step 3: Verify logo renders**

Navigate to the IDP homepage. The bank logo should appear in the top-left of the sidenav. Verify it looks correct at normal browser zoom (100%).

> **Production swap:** When going to prod, repeat this step on the prod Harness account. The YAML files do not need to change.

---

### Task 3: Apply sidenav.yaml in Harness IDP Admin UI

**Files:**
- Reference: `harness-catalog/sidenav.yaml`

- [ ] **Step 1: Copy the sidenav YAML content**

Open `harness-catalog/sidenav.yaml`. The content is:

```yaml
page:
  name: SidebarLayout
  items:
    - title: Home
      url: /
      icon: home

    - title: Catalog
      url: /catalog
      icon: catalog

    - title: Self-Service
      url: /create
      icon: scaffolder

    - title: Scorecards
      url: /scorecards
      icon: quality

    - title: TechDocs
      url: /docs
      icon: docs

    - title: My Starred
      url: /catalog?filters[starred]=true
      icon: star
```

- [ ] **Step 2: Apply in Admin UI**

1. Log in to Harness IDP as a Platform Admin
2. Go to **Admin > Layout**
3. Select **Sidenav** from the layout list
4. Click **Edit YAML**
5. Replace the existing content with the YAML above
6. Click **Save**

- [ ] **Step 3: Verify sidenav renders correctly**

Navigate to the IDP homepage. The left sidebar should show exactly 6 items in this order:
- Home
- Catalog
- Self-Service
- Scorecards
- TechDocs
- My Starred

Click each item and confirm it navigates to the correct page.

---

### Task 4: Apply homepage.yaml in Harness IDP Admin UI

**Files:**
- Reference: `harness-catalog/homepage.yaml` (finalised in Task 1)

- [ ] **Step 1: Copy the homepage YAML content**

Open `harness-catalog/homepage.yaml`. Copy the entire file content.

- [ ] **Step 2: Apply in Admin UI**

1. Log in to Harness IDP as a Platform Admin
2. Go to **Admin > Layout**
3. Select **Homepage** from the layout list
4. Click **Edit YAML**
5. Replace the existing content with the copied YAML
6. Click **Save**

- [ ] **Step 3: Verify homepage renders — header**

Log in as a regular developer account (not admin). Navigate to the IDP homepage. Verify:
- `<+greeting>` renders as `Good morning` / `Good afternoon` / `Good evening` based on current time
- `<+first_name>` renders as the logged-in developer's first name (not the literal string `<+first_name>`)
- Subtitle text is visible
- Quick links (Catalog, Self-Service, Scorecards, TechDocs) appear in the header

- [ ] **Step 4: Verify homepage renders — banner**

Verify the banner strip shows:
```
✅ No active production freezes. For platform support contact <your-email>
```

- [ ] **Step 5: Verify homepage renders — all rows**

Scroll through the homepage and confirm all 4 rows render:

| Row | Expected content |
|---|---|
| Row 1 | Platform Status card — full width — 6 links in a table |
| Row 2 | Recently Visited (left) + Top Visited (right) — side by side |
| Row 3 | Starred Entities (left) + What's New (right) — side by side |
| Row 4 | Getting Started (left) + IDP 101 (right) — side by side |

- [ ] **Step 6: Test Recently Visited populates**

Visit 3-5 catalog entities. Navigate back to the homepage. Confirm the Recently Visited card now shows those entities.

- [ ] **Step 7: Test Starred Entities populates**

Star 1-2 catalog entities (click ⭐ on any entity page). Navigate back to the homepage. Confirm the Starred Entities card shows them.

---

### Task 5: Banner update procedure — brief Platform Engineering team

**Files:** None — this is a process task.

The banner must be updated manually whenever a production freeze is active.

- [ ] **Step 1: Document the freeze banner procedure**

Share this procedure with the Platform Engineering team:

**To activate a freeze banner:**
1. Log in to Harness IDP as a Platform Admin
2. Go to **Admin > Layout > Homepage > Edit YAML**
3. Find the `banner:` section
4. Replace the text with:
   ```yaml
   banner:
     text: "🔴 PROD FREEZE ACTIVE — no deployments to production until [DATE]. Contact <email>"
   ```
5. Replace `[DATE]` with the actual freeze end date (e.g. `Friday 1 May 2026 17:00 UTC`)
6. Click **Save**

**To deactivate the freeze banner:**
1. Repeat steps 1-3 above
2. Replace the text with:
   ```yaml
   banner:
     text: "✅ No active production freezes. For platform support contact <email>"
   ```
3. Click **Save**

> Total time: ~2 minutes. No code change, no PR, no deployment required.

- [ ] **Step 2: Commit the finalised plan**

```bash
cd /Users/chandinidev/Documents/rohid-code/harness-catalog
git add docs/superpowers/plans/2026-04-24-homepage-implementation.md
git add docs/superpowers/specs/2026-04-24-homepage-design.md
git add sidenav.yaml
git commit -m "feat: add homepage and sidenav YAML with implementation plan"
```

---

## Verification Summary

After all tasks are complete, run through this final check:

- [ ] Bank logo visible top-left on every page
- [ ] Sidenav shows: Home, Catalog, Self-Service, Scorecards, TechDocs, My Starred
- [ ] Homepage header shows personalised greeting + first name
- [ ] Banner shows no-freeze status message
- [ ] Platform Status card shows 6 links (all real URLs, not bank.com placeholders)
- [ ] Recently Visited card populates after browsing
- [ ] Top Visited card shows entries
- [ ] Starred Entities card populates after starring
- [ ] What's New card shows real announcements
- [ ] Getting Started card shows 5-step checklist
- [ ] IDP 101 card shows concept table
- [ ] All sidenav links navigate to correct pages
- [ ] Tested on a non-admin developer account
