package minimockbob

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const (
	packageName    = "minimockbob"
	packageVersion = "0.0.1"
	melangeYAML    = "melange.yaml"
)

// getNativeArch returns the native architecture for melange builds
func getNativeArch() string {
	cmd := exec.Command("uname", "-m")
	output, err := cmd.Output()
	if err != nil {
		return "x86_64" // default fallback
	}
	arch := strings.TrimSpace(string(output))
	// Map common arch names to melange arch names
	switch arch {
	case "arm64", "aarch64":
		return "aarch64"
	case "x86_64", "amd64":
		return "x86_64"
	default:
		return "x86_64"
	}
}

// TestMelangePrerequisites checks that melange and docker are available
func TestMelangePrerequisites(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping melange prerequisite tests in short mode")
	}

	t.Run("melange_installed", func(t *testing.T) {
		_, err := exec.LookPath("melange")
		if err != nil {
			t.Skip("melange is not installed. Install from: https://github.com/chainguard-dev/melange")
		}
	})

	t.Run("docker_running", func(t *testing.T) {
		cmd := exec.Command("docker", "info")
		err := cmd.Run()
		if err != nil {
			t.Skip("Docker is not running")
		}
	})

	t.Run("melange_yaml_exists", func(t *testing.T) {
		if _, err := os.Stat(melangeYAML); os.IsNotExist(err) {
			t.Fatalf("%s not found", melangeYAML)
		}
	})
}

// TestMelangeYAMLValid validates that melange.yaml can be parsed
func TestMelangeYAMLValid(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping melange YAML validation test in short mode")
	}

	// Check prerequisites
	if _, err := exec.LookPath("melange"); err != nil {
		t.Skip("melange is not installed")
	}

	if _, err := os.Stat(melangeYAML); os.IsNotExist(err) {
		t.Fatalf("%s not found", melangeYAML)
	}

	// Use 'melange query' to validate the YAML can be parsed
	cmd := exec.Command("melange", "query", melangeYAML, "{{ .Package.Name }}")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("melange YAML validation failed: %v\nOutput: %s", err, string(output))
	}

	result := strings.TrimSpace(string(output))
	if result != packageName {
		t.Errorf("Expected package name %s, got %s", packageName, result)
	}
	t.Logf("melange.yaml is valid, package name: %s", result)
}

// getDockerSocket returns the Docker socket path from the current context
func getDockerSocket(t *testing.T) string {
	cmd := exec.Command("docker", "context", "inspect", "--format", "{{.Endpoints.docker.Host}}")
	output, err := cmd.Output()
	if err == nil && strings.TrimSpace(string(output)) != "" {
		return strings.TrimSpace(string(output))
	}

	// Fall back to DOCKER_HOST environment variable
	if dockerHost := os.Getenv("DOCKER_HOST"); dockerHost != "" {
		return dockerHost
	}

	// Check for colima socket
	homeDir, err := os.UserHomeDir()
	if err == nil {
		colimaSocket := filepath.Join(homeDir, ".colima", "default", "docker.sock")
		if _, err := os.Stat(colimaSocket); err == nil {
			return "unix://" + colimaSocket
		}
	}

	// Check for podman socket (default user location)
	if homeDir, err := os.UserHomeDir(); err == nil {
		podmanSocket := filepath.Join(homeDir, ".local", "share", "containers", "podman", "machine", "podman.sock")
		if _, err := os.Stat(podmanSocket); err == nil {
			return "unix://" + podmanSocket
		}

		// Check for podman socket (alternative location)
		podmanAltSocket := filepath.Join(homeDir, ".local", "share", "containers", "podman.sock")
		if _, err := os.Stat(podmanAltSocket); err == nil {
			return "unix://" + podmanAltSocket
		}
	}

	return ""
}

// TestMelangeBuild builds the APK package
func TestMelangeBuild(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping melange build test in short mode")
	}

	// Check prerequisites
	if _, err := exec.LookPath("melange"); err != nil {
		t.Skip("melange is not installed")
	}

	cmd := exec.Command("docker", "info")
	if err := cmd.Run(); err != nil {
		t.Skip("Docker is not running")
	}

	// Get native architecture for the build
	arch := getNativeArch()
	t.Logf("Building for native architecture: %s", arch)

	// Get Docker socket path
	dockerHost := getDockerSocket(t)
	if dockerHost == "" {
		t.Log("Warning: Could not determine Docker socket, using default")
	} else {
		t.Logf("Using Docker socket: %s", dockerHost)
	}

	// Generate signing key if it doesn't exist
	if _, err := os.Stat("melange.rsa"); os.IsNotExist(err) {
		t.Log("Generating signing key...")
		cmd := exec.Command("melange", "keygen")
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("Failed to generate signing key: %v\nOutput: %s", err, string(output))
		}
		t.Log("Signing key generated")
	}

	// Clean previous build artifacts
	t.Log("Cleaning previous build artifacts...")
	if err := os.RemoveAll("packages"); err != nil && !os.IsNotExist(err) {
		t.Fatalf("Failed to clean packages directory: %v", err)
	}

	// Clean cached melange Docker image to avoid stale workspace path issues
	t.Log("Cleaning cached melange Docker images...")
	cleanCmd := exec.Command("docker", "rmi", "-f", "melange:latest")
	if output, err := cleanCmd.CombinedOutput(); err != nil {
		// Ignore errors if image doesn't exist
		if !strings.Contains(string(output), "No such image") {
			t.Logf("Warning: Failed to remove melange image: %v (output: %s)", err, string(output))
		}
	}

	// Create a test-specific melange.yaml that builds from local sources
	// instead of using git-checkout
	testMelangeYAML := "melange-test.yaml"
	content := `package:
  name: minimockbob
  version: 0.0.1
  epoch: 0
  description: "A sarcastic text generator - transforms text with alternating capitalization"
  copyright:
    - license: Apache-2.0
  target-architecture:
    - x86_64
    - aarch64

environment:
  contents:
    repositories:
      - https://packages.wolfi.dev/os
    keyring:
      - https://packages.wolfi.dev/os/wolfi-signing.rsa.pub
    packages:
      - build-base
      - busybox
      - ca-certificates-bundle
      - go

pipeline:
  - uses: go/build
    with:
      modroot: .
      packages: ./cmd/minimockbob
      output: minimockbob
      ldflags: -buildid= -s -w

  - uses: strip

  - runs: |
      install -Dm644 README.md ${{targets.destdir}}/usr/share/doc/minimockbob/README.md
      install -Dm644 LICENCE ${{targets.destdir}}/usr/share/licenses/minimockbob/LICENCE

subpackages:
  - name: minimockbob-doc
    description: "Documentation for minimockbob"
    pipeline:
      - runs: |
          mkdir -p ${{targets.subpkgdir}}/usr/share/doc/minimockbob
          cp README.md ${{targets.subpkgdir}}/usr/share/doc/minimockbob/
`
	if err := os.WriteFile(testMelangeYAML, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test melange.yaml: %v", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("Failed to remove test melange.yaml: %v", err)
		}
	}(testMelangeYAML)

	// Get current working directory to use as source
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Create workspace directory using OS temp dir
	// Using system temp avoids Docker Desktop path mapping issues
	workspaceDir, err := os.MkdirTemp("", "melange-workspace-*")
	if err != nil {
		t.Fatalf("Failed to create temp workspace directory: %v", err)
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Failed to remove workspace directory: %v", err)
		}
	}(workspaceDir)
	t.Logf("Using workspace directory: %s", workspaceDir)

	// Build the package
	t.Log("Building APK package...")
	cmd = exec.Command("melange", "build", testMelangeYAML,
		"--signing-key", "melange.rsa",
		"--runner", "docker",
		"--arch", arch,
		"--workspace-dir", workspaceDir,
		"--source-dir", cwd,
		"--log-level", "info")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set DOCKER_HOST environment variable if we found a socket
	if dockerHost != "" {
		cmd.Env = append(os.Environ(), "DOCKER_HOST="+dockerHost)
	}

	if err := cmd.Run(); err != nil {
		t.Logf("melange build failed: %v", err)

		// Check if using colima and provide helpful guidance
		if strings.Contains(dockerHost, "colima") {
			t.Log("Detected colima - checking resource allocation...")
			if colimaCmd := exec.Command("colima", "list"); colimaCmd.Run() == nil {
				if output, err := colimaCmd.Output(); err == nil {
					t.Logf("Colima configuration:\n%s", string(output))
				}
			}
			t.Fatal("Build failed with colima. This may be due to insufficient memory. " +
				"Try increasing colima memory: colima stop && colima start --memory 4")
		}

		t.Fatalf("melange build failed: %v", err)
	}
	t.Log("APK package built successfully")
}

// TestPackageArtifacts verifies that the expected package files were created
func TestPackageArtifacts(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping package artifacts test in short mode")
	}

	arch := getNativeArch()

	// This test depends on TestMelangeBuild having run
	packageFile := filepath.Join("packages", arch, fmt.Sprintf("%s-%s-r0.apk", packageName, packageVersion))

	t.Run("package_file_exists", func(t *testing.T) {
		if _, err := os.Stat(packageFile); os.IsNotExist(err) {
			t.Skipf("Package file not found: %s. Run TestMelangeBuild first.", packageFile)
		}

		info, err := os.Stat(packageFile)
		if err != nil {
			t.Fatalf("Failed to stat package file: %v", err)
		}
		t.Logf("Package size: %d bytes (%.2f KB)", info.Size(), float64(info.Size())/1024)
	})

	t.Run("apk_index_exists", func(t *testing.T) {
		indexFile := filepath.Join("packages", arch, "APKINDEX.tar.gz")
		if _, err := os.Stat(indexFile); os.IsNotExist(err) {
			t.Skipf("APKINDEX.tar.gz not found: %s", indexFile)
		}
	})

	t.Run("doc_subpackage_exists", func(t *testing.T) {
		docPackage := filepath.Join("packages", arch, fmt.Sprintf("%s-doc-%s-r0.apk", packageName, packageVersion))
		if _, err := os.Stat(docPackage); os.IsNotExist(err) {
			t.Log("Documentation subpackage not found (optional)")
			return
		}

		info, err := os.Stat(docPackage)
		if err != nil {
			t.Fatalf("Failed to stat doc package: %v", err)
		}
		t.Logf("Documentation package size: %d bytes (%.2f KB)", info.Size(), float64(info.Size())/1024)
	})
}

// TestPackageContents verifies the contents of the built package
func TestPackageContents(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping package contents test in short mode")
	}

	cmd := exec.Command("docker", "info")
	if err := cmd.Run(); err != nil {
		t.Skip("Docker is not running")
	}

	arch := getNativeArch()

	packageFile := filepath.Join("packages", arch, fmt.Sprintf("%s-%s-r0.apk", packageName, packageVersion))
	if _, err := os.Stat(packageFile); os.IsNotExist(err) {
		t.Skipf("Package file not found: %s. Run TestMelangeBuild first.", packageFile)
	}

	// Get package contents
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	containerPackagePath := fmt.Sprintf("/packages/%s/%s-%s-r0.apk", arch, packageName, packageVersion)
	cmd = exec.Command("docker", "run", "--rm",
		"-v", fmt.Sprintf("%s/packages:/packages", cwd),
		"cgr.dev/chainguard/wolfi-base", "sh", "-c",
		fmt.Sprintf("apk add --allow-untrusted %s > /dev/null 2>&1 && apk info --contents %s", containerPackagePath, packageName))

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to get package contents: %v\nOutput: %s", err, string(output))
	}

	contents := string(output)

	t.Run("binary_included", func(t *testing.T) {
		if !strings.Contains(contents, "usr/bin/minimockbob") {
			t.Errorf("Binary not found in package.\nPackage contents:\n%s", contents)
		}
	})

	t.Run("documentation_included", func(t *testing.T) {
		if !strings.Contains(contents, "usr/share/doc/minimockbob") {
			t.Log("Documentation not found in package (may be in subpackage)")
		}
	})

	t.Run("license_included", func(t *testing.T) {
		if !strings.Contains(contents, "usr/share/licenses/minimockbob") {
			t.Log("License not found in package (may be in subpackage)")
		}
	})
}

// TestPackageInstallation tests installing the package in a Wolfi container
func TestPackageInstallation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping package installation test in short mode")
	}

	cmd := exec.Command("docker", "info")
	if err := cmd.Run(); err != nil {
		t.Skip("Docker is not running")
	}

	arch := getNativeArch()

	packageFile := filepath.Join("packages", arch, fmt.Sprintf("%s-%s-r0.apk", packageName, packageVersion))
	if _, err := os.Stat(packageFile); os.IsNotExist(err) {
		t.Skipf("Package file not found: %s. Run TestMelangeBuild first.", packageFile)
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	containerPackagePath := fmt.Sprintf("/packages/%s/%s-%s-r0.apk", arch, packageName, packageVersion)
	cmd = exec.Command("docker", "run", "--rm",
		"-v", fmt.Sprintf("%s/packages:/packages", cwd),
		"cgr.dev/chainguard/wolfi-base", "sh", "-c",
		fmt.Sprintf("apk add --allow-untrusted %s > /dev/null 2>&1 && echo SUCCESS || echo FAILED", containerPackagePath))

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run installation test: %v\nOutput: %s", err, string(output))
	}

	result := strings.TrimSpace(string(output))
	if result != "SUCCESS" {
		t.Errorf("Package installation failed. Output: %s", result)
	} else {
		t.Log("Package installed successfully in Wolfi container")
	}
}

// TestBinaryExecution tests running the binary from the installed package
func TestBinaryExecution(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping binary execution test in short mode")
	}

	cmd := exec.Command("docker", "info")
	if err := cmd.Run(); err != nil {
		t.Skip("Docker is not running")
	}

	arch := getNativeArch()

	packageFile := filepath.Join("packages", arch, fmt.Sprintf("%s-%s-r0.apk", packageName, packageVersion))
	if _, err := os.Stat(packageFile); os.IsNotExist(err) {
		t.Skipf("Package file not found: %s. Run TestMelangeBuild first.", packageFile)
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	containerPackagePath := fmt.Sprintf("/packages/%s/%s-%s-r0.apk", arch, packageName, packageVersion)

	t.Run("execution_with_arguments", func(t *testing.T) {
		cmd := exec.Command("docker", "run", "--rm",
			"-v", fmt.Sprintf("%s/packages:/packages", cwd),
			"cgr.dev/chainguard/wolfi-base", "sh", "-c",
			fmt.Sprintf("apk add --allow-untrusted %s > /dev/null 2>&1 && minimockbob 'Hello World'", containerPackagePath))

		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Binary execution failed: %v\nOutput: %s", err, string(output))
		}

		result := strings.TrimSpace(string(output))
		expected := "hElLo WoRlD"
		if result != expected {
			t.Errorf("Binary output mismatch.\nExpected: %s\nGot: %s", expected, result)
		} else {
			t.Logf("Binary executed successfully. Input: 'Hello World', Output: '%s'", result)
		}
	})

	t.Run("execution_with_piped_input", func(t *testing.T) {
		cmd := exec.Command("docker", "run", "--rm",
			"-v", fmt.Sprintf("%s/packages:/packages", cwd),
			"cgr.dev/chainguard/wolfi-base", "sh", "-c",
			fmt.Sprintf("apk add --allow-untrusted %s > /dev/null 2>&1 && echo 'Testing Pipe' | minimockbob", containerPackagePath))

		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Binary execution with piped input failed: %v\nOutput: %s", err, string(output))
		}

		result := strings.TrimSpace(string(output))
		expected := "tEsTiNg PiPe"
		if result != expected {
			t.Errorf("Binary output mismatch with piped input.\nExpected: %s\nGot: %s", expected, result)
		} else {
			t.Logf("Binary with piped input executed successfully. Input: 'Testing Pipe', Output: '%s'", result)
		}
	})
}
