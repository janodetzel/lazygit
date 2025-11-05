package helpers

import (
	"path/filepath"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
)

func TestDefaultWorktreePath(t *testing.T) {
	type testCase struct {
		name              string
		base              string
		expected          string
		worktreeParentDir string
		remoteNames       []string
	}

	sep := string(filepath.Separator)

	tests := []testCase{
		{
			name:     "simple branch",
			base:     "feature/EMP-1234",
			expected: filepath.FromSlash("feature/EMP-1234") + sep,
		},
		{
			name:     "branch with refs heads prefix",
			base:     "refs/heads/feature/EMP-1234",
			expected: filepath.FromSlash("feature/EMP-1234") + sep,
		},
		{
			name:        "remote branch with refs prefix",
			base:        "refs/remotes/origin/feature/EMP-1234",
			expected:    filepath.FromSlash("feature/EMP-1234") + sep,
			remoteNames: []string{"origin"},
		},
		{
			name:        "remote branch without refs prefix",
			base:        "origin/feature/EMP-2222",
			expected:    filepath.FromSlash("feature/EMP-2222") + sep,
			remoteNames: []string{"origin"},
		},
		{
			name:     "commit hash fallback",
			base:     "abc123",
			expected: "abc123" + sep,
		},
		{
			name:     "empty base",
			base:     "",
			expected: "." + sep,
		},
		{
			name:              "custom parent directory",
			base:              "feature/EMP-5678",
			worktreeParentDir: ".worktrees",
			expected:          filepath.Join(".worktrees", filepath.FromSlash("feature/EMP-5678")) + sep,
		},
		{
			name:              "custom parent directory with trailing slash",
			base:              "feature/EMP-9012",
			worktreeParentDir: "../custom/",
			expected:          filepath.Join("..", "custom", filepath.FromSlash("feature/EMP-9012")) + sep,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var helper *WorktreeHelper
			if test.worktreeParentDir != "" {
				cfg := config.GetDefaultConfig()
				cfg.Git.WorktreeParentDir = test.worktreeParentDir
				common := &common.Common{}
				common.SetUserConfig(cfg)
				helper = &WorktreeHelper{c: &HelperCommon{Common: common}, remoteNamesCache: test.remoteNames}
			} else {
				helper = &WorktreeHelper{remoteNamesCache: test.remoteNames}
			}

			actual := helper.defaultWorktreePath(test.base)
			if actual != test.expected {
				t.Fatalf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}

func TestSanitizeBranchName(t *testing.T) {
	tests := []struct {
		name     string
		base     string
		remotes  []string
		expected string
	}{
		{
			name:     "local branch",
			base:     "feature/EMP-1234",
			expected: "feature/EMP-1234",
		},
		{
			name:     "refs heads",
			base:     "refs/heads/feature/EMP-1234",
			expected: "feature/EMP-1234",
		},
		{
			name:     "remote branch with refs prefix",
			base:     "refs/remotes/origin/feature/EMP-1234",
			remotes:  []string{"origin"},
			expected: "feature/EMP-1234",
		},
		{
			name:     "remote branch without refs prefix",
			base:     "origin/feature/EMP-1234",
			remotes:  []string{"origin"},
			expected: "feature/EMP-1234",
		},
		{
			name:     "remote branch unknown remote",
			base:     "upstream/feature/EMP-1234",
			remotes:  []string{"origin"},
			expected: "upstream/feature/EMP-1234",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := sanitizeBranchName(test.base, test.remotes)
			if actual != test.expected {
				t.Fatalf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}
