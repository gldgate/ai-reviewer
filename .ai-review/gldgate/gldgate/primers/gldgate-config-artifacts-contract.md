---
id: gldgate-config-artifacts-contract
type: caution
path_filters:
  - "./config/**"
  - "./scripts/**"
  - "./typescript/icr-cli/**"
  - "./artifacts/**"
---
Deployment config contract:

- `config/<target>/<chain>/...` files are manual inputs.
- `artifacts/<target>/...` files are generated outputs and should not be hand-edited.
- Secrets must stay under `secrets/` and out of tracked config/artifacts.
- Deploy/resume/update modes must avoid accidental `factory.spawn` re-creation when updating existing apps.
- Chain-support-only deploys must preserve existing app IDs and service mappings.

Prioritize findings that can cause environment drift, wrong service wiring, or secret exposure.
