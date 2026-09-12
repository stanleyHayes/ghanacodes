# Contributing to GhanaCodes

GhanaCodes resolves Ghanaian identifiers exactly, within a declared namespace, dataset version and effective date. Contributions are welcome in code, contracts, documentation and — especially — evidence-backed data corrections.

Read [`AGENTS.md`](AGENTS.md), [`agent_plan.md`](agent_plan.md) and [`docs/adr/0001-product-boundary.md`](docs/adr/0001-product-boundary.md) before making a substantive change. They record the product boundary, the live task board and the rules that automated contributors follow too.

## Before you start

1. Confirm the work belongs here. Portfolio-wide governance, standards and the umbrella site live in the [`digitalghana`](https://github.com/stanleyHayes/digitalghana) repository; product implementation lives here.
2. For anything touching [`data/codes.json`](data/codes.json), identify the source authority, its licence, the retrieval date and the permitted use **before** importing. The register of accepted sources and publication decisions is [`docs/governance/source-register.json`](docs/governance/source-register.json).
3. Check the task board in [`agent_plan.md`](agent_plan.md). Claim one path-bounded task rather than editing broadly.

## Prerequisites

| Tool | Version | Where it is pinned |
|---|---|---|
| Node.js | 24 | [`.github/workflows/quality.yml`](.github/workflows/quality.yml) |
| pnpm | 11.11.0 | `packageManager` in [`package.json`](package.json) |
| Go | 1.23 | [`go.mod`](go.mod), [`Dockerfile`](Dockerfile), [`.github/workflows/quality.yml`](.github/workflows/quality.yml) — note that [`render.yaml`](render.yaml) declares `runtime: go` without a version |
| Ruby | 3.4 | [`.github/workflows/quality.yml`](.github/workflows/quality.yml) |
| `pdftotext` (poppler) | any | Only needed to re-run [`scripts/import_gss_2021.rb`](scripts/import_gss_2021.rb) |

## Local setup

```sh
git clone https://github.com/stanleyHayes/ghanacodes.git
cd ghanacodes
pnpm install --frozen-lockfile

pnpm dev          # public register and browser resolver on http://localhost:3000
go run ./cmd/api  # API on http://localhost:8080 — run from the repository root
```

The Go service reads `data/codes.json` relative to the working directory, so start it from the repository root. Set `PORT` to use a different port.

Smoke-test a change to the resolver:

```sh
curl -s "http://localhost:8080/v1/resolve/gss-phc-2021-district/0101?at=2021-06-27"
curl -s "http://localhost:8080/v1/crosswalk?fromNamespace=gss-phc-2021-region&code=03&toNamespace=gss-phc-2021-region-abbreviation"
```

## Verification

Run the full set before opening a pull request. These are exactly the commands CI runs, in order:

```sh
ruby scripts/validate.rb
go test ./...
go vet ./...
pnpm typecheck
pnpm test
pnpm build
```

> If `ruby scripts/validate.rb` aborts with `invalid byte sequence in US-ASCII`, your shell's default external encoding is not UTF-8. Run it as `RUBYOPT="-E UTF-8" ruby scripts/validate.rb`, or set `LANG=en_US.UTF-8`. CI runs under a UTF-8 locale and is unaffected.

`ruby scripts/validate.rb` checks that the required governance files exist, that the source register is non-empty, and that no unresolved template token or private key has been committed. `pnpm test` asserts the dataset invariants — 3 namespaces, 261 districts, 16 region codes, uniqueness per namespace, and `sourceId` plus a well-formed `validFrom` on every record.

## Commit and pull request conventions

This repository uses short, lower-case Conventional Commit subjects in the imperative mood, matching the existing history:

```
chore: establish GhanaCodes foundation
feat: build GhanaCodes beta foundation
docs: record GhanaCodes beta release
feat: align Codes UI and social metadata
docs: record Codes UI and SEO evidence [skip ci]
```

Use `feat`, `fix`, `docs`, `chore` or `refactor`. Keep the subject under about 72 characters and do not end it with a full stop. Append `[skip ci]` only for documentation-only commits that cannot affect a build.

A pull request should state:

- the scope and the affected paths;
- the source, ADR or task-board ID it relates to (for example `CODE-3.2`);
- the verification actually performed, with output where it matters;
- migration and rollback impact for any contract, dataset-version or endpoint change;
- any external gate that still blocks it.

Breaking changes to [`contracts/openapi.yaml`](contracts/openapi.yaml), [`contracts/graphql/schema.graphql`](contracts/graphql/schema.graphql) or the shape of a resolution need a version decision and a stated migration path. Stable identifiers — `entityId`, `namespace`, `code` — must not be silently repurposed.

## Proposing a data correction

Canonical data is generated, not hand-edited. **Do not edit [`data/codes.json`](data/codes.json) directly in a pull request**; it is produced by [`scripts/import_gss_2021.rb`](scripts/import_gss_2021.rb) from the pinned source, and a manual edit breaks reproducibility.

Open an issue, or a pull request against the importer or the source register, containing:

1. the affected `namespace` and `code`, plus the `entityId`;
2. the current published value and the proposed value;
3. the authoritative source — a URL or document title, its issuing authority and its publication date;
4. whether the change is a **correction** to the 2021 record, or a **later fact** that belongs in a new dataset version;
5. whether historical records change, and if so which dates are affected;
6. the licence status of the evidence you are citing.

Rules that apply to every correction:

- The 2021 namespace is published as the 2021 namespace. A district renamed after the census is a new version or a new namespace, never an in-place rewrite of the 2021 record.
- Aliases must be explicit, versioned records. Nothing is ever resolved by approximate name matching.
- A new source needs an entry in [`docs/governance/source-register.json`](docs/governance/source-register.json) with its authority, licence and publication decision, and it must satisfy [`docs/governance/source-register.schema.json`](docs/governance/source-register.schema.json).
- Automation may draft a correction; a human reviewer approves canonical publication.

## Review expectations

- Changes are reviewed against the product boundary in [`docs/adr/0001-product-boundary.md`](docs/adr/0001-product-boundary.md), not only against whether the code works.
- Any change to resolution semantics needs a Go test in [`internal/codes/resolver_test.go`](internal/codes/resolver_test.go) covering the new behaviour, including the `not_found` and `ambiguous` paths.
- Lifecycle claims must be evidenced. Use the portfolio vocabulary: proposed, building, beta, stable, externally blocked, retired, deferred. Do not upgrade a status in documentation without recording the evidence in [`docs/runbooks/release-evidence.md`](docs/runbooks/release-evidence.md).
- Never commit credentials, provider environment files, private exports or personal data. Report security issues privately per [`SECURITY.md`](SECURITY.md) — never in a public issue.
- Do not imply government endorsement, affiliation or official status anywhere in code, data or copy.

All participation is covered by [`CODE_OF_CONDUCT.md`](CODE_OF_CONDUCT.md).

## Good first contributions

Each of these is a real gap in the current beta, derived from [`ROADMAP.md`](ROADMAP.md) and [`agent_plan.md`](agent_plan.md):

1. **Flesh out the OpenAPI contract.** [`contracts/openapi.yaml`](contracts/openapi.yaml) declares the four REST paths with descriptions but no response schemas. Add `Resolution` and `Code` component schemas matching [`internal/codes/resolver.go`](internal/codes/resolver.go), plus the `at`, `fromNamespace`, `code` and `toNamespace` parameters. No behaviour change.
2. **Add HTTP handler tests.** [`internal/codes/resolver_test.go`](internal/codes/resolver_test.go) covers the resolver well; [`cmd/api/main.go`](cmd/api/main.go) has no tests. Cover the malformed `/v1/resolve/` path (400), the 1000-query and 1 MiB bulk limits, and the GraphQL rejection of an unsupported operation.
3. **Document the GraphQL subset accurately.** The `/graphql` handler dispatches on a substring of the query rather than parsing it, so only `resolveCode` and `crosswalk` are supported and no field selection is honoured. Record that constraint alongside [`contracts/graphql/schema.graphql`](contracts/graphql/schema.graphql) so consumers are not surprised.
4. **Extend the dataset invariant tests.** [`tests/dataset.test.mjs`](tests/dataset.test.mjs) could additionally assert that every `entityId` matches `geo:(region|district):<code>`, that every district `regionCode` exists as a region code, and that every `sourceId` appears in the source register.
5. **Draft the admin conflict and import-approval design.** Task `CODE-3.2` is pending and unassigned, and is a stable gate. A documentation-only design note in `docs/` covering conflict review, audited import approval and the reviewer role would unblock the implementation.
6. **Prepare the TypeScript client for publication.** [`sdk/typescript/index.ts`](sdk/typescript/index.ts) is complete as source but unpublished. Package metadata, a build step, a README and usage examples are needed before `CODE-2.2` can close.

## Licensing

Unless you state otherwise, contributions intentionally submitted for inclusion in GhanaCodes are provided under the Apache License 2.0 on the licence's inbound=outbound terms. See [`LICENSE`](LICENSE) and [`NOTICE`](NOTICE).

Contributions must not include third-party data or documents without recorded permission. Referenced source documents retain their own rights and are not relicensed by being cited here.
