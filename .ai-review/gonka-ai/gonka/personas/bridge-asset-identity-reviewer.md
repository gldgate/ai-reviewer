---
id: bridge-asset-identity
model_category: best_code
path_filters:
  - "inference-chain/contracts/wrapped-token/src/**.rs"
  - "inference-chain/contracts/liquidity-pool/src/**.rs"
  - "inference-chain/contracts/community-sale/src/**.rs"
exclude_filters:
  - "**/lib.rs"
  - "**/error.rs"
---
Review bridge token identity and acceptance logic for failures that could let the wrong external asset pass as legitimate value.

Focus on:
- inconsistent normalization of chain IDs, EVM contract addresses, `cw20:` prefixes, or IBC denoms
- trusting user-provided identifiers where the contract should rely on queried bridge metadata or module-approved token lists
- mismatches between stored bridge metadata, query results, and the token actually transferred or burned
- migrations or config updates that can silently desynchronize acceptance checks from previously issued wrapped assets

Prioritize only issues that could cause spoofed deposits, acceptance of the wrong asset, or permanent bridge-accounting divergence.
