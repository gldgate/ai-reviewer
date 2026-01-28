package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type PRInfo struct {
	Title       string `json:"title"`
	Body        string `json:"body"`
	BaseRefName string `json:"baseRefName"`
	BaseRefOid  string `json:"baseRefOid"`
	HeadRefName string `json:"headRefName"`
	HeadRefOid  string `json:"headRefOid"`
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
	// Use ... for triple-dot diff (find common ancestor)
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
	return string(output), nil
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

func buildPrompt(p Persona, ctx *PRContext, globalInstructions string, preRunAnalyses map[string][]string) string {
	var fileList strings.Builder
	for _, file := range ctx.ChangedFiles {
		fileList.WriteString(file)
		if analyses, ok := preRunAnalyses[file]; ok {
			for _, analysis := range analyses {
				fileList.WriteString(fmt.Sprintf("\n  - Explainer Analysis: %s", analysis))
			}
		}
		fileList.WriteString("\n")
	}

	prompt := fmt.Sprintf(`%s

PR Title: %s
PR Description: %s

Changed Files:
%s

Unified Diff:
%s
`, p.Instructions, ctx.Title, ctx.Description, fileList.String(), ctx.Diff)

	if globalInstructions != "" {
		prompt += fmt.Sprintf("\nStandard Instructions:\n%s\n", globalInstructions)
	}

	if p.Role == "explainer" && p.Stage == "pre" {
		prompt = PreRunExplainerSystemPrompt + "\n\n" + prompt
	}

	return prompt
}
