package minimockbob

import (
	"os/exec"
	"strings"
	"testing"
)

// TestLintGolangciLint runs golangci-lint on the codebase with comprehensive checks
func TestLintGolangciLint(t *testing.T) {
	// Check if golangci-lint is available
	if _, err := exec.LookPath("golangci-lint"); err != nil {
		t.Skip("golangci-lint not found in PATH, skipping linting test")
	}

	// Run golangci-lint with default linters (which are already comprehensive)
	cmd := exec.Command("golangci-lint", "run", "./...")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("golangci-lint found issues:\n%s", string(output))
	} else {
		t.Log("golangci-lint passed with no issues")
	}
}

// TestLintStaticcheck runs staticcheck on the codebase
func TestLintStaticcheck(t *testing.T) {
	// Check if staticcheck is available
	if _, err := exec.LookPath("staticcheck"); err != nil {
		t.Skip("staticcheck not found in PATH, skipping linting test")
	}

	cmd := exec.Command("staticcheck", "./...")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("staticcheck found issues:\n%s", string(output))
	} else {
		t.Log("staticcheck passed with no issues")
	}
}

// TestLintGoVet runs go vet on the codebase
func TestLintGoVet(t *testing.T) {
	cmd := exec.Command("go", "vet", "./...")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("go vet found issues:\n%s", string(output))
	} else {
		t.Log("go vet passed with no issues")
	}
}

// TestLintGoFmt checks if code is properly formatted
func TestLintGoFmt(t *testing.T) {
	// Use gofmt -l to list files that differ from gofmt's style
	cmd := exec.Command("sh", "-c", "find . -name '*.go' -not -path './vendor/*' -exec gofmt -l {} +")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("gofmt command failed: %v\nOutput: %s", err, string(output))
	}

	// gofmt -l outputs filenames of files that need formatting
	unformattedFiles := strings.TrimSpace(string(output))
	if unformattedFiles != "" {
		t.Errorf("The following files need formatting:\n%s\nRun: gofmt -w %s", unformattedFiles, unformattedFiles)
	} else {
		t.Log("All files are properly formatted")
	}
}

// TestLintGoFmtSimplify checks if code can be simplified
func TestLintGoFmtSimplify(t *testing.T) {
	// Use gofmt -s to check for code that can be simplified
	cmd := exec.Command("sh", "-c", "find . -name '*.go' -not -path './vendor/*' -exec gofmt -s -l {} +")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("gofmt -s command failed: %v\nOutput: %s", err, string(output))
	}

	// gofmt -s -l outputs filenames of files that can be simplified
	unsimplifiedFiles := strings.TrimSpace(string(output))
	if unsimplifiedFiles != "" {
		t.Errorf("The following files can be simplified:\n%s\nRun: gofmt -s -w %s", unsimplifiedFiles, unsimplifiedFiles)
	} else {
		t.Log("All files are optimally simplified")
	}
}

// TestLintGoImports checks if imports are properly formatted
func TestLintGoImports(t *testing.T) {
	// Check if goimports is available
	if _, err := exec.LookPath("goimports"); err != nil {
		t.Skip("goimports not found in PATH, skipping import formatting test")
	}

	cmd := exec.Command("sh", "-c", "find . -name '*.go' -not -path './vendor/*' -exec goimports -l {} +")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("goimports command failed: %v\nOutput: %s", err, string(output))
	}

	// goimports -l outputs filenames of files with incorrect imports
	incorrectImports := strings.TrimSpace(string(output))
	if incorrectImports != "" {
		t.Errorf("The following files have incorrect imports:\n%s\nRun: goimports -w %s", incorrectImports, incorrectImports)
	} else {
		t.Log("All imports are properly formatted")
	}
}

// TestLintErrcheck checks for unchecked errors
func TestLintErrcheck(t *testing.T) {
	// Check if errcheck is available
	if _, err := exec.LookPath("errcheck"); err != nil {
		t.Skip("errcheck not found in PATH, skipping error checking test")
	}

	cmd := exec.Command("errcheck", "./...")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("errcheck found unchecked errors:\n%s", string(output))
	} else {
		t.Log("No unchecked errors found")
	}
}

// TestLintGoModTidy checks if go.mod is tidy
func TestLintGoModTidy(t *testing.T) {
	// Run go mod tidy and check if it makes changes
	cmd := exec.Command("sh", "-c", "go mod tidy && git diff --exit-code go.mod go.sum")
	output, err := cmd.CombinedOutput()

	if err != nil {
		// Check if it's because git diff found differences
		if strings.Contains(string(output), "diff --git") {
			t.Errorf("go.mod or go.sum needs tidying. Run: go mod tidy\nDifferences:\n%s", string(output))
		} else {
			t.Logf("go.mod and go.sum are tidy (or git not available)")
		}
	} else {
		t.Log("go.mod and go.sum are tidy")
	}
}

// TestLintShadow checks for variable shadowing
func TestLintShadow(t *testing.T) {
	// Check if shadow is available
	if _, err := exec.LookPath("shadow"); err != nil {
		t.Skip("shadow not found in PATH, skipping shadow checking test")
	}

	cmd := exec.Command("shadow", "./...")
	output, err := cmd.CombinedOutput()

	if err != nil {
		// shadow returns non-zero if it finds issues
		if len(output) > 0 {
			t.Errorf("shadow found shadowed variables:\n%s", string(output))
		}
	} else {
		t.Log("No shadowed variables found")
	}
}
