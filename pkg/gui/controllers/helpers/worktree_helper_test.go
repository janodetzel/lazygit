package helpers

import (
	"path/filepath"
	"testing"
)

func TestDefaultWorktreePath(t *testing.T) {
	type testCase struct {
		name     string
		base     string
		expected string
	}

	sep := string(filepath.Separator)

	tests := []testCase{
		{
			name:     "simple branch",
			base:     "feature/EMP-1234",
			expected: filepath.Join("..", "feature/EMP-1234") + sep,
		},
		{
			name:     "branch with refs heads prefix",
			base:     "refs/heads/feature/EMP-1234",
			expected: filepath.Join("..", "feature/EMP-1234") + sep,
		},
		{
			name:     "remote branch",
			base:     "refs/remotes/origin/feature/EMP-1234",
			expected: filepath.Join("..", "origin/feature/EMP-1234") + sep,
		},
		{
			name:     "commit hash fallback",
			base:     "abc123",
			expected: filepath.Join("..", "abc123") + sep,
		},
		{
			name:     "empty base",
			base:     "",
			expected: ".." + sep,
		},
	}

	helper := &WorktreeHelper{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := helper.defaultWorktreePath(test.base)
			if actual != test.expected {
				t.Fatalf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}
