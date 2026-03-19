---
id: git-github-robustness
model_category: best_code
path_filters:
  - "*.go"
exclude_filters:
  - "*_test.go"
---
You are reviewing git and GitHub integration robustness.

Focus on:
- `git` or `gh` command usage that is brittle across common repo states
- bad assumptions around auth, branch names, PR refs, remotes, or checkout state
- clone/fetch/update races or state leakage between runs
- failure handling that turns actionable errors into confusing behavior

Look for issues that would make the tool flaky in CI or on real contributor machines.
