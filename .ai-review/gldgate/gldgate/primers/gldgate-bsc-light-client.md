---
id: gldgate-bsc-light-client
type: implementation
path_filters:
  - "./canisters/adapter/**"
  - "./typescript/icr-cli/**"
  - "./typescript/transport-node/**"
  - "./docs/bsc-light-client.md"
---
BSC light-client review cues:

- Bootstrap artifact is the trust anchor; post-bootstrap claims must be proof-verified.
- Post-Bohr validator derivation must come from contract state proofs against epoch `stateRoot`.
- Pre-Bohr vs post-Bohr branching must remain unambiguous (avoid mixed-source acceptance).
- Header RLP encoding order changes are consensus-critical and must preserve hash parity.
- Archive-RPC assumptions (`eth_getProof` at historical blocks) are mandatory for zero-trust updates.

Prioritize high-impact issues (invalid proof acceptance, validator set corruption, epoch progression bypass), but also report medium-impact correctness risks that could evolve into safety failures.
