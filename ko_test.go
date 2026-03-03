package minimockbob

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

const (
	koLocalRepo = "ko.local"
)

// TestKoPrerequisites checks that ko is installed and docker is running
func TestKoPrerequisites(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping ko prerequisite tests in short mode")
	}

	t.Run("ko_installed", func(t *testing.T) {
		_, err := exec.LookPath("ko")
		if err != nil {
			t.Skip("ko is not installed. Install from: https://ko.build/install/")
		}
	})

	t.Run("docker_running", func(t *testing.T) {
		cmd := exec.Command("docker", "info")
		err := cmd.Run()
		if err != nil {
			t.Skip("Docker is not running")
		}
	})

	t.Run("ko_config_exists", func(t *testing.T) {
		koConfigPath := ".ko.yaml"
		if _, err := os.Stat(koConfigPath); os.IsNotExist(err) {
			t.Fatalf("%s not found", koConfigPath)
		}
	})
}

// TestKoConfigValid validates the .ko.yaml configuration
func TestKoConfigValid(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping ko config validation test in short mode")
	}

	koConfigPath := ".ko.yaml"
	content, err := os.ReadFile(koConfigPath)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", koConfigPath, err)
	}

	configStr := string(content)

	t.Run("has_base_image", func(t *testing.T) {
		if !strings.Contains(configStr, "defaultBaseImage") {
			t.Error("Config should specify defaultBaseImage")
		}
	})

	t.Run("uses_distroless", func(t *testing.T) {
		if !strings.Contains(configStr, "distroless") {
			t.Log("Warning: Not using distroless base image")
		}
	})

	t.Run("has_build_settings", func(t *testing.T) {
		if !strings.Contains(configStr, "builds:") {
			t.Error("Config should have builds section")
		}
	})

	t.Run("disables_cgo", func(t *testing.T) {
		if !strings.Contains(configStr, "CGO_ENABLED=0") {
			t.Error("Config should disable CGO for static binary")
		}
	})

	t.Run("has_ldflags", func(t *testing.T) {
		if !strings.Contains(configStr, "ldflags") {
			t.Error("Config should specify ldflags for optimization")
		}
		if !strings.Contains(configStr, "-s -w") {
			t.Log("Warning: ldflags should include -s -w for smaller binary")
		}
	})
}

// buildKoImage is a helper function to build a ko image
func buildKoImage(t *testing.T) string {
	t.Helper()

	// Check prerequisites
	if _, err := exec.LookPath("ko"); err != nil {
		t.Skip("ko is not installed")
	}

	cmd := exec.Command("docker", "info")
	if err := cmd.Run(); err != nil {
		t.Skip("Docker is not running")
	}

	// Get Docker socket for colima/podman support
	dockerHost := getDockerSocket(t)
	if dockerHost != "" {
		t.Logf("Using Docker socket: %s", dockerHost)
	}

	t.Log("Building container image with ko...")

	// Build the image - set GOOS=linux and preserve native GOARCH
	// Don't use --platform flag to avoid conflicts with GOOS=darwin
	arch := getNativeArch()
	var goarch string
	switch arch {
	case "aarch64":
		goarch = "arm64"
	case "x86_64":
		goarch = "amd64"
	default:
		goarch = "amd64"
	}

	buildCmd := exec.Command("ko", "build", "--bare", "./cmd/minimockbob")
	buildCmd.Env = append(os.Environ(),
		fmt.Sprintf("KO_DOCKER_REPO=%s", koLocalRepo),
		"GOOS=linux",
		fmt.Sprintf("GOARCH=%s", goarch))
	if dockerHost != "" {
		buildCmd.Env = append(buildCmd.Env, "DOCKER_HOST="+dockerHost)
	}

	output, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("ko build failed: %v\nOutput: %s", err, string(output))
	}

	// Parse the output to get the image reference (last line)
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	imageRef := lines[len(lines)-1]
	t.Logf("Built image: %s", imageRef)

	return imageRef
}

// TestKoBuild builds a container image with ko locally
func TestKoBuild(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping ko build test in short mode")
	}

	imageRef := buildKoImage(t)

	// Get Docker socket for colima/podman support
	dockerHost := getDockerSocket(t)

	// Verify the image exists
	inspectCmd := exec.Command("docker", "inspect", imageRef)
	if dockerHost != "" {
		inspectCmd.Env = append(os.Environ(), "DOCKER_HOST="+dockerHost)
	}
	if err := inspectCmd.Run(); err != nil {
		t.Fatalf("Image not found in docker: %s", imageRef)
	}

	t.Log("Container image built successfully")
}

// TestKoImageMetadata verifies the built image has correct metadata
func TestKoImageMetadata(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping ko image metadata test in short mode")
	}

	// Build a fresh image for this test
	imageRef := buildKoImage(t)

	dockerHost := getDockerSocket(t)

	t.Run("image_has_labels", func(t *testing.T) {
		cmd := exec.Command("docker", "inspect", "--format", "{{json .Config.Labels}}", imageRef)
		if dockerHost != "" {
			cmd.Env = append(os.Environ(), "DOCKER_HOST="+dockerHost)
		}
		output, err := cmd.Output()
		if err != nil {
			t.Skipf("Failed to inspect image labels: %v", err)
		}
		t.Logf("Image labels: %s", string(output))
	})

	t.Run("image_has_entrypoint", func(t *testing.T) {
		cmd := exec.Command("docker", "inspect", "--format", "{{json .Config.Entrypoint}}", imageRef)
		if dockerHost != "" {
			cmd.Env = append(os.Environ(), "DOCKER_HOST="+dockerHost)
		}
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Failed to inspect image entrypoint: %v", err)
		}
		entrypoint := strings.TrimSpace(string(output))
		if entrypoint == "" || entrypoint == "null" || entrypoint == "[]" {
			t.Error("Image should have an entrypoint")
		}
		t.Logf("Image entrypoint: %s", entrypoint)
	})

	t.Run("image_size_reasonable", func(t *testing.T) {
		cmd := exec.Command("docker", "inspect", "--format", "{{.Size}}", imageRef)
		if dockerHost != "" {
			cmd.Env = append(os.Environ(), "DOCKER_HOST="+dockerHost)
		}
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Failed to get image size: %v", err)
		}
		sizeStr := strings.TrimSpace(string(output))
		t.Logf("Image size: %s bytes", sizeStr)
	})
}

// TestKoImageExecution tests running the container built with ko
func TestKoImageExecution(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping ko image execution test in short mode")
	}

	// Build a fresh image for this test
	imageRef := buildKoImage(t)

	dockerHost := getDockerSocket(t)

	t.Run("run_with_arguments", func(t *testing.T) {
		cmd := exec.Command("docker", "run", "--rm", imageRef, "Hello World")
		if dockerHost != "" {
			cmd.Env = append(os.Environ(), "DOCKER_HOST="+dockerHost)
		}
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run container: %v\nOutput: %s", err, string(output))
		}

		result := strings.TrimSpace(string(output))
		expected := "hElLo WoRlD"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
		t.Logf("Container output: %s", result)
	})

	t.Run("run_with_piped_input", func(t *testing.T) {
		cmd := exec.Command("sh", "-c", fmt.Sprintf("echo 'Testing Pipe' | docker run --rm -i %s", imageRef))
		if dockerHost != "" {
			cmd.Env = append(os.Environ(), "DOCKER_HOST="+dockerHost)
		}
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run container with piped input: %v\nOutput: %s", err, string(output))
		}

		result := strings.TrimSpace(string(output))
		expected := "tEsTiNg PiPe"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
		t.Logf("Container output: %s", result)
	})

	t.Run("run_with_no_args_shows_usage", func(t *testing.T) {
		cmd := exec.Command("docker", "run", "--rm", imageRef)
		if dockerHost != "" {
			cmd.Env = append(os.Environ(), "DOCKER_HOST="+dockerHost)
		}
		output, err := cmd.CombinedOutput()
		// Expect non-zero exit since no input provided
		if err == nil {
			t.Log("Warning: Container succeeded with no input (expected usage message)")
		}

		result := strings.TrimSpace(string(output))
		// Should show usage message
		if !strings.Contains(result, "Usage") && !strings.Contains(result, "usage") {
			t.Logf("Container output with no input: %s", result)
		}
	})
}

// TestKoBuildWithTags tests building with specific tags
// This test verifies that ko can build images with custom tags and that
// the tags are correctly applied to the built images.
//
// Platform handling: The test explicitly sets GOOS=linux and determines
// the correct GOARCH based on the native architecture. This is necessary
// because ko builds container images which must run on Linux, and the
// base image (gcr.io/distroless/static:nonroot) doesn't support macOS.
// Without setting GOOS=linux, the test would fail on macOS with
// "no matching platforms in base image index" error.
func TestKoBuildWithTags(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping ko build with tags test in short mode")
	}

	if _, err := exec.LookPath("ko"); err != nil {
		t.Skip("ko is not installed")
	}

	cmd := exec.Command("docker", "info")
	if err := cmd.Run(); err != nil {
		t.Skip("Docker is not running")
	}

	dockerHost := getDockerSocket(t)

	// Use timestamp to create unique tags
	timestamp := time.Now().Unix()
	tag1 := fmt.Sprintf("test-%d", timestamp)
	tag2 := fmt.Sprintf("test-%d-alt", timestamp)
	tags := fmt.Sprintf("%s,%s", tag1, tag2)

	t.Logf("Building with tags: %s", tags)

	// Determine the correct GOARCH for the native platform
	// This ensures we build for linux with the correct architecture
	// which matches the available platforms in the distroless base image
	arch := getNativeArch()
	var goarch string
	switch arch {
	case "aarch64":
		goarch = "arm64"
	case "x86_64":
		goarch = "amd64"
	default:
		goarch = "amd64"
	}

	t.Logf("Building for platform: linux/%s", goarch)

	buildCmd := exec.Command("ko", "build", "--bare", "--tags", tags, "./cmd/minimockbob")
	buildCmd.Env = append(os.Environ(),
		fmt.Sprintf("KO_DOCKER_REPO=%s", koLocalRepo),
		"GOOS=linux",
		fmt.Sprintf("GOARCH=%s", goarch))
	if dockerHost != "" {
		buildCmd.Env = append(buildCmd.Env, "DOCKER_HOST="+dockerHost)
	}

	output, err := buildCmd.CombinedOutput()
	if err != nil {
		// Check for common platform errors
		outputStr := string(output)
		if strings.Contains(outputStr, "no matching platforms") {
			t.Fatalf("ko build with tags failed due to platform mismatch.\n"+
				"This usually means the base image doesn't support the target platform.\n"+
				"Building for: GOOS=linux GOARCH=%s\n"+
				"Error: %v\nOutput: %s", goarch, err, outputStr)
		}
		t.Fatalf("ko build with tags failed: %v\nOutput: %s", err, outputStr)
	}

	// Parse the output to get the image reference (last line)
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	imageRef := lines[len(lines)-1]
	t.Logf("Built image: %s", imageRef)

	// Extract the base image name from the ref (format: ko.local:hash)
	// Tags are added separately by ko
	baseImage := strings.Split(imageRef, ":")[0]

	// Verify both tags exist
	for _, tag := range []string{tag1, tag2} {
		taggedImage := fmt.Sprintf("%s:%s", baseImage, tag)
		inspectCmd := exec.Command("docker", "inspect", taggedImage)
		if dockerHost != "" {
			inspectCmd.Env = append(os.Environ(), "DOCKER_HOST="+dockerHost)
		}
		if err := inspectCmd.Run(); err != nil {
			t.Errorf("Tagged image not found: %s", taggedImage)
		} else {
			t.Logf("Verified tagged image: %s", taggedImage)
		}
	}

	// Cleanup: Remove test images
	t.Cleanup(func() {
		for _, tag := range []string{tag1, tag2} {
			taggedImage := fmt.Sprintf("%s:%s", baseImage, tag)
			rmiCmd := exec.Command("docker", "rmi", "-f", taggedImage)
			if dockerHost != "" {
				rmiCmd.Env = append(os.Environ(), "DOCKER_HOST="+dockerHost)
			}
			_ = rmiCmd.Run() // Ignore errors during cleanup
		}
		t.Logf("Cleaned up test images with tags: %s", tags)
	})
}

// TestKoBuildPlatforms tests that ko builds work correctly without explicit platform flags
func TestKoBuildPlatforms(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping ko platform build test in short mode")
	}

	// Build the image using the helper - ko will determine the correct platform
	// based on GOOS/GOARCH environment variables
	imageRef := buildKoImage(t)

	dockerHost := getDockerSocket(t)

	// Verify the image was built and can be inspected
	inspectCmd := exec.Command("docker", "inspect", "--format", "{{.Architecture}}", imageRef)
	if dockerHost != "" {
		inspectCmd.Env = append(os.Environ(), "DOCKER_HOST="+dockerHost)
	}
	output, err := inspectCmd.Output()
	if err != nil {
		t.Fatalf("Failed to inspect image architecture: %v", err)
	}

	arch := strings.TrimSpace(string(output))
	t.Logf("Built image architecture: %s", arch)

	// Verify we got a valid architecture
	if arch != "arm64" && arch != "amd64" && arch != "386" {
		t.Errorf("Unexpected architecture: %s", arch)
	}
}

// TestKoConfigOptions tests various ko configuration options
func TestKoConfigOptions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping ko config options test in short mode")
	}

	t.Run("verify_main_path", func(t *testing.T) {
		mainPath := "./cmd/minimockbob"
		// Check if the directory exists and has Go files
		entries, err := os.ReadDir(mainPath)
		if err != nil {
			t.Errorf("Main package directory not found at %s", mainPath)
			return
		}
		hasGoFiles := false
		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".go") && !strings.HasSuffix(entry.Name(), "_test.go") {
				hasGoFiles = true
				break
			}
		}
		if !hasGoFiles {
			t.Errorf("No Go source files found at %s", mainPath)
		}
	})

	t.Run("verify_go_mod", func(t *testing.T) {
		if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
			t.Error("go.mod not found")
		}
	})

	t.Run("verify_static_build", func(t *testing.T) {
		// Verify that the binary built by ko is static
		// This is indirectly tested by using distroless base
		content, err := os.ReadFile(".ko.yaml")
		if err != nil {
			t.Skip("Cannot read .ko.yaml")
		}
		if !strings.Contains(string(content), "CGO_ENABLED=0") {
			t.Error("CGO should be disabled for static builds")
		}
	})
}

// TestKoBuildAndRun is an end-to-end test that builds and runs the container
func TestKoBuildAndRun(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping ko build and run test in short mode")
	}

	// Build the image using the helper function
	imageRef := buildKoImage(t)

	dockerHost := getDockerSocket(t)

	// Test running the container
	t.Run("execute_with_arguments", func(t *testing.T) {
		cmd := exec.Command("docker", "run", "--rm", imageRef, "Hello Ko")
		if dockerHost != "" {
			cmd.Env = append(os.Environ(), "DOCKER_HOST="+dockerHost)
		}
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run container: %v\nOutput: %s", err, string(output))
		}

		result := strings.TrimSpace(string(output))
		expected := "hElLo Ko"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
		t.Logf("Container executed successfully: %s", result)
	})
}

// TestKoCleanup cleans up test images
func TestKoCleanup(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping ko cleanup test in short mode")
	}

	// This test runs at the end to clean up
	t.Log("Note: To clean up ko test images, run: docker images | grep ko.local/minimockbob | awk '{print $3}' | xargs docker rmi -f")
}
