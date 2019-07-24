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
	windowsArgs = args{"-ldflags", "-H=windowsgui"}

	// Darwin build options
	darwinExecutable = "./build/darwin/tecla"
	darwinEnv        = env{
		"GOOS":   "darwin",
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

// args represents custom arguments.
type args []string

// Run namespace
type Run mg.Namespace

// Builds Tecla for Windows and runs it.
func (Run) Windows() error {
	mg.Deps(Build.Windows)

	fmt.Println("Running Tecla for Windows...")
	return sh.Run(windowsExecutable)
}

// Builds Tecla for Darwin and runs it (NOT TESTED!).
func (Run) Darwin() error {
	mg.Deps(Build.Darwin)

	fmt.Println("Running Tecla for Darwin...")
	return sh.Run(darwinExecutable)
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
	mg.Deps(Build.Windows, Build.Darwin, Build.Linux)
}

// Builds Tecla for Windows.
func (Build) Windows() error {
	mg.Deps(Server.InstallDeps)

	fmt.Println("Building Tecla for Windows...")
	return build(windowsEnv, windowsExecutable, windowsArgs...)
}

// Builds Tecla for Darwin (NOT TESTED!).
func (Build) Darwin() error {
	mg.Deps(Server.InstallDeps)

	fmt.Println("Building Tecla for Darwin...")
	return build(darwinEnv, darwinExecutable)
}

// Builds Tecla for Linux.
func (Build) Linux() error {
	mg.Deps(Server.InstallDeps)

	fmt.Println("Building Tecla for Linux...")
	return build(linuxEnv, linuxExecutable)
}

func build(env env, executable string, args ...string) error {
	goArgs := []string{
		"build", "-v", "-o", executable, teclaCmd,
	}
	// Insert custom args between "build" and "-v"
	goArgs = append(goArgs[:1], append(args, goArgs[1:]...)...)
	return sh.RunWith(env, "go", goArgs...)
}

// Server namespace
type Server mg.Namespace

// Runs the server tests.
func (Server) Test() error {
	mg.Deps(Server.InstallDeps)

	fmt.Println("Running server tests...")
	return sh.Run("go", "test", "-v", "-race", "./server/pkg/...")
}

// Lists all server production dependencies.
func (Server) ProdDeps() error {
	fmt.Println("Server production dependencies:")
	return sh.RunV("go", "list", "-f", `{{ join .Deps "\n" }}`, teclaCmd)
}

// Installs the server dependencies.
func (Server) InstallDeps() error {
	fmt.Println("Installing server dependencies...")
	return sh.Run("go", "mod", "download")
}
