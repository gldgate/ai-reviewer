---
id: params-files-testermint
model_category: balanced
path_filters:
  - inference-chain/proto/inference/inference/params.proto
  - testermint/src/main/kotlin/data/AppExport.kt
---
If parameters are changed in `inference-chain/proto/inference/inference/params.proto`, look to make sure that the params Kotlin classes for testermint are updated accordingly. The classes file is located at `testermint/src/main/kotlin/data/AppExport.kt`. 

ALL parameters from the proto need to be in the Kotlin classes or Governance tests will fail.

Instructions:
1. Check if `inference-chain/proto/inference/inference/params.proto` is among the changed files.
2. If it is, verify if `testermint/src/main/kotlin/data/AppExport.kt` is also changed and if it correctly reflects the changes in the proto.
3. If `params.proto` changed but `AppExport.kt` did not, flag this as a potential issue.
4. If only `AppExport.kt` changed, you can generally ignore it unless you see something obviously wrong related to parameters.
