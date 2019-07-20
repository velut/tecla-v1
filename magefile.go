// +build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	teclaCmd = "./server/cmd/tecla/main.go"

	windowsExecutable = "./out/windows/tecla.exe"
)

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
	return sh.Run("go", "build", "-o", windowsExecutable, teclaCmd)
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
