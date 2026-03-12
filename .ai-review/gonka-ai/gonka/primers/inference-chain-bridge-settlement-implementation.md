---
id: inference-chain-bridge-settlement-implementation
type: implementation
path_filters:
  - "inference-chain/contracts/wrapped-token/src/**.rs"
  - "inference-chain/contracts/liquidity-pool/src/**.rs"
  - "inference-chain/contracts/community-sale/src/**.rs"
---
There are several unit and trust-boundary conventions reviewers should keep in mind.

USD values are generally represented with 6 decimals. Native-token sale quantities are generally represented with 9 decimals. Wrapped CW20 tokens and IBC tokens may use other precisions and are normalized before pricing logic runs. Changes to multiplication, division order, or integer truncation can materially change who wins the rounding.

These contracts rely on chain-module gRPC queries such as `ValidateWrappedTokenForTrade`, `ValidateIbcTokenForTrade`, and `ApprovedTokensForTrade`. Those queries are part of the security boundary: review whether failures are handled conservatively and whether query outputs are cross-checked against the asset actually received.

`wrapped-token` also emits a stargate `MsgRequestBridgeWithdrawal` after burning supply. Review any change to the burn-before-message or message-construction sequence as a bridge-settlement change, not as ordinary metadata plumbing.
