---
id: gldgate-dispatcher-registry-invariants
type: caution
path_filters:
  - "./canisters/dispatcher/**"
  - "./canisters/token_adapter/**"
  - "./typescript/icr-cli/scripts/dispatcher.ts"
  - "./config/**"
---
Dispatcher is the single source of truth for chain/token routing.

Review invariants:
- `(chain_id, token)` must map deterministically to a ledger canister.
- Chain-level and token-level pause flags must be enforced by token routing paths.
- Controller-only admin paths (`upsert_chain`, `register_token`, pause/resume) must remain tightly scoped.
- Disabled chain behavior must be equivalent to non-operational routing.

Findings should prioritize wrong-ledger routing, pause bypass, or privilege mistakes.
