# GhanaCodes roadmap

This roadmap is directional, not a commitment. It is a readable summary of the execution ledger; the machine-readable state, task ownership and evidence live in [`agent_plan.md`](agent_plan.md). No dates are promised here, and none should be inferred.

**Current lifecycle: public beta.** Web and API are live on their canonical hosts with recorded smoke and rollback evidence in [`docs/runbooks/release-evidence.md`](docs/runbooks/release-evidence.md). Stable gates remain open.

## Now — shipped in the current beta

- [x] **GSS source and licence review** (`CODE-0.1`). Source register entries with authority, licence, publication decision and a pinned SHA-256; the official PDF is linked, never redistributed.
- [x] **Reproducible namespace dataset** (`CODE-1.1`). [`scripts/import_gss_2021.rb`](scripts/import_gss_2021.rb) regenerates `gss-phc-2021.1` — 261 unique district codes extracted from Appendix 2 of the official PDF under a verified SHA-256 checksum, plus 16 region codes and 16 region abbreviations transcribed into the importer as a reviewed constant, 293 records in total — and aborts on a checksum, count, uniqueness or empty-name failure.
- [x] **Date-aware resolver and crosswalk** (`CODE-1.2`). Exact resolution within a declared namespace, `validFrom`/`validTo` enforcement, crosswalks through canonical entity IDs, and `resolved` / `not_found` / `ambiguous` outcomes that never guess. Covered by [`internal/codes/resolver_test.go`](internal/codes/resolver_test.go).
- [x] **REST, constrained GraphQL and bulk interfaces** (`CODE-2.1`) with committed contracts in [`contracts/`](contracts/).
- [x] **TypeScript client source** (`CODE-2.2`, locally complete). [`sdk/typescript/index.ts`](sdk/typescript/index.ts) — written and typed, not yet published as a package.
- [x] **Public register and browser resolver** (`CODE-3.1`). Accessible namespace picker, portfolio typography, favicon, manifest and full canonical/Open Graph metadata.
- [x] **Web production release** (`CODE-4.1`) on `codes.digitalghana.dev`, with an immutable deployment recorded for rollback.
- [x] **API production release** (`CODE-4.2`) on `api-codes.digitalghana.dev`, with a verified custom-domain certificate.
- [x] **Production smoke and rollback evidence** (`CODE-4.3`). Exact, historical, crosswalk, bulk, GraphQL, health, TLS and web smoke all passed against the canonical hosts.

## Next — the open stable gates

The release rule in [`agent_plan.md`](agent_plan.md) states what stable additionally requires beyond beta: reviewed admin publication, broader version history, package releases, security and load evidence, and named ongoing source stewardship.

| Gate | State | Blocking dependency |
|---|---|---|
| Namespace and source administration, conflict review, audited import approval (`CODE-3.2`) | Pending | Source review `CODE-0.1` is done; the task is unassigned and needs an owner and a design |
| Package publication of the TypeScript client (`CODE-2.2`) | Done locally, publication pending | Interfaces `CODE-2.1` are done; needs package metadata, a build step and a release path |
| Broader version history | Not started | A second dataset version admitted through source review, so `validTo` and supersession are exercised on real data |
| Security and load evidence | Not started | Named in the stable release rule; no evidence recorded yet |
| Named ongoing source stewardship | Not started | Requires a named human owner for re-verification of GSS sources |

## Later — deferred scope

These are plausible and not rejected, but nothing is scheduled and each is gated on its own review.

- Additional identifier namespaces, admitted one at a time and only after that namespace's own authority and licence review — never bundled in.
- Historical namespace coverage predating the 2021 census codes, so a code can be resolved against the structure in force at an earlier date.
- Richer conflict and supersession records once more than one dataset version exists.
- Operational maturity described in [`docs/runbooks/operations.md`](docs/runbooks/operations.md): structured logs with correlation identifiers, latency and correctness signals, alert routes with an acknowledged owner, and rate limits and audit trails for privileged publication.

## Explicitly out of scope

From the stated non-goals in [`agent_plan.md`](agent_plan.md) and [`docs/adr/0001-product-boundary.md`](docs/adr/0001-product-boundary.md):

- **A perpetually current MMDA list.** The 2021 namespace is published as the 2021 namespace and will not be quietly maintained as "today's districts".
- **An automatically published scrape of the Ministry MMDA directory.** It is registered as a review reference only.
- **Any institution, education, tax, health or private identifier namespace** admitted without its own authority and licence review.
- **Fuzzy resolution of any kind.** Aliases must be explicit versioned records; approximate name matching will not be added.
- **Redistribution of the source PDF.** The official document stays linked, not re-hosted.
- **Geometry, boundaries, population or coordinates.** GhanaCodes publishes identifier facts only.
- **Any government transaction, credential, form or payment proxy**, and any claim of government affiliation or endorsement.
- **Shared databases or secrets with another Digital Ghana product.** Integration happens only through versioned contracts or pinned dataset artifacts, per ADR-0001.
