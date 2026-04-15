---
id: gldgate-ledger-mint-burn
type: caution
path_filters:
  - "./canisters/token_adapter/**"
  - "./canisters/ledgers/**"
  - "./docs/usdt-ledger.md"
---
Ledger and token-adapter safety checks:

- Minting account authority must remain restricted to token_adapter workflow.
- `balance_up` and `withdraw` must validate caller permissions and enforce unique `tx_id`.
- Mint/burn operations must respect dispatcher pause state before touching ledger balances.
- Memo handling should preserve traceability without violating ICRC limits.
- U256 ledger assumptions should not be weakened by type narrowing in adapters or scripts.
- Mint/burn paths must be overflow-safe end-to-end (no unchecked arithmetic, no unsafe narrowing casts, no aggregate amount wraparound).
- Cross-canister balance flows should be transactional from an observable state perspective: if a downstream canister call fails, there must be a safe rollback/compensation path so partial state does not persist.

Prioritize issues that can mis-account supply, duplicate mint/burn, or bypass authorization.
