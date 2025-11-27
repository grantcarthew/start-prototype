# Phase 6: Asset Catalog & Lazy Loading

**Status:** Not Started
**Dependencies:** Phase 5
**Estimated Effort:** 8-10 hours

---

## Required Reading

Before starting this phase, review these documents:

**Design Records (Critical):**

- [DR-031: Catalog-Based Assets](../design/design-records/dr-031-catalog-based-assets.md) - Catalog architecture overview
- [DR-032: Metadata Schema](../design/design-records/dr-032-asset-metadata-schema.md) - Asset metadata structure
- [DR-033: Asset Resolution Algorithm](../design/design-records/dr-033-asset-resolution-algorithm.md) - How assets are resolved
- [DR-034: GitHub Catalog API](../design/design-records/dr-034-github-catalog-api.md) - GitHub API integration
- [DR-039: Catalog Index](../design/design-records/dr-039-catalog-index.md) - Index file structure

**Design Records (Supporting):**

- [DR-036: Cache Management](../design/design-records/dr-036-cache-management.md) - Caching strategy
- [DR-026: Offline Behavior](../design/design-records/dr-026-offline-behavior.md) - Network unavailable handling
- [DR-027: Security Trust Model](../design/design-records/dr-027-security-trust-model.md) - Security considerations

**CLI Documentation:**

- [start-assets.md](../cli/start-assets.md) - Asset management overview
- [start-assets-search.md](../cli/start-assets-search.md) - Search command
- [start-assets-add.md](../cli/start-assets-add.md) - Add command

---

## Goal

GitHub integration for on-demand asset downloads.

---

## Deliverables

- [ ] GitHub catalog client (index.csv)
- [ ] Asset resolution algorithm (DR-033)
- [ ] Cache management
- [ ] Lazy loading on first use
- [ ] `start assets search/add` commands

---

## Testing Criteria

- [ ] Can search catalog
- [ ] Can download assets
- [ ] Assets cached correctly
- [ ] Lazy loading works
- [ ] Works offline if cached

---

_Next: [Phase 7](phase-7.md)_
