// +build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var (
	// Tecla's main package
	teclaCmd = "./server/cmd/tecla/main.go"

	// Windows build options
	windowsExecutable = "./out/windows/tecla.exe"
	windowsEnv        = env{
		"GOOS":   "windows",
		"GOARCH": "amd64",
	}

	// Linux build options
	linuxExecutable = "./out/linux/tecla"
	linuxEnv        = env{
		"GOOS":   "linux",
		"GOARCH": "amd64",
	}
)

// env represents environment variables.
type env map[string]string

// BuildAll builds the program's executable for all platforms.
func BuildAll() {
	mg.Deps(BuildWindows)
	mg.Deps(BuildLinux)
}

// RunWindows runs the Windows executable.
func RunWindows() error {
	mg.Deps(BuildWindows)

	fmt.Println("Running the executable for Windows...")
	return sh.Run(windowsExecutable)
}

// BuildWindows builds the program's executable for Windows.
func BuildWindows() error {
	mg.Deps(InstallServerDependencies)

	fmt.Println("Building executable for Windows...")
	return build(windowsEnv, windowsExecutable)
}

// RunLinux runs the Linux executable.
func RunLinux() error {
	mg.Deps(BuildLinux)

	fmt.Println("Running the executable for Linux...")
	return sh.Run(linuxExecutable)
}

// BuildLinux builds the program's executable for Linux.
func BuildLinux() error {
	mg.Deps(InstallServerDependencies)

	fmt.Println("Building executable for Linux...")
	return build(linuxEnv, linuxExecutable)
}

func build(env env, executable string) error {
	return sh.RunWith(
		env, "go", "build", "-v", "-o", executable, teclaCmd,
	)
}

// TestServer tests the server packages.
func TestServer() error {
	mg.Deps(InstallServerDependencies)

	fmt.Println("Testing server packages...")
	return sh.Run("go", "test", "-v", "-race", "./server/pkg/...")
}

// InstallServerDependencies installs the server dependencies.
func InstallServerDependencies() error {
	fmt.Println("Installing server dependencies...")
	return sh.Run("go", "mod", "download")
}
