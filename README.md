# GhanaCodes

`GhanaCodes` is an independent Digital Ghana interoperability resolver for versioned, source-linked Ghanaian identifiers.

- Web and resolver: <https://codes.digitalghana.dev>
- REST, GraphQL and bulk API: <https://api-codes.digitalghana.dev>
- Initial dataset: GSS 2021 PHC regions and all 261 district codes

## Safety model

Every request names a namespace and may name an effective date. Crosswalks use canonical entity IDs and return all candidates when ambiguous. GhanaCodes never silently fuzzy-matches a district or institution name.

## Verification

Run `ruby scripts/validate.rb`, `go test ./...`, `go vet ./...`, `pnpm typecheck`, `pnpm test`, and `pnpm build`.
