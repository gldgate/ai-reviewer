---
id: gldgate-service-resolver-ttl
type: implementation
path_filters:
  - "./typescript/icr-service-resolver/**"
  - "./canisters/factory/**"
  - "./canisters/gateway/**"
  - "./canisters/adapter/**"
  - "./canisters/token_adapter/**"
---
Service resolution must preserve factory-driven cache semantics:

- Runtime source of truth is `factory.resolve_services(app_id)`.
- Cache entries must respect factory `ttl_seconds` and expire predictably.
- Query handlers must not attempt inter-canister refreshes; update calls should refresh.
- Manual controller overrides are exceptional and intentionally non-expiring.

Flag any change that can cause stale routing, split-brain service IDs, or query/update inconsistency.
