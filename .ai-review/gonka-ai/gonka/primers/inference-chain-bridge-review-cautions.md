---
id: inference-chain-bridge-review-cautions
type: caution
path_filters:
  - "inference-chain/contracts/**.rs"
---
Be conservative about anything that changes token identity, authority, or value conversion.

High-risk areas in this folder include:
- `Receive`, `PurchaseWithNative`, and custom withdraw flows
- bridge metadata storage and queries like `BridgeInfo`, `ValidateWrappedTokenForTrade`, and `ValidateIbcTokenForTrade`
- admin or creator updates, emergency withdrawals, and migrations
- decimal normalization and tier math

Prefer findings only when the consequence is major: wrong-asset acceptance, unauthorized mint/burn/withdraw, broken bridge withdrawal requests, sale mispricing, or state that can no longer reconcile with the assets actually held.
