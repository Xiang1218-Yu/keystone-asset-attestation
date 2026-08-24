# Keystone Asset Attestation

Keystone coordinates supplier lineage, industrial inspection campaigns, maintenance evidence, risk posture and release decisions for complex assets.

This is a standalone Go service using only the standard library. It keeps a
small in-memory domain store so the complete lifecycle can be explored without
external infrastructure.

## Run

```bash
go run ./cmd/keystone-asset-attestation
```

The default address is `:8282`; set `SERVICE_ADDR` to override it.

## Quick flow

```bash
curl http://localhost:8282/healthz
curl -X POST http://localhost:8282/v1/records \
  -H 'content-type: application/json' \
  -d '{"id":"demo-1","payload":"asset turbine supplier lineage inspection baseline"}'
curl -X POST http://localhost:8282/v1/records/demo-1/advance \
  -H 'content-type: application/json' \
  -d '{"stage":"verified"}'
```

The project has no test files by design; package compilation, vetting and
runtime API checks are used for validation.
