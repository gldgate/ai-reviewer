---
id: waive-branch-subnet-only-issues
model_category: fastest_good
branch_filters:
  - "gm/inference"
path_filters:
  - "**subnet**"
---
Waive findings on this branch when all of the following are true:

- the issue is confined to subnet-specific code in the matched files
- the consequence is limited to the behavior, correctness, availability, accounting, or operator workflow of the subnet feature itself
- the issue does not create a plausible impact on overall chain operation, global chain safety, consensus behavior, shared bank or module accounting outside subnet flows, unrelated inference features, or general transaction processing
- the analysis assumes subnet features are protected by allow lists, so findings that depend on unrestricted public access to subnet creation, subnet settlement, or other subnet-only entry points should be suppressed on that basis

Do not waive a finding if the changed code might affect any non-subnet path, shared keeper behavior, shared params, shared stores, message routing, authz, fee handling, token movement outside subnet-scoped funds, state growth or pruning with chain-wide consequences, node stability, or any behavior that could degrade the chain even when subnet usage is allow-listed.

When in doubt, raise the finding instead of waiving it.
