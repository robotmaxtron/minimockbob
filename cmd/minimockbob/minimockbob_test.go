package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestCLI tests the CLI functionality binary before running any tests
func TestCLI(t *testing.T) {
	cmd := exec.Command("go", "build", "-o", "minimockbob")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build binary: %v\nOutput: %s", err, output)
	}
	t.Log("Binary built successfully")
}

// TestCLIQuotedArgument tests the CLI with a quoted argument
func TestCLIQuotedArgument(t *testing.T) {
	cmd := exec.Command("./minimockbob", "Hello, World!")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}
	got := strings.TrimSpace(string(output))
	want := "hElLo, WoRlD!"
	if got != want {
		t.Errorf("./minimockbob \"Hello, World!\" = %q, want %q", got, want)
	}
}

// TestCLIMultipleArguments tests the CLI with multiple unquoted arguments
func TestCLIMultipleArguments(t *testing.T) {
	cmd := exec.Command("./minimockbob", "Hello", "World")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}
	got := strings.TrimSpace(string(output))
	want := "hElLo WoRlD"
	if got != want {
		t.Errorf("./minimockbob Hello World = %q, want %q", got, want)
	}
}

// TestCLIPipeInput tests the CLI with pipe input
func TestCLIPipeInput(t *testing.T) {
	cmd := exec.Command("sh", "-c", "echo \"Hello, World!\" | ./minimockbob")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}
	got := strings.TrimSpace(string(output))
	want := "hElLo, WoRlD!"
	if got != want {
		t.Errorf("echo \"Hello, World!\" | ./minimockbob = %q, want %q", got, want)
	}
}

// TestCLIPipeInputVariations tests various piped inputs
func TestCLIPipeInputVariations(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple text", "hello", "hElLo"},
		{"with spaces", "foo bar", "fOo BaR"},
		{"with punctuation", "foo, bar!", "fOo, BaR!"},
		{"with numbers", "abc123def", "aBc123DeF"},
		{"empty line handling", "line1\nline2", "lInE1\nlInE2"},
		{"multiple spaces", "a  b  c", "a  B  c"},
		{"special chars basic", "test!@#$%", "tEsT!@#$%"},
		{"special chars extended", "test!@#$%^&*()", "tEsT!@#$%^&*()"},
		{"question and exclamation", "hello?world!", "hElLo?WoRlD!"},
		{"hyphens and underscores", "a-b_c=d+e", "a-B_c=D+e"},
		{"brackets and braces", "test[abc]def", "tEsT[aBc]DeF"},
		{"colon and semicolon", "a:b;c", "a:B;c"},
		{"tilde", "a~b", "a~B"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("sh", "-c", "printf '%s' '"+tt.input+"' | ./minimockbob")
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("Command failed: %v\nOutput: %s", err, output)
			}
			got := strings.TrimSpace(string(output))
			if got != tt.want {
				t.Errorf("piped %q = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestCLIPipeEmptyInput tests piping empty input
func TestCLIPipeEmptyInput(t *testing.T) {
	// Note: The current implementation handles empty pipe gracefully
	// This test documents that behavior - empty input shows usage
	cmd := exec.Command("sh", "-c", "printf '' | ./minimockbob")
	output, err := cmd.CombinedOutput()
	got := string(output)

	// Empty piped input should either fail or show usage
	if err != nil {
		// Expected: command exits with error
		if !strings.Contains(got, "Usage:") {
			t.Errorf("Expected usage message for empty pipe, got: %s", got)
		}
	} else {
		// If it doesn't error, output should be empty or usage
		if got != "" && !strings.Contains(got, "Usage:") {
			t.Errorf("Expected empty output or usage for empty pipe, got: %s", got)
		}
	}
}

// TestCLIPipeMultiline tests piping multiline input
func TestCLIPipeMultiline(t *testing.T) {
	cmd := exec.Command("sh", "-c", "printf 'Hello\\nWorld' | ./minimockbob")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}
	got := strings.TrimSpace(string(output))
	want := "hElLo\nWoRlD"
	if got != want {
		t.Errorf("multiline pipe = %q, want %q", got, want)
	}
}

// TestCLINoArguments tests the CLI with no arguments (should show usage)
func TestCLINoArguments(t *testing.T) {
	cmd := exec.Command("./minimockbob")
	output, err := cmd.CombinedOutput()
	// This should exit with error code 1
	if err == nil {
		t.Fatal("Expected command to fail with no arguments, but it succeeded")
	}
	got := string(output)
	if !strings.Contains(got, "Usage:") {
		t.Errorf("Expected usage message, got: %s", got)
	}
}

// TestCLIEmptyString tests the CLI with an empty string
func TestCLIEmptyString(t *testing.T) {
	cmd := exec.Command("./minimockbob", "")
	output, err := cmd.CombinedOutput()
	// This should exit with error code 1
	if err == nil {
		t.Fatal("Expected command to fail with empty string, but it succeeded")
	}
	got := string(output)
	if !strings.Contains(got, "Usage:") {
		t.Errorf("Expected usage message, got: %s", got)
	}
}

// TestCLISpecialCharacters tests the CLI with special characters via arguments
func TestCLISpecialCharacters(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"foo bar", "fOo BaR"},
		{"foo, bar?", "fOo, BaR?"},
		{"123", "123"},
		{"AB", "aB"},
		{"hello?", "hElLo?"},
		{"world!", "wOrLd!"},
		{"hello? world!", "hElLo? WoRlD!"},
		{"test@#$%", "tEsT@#$%"},
		{"a+b=c", "a+B=c"},
		{"x[y]z", "x[Y]z"},
		{"m&n*o", "m&N*o"},
		{"a:b;c", "a:B;c"},
		{"p-q_r", "p-Q_r"},
		{"s/t\\u", "s/T\\u"},
		{"v<w>x", "v<W>x"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			cmd := exec.Command("./minimockbob", tt.input)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("Command failed: %v\nOutput: %s", err, output)
			}
			got := strings.TrimSpace(string(output))
			if got != tt.want {
				t.Errorf("./minimockbob %q = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestCLICleanup tests cleaning up the binary after tests
func TestCLICleanup(t *testing.T) {
	// Remove the binary if it exists
	err := os.Remove("./minimockbob")
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("Failed to remove binary: %v", err)
	}
	t.Log("Binary cleaned up successfully")
}
