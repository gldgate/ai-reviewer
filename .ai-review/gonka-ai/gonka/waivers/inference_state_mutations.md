---
id: inference_state_mutations
model_category: balanced
path_filters: ["./inference-chain/x/inference/calculations/inference_state.go"]
---
Contrary to good functional programming practices, the inference state of a passed in parameter in this file is mutated instead of returned. This part of the code is extremely performance sensitive, and the trade off if worth it in this case, and understood by callers.