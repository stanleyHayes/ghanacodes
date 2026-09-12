# Operations baseline

Web and API are live in public beta on `codes.digitalghana.dev` and `api-codes.digitalghana.dev`; deployment, smoke, TLS and rollback evidence is recorded in [`release-evidence.md`](release-evidence.md).

Open before stable — define and verify:

- health, readiness and dependency checks;
- structured logs with request/correlation identifiers and secret redaction;
- latency, error, saturation and domain-correctness signals;
- alert routes with an acknowledged owner;
- backup/restore where state exists;
- least-privilege provider credentials and rotation;
- rate limits, abuse controls and audit trails for privileged publication;
- rollback and incident communication steps.
