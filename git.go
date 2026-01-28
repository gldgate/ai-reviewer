package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func EnsureRepo(repo string) error {
	repoParts := strings.Split(repo, "/")
	repoName := repoParts[len(repoParts)-1]

	// 1. Check if current directory is a git repo and seems related
	if isGitRepo() {
		if isRelatedRepo(repo, repoName) {
			return nil
		}
	}

	// 2. Check if it exists in .repos/repo
	reposDir := ".repos"
	targetDir := filepath.Join(reposDir, repo)

	if info, err := os.Stat(targetDir); err == nil && info.IsDir() {
		if err := os.Chdir(targetDir); err == nil {
			if isGitRepo() {
				return nil
			}
			// If it's a directory but not a git repo, back out
			// We need to know where we came from.
			// Actually, EnsureRepo is called early, so we can probably just assume we want to be in targetDir if it's a git repo.
			// If not, we might have an issue.
		}
	}

	// 3. Clone
	fmt.Printf("--- Cloning %s into %s...\n", repo, targetDir)
	if err := os.MkdirAll(filepath.Dir(targetDir), 0755); err != nil {
		return fmt.Errorf("failed to create directory for clone: %w", err)
	}

	// gh repo clone <repo> <directory>
	cmd := exec.Command("gh", "repo", "clone", repo, targetDir)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to clone repo: %w, output: %s", err, string(out))
	}

	if err := os.Chdir(targetDir); err != nil {
		return fmt.Errorf("failed to change directory to %s: %w", targetDir, err)
	}
	return nil
}

func isGitRepo() bool {
	return exec.Command("git", "rev-parse", "--is-inside-work-tree").Run() == nil
}

func isRelatedRepo(repo, repoName string) bool {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	url := strings.TrimSpace(string(out))
	// Matches if it's the exact repo or at least has the same repo name (likely a fork)
	return strings.Contains(url, repo) || strings.Contains(url, "/"+repoName)
}

func FetchRefs(repo, prNumber, baseRef string) error {
	remote := fmt.Sprintf("https://github.com/%s.git", repo)

	// Fetch base branch
	fmt.Printf("    -> Fetching base branch %s...\n", baseRef)
	cmd := exec.Command("git", "fetch", remote, baseRef)
	if _, err := cmd.CombinedOutput(); err != nil {
		// If it fails, try fetching origin as fallback
		exec.Command("git", "fetch", "origin", baseRef).Run()
	}

	// Fetch PR head
	// GitHub exposes PRs at refs/pull/ID/head
	fmt.Printf("    -> Fetching PR head refs/pull/%s/head...\n", prNumber)
	cmd = exec.Command("git", "fetch", remote, fmt.Sprintf("refs/pull/%s/head", prNumber))
	if out, err := cmd.CombinedOutput(); err != nil {
		// Fallback to origin if remote URL fails
		cmd = exec.Command("git", "fetch", "origin", fmt.Sprintf("refs/pull/%s/head", prNumber))
		if out2, err2 := cmd.CombinedOutput(); err2 != nil {
			return fmt.Errorf("failed to fetch PR head from %s or origin: %v (out1: %s, out2: %s)", remote, err2, string(out), string(out2))
		}
	}
	return nil
}
