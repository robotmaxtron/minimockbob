package minimockbob

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// errorReader always returns an error on Read
type errorReader struct{}

func (e errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated read error")
}

func TestGen(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// Basic cases
		{"Hello, World!", "hElLo, WoRlD!"},
		{"foobar", "fOoBaR"},
		{"Foobar", "fOoBaR"},
		{"foo bar", "fOo BaR"},
		{"foo, bar?", "fOo, BaR?"},
		{"f", "f"},
		{"", ""},
		{"AB", "aB"},
		{"123", "123"},

		// Edge cases
		{"a", "a"},
		{"A", "a"},
		{"aB", "aB"},
		{"Ab", "aB"},

		// Multiple spaces
		{"a  b  c", "a  B  c"},
		{"   leading", "   lEaDiNg"},
		{"trailing   ", "tRaIlInG   "},

		// Special characters
		{"!@#$%^&*()", "!@#$%^&*()"},
		{"test!@#", "tEsT!@#"},
		{"a-b_c", "a-B_c"},
		{"x[y]z", "x[Y]z"},
		{"m&n*o", "m&N*o"},
		{"a:b;c", "a:B;c"},
		{"p+q=r", "p+Q=r"},
		{"s/t\\u", "s/T\\u"},
		{"v<w>x", "v<W>x"},
		{"y{z}a", "y{Z}a"},
		{"test'quote", "tEsT'qUoTe"},
		{"test\"quote", "tEsT\"qUoTe"},

		// Numbers mixed with letters
		{"abc123def456ghi", "aBc123DeF456gHi"},
		{"123abc", "123aBc"},
		{"abc123", "aBc123"},
		{"1a2b3c", "1a2B3c"},

		// Punctuation
		{"hello!", "hElLo!"},
		{"hello?", "hElLo?"},
		{"hello.", "hElLo."},
		{"hello,world", "hElLo,WoRlD"},
		{"a.b.c.d", "a.B.c.D"},

		// Unicode and international characters
		{"café", "cAfÉ"},
		{"naïve", "nAïVe"},
		{"Ñoño", "ñOñO"},

		// Newlines and tabs
		{"hello\nworld", "hElLo\nWoRlD"},
		{"hello\tworld", "hElLo\tWoRlD"},
		{"a\nb\nc", "a\nB\nc"},

		// Long strings
		{"abcdefghijklmnopqrstuvwxyz", "aBcDeFgHiJkLmNoPqRsTuVwXyZ"},
		{"ABCDEFGHIJKLMNOPQRSTUVWXYZ", "aBcDeFgHiJkLmNoPqRsTuVwXyZ"},

		// Mixed case scenarios
		{"HeLLo WoRLd", "hElLo WoRlD"},
		{"ALLCAPS", "aLlCaPs"},
		{"alllower", "aLlLoWeR"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := Gen(tt.input); got != tt.want {
				t.Errorf("Gen(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestGenConsistency ensures Gen produces consistent results
func TestGenConsistency(t *testing.T) {
	input := "Hello, World!"
	want := "hElLo, WoRlD!"

	// Run multiple times to ensure consistency
	for i := 0; i < 100; i++ {
		got := Gen(input)
		if got != want {
			t.Errorf("Gen(%q) iteration %d = %q, want %q", input, i, got, want)
		}
	}
}

// TestGenAlternation ensures a proper alternation pattern
func TestGenAlternation(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"even letters", "abcd", "aBcD"},
		{"odd letters", "abcde", "aBcDe"},
		{"with non-letters", "a1b2c3d", "a1B2c3D"},
		{"letters at start and end", "xabc123defx", "xAbC123dEfX"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Gen(tt.input)
			if got != tt.want {
				t.Errorf("Gen(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestGenStartsWithLower ensures the first letter is always lowercase
func TestGenStartsWithLower(t *testing.T) {
	tests := []string{
		"Hello",
		"HELLO",
		"hELLO",
		"A",
		"Z",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			got := Gen(input)
			if got == "" {
				t.Fatalf("Gen(%q) returned empty string", input)
			}
			firstChar := rune(got[0])
			if firstChar >= 'A' && firstChar <= 'Z' {
				t.Errorf("Gen(%q) = %q, first letter should be lowercase", input, got)
			}
		})
	}
}

// BenchmarkGenShort benchmarks Gen with short strings
func BenchmarkGenShort(b *testing.B) {
	input := "Hello, World!"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Gen(input)
	}
}

// BenchmarkGenMedium benchmarks Gen with medium-length strings
func BenchmarkGenMedium(b *testing.B) {
	input := "The quick brown fox jumps over the lazy dog. Pack my box with five dozen liquor jugs!"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Gen(input)
	}
}

// BenchmarkGenLong benchmarks Gen with long strings
func BenchmarkGenLong(b *testing.B) {
	input := strings.Repeat("abcdefghijklmnopqrstuvwxyz ", 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Gen(input)
	}
}

// BenchmarkGenUnicode benchmarks Gen with unicode strings
func BenchmarkGenUnicode(b *testing.B) {
	input := "café naïve Ñoño résumé"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Gen(input)
	}
}

// BenchmarkGenAllCaps benchmarks Gen with all uppercase input
func BenchmarkGenAllCaps(b *testing.B) {
	input := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Gen(input)
	}
}

// BenchmarkGenAllLower benchmarks Gen with all lowercase input
func BenchmarkGenAllLower(b *testing.B) {
	input := "abcdefghijklmnopqrstuvwxyz"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Gen(input)
	}
}

// BenchmarkGenMixedContent benchmarks Gen with mixed letters, numbers, and symbols
func BenchmarkGenMixedContent(b *testing.B) {
	input := "abc123def456ghi!@#$%^&*()xyz"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Gen(input)
	}
}

// BenchmarkGenParallel benchmarks Gen with parallel execution
func BenchmarkGenParallel(b *testing.B) {
	input := "Hello, World!"
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Gen(input)
		}
	})
}

// TestGenEdgeCasesExtended tests additional edge cases
func TestGenEdgeCasesExtended(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Emoji and special unicode
		{"emoji", "hello 👋 world", "hElLo 👋 WoRlD"},
		{"mixed emoji", "test🔥code", "tEsT🔥cOdE"},

		// Various whitespace
		{"tab separated", "a\tb\tc", "a\tB\tc"},
		{"carriage return", "a\rb\rc", "a\rB\rc"},
		{"form feed", "a\fb\fc", "a\fB\fc"},
		{"vertical tab", "a\vb\vc", "a\vB\vc"},

		// Repeated punctuation
		{"multiple exclamation", "wow!!!", "wOw!!!"},
		{"multiple question", "what???", "wHaT???"},
		{"ellipsis", "wait...", "wAiT..."},

		// Quotes and brackets
		{"single quotes", "'hello'", "'hElLo'"},
		{"double quotes", "\"hello\"", "\"hElLo\""},
		{"backticks", "`code`", "`cOdE`"},
		{"curly braces", "{hello}", "{hElLo}"},
		{"angle brackets", "<hello>", "<hElLo>"},
		{"parentheses", "(hello)", "(hElLo)"},

		// Currency and symbols
		{"dollar sign", "$100", "$100"},
		{"euro sign", "€50", "€50"},
		{"pound sign", "£30", "£30"},
		{"percent", "50%", "50%"},
		{"at symbol", "@user", "@uSeR"},
		{"hash tag", "#tag", "#tAg"},
		{"ampersand", "A&B", "a&B"},
		{"asterisk", "a*b*c", "a*B*c"},

		// Math symbols
		{"plus minus", "a+b-c", "a+B-c"},
		{"multiply divide", "a*b/c", "a*B/c"},
		{"equals", "a=b", "a=B"},
		{"less greater", "a<b>c", "a<B>c"},

		// Only non-letters
		{"only numbers", "123456", "123456"},
		{"only punctuation", "!@#$%^", "!@#$%^"},
		{"only spaces", "    ", "    "},
		{"mixed non-letters", "123 !@# 456", "123 !@# 456"},

		// Alternating patterns
		{"single letter words", "I am a go dev", "i Am A gO dEv"},
		{"all single chars", "a b c d e f", "a B c D e F"},

		// Very long consecutive letters
		{"long word", "antidisestablishmentarianism", "aNtIdIsEsTaBlIsHmEnTaRiAnIsM"},

		// Mixed scripts (if supported)
		{"cyrillic", "Привет", "пРиВеТ"},
		{"greek", "Γειά", "γΕιΆ"},

		// Leading/trailing special chars
		{"leading special", "!!!hello", "!!!hElLo"},
		{"trailing special", "hello!!!", "hElLo!!!"},
		{"surrounded", "***test***", "***tEsT***"},

		// URL-like strings
		{"url pattern", "https://example.com", "hTtPs://ExAmPlE.cOm"},
		{"email pattern", "test@example.com", "tEsT@eXaMpLe.CoM"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Gen(tt.input)
			if got != tt.want {
				t.Errorf("Gen(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestGenNonASCII tests behavior with various non-ASCII characters
func TestGenNonASCII(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"accented chars", "crème brûlée", "cRèMe BrÛlÉe"},
		{"german umlaut", "Über Größe", "üBeR gRöße"}, // Note: ß lowercase is ß, not ẞ
		{"scandinavian", "Åse Ørsted", "åSe ØrStEd"},
		{"turkish", "İstanbul", "iStAnBuL"}, // Note: İ (capital dotted I) lowercases to i
		{"polish", "Łódź", "łÓdŹ"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Gen(tt.input)
			if got != tt.want {
				t.Errorf("Gen(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestGenRobustness tests robustness with unusual inputs
func TestGenRobustness(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"very long string", strings.Repeat("abcdefghijklmnopqrstuvwxyz", 1000)},
		{"many spaces", strings.Repeat(" ", 1000)},
		{"many newlines", strings.Repeat("\n", 100)},
		{"alternating space letter", strings.Repeat("a ", 500)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just ensure it doesn't panic or hang
			result := Gen(tt.input)
			if len(result) != len(tt.input) {
				t.Errorf("Gen() changed input length: got %d, want %d", len(result), len(tt.input))
			}
		})
	}
}

// TestRunWithArgs tests the run function with command-line arguments
func TestRunWithArgs(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantOutput string
		wantCode   int
	}{
		{"single arg", []string{"Hello"}, "hElLo\n", 0},
		{"multiple args", []string{"Hello", "World"}, "hElLo WoRlD\n", 0},
		{"with punctuation", []string{"Hello,", "World!"}, "hElLo, WoRlD!\n", 0},
		{"empty args", []string{}, "Usage: minimockbob \"<text>\"\nOr pipe text to it: echo \"<text>\" | minimockbob\n", 1},
		{"empty string arg", []string{""}, "Usage: minimockbob \"<text>\"\nOr pipe text to it: echo \"<text>\" | minimockbob\n", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			stdin := strings.NewReader("")
			code := Run(tt.args, stdin, &stdout, &stderr)
			if code != tt.wantCode {
				t.Errorf("run() = %d, want %d", code, tt.wantCode)
			}
			if got := stdout.String(); got != tt.wantOutput {
				t.Errorf("run() output = %q, want %q", got, tt.wantOutput)
			}
		})
	}
}

// TestRunWithStdin tests the run function reading from stdin
func TestRunWithStdin(t *testing.T) {
	tests := []struct {
		name       string
		stdin      string
		wantOutput string
		wantCode   int
	}{
		{"simple input", "Hello World", "hElLo WoRlD\n", 0},
		{"with newline", "Hello World\n", "hElLo WoRlD\n", 0},
		{"multiline", "Hello\nWorld", "hElLo\nWoRlD\n", 0},
		{"empty input", "", "Usage: minimockbob \"<text>\"\nOr pipe text to it: echo \"<text>\" | minimockbob\n", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			stdin := strings.NewReader(tt.stdin)
			code := Run([]string{}, stdin, &stdout, &stderr)
			if code != tt.wantCode {
				t.Errorf("run() = %d, want %d", code, tt.wantCode)
			}
			if got := stdout.String(); got != tt.wantOutput {
				t.Errorf("run() output = %q, want %q", got, tt.wantOutput)
			}
		})
	}
}

// TestRunStdinReadError tests error handling when reading from stdin fails
func TestRunStdinReadError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	stdin := errorReader{}
	code := Run([]string{}, stdin, &stdout, &stderr)
	if code != 1 {
		t.Errorf("run() with read error = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "Error reading from STDIN") {
		t.Errorf("run() stderr = %q, want error message", stderr.String())
	}
}

// TestRunWithPipe tests the run function with actual pipe input (os.File)
func TestRunWithPipe(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple pipe", "Hello Pipe", "hElLo PiPe\n"},
		{"empty pipe", "", "Usage: minimockbob \"<text>\"\nOr pipe text to it: echo \"<text>\" | minimockbob\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("Failed to create pipe: %v", err)
			}

			// Write test data to pipe in a goroutine
			go func() {
				if tt.input != "" {
					_, _ = w.Write([]byte(tt.input))
				}
				_ = w.Close()
			}()

			var stdout, stderr bytes.Buffer
			code := Run([]string{}, r, &stdout, &stderr)
			_ = r.Close()

			if tt.input == "" && code != 1 {
				t.Errorf("run() with empty pipe = %d, want 1", code)
			} else if tt.input != "" && code != 0 {
				t.Errorf("run() with pipe = %d, want 0", code)
			}
			if got := stdout.String(); got != tt.want {
				t.Errorf("run() output = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestCLI tests the CLI functionality binary before running any tests
func TestCLI(t *testing.T) {
	cmd := exec.Command("go", "build", "-o", "minimockbob", "./cmd/minimockbob")
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
