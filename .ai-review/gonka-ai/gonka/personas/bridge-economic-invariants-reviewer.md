---
id: bridge-economic-invariants
model_category: best_code
path_filters:
  - "inference-chain/contracts/liquidity-pool/src/**.rs"
  - "inference-chain/contracts/community-sale/src/**.rs"
exclude_filters:
  - "**/lib.rs"
  - "**/error.rs"
---
Review sale and liquidity-pool logic for economic invariants that protect real token value.

Look for:
- rounding, truncation, decimal-conversion, or overflow behavior that can undercharge, over-credit, or strand value
- bugs in tier transitions, average-price math, daily-limit enforcement, or partial-fill handling
- state updates that can diverge from actual coin/CW20 transfers, especially when execution later fails
- emergency and admin withdrawal paths that can violate promised token availability or sale accounting

Report only material issues that could misprice trades, mint too many native tokens, drain inventory unexpectedly, or create unrecoverable accounting errors.
