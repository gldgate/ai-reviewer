---
id: gldgate-deposit-flow-implementation
type: implementation
path_filters:
  - "./canisters/adapter/**"
  - "./canisters/token_adapter/**"
  - "./canisters/dispatcher/**"
  - "./typescript/transport-node/**"
  - "./typescript/web-ui/**"
  - "./evm/contracts/**"
---
Deposit flow invariants (release 0.1):

1. `adapter.submit_message` must validate and deduplicate before persisting.
2. Mint path must go through `token_adapter.balance_up` with authorized caller checks.
3. `token_adapter` must resolve `(chain_id, token)` via dispatcher and reject paused/disabled routes.
4. Ledger minting must preserve replay protection using unique `tx_id` records.
5. UI history and balance updates are eventually consistent; avoid assumptions of immediate settlement visibility.
6. EVM-side deposit intent must be verifiably transported to ICP by `transport-node`; missing delivery is a correctness failure.
7. Transport submission failures must be retried safely until terminal success/failure, using idempotent semantics (`tx_id` dedup) to prevent double mint.
8. ICP-side processing should be atomic from an observable state perspective: no partial success where one canister persists and another fails; on failure, state must roll back to pre-transaction state.
9. The main business result is exact balance-up on ICP ledger/accounting state (no loss, duplication, or drift from the validated EVM deposit amount).

When changes touch this flow, prioritize correctness of ordering, dedup, authorization, retry/error handling, and atomic settlement guarantees over style concerns.
