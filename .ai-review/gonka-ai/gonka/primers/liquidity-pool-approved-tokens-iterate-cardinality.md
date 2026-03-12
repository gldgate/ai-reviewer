---
id: liquidity-pool-approved-tokens-iterate-cardinality
type: caution
regex_filters:
  - "LiquidityPoolApprovedTokensMap\\.Iterate"
---
`LiquidityPoolApprovedTokensMap.Iterate` is only acceptable because this allowlist is expected to stay very small in production, roughly a dozen entries or fewer.

Reviewers should treat every use of this iterator as a linear scan over a tiny governance-curated set, not as a scalable lookup pattern. Findings are warranted if a change:

- makes the approved-token set likely to grow materially beyond that rough size bound
- adds new `Iterate` call sites on hot execution paths
- replaces direct keyed lookups with iteration when the caller already has the full `(chainId, contractAddress)` key
- relies on iteration without preserving the assumption that approvals remain rare and operationally curated

If a feature needs this map to hold many more entries, the implementation should move to an indexed or directly keyed lookup strategy rather than depending on repeated full scans.
