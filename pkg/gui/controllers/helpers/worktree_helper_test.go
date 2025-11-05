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
			name:     "remote branch",
			base:     "refs/remotes/origin/feature/EMP-1234",
			expected: filepath.FromSlash("origin/feature/EMP-1234") + sep,
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
				helper = &WorktreeHelper{c: &HelperCommon{Common: common}}
			} else {
				helper = &WorktreeHelper{}
			}

			actual := helper.defaultWorktreePath(test.base)
			if actual != test.expected {
				t.Fatalf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}
