# GhanaCodes beta release evidence

Verified: 2026-09-01<br>
Release commit: `d576403f6b0cab1ae7dc0652b764b6dd689a515c`<br>
Dataset: `gss-phc-2021.1`

## Web

- Vercel project: `hayfordstanleys-projects/ghanacodes`
- Deployment: `dpl_8m9yeNpWPt6hFf2SDV7AZZEqf31u`
- Immutable URL: `ghanacodes-msvmf7h8a-hayfordstanleys-projects.vercel.app`
- Canonical URL: <https://codes.digitalghana.dev>
- `/`, `/sitemap.xml`, and `/robots.txt`: HTTP 200

## API

- Render service: `srv-dabdm0ss728c73adbong`
- Deploy: `dep-dabdm1ks728c73adbqeg`
- Provider URL: `ghanacodes-api.onrender.com`
- Canonical URL: <https://api-codes.digitalghana.dev>
- Health reports 293 codes and dataset version `gss-phc-2021.1`.
- Custom-domain certificate: `CN=api-codes.digitalghana.dev`, Google Trust Services `WE1`, valid at verification.

## Production smoke

- GSS district `0101` at `2021-06-27` resolves exactly to Jomoro and `geo:district:0101`.
- The same code at `2020-01-01` returns `not_found`, proving effective-date enforcement.
- Region code `03` crosswalks to abbreviation `GAR` through `geo:region:03`.
- A two-query bulk request preserved input order and returned Ellembelle followed by Greater Accra.
- GraphQL `resolveCode` matched REST status and dataset version.
- GitHub Quality run `33516970581` passed for the release commit.

## Source and rollback

The generated dataset is reproducible only from the official GSS PDF with SHA-256 `06bb5ffc8cfeb59df612a10432209de82ddd8988d4e2fbb8c75e359dc2968560`; the source artifact is linked and not redistributed.

Rollback uses the immutable Vercel deployment and Render deploy/provider URL recorded above. A replacement must pass the same exact, historical, crosswalk, bulk, GraphQL, health, and TLS smoke before DNS changes.

## Beta limitations

Namespace/source administration, conflict review, audited import approval, broader historical namespace coverage, package publication, and mature security/load/source-stewardship evidence remain stable gates.

## UI and discoverability addendum

Commit `7ef19244935b4c6c3ab10b6213822dd5172746c1`, CI `33527986949`, and deployment `dpl_CLH7w8FDL1aLGVuhJUzxSjvRWh2d` replaced the native namespace select with an accessible three-option Radix listbox and aligned typography to Outfit, Geist Mono and Newsreader. Browser inspection found no prohibited native control or overflow. The canonical page now exposes a product favicon, manifest, canonical/Open Graph/Twitter metadata, and a verified 1200x630 PNG.
