# Vision

## Why This Exists

AI code review is one of the most promising uses of language models in software engineering.

That opportunity is getting bigger, not smaller. More code is being written with AI assistance, which means human reviewers are increasingly asked to evaluate changes that were produced faster than they can comfortably inspect in depth. The bottleneck is no longer only writing code. It is reviewing enough of it, well enough, before it ships.

The goal of `ai-reviewer` is not to replace human judgment. The goal is to make human review scale better, with more consistency and better project awareness.

## The Problem With Many AI Review Tools

There are many strong commercial and open source tools in this space, and they have helped prove that AI review can be genuinely useful. But many approaches also fall into one of a few patterns:

- one large generic reviewer prompt that tries to do everything at once
- AI-flavored linting that is good at local issues but weak at deeper domain reasoning
- very broad repository context injection that can become slow, expensive, and difficult to control
- opaque review behavior that is hard to inspect, tune, or debug

Those approaches can be helpful, but they often leave teams with one of two tradeoffs:

- high cost and large token usage for broad context
- lower cost but weaker signal because the model is not given the right project-specific lens

For systems with meaningful domain complexity, those tradeoffs become more painful.

## The Core Bet

`ai-reviewer` is built on a different bet:

**Code review is not one job. It is an orchestration problem.**

A good review is usually the product of multiple specialized perspectives:

- security
- correctness
- architecture
- performance
- testing
- operational safety
- project-specific invariants
- historical and version-specific context

Instead of asking one generic reviewer to do all of that at once, `ai-reviewer` breaks review into smaller, more focused parts and composes them into a pipeline.

## What That Means In Practice

The system is designed around a few core ideas.

### Specialized personas

Each persona is meant to review with a narrow purpose. A persona should answer a specific question well, not perform a vague imitation of “senior engineer review.”

This makes prompts clearer and gives the system better odds of producing findings that are actually actionable.

### Repo-aware context

Review policy should live with the repository being reviewed.

That is why the system favors:

- repo-scoped local artifacts under `.ai-review/<owner>/<repo>/...`
- versioned Markdown artifacts in the repository itself via `ai_review` frontmatter

This allows primers, personas, waivers, and operative documentation to evolve with the codebase instead of existing as ambient machine-local behavior.

### Structured review pipeline

The pipeline is intentionally multi-stage:

1. gather context
2. run explainers and reviewers
3. normalize raw output into structured findings
4. apply waivers as review policy
5. aggregate the remaining findings into a concise report

The point is not complexity for its own sake. The point is to separate concerns so each step can be simpler, cheaper, and more inspectable.

### Cost consciousness

Large context windows are useful, but they are not free.

The architecture tries to keep token use under control by:

- filtering personas to relevant files
- filtering primers to relevant changes
- allowing different model categories for different jobs
- using cheaper normalization and policy steps where possible

The goal is not “use the smallest model possible.” The goal is to spend tokens where they create the most review value.

### Auditability

AI review systems are only trustworthy if teams can inspect what happened.

That is why `ai-reviewer` persists prompts, raw outputs, normalized findings, summaries, and reports as artifacts. If a review is confusing, expensive, or wrong, there should be a paper trail.

## What We Want From AI Review

We want AI review to be:

- high-signal
- explicit
- configurable
- repo-aware
- cost-conscious
- inspectable

We do not want it to be:

- ambient
- magical
- monolithic
- difficult to debug
- detached from the codebase’s actual rules

## Why Repo-Scoped Review Logic Matters

Different repositories have different concerns. A blockchain repository, a web app, an SDK, and a data pipeline should not all be reviewed as if the same checklist applies equally well to each.

That is why `ai-reviewer` treats review logic as something that can be authored close to the code:

- primers can explain architecture or constraints
- personas can focus on project-specific risks
- waivers can encode policy about known exceptions

This is especially important for projects like Gonka, where architectural invariants and domain rules matter more than generic style advice.

## Long-Term Direction

The long-term idea is not just “run LLMs on diffs.”

The long-term idea is to build a review system where:

- review intelligence can evolve with the codebase
- different concerns are handled by the right specialized lenses
- output is structured enough to be aggregated and audited
- teams can choose how much cost and depth they want

There is also room for shared review logic in the future, but it should be explicit and intentional, not ambient. A possible future direction is importable shared persona packs rather than globally loaded local files.

## Bottom Line

The vision is simple:

**AI review should help human reviewers keep up with what is shipping, without turning review into an expensive black box.**

That means focused reviewers, project-aware context, explicit policy, and artifacts that make the system understandable when it succeeds and debuggable when it fails.
