---
id: gldgate-architecture-overview
type: overview
path_filters:
  - "./canisters/**"
  - "./typescript/**"
  - "./evm/**"
  - "./config/**"
---
# GldGate architecture overview

GldGate is split into four planes:

- UX plane: `web_ui` + Internet Identity + wallet integrations.
- IC execution plane: `adapter` -> `token_adapter` -> `dispatcher` -> ICRC ledgers.
- Off-chain verification plane: `transport-node` + EVM RPC + chain contracts.
- Control plane: `factory` + resolver cache for runtime service discovery.

The core deposit path is EVM deposit event -> transport-node claim -> `adapter` validation/dedup -> `token_adapter.balance_up` -> ledger mint transfer.

When reviewing PRs, anchor findings to this boundary model so issues are evaluated in the correct plane.
