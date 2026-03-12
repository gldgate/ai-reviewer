---
id: inference-chain-bridge-contracts-overview
type: overview
path_filters:
  - "inference-chain/contracts/**.rs"
---
`inference-chain/contracts` contains three distinct CosmWasm contracts that participate in bridge-related token flows.

`wrapped-token` is a CW20 wrapper around an external asset. It stores bridge origin metadata (`chain_id`, external `contract_address`), delegates standard token behavior to `cw20-base`, and adds a custom withdraw path that burns wrapped supply and emits a stargate bridge-withdrawal message.

`liquidity-pool` accepts approved bridge assets and sells native Gonka tokens using tiered pricing and daily limits. It depends on external gRPC queries to decide whether a CW20 or IBC token is approved and to learn token decimal precision for USD normalization.

`community-sale` is a more tightly scoped sale contract for a designated buyer. It also accepts bridge-originated value, but its pricing is fixed rather than tiered and its acceptance rules can be configured to allow either a single expected asset or any module-approved trade token.

Reviewers should treat these contracts as custody and settlement code. Small logic mistakes here can turn into wrong-asset acceptance, permanent supply divergence, or operator-only recovery paths.
