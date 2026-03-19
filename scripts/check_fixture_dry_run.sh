#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

if ! command -v gh >/dev/null 2>&1; then
  echo "gh is required" >&2
  exit 1
fi

if [[ -z "${GH_TOKEN:-}" ]]; then
  echo "GH_TOKEN must be set" >&2
  exit 1
fi

BIN="${ROOT_DIR}/ai-reviewer"

echo "Building ai-reviewer..."
go build -o "$BIN" .

strip_ansi() {
  LC_ALL=C LANG=C perl -pe 's/\e\[[0-9;]*[A-Za-z]//g'
}

assert_contains() {
  local haystack="$1"
  local needle="$2"
  if ! grep -Fq -- "$needle" <<<"$haystack"; then
    echo "Expected output to contain: $needle" >&2
    exit 1
  fi
}

assert_not_contains() {
  local haystack="$1"
  local needle="$2"
  if grep -Fq -- "$needle" <<<"$haystack"; then
    echo "Expected output NOT to contain: $needle" >&2
    exit 1
  fi
}

run_dry_run() {
  local pr_number="$1"
  echo
  echo "Running dry-run for fixture PR #$pr_number..."
  "$BIN" pr gonka-ai/ai-reviewer-fixture "$pr_number" --dry-run | strip_ansi
}

pr2_output="$(run_dry_run 2)"
printf '%s\n' "$pr2_output"

assert_contains "$pr2_output" "Loading config from: ${ROOT_DIR}/.ai-review/gonka-ai/ai-reviewer-fixture/config.yaml"
assert_not_contains "$pr2_output" ".repos/gonka-ai/ai-reviewer-fixture/.ai-review/config.yaml"
assert_contains "$pr2_output" "- branch-reviewer (reviewer)"
assert_contains "$pr2_output" "- date-reviewer (reviewer)"
assert_contains "$pr2_output" "- docs-reviewer (reviewer)"
assert_contains "$pr2_output" "- frontend-reviewer (reviewer)"
assert_contains "$pr2_output" "- function-reviewer (reviewer)"
assert_contains "$pr2_output" "- go-reviewer (reviewer)"
assert_contains "$pr2_output" "- legacy-reviewer (reviewer)"
assert_contains "$pr2_output" "- line-window-reviewer (reviewer)"
assert_contains "$pr2_output" "- regex-reviewer (reviewer)"
assert_contains "$pr2_output" "with primer: api-guardrails (matches 1 files)"
assert_contains "$pr2_output" "with primer: legacy-warning (matches 1 files)"

pr1_output="$(run_dry_run 1)"
printf '%s\n' "$pr1_output"

assert_contains "$pr1_output" "- docs-reviewer (reviewer)"
assert_contains "$pr1_output" "- frontend-reviewer (reviewer)"
assert_contains "$pr1_output" "- legacy-reviewer (reviewer)"
assert_contains "$pr1_output" "To be skipped (no matching files):"
assert_contains "$pr1_output" "- branch-reviewer"
assert_contains "$pr1_output" "- date-reviewer"
assert_contains "$pr1_output" "- function-reviewer"
assert_contains "$pr1_output" "- go-reviewer"
assert_contains "$pr1_output" "- line-window-reviewer"
assert_contains "$pr1_output" "- regex-reviewer"
assert_contains "$pr1_output" "with primer: legacy-warning (matches 1 files)"
assert_not_contains "$pr1_output" "- branch-reviewer (reviewer)"
assert_not_contains "$pr1_output" "- function-reviewer (reviewer)"

echo
echo "Fixture dry-run checks passed."
