# GhanaCodes execution ledger

Last updated: 2026-09-01
Status: Beta implementation complete locally; production release in progress
Canonical hosts: `codes.digitalghana.dev`, `api-codes.digitalghana.dev`

## Product definition

GhanaCodes resolves identifiers only within an explicit namespace, dataset version, and effective date. It never joins records by fuzzy name matching and never collapses multiple valid targets into a guessed answer.

### Beta scope

- GSS 2021 PHC region codes, region abbreviations, and all 261 district codes.
- Exact REST and constrained GraphQL resolution.
- Crosswalk by canonical entity ID with `resolved`, `not_found`, or `ambiguous` outcomes.
- Effective-date filtering and deterministic ordered bulk resolution.
- Reproducible importer pinned to the official PDF checksum.
- Public register, browser resolver, contracts, and TypeScript client source.

### Non-goals and limitations

- The 2021 namespace is not represented as a perpetually current MMDA list.
- The Ministry directory is a review source, not an automatically published scrape.
- No institution, education, tax, health, or private identifier namespace is admitted without its own authority/licence review.
- No fuzzy resolution; aliases must be explicit versioned records.

## Acceptance evidence

- [x] Exactly 16 region codes, 16 abbreviations, and 261 unique district codes generated.
- [x] Every code carries source ID, entity ID, effective date, and version.
- [x] Code uniqueness is enforced per namespace/effective result set.
- [x] Ambiguous mappings return sorted candidates and never guess.
- [x] Effective dates prevent a 2021 code from resolving in 2020.
- [x] REST, constrained GraphQL, bulk resolution, contracts, and source TypeScript client implemented.
- [x] Responsive public register and exact browser resolver built.
- [ ] Namespace/source admin, conflict review, audited import approval, and package publication remain stable gates.

## Live task board

| ID | Task | Status | Owner | Dependency | Evidence |
|---|---|---|---|---|---|
| CODE-0.1 | GSS source/licence review | Done | Codex | — | Source register, pinned SHA-256, no PDF redistribution |
| CODE-1.1 | Reproducible namespace dataset | Done | Codex | CODE-0.1 | Importer asserts 261 unique districts; 293 total codes |
| CODE-1.2 | Date-aware resolver/crosswalk | Done | Codex | CODE-1.1 | Go tests cover exact, historical, ambiguity and crosswalk behavior |
| CODE-2.1 | REST/GraphQL/bulk interfaces | Done | Codex | CODE-1.2 | Handlers and committed contracts |
| CODE-2.2 | TypeScript client | Done locally | Codex | CODE-2.1 | Source client complete; package publication pending |
| CODE-3.1 | Public website and sandbox | Done locally | Codex | CODE-1.2 | Typecheck, tests and production build pass |
| CODE-3.2 | Admin conflict/import workflow | Pending | Unassigned | CODE-0.1 | Required before stable |
| CODE-4.1 | Web production release | In progress | Codex | CODE-3.1 | Vercel deployment next |
| CODE-4.2 | API production release | In progress | Codex | CODE-2.1 | Render deployment next |
| CODE-4.3 | Production smoke/rollback | Blocked | Codex | CODE-4.1, CODE-4.2 | Verify exact, ambiguous-safe, crosswalk, bulk and TLS after attach |

## Release rule

Beta requires green CI, immutable provider deployments, canonical-domain TLS, exact/crosswalk/bulk parity smoke, and rollback evidence. Stable additionally requires reviewed admin publication, broader version history, package releases, security/load evidence, and named ongoing source stewardship.
