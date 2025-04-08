package git

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const (
	DefaultBranch     = "main"
	DefaultRemote     = "origin"
	DefaultTimeFormat = "02-01-2006_15:04:05"
)

type Git struct {
	repoPath string
}

func New(repoPath string) *Git {
	return &Git{
		repoPath: repoPath,
	}
}

func (g *Git) HasChanges() (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = g.repoPath
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("error while checking for changes: %w", err)
	}
	return len(output) > 0, nil
}

func (g *Git) AddAll() error {
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = g.repoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error while adding changes: %w", err)
	}
	return nil
}

func (g *Git) Commit() error {
	now := time.Now().Format(DefaultTimeFormat)
	cmd := exec.Command("git", "commit", "--quiet", "-m", now)
	cmd.Dir = g.repoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error while creating commit: %w", err)
	}
	return nil
}

func (g *Git) Push() error {
	cmd := exec.Command("git", "push", "--quiet", DefaultRemote, DefaultBranch)
	cmd.Dir = g.repoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error while pushing to GitHub: %w", err)
	}
	return nil
}

func (g *Git) GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = g.repoPath
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error while getting current branch: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func (g *Git) Fetch() error {
	cmd := exec.Command("git", "fetch", DefaultRemote, DefaultBranch)
	cmd.Dir = g.repoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error while fetching: %w", err)
	}
	return nil
}

func (g *Git) Pull() error {
	cmd := exec.Command("git", "pull", DefaultRemote, DefaultBranch)
	cmd.Dir = g.repoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error while pulling: %w", err)
	}
	return nil
}

func (g *Git) GetConflicts() ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", "--diff-filter=U")
	cmd.Dir = g.repoPath
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error while getting conflicts: %w", err)
	}

	if len(output) == 0 {
		return []string{}, nil
	}

	return strings.Split(strings.TrimSpace(string(output)), "\n"), nil
}
