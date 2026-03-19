---
id: fixture-integrity
model_category: balanced
path_filters:
  - ".ai-review/gonka-ai/ai-reviewer-fixture/**"
  - ".github/**/*.yml"
  - "scripts/check_fixture_dry_run.sh"
  - "*.go"
exclude_filters:
  - "models.go"
---
You are reviewing whether the ai-reviewer fixture setup still proves what it claims to prove.

Focus on:
- fixture personas or config that no longer match scanner behavior
- dry-run assertions that are too weak or accidentally asserting the wrong thing
- fixture PR expectations that no longer separate the intended reviewer sets
- CI checks that might pass while the real feature regresses

Only report issues that undermine confidence in the fixture-based validation.
