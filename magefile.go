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
	windowsExecutable = "./build/windows/tecla.exe"
	windowsEnv        = env{
		"GOOS":   "windows",
		"GOARCH": "amd64",
	}

	// Linux build options
	linuxExecutable = "./build/linux/tecla"
	linuxEnv        = env{
		"GOOS":   "linux",
		"GOARCH": "amd64",
	}
)

// env represents environment variables.
type env map[string]string

// Run namespace
type Run mg.Namespace

// Builds Tecla for Windows and runs it.
func (Run) Windows() error {
	mg.Deps(Build.Windows)

	fmt.Println("Running Tecla for Windows...")
	return sh.Run(windowsExecutable)
}

// Builds Tecla for Linux and runs it.
func (Run) Linux() error {
	mg.Deps(Build.Linux)

	fmt.Println("Running Tecla for Linux...")
	return sh.Run(linuxExecutable)
}

// Build namespace
type Build mg.Namespace

// Builds Tecla for all platforms.
func (Build) All() {
	mg.Deps(Build.Windows, Build.Linux)
}

// Builds Tecla for Windows.
func (Build) Windows() error {
	mg.Deps(installServerDependencies)

	fmt.Println("Building Tecla for Windows...")
	return build(windowsEnv, windowsExecutable)
}

// Builds Tecla for Linux.
func (Build) Linux() error {
	mg.Deps(installServerDependencies)

	fmt.Println("Building Tecla for Linux...")
	return build(linuxEnv, linuxExecutable)
}

func build(env env, executable string) error {
	return sh.RunWith(
		env, "go", "build", "-v", "-o", executable, teclaCmd,
	)
}

// Test namespace
type Test mg.Namespace

// Tests the server packages.
func (Test) Server() error {
	mg.Deps(installServerDependencies)

	fmt.Println("Testing server packages...")
	return sh.Run("go", "test", "-v", "-race", "./server/pkg/...")
}

// Installs the server dependencies.
func installServerDependencies() error {
	fmt.Println("Installing server dependencies...")
	return sh.Run("go", "mod", "download")
}
