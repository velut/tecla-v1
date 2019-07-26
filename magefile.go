// +build mage

package main

import (
	"fmt"
	"os"

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

func init() {
	// Always run Mage in verbose mode.
	os.Setenv("MAGEFILE_VERBOSE", "true")
}

// Run namespace
type Run mg.Namespace

// Builds Tecla for Windows and runs it.
func (Run) Windows() error {
	mg.Deps(Build.Windows)

	fmt.Println("Running Tecla for Windows...")
	return sh.RunV(windowsExecutable)
}

// Builds Tecla for Darwin and runs it (NOT TESTED!).
func (Run) Darwin() error {
	mg.Deps(Build.Darwin)

	fmt.Println("Running Tecla for Darwin...")
	return sh.RunV(darwinExecutable)
}

// Builds Tecla for Linux and runs it.
func (Run) Linux() error {
	mg.Deps(Build.Linux)

	fmt.Println("Running Tecla for Linux...")
	return sh.RunV(linuxExecutable)
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
	return sh.RunV("go", "test", "-v", "-race", "./server/pkg/...")
}

// Lints the server.
// Requires golangci-lint.
func (Server) Lint() error {
	fmt.Println("Linting server...")
	return sh.RunV("golangci-lint", "run", "--enable-all")
}

// Lists all server production dependencies.
func (Server) ProdDeps() error {
	fmt.Println("Server production dependencies:")
	return sh.RunV("go", "list", "-f", `{{ join .Deps "\n" }}`, teclaCmd)
}

// Installs the server dependencies.
func (Server) InstallDeps() error {
	fmt.Println("Installing server dependencies...")
	return sh.RunV("go", "mod", "download")
}

// Client namespace
type Client mg.Namespace

// Serves the client for development.
// This command assumes that client dependencies are already installed.
func (Client) Serve() error {
	mg.Deps(Client.chdir)

	fmt.Println("Serving client...")
	return sh.RunV("npm", "run", "serve")
}

// Builds the client for deployment.
func (Client) Build() error {
	mg.Deps(Client.chdir, Client.InstallDeps)

	fmt.Println("Building client...")
	return sh.RunV("npm", "run", "build")
}

// Lints the client.
func (Client) Lint() error {
	mg.Deps(Client.chdir)

	fmt.Println("Linting client...")
	return sh.RunV("npm", "run", "lint")
}

// Lists all client production dependencies.
func (Client) ProdDeps() error {
	mg.Deps(Client.chdir)

	fmt.Println("Client production dependencies:")
	return sh.RunV("npx", "license-checker", "--production")
}

// Installs the client dependencies.
func (Client) InstallDeps() error {
	mg.Deps(Client.chdir)

	fmt.Println("Installing client dependencies...")
	return sh.RunV("npm", "install")
}

func (Client) chdir() error {
	return os.Chdir("./client")
}
