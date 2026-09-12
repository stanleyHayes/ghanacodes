# GhanaCodes

Exact, date-aware resolution and crosswalking of Ghanaian region and district identifiers, against an explicit namespace, dataset version and effective date — never a fuzzy name match.

[![Licence](https://img.shields.io/badge/licence-Apache--2.0-blue.svg)](LICENSE)
[![Web](https://img.shields.io/badge/web-codes.digitalghana.dev-brightgreen.svg)](https://codes.digitalghana.dev)
[![API](https://img.shields.io/badge/API-api--codes.digitalghana.dev-brightgreen.svg)](https://api-codes.digitalghana.dev/health)
[![Quality](https://img.shields.io/github/actions/workflow/status/stanleyHayes/ghanacodes/quality.yml?branch=main&label=quality)](https://github.com/stanleyHayes/ghanacodes/actions/workflows/quality.yml)
[![Runtime](https://img.shields.io/badge/runtime-Go%201.23%20%C2%B7%20Node%2024-informational.svg)](go.mod)

## Live now

| Surface | URL | Notes |
|---|---|---|
| Public register and browser resolver | <https://codes.digitalghana.dev> | Next.js page; resolves in the browser against the pinned dataset |
| REST, GraphQL and bulk API | <https://api-codes.digitalghana.dev> | Go service on Render's free plan |
| Health | <https://api-codes.digitalghana.dev/health> | Dataset version and loaded code count |

The API runs on a free Render instance that sleeps when idle. **The first request after a quiet period can take 30–60 seconds**; later ones are fast. Hit the health check first if you are demonstrating it.

```sh
# Resolve GSS 2021 district code 0101 as it stood on census day.
curl -s "https://api-codes.digitalghana.dev/v1/resolve/gss-phc-2021-district/0101?at=2021-06-27"
```

```json
{
  "status": "resolved",
  "query": {"at": "2021-06-27", "code": "0101", "namespace": "gss-phc-2021-district"},
  "candidates": [{
    "namespace": "gss-phc-2021-district", "code": "0101", "name": "Jomoro",
    "entityId": "geo:district:0101", "entityType": "district",
    "districtType": "municipal", "regionCode": "01",
    "validFrom": "2021-06-27", "validTo": null,
    "status": "active", "sourceId": "gss-2021-phc-field-manual"
  }],
  "version": "gss-phc-2021.1"
}
```

Ask for the same code at a date before the namespace existed and you get an honest `not_found` rather than a lucky hit — change `at` to `2020-01-01` and `candidates` comes back empty with `"status": "not_found"`.

## The problem this solves

### For developers

You have two Ghanaian datasets and you need to join them. One is keyed on a four-digit GSS district code. The other has a column called `District` with strings in it. Before GhanaCodes, the path from A to B looked like this:

- **Find the codes at all.** The authoritative 2021 region and district lists live inside Appendix 2 and Table 10.1 of a Ghana Statistical Service PDF field manual. Not a CSV, not an API — a PDF. You run `pdftotext`, eyeball the columns, and fix the rows that wrapped.
- **Hand-maintain the result.** That extraction becomes a spreadsheet in someone's Drive or a hard-coded array in one service. Six months later a second service grows its own copy, the two disagree about three districts, and nobody knows which is right.
- **Guess at ambiguity.** `03` is Greater Accra as a GSS region code. It is also a plausible prefix, a month and an internal ID. A bare value carries no namespace, so joins silently bind the wrong rows.
- **Fight the names.** District 0201 is published as `Komenda Edina Eguafo Abirem Municipal`; your other dataset says `KEEA`. Approximate matching will confidently hand you the wrong district without telling you.
- **Ignore time.** A code valid in a 2021 census extract is not necessarily valid in a table compiled under an earlier regional structure. Most pipelines have no concept of *when* a code was valid, so historical records get re-labelled with present-day identifiers and the error becomes invisible.

Each is a few hours of work and a permanent source of quiet defects — the kind that survive review because the join succeeded and returned rows.

What GhanaCodes gives you instead:

- One pinned, reproducible dataset — `gss-phc-2021.1`, 293 codes across 3 declared namespaces. The 261 district codes are extracted from Appendix 2 of the official PDF under a verified SHA-256 checksum; the 16 region codes and 16 abbreviations are transcribed into the importer as a reviewed constant.
- Resolution that **requires** a namespace. There is no "just look up `03`" call to misuse.
- An `at` parameter on every lookup. A code resolves only within its own validity window.
- Crosswalks that travel through a canonical entity ID (`geo:region:03`), never through a name.
- Three explicit outcomes — `resolved`, `not_found`, `ambiguous` — with all candidates returned when more than one record is valid. The service never picks a winner for you.
- Deterministic, order-preserving bulk resolution for pipeline work, plus REST, a constrained GraphQL surface, a committed OpenAPI contract and a TypeScript client in [`sdk/typescript/index.ts`](sdk/typescript/index.ts).

### For the community

Identifier resolution is boring infrastructure, and that is exactly why it should be public. When every analyst, newsroom, ministry contractor and startup maintains a private district list, the same errors are re-made independently and there is no way to compare two published figures with confidence.

- **Reproducibility.** [`scripts/import_gss_2021.rb`](scripts/import_gss_2021.rb) refuses to run unless the source PDF matches a pinned checksum, and asserts 261 unique districts before writing anything; the 16 region codes and 16 abbreviations are not extracted from the PDF but transcribed into the importer as a reviewed constant. Anyone holding the official document can regenerate the identical file.
- **Provenance.** Every record carries `sourceId`, `validFrom` and dataset `version`, so "where did this come from, and when was it true?" always has an answer.
- **No lock-in.** Apache-2.0, plain JSON in [`data/codes.json`](data/codes.json), self-hostable from the committed [`Dockerfile`](Dockerfile). If this project stops, the data does not become unreachable.
- **Independence.** Built and operated outside government, so it can be honest about gaps rather than presenting an incomplete list as a complete register.

### What this is not

Drawn from the stated non-goals in [`agent_plan.md`](agent_plan.md), the publication decisions in [`docs/governance/source-register.json`](docs/governance/source-register.json), and the product boundary in [`docs/adr/0001-product-boundary.md`](docs/adr/0001-product-boundary.md):

- **Not a perpetually current MMDA list.** The 2021 namespace is published as the 2021 namespace. It is not maintained as "today's districts".
- **Not a scrape of the Ministry directory.** The Ministry of Local Government, Chieftaincy and Religious Affairs MMDA directory is registered as a *review reference only* and is never silently imported into canonical data.
- **Not a fuzzy matcher.** There is no name search and no approximate matching. Aliases must exist as explicit, versioned records.
- **Not an open namespace.** No institution, education, tax, health or private identifier namespace is admitted without its own authority and licence review.
- **Not a redistribution of the source PDF.** The official document is linked, never re-hosted.
- **Not a geography service.** The dataset carries no boundaries, geometry, population or coordinates.
- **Not a government system.** See the independence note at the end of this file.

## Quickstart

Prerequisites: Node 24, [pnpm](https://pnpm.io) 11.11.0, Go 1.23, and Ruby 3.4 if you want to run the foundation validator. These are the exact versions pinned in [`package.json`](package.json), [`go.mod`](go.mod) and [`.github/workflows/quality.yml`](.github/workflows/quality.yml).

```sh
git clone https://github.com/stanleyHayes/ghanacodes.git
cd ghanacodes
pnpm install --frozen-lockfile

# 1. Public register and browser resolver on http://localhost:3000
pnpm dev

# 2. API on http://localhost:8080 (run from the repository root; the
#    service loads data/codes.json relative to the working directory)
go run ./cmd/api
curl -s http://localhost:8080/health
```

Set `PORT` to move the API off 8080. Production builds the same binary using the command in [`render.yaml`](render.yaml).

## Usage

### Namespaces

Every call names one of the three declared namespaces. `GET /v1/namespaces` returns them with the dataset version.

| Namespace ID | Entity | Codes | Example |
|---|---|---|---|
| `gss-phc-2021-region` | region | 16 | `03` → Greater Accra |
| `gss-phc-2021-region-abbreviation` | region | 16 | `GAR` → Greater Accra |
| `gss-phc-2021-district` | district | 261 | `0101` → Jomoro |

### Endpoints

| Method | Path | Purpose |
|---|---|---|
| `GET` | `/health` | Dataset version and loaded code count |
| `GET` | `/v1/namespaces` | Declared namespaces for the loaded dataset |
| `GET` | `/v1/resolve/{namespace}/{code}?at=YYYY-MM-DD` | Exact, date-aware resolution |
| `GET` | `/v1/crosswalk?fromNamespace=&code=&toNamespace=&at=` | Crosswalk through the canonical entity ID |
| `POST` | `/v1/bulk-resolve` | Ordered bulk resolution; max 1000 queries, 1 MiB body |
| `POST` | `/graphql` | Constrained GraphQL: `resolveCode` and `crosswalk` only |

Omitting `at` resolves as at today (UTC). The committed contracts are [`contracts/openapi.yaml`](contracts/openapi.yaml) — currently path declarations only, without response schemas, and not covering `/health` or `/graphql` — and [`contracts/graphql/schema.graphql`](contracts/graphql/schema.graphql), whose `Resolution` type does not yet declare the `query` field that the JSON responses return.

### Response shape

Every resolve and crosswalk response is a `Resolution`:

| Field | Meaning |
|---|---|
| `status` | `resolved` (exactly one active record), `not_found` (none), or `ambiguous` (more than one) |
| `query` | The parameters as received, echoed back |
| `candidates` | Every matching active record, sorted deterministically. Empty when `not_found` |
| `version` | The dataset version that answered the query |

Each candidate carries `namespace`, `code`, `name`, `entityId`, `entityType`, `validFrom`, `validTo`, `status` and `sourceId`; districts add `districtType` and `regionCode`, and regions may carry `aliases`. The Go types are in [`internal/codes/resolver.go`](internal/codes/resolver.go) and a partial TypeScript type in [`sdk/typescript/index.ts`](sdk/typescript/index.ts), which currently omits `districtType`, `regionCode` and `aliases`.

`ambiguous` is a first-class outcome, not an error path: when several records are valid for a key at the queried date, all are returned and none is preferred. Within the single pinned 2021 dataset every namespace key has exactly one active record, so no query against the published beta returns `ambiguous` today. The behaviour exists for future dataset versions and is covered by `TestAmbiguityReturnsCandidates` in [`internal/codes/resolver_test.go`](internal/codes/resolver_test.go).

### Crosswalk

```sh
curl -s "https://api-codes.digitalghana.dev/v1/crosswalk?fromNamespace=gss-phc-2021-region&code=03&toNamespace=gss-phc-2021-region-abbreviation"
```

Both `03` and `GAR` resolve to `geo:region:03`, so the crosswalk is a walk through that entity ID rather than a string comparison of "Greater Accra".

### Bulk

```sh
curl -s -X POST https://api-codes.digitalghana.dev/v1/bulk-resolve \
  -H 'content-type: application/json' \
  -d '{"queries":[
        {"namespace":"gss-phc-2021-district","code":"0102","at":"2021-06-27"},
        {"namespace":"gss-phc-2021-region","code":"03","at":"2021-06-27"}]}'
```

Results are returned in `results`, in the order the queries were submitted, alongside the dataset `version`.

### Calling from a browser

The API sends a fixed `Access-Control-Allow-Origin` of `https://codes.digitalghana.dev`. Server-side and CLI callers are unaffected; browser code on other origins should call through your own backend, or embed [`data/codes.json`](data/codes.json) directly as the public site does.

## Data and provenance

| Property | Value |
|---|---|
| Dataset version | `gss-phc-2021.1` |
| Effective from | `2021-06-27` |
| Records | 293 codes — 16 region codes, 16 region abbreviations, 261 district codes |
| Primary source | Ghana Statistical Service, *2021 Population and Housing Census Field Officer's Manual* |
| Source SHA-256 | `06bb5ffc8cfeb59df612a10432209de82ddd8988d4e2fbb8c75e359dc2968560` |
| Source review date | 2026-09-01 |
| Publication decision | Derived code facts with attribution; the PDF itself is linked, not redistributed |

The full register of sources, licences and publication decisions is [`docs/governance/source-register.json`](docs/governance/source-register.json), validated against [`docs/governance/source-register.schema.json`](docs/governance/source-register.schema.json).

**Licence boundary.** The original code and configuration in this repository are Apache-2.0. The referenced source documents are not: the GSS manual states no open-data licence, so only derived code facts are published, with attribution, and the PDF is never re-hosted. Nothing referenced here is relicensed by inclusion.

**Regeneration.** Obtain the official PDF, then:

```sh
ruby scripts/import_gss_2021.rb /path/to/official.pdf
```

The importer aborts on a checksum mismatch, on a district count other than 261, on duplicate codes, and on an empty district name. It requires `pdftotext` on the path.

One district name (`0201`) is normalised by an explicit `NAME_OVERRIDES` entry in the importer because the source layout does not extract cleanly; every other name is taken verbatim from Appendix 2.

**Corrections.** Open an issue or pull request with the affected code, the current and proposed value, the authoritative source and its publication date, and whether historical records change — see [`CONTRIBUTING.md`](CONTRIBUTING.md). A human reviewer approves canonical publication; corrections are never auto-merged.

## Project layout

```
app/         Next.js public register and browser resolver (Outfit/Geist Mono/Newsreader)
cmd/api/     Go HTTP service: REST, bulk and the constrained GraphQL endpoint
internal/    Resolver, crosswalk, effective-date logic and its Go tests
contracts/   Committed OpenAPI and GraphQL contracts
data/        codes.json — the generated, pinned dataset
docs/        ADR, governance source register, operations and release-evidence runbooks
infra/       Vercel configuration and security headers
scripts/     Reproducible GSS importer and the foundation validator
sdk/         TypeScript client source (not yet published as a package)
tests/       Node dataset invariant tests
```

## Verification

These are the commands CI runs, in order, in [`.github/workflows/quality.yml`](.github/workflows/quality.yml):

```sh
ruby scripts/validate.rb   # required governance files, source register, template/secret scan
go test ./...              # resolver, crosswalk, effective-date and ambiguity behaviour
go vet ./...
pnpm typecheck
pnpm test                  # dataset invariants: counts, uniqueness, provenance
pnpm build
```

If `ruby scripts/validate.rb` aborts with `invalid byte sequence in US-ASCII`, run it as `RUBYOPT="-E UTF-8" ruby scripts/validate.rb`; see [`CONTRIBUTING.md`](CONTRIBUTING.md).

Immutable deployment, smoke and rollback evidence is recorded in [`docs/runbooks/release-evidence.md`](docs/runbooks/release-evidence.md).

## Status and roadmap

**Public beta.** Web and API are live on their canonical hosts with recorded smoke and rollback evidence. Namespace and source administration, conflict review, audited import approval, broader version history, package publication, and security/load evidence remain stable gates.

See [`ROADMAP.md`](ROADMAP.md) for what is shipped, what is next and what is deliberately out of scope. The machine-readable ledger is [`agent_plan.md`](agent_plan.md).

## Contributing, conduct, security and licence

- [`CONTRIBUTING.md`](CONTRIBUTING.md) — setup, verification, commit conventions, data corrections
- [`CODE_OF_CONDUCT.md`](CODE_OF_CONDUCT.md) — Contributor Covenant v2.1; reports to the private channel documented there
- [`SECURITY.md`](SECURITY.md) — report privately to the private channel documented there, never in a public issue
- [`AGENTS.md`](AGENTS.md) — coordination rules for automated contributors
- [`LICENSE`](LICENSE) and [`NOTICE`](NOTICE) — Apache License 2.0

## Independence

GhanaCodes is an independent open-source project in the [Digital Ghana](https://digitalghana.dev) portfolio. It is **not** operated by, affiliated with, or endorsed by the Government of Ghana, the Ghana Statistical Service, or any ministry, department or agency. It proxies no government transaction, credential, form or payment. Official records remain with their issuing authorities; this project publishes only derived, source-linked identifier facts.
