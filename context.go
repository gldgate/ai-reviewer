package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type PRInfo struct {
	Title       string `json:"title"`
	Body        string `json:"body"`
	BaseRefName string `json:"baseRefName"`
	BaseRefOid  string `json:"baseRefOid"`
	HeadRefName string `json:"headRefName"`
	HeadRefOid  string `json:"headRefOid"`
	IsCommit    bool
}

type PRContext struct {
	Title        string
	Description  string
	ChangedFiles []string
	Diff         string
}

func GetPRInfo(repo, prNumber string) (*PRInfo, error) {
	fmt.Printf("    -> Running gh pr view %s...\n", prNumber)
	cmd := exec.Command("gh", "pr", "view", prNumber, "-R", repo, "--json", "title,body,baseRefName,baseRefOid,headRefName,headRefOid")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error running gh pr view: %w, output: %s", err, string(output))
	}

	var pr PRInfo
	if err := json.Unmarshal(output, &pr); err != nil {
		return nil, fmt.Errorf("error unmarshaling gh output: %w", err)
	}

	return &pr, nil
}

func GetCommitInfo(commitHash, compareTo string) (*PRInfo, error) {
	fmt.Printf("    -> Getting info for commit %s...\n", commitHash)

	// Get commit message (title and body)
	cmd := exec.Command("git", "show", "-s", "--format=%s%n%n%b", commitHash)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error getting commit message: %w, output: %s", err, string(output))
	}

	fullMsg := strings.TrimSpace(string(output))
	parts := strings.SplitN(fullMsg, "\n", 2)
	title := parts[0]
	body := ""
	if len(parts) > 1 {
		body = strings.TrimSpace(parts[1])
	}

	// Get full SHA for head
	cmd = exec.Command("git", "rev-parse", commitHash)
	headSHA, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error resolving commit hash: %w", err)
	}

	baseSHA := compareTo
	if baseSHA == "" {
		// Default to parent commit
		cmd = exec.Command("git", "rev-parse", commitHash+"^")
		out, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("error getting parent commit: %w", err)
		}
		baseSHA = strings.TrimSpace(string(out))
	} else {
		// Resolve compareTo to full SHA
		cmd = exec.Command("git", "rev-parse", compareTo)
		out, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("error resolving comparison commit: %w", err)
		}
		baseSHA = strings.TrimSpace(string(out))
	}

	return &PRInfo{
		Title:       title,
		Body:        body,
		BaseRefOid:  baseSHA,
		HeadRefOid:  strings.TrimSpace(string(headSHA)),
		IsCommit:    true,
		BaseRefName: baseSHA[:8], // Short SHA for display
		HeadRefName: strings.TrimSpace(string(headSHA))[:8],
	}, nil
}

func GetPRContext(prInfo *PRInfo, includeFilters, excludeFilters []string) (*PRContext, error) {
	diff, err := GetDiff(prInfo.BaseRefOid, prInfo.HeadRefOid, includeFilters, excludeFilters)
	if err != nil {
		return nil, err
	}

	files, err := GetChangedFiles(prInfo.BaseRefOid, prInfo.HeadRefOid, includeFilters, excludeFilters)
	if err != nil {
		return nil, err
	}

	return &PRContext{
		Title:        prInfo.Title,
		Description:  prInfo.Body,
		ChangedFiles: files,
		Diff:         diff,
	}, nil
}

func GetDiff(baseSHA, headSHA string, includeFilters, excludeFilters []string) (string, error) {
	// Triple-dot (A...B) means diff from common ancestor of A and B to B.
	// Double-dot (A..B) means diff from A to B.
	// For PRs we usually want triple-dot.
	// For "ai-review commit" we probably want double-dot if a base is specified,
	// or triple-dot if comparing to parent (which ends up being same as double-dot).
	// Let's use triple-dot as it's generally safer for PR-like workflows.
	args := []string{"diff", fmt.Sprintf("%s...%s", baseSHA, headSHA)}
	if len(includeFilters) > 0 || len(excludeFilters) > 0 {
		args = append(args, "--")
		for _, f := range includeFilters {
			args = append(args, f)
		}
		for _, f := range excludeFilters {
			args = append(args, ":(exclude)"+f)
		}
	}

	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error running git diff: %w, output: %s", err, string(output))
	}
	return AnnotateDiff(string(output)), nil
}

var hunkHeaderRegexp = regexp.MustCompile(`^@@ -\d+(?:,\d+)? \+(\d+)(?:,\d+)? @@`)

func AnnotateDiff(diff string) string {
	var result strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(diff))
	currentLine := 0

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "@@ ") {
			matches := hunkHeaderRegexp.FindStringSubmatch(line)
			if len(matches) > 1 {
				startLine, _ := strconv.Atoi(matches[1])
				currentLine = startLine
			}
			result.WriteString(line + "\n")
		} else if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++ ") {
			result.WriteString(fmt.Sprintf("%d:%s\n", currentLine, line))
			currentLine++
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "--- ") {
			result.WriteString(fmt.Sprintf("%d:%s\n", currentLine, line))
		} else if strings.HasPrefix(line, " ") {
			result.WriteString(fmt.Sprintf("%d:%s\n", currentLine, line))
			currentLine++
		} else {
			result.WriteString(line + "\n")
		}
	}

	return result.String()
}

func GetChangedFiles(baseSHA, headSHA string, includeFilters, excludeFilters []string) ([]string, error) {
	args := []string{"diff", "--name-only", fmt.Sprintf("%s...%s", baseSHA, headSHA)}
	if len(includeFilters) > 0 || len(excludeFilters) > 0 {
		args = append(args, "--")
		for _, f := range includeFilters {
			args = append(args, f)
		}
		for _, f := range excludeFilters {
			args = append(args, ":(exclude)"+f)
		}
	}

	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error running git diff --name-only: %w, output: %s", err, string(output))
	}
	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(files) == 1 && files[0] == "" {
		return []string{}, nil
	}
	return files, nil
}

func buildPrompt(p Persona, ctx *PRContext, globalInstructions string, preRunAnalyses map[string][]string, summary string) string {
	var fileList strings.Builder
	for _, file := range ctx.ChangedFiles {
		fileList.WriteString(file)
		if len(p.IncludeExplainers) > 0 {
			if analyses, ok := preRunAnalyses[file]; ok {
				for _, analysis := range analyses {
					// Check if this analysis is from one of the included explainers
					// Analysis format is "PersonaID: description" (based on pipeline.go and how preRunAnalyses is populated)
					parts := strings.SplitN(analysis, ": ", 2)
					if len(parts) > 0 {
						explainerID := parts[0]
						included := false
						for _, id := range p.IncludeExplainers {
							if id == explainerID {
								included = true
								break
							}
						}
						if included {
							fileList.WriteString(fmt.Sprintf("\n  - Explainer Analysis: %s", analysis))
						}
					}
				}
			}
		}
		fileList.WriteString("\n")
	}

	findingsText := ""
	if p.IncludeFindings && summary != "" {
		findingsText = fmt.Sprintf("\n\n--- AGGREGATED REPORT ---\n%s\n", summary)
	}

	diffSection := ""
	if p.ExcludeDiff {
		// Calculate diff stats
		addedLines := 0
		deletedLines := 0
		scanner := bufio.NewScanner(strings.NewReader(ctx.Diff))
		for scanner.Scan() {
			line := scanner.Text()
			// Extract original diff line (after line number and colon)
			parts := strings.SplitN(line, ":", 2)
			if len(parts) < 2 {
				continue
			}
			diffLine := parts[1]
			if strings.HasPrefix(diffLine, "+") && !strings.HasPrefix(diffLine, "+++ ") {
				addedLines++
			} else if strings.HasPrefix(diffLine, "-") && !strings.HasPrefix(diffLine, "--- ") {
				deletedLines++
			}
		}
		diffSection = fmt.Sprintf("\nDiff Stats: %d files changed, %d lines added, %d lines deleted. (Full diff excluded by configuration)\n", len(ctx.ChangedFiles), addedLines, deletedLines)
	} else {
		fence := "```"
		diffSection = fmt.Sprintf(`
This is a unified diff format where each line is prefixed with a line number (e.g., "123:"). Lines starting with "-" after the line number indicate removed lines, lines starting with "+" after the line number indicate added lines, and lines starting with " " (space) are unchanged context lines. Diff chunks begin with "@@ -old_start,old_count +new_start,new_count @@" headers that may include function/class context, and are preceded by +++ lines indicating the file being modified.

Unified Diff:
%sdiff
%s
%s
`, fence, ctx.Diff, fence)
	}

	prompt := fmt.Sprintf(`%s
%s
PR Title: %s
PR Description: %s

Changed Files:
%s
%s
`, p.Instructions, findingsText, ctx.Title, ctx.Description, fileList.String(), diffSection)

	if globalInstructions != "" {
		prompt += fmt.Sprintf("\nStandard Instructions:\n%s\n", globalInstructions)
	}

	if p.Role == "explainer" && p.Stage == "pre" {
		prompt = PreRunExplainerSystemPrompt + "\n\n" + prompt
	}

	return prompt
}
