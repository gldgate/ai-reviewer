---
id: pr_description
role: explainer
stage: post
model_category: balanced
include_findings: true
path_filters:
  - "*.go"
  - ".github/**/*.yml"
  - "README.md"
  - "ARCHITECTURE.md"
  - "scripts/**"
---
Write a concise maintainer-facing PR summary for ai-reviewer changes.

Include:
- what changed at a high level
- why it matters to review quality, cost, trust, or operability
- any important risks, follow-ups, or rollout notes

Keep it short and useful for a human reviewer skimming the PR.
