---
id: bridge-authority-and-escape-hatches
model_category: best_code
path_filters:
  - "inference-chain/contracts/wrapped-token/src/**.rs"
  - "inference-chain/contracts/liquidity-pool/src/**.rs"
  - "inference-chain/contracts/community-sale/src/**.rs"
exclude_filters:
  - "**/lib.rs"
  - "**/state.rs"
---
Review privileged control paths and escape hatches.

Check:
- admin, creator, governance, and caller authorization on pause, resume, update, withdraw, emergency-withdraw, metadata, mint, burn, and migrate flows
- fallback initialization paths that can assign authority incorrectly during instantiate or migration
- delegation into `cw20-base`, bank sends, and stargate/Any messages where a wrapper may unintentionally widen permissions or bypass local checks
- privilege splits that are ambiguous enough to cause operator mistakes, permanent lockout, or unauthorized asset movement

Only surface issues with major operational impact: privilege escalation, unintended custody loss, or dangerous authority ambiguity.
