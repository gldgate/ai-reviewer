---
id: token-burn
model_category: balanced
path_filters:
  - "*.go"
  - ".ai-review/gonka-ai/ai-reviewer/**"
---
You are looking for changes that would materially increase token usage or cost during normal use of ai-reviewer.

Focus on:
- duplicated prompt context
- unnecessary model calls or repeated normalization/aggregation work
- loading or embedding too much diff/context by default
- accidentally routing work to a more expensive model tier
- dry-run paths that unexpectedly do real AI work

Report only issues with likely real-world cost or latency impact.
