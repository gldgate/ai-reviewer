---
id: gldgate-trust-boundaries
type: caution
path_filters:
  - "./canisters/adapter/**"
  - "./canisters/token_adapter/**"
  - "./typescript/transport-node/**"
  - "./typescript/icr-cli/**"
---
Treat the transport node as untrusted input.

Critical boundaries:
- Adapter should trust only bootstrap anchors and cryptographic proofs, not operator claims.
- Authorization lists are primarily economic policy boundaries under the IC reverse-gas model: authorized transport nodes get free-call behavior, while non-authorized callers remain allowed but must pay gas/fees correctly.
- Review changes to ensure open-call access is preserved where intended, and gas charging is applied correctly and cannot be bypassed by unauthorized callers.
- Replay protection by `tx_id` dedup is mandatory for both mint and withdraw paths.
- Any anonymous mode toggle (`allow_anonymous`) must be scrutinized for abuse paths.

Only report issues with clear impact: unauthorized mint/burn, replay, broken proof validation, or route bypass.
