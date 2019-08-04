// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Go commands
var (
	goGet          = sh.RunCmd("go", "get", "-v")
	goBuild        = sh.RunCmd("go", "build", "-v")
	goTest         = sh.RunCmd("go", "test", "-v", "-race")
	goLint         = sh.RunCmd("golangci-lint", "run", "--enable-all")
	goListProdDeps = sh.RunCmd("go", "list", "-f", `{{ join .Deps "\n" }}`)
	goModDownload  = sh.RunCmd("go", "mod", "download")
)

// Npm commands
var (
	npmRunServe     = sh.RunCmd("npm", "run", "serve")
	npmRunBuild     = sh.RunCmd("npm", "run", "build")
	npmRunLint      = sh.RunCmd("npm", "run", "lint")
	npmListProdDeps = sh.RunCmd("npx", "license-checker", "--production")
	npmInstall      = sh.RunCmd("npm", "install")
)

// Tools
var (
	statik = sh.RunCmd("statik")
	esc    = sh.RunCmd("esc")
)

// Project directories
var (
	projectRootDir = func() string {
		dir, err := os.Getwd()
		check(err)
		return dir
	}()

	buildDir = filepath.Join(projectRootDir, "build")

	staticDir = filepath.Join(projectRootDir, "static")

	clientDir            = filepath.Join(projectRootDir, "client")
	clientDistDir        = filepath.Join(clientDir, "dist")
	clientNodeModulesDir = filepath.Join(clientDir, "node_modules")
)

// Main packages
var (
	// Tecla's main package
	teclaCmd = "./server/cmd/tecla/main.go"
)

// Build environments
var (
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
	mg.SerialDeps(
		Build.Windows,
		Build.Darwin,
		Build.Linux,
	)
}

// Builds Tecla for Windows.
func (Build) Windows() error {
	mg.Deps(Build.preBuild)

	fmt.Println("Building Tecla for Windows...")
	return build(windowsEnv, windowsExecutable, windowsArgs...)
}

// Builds Tecla for Darwin (NOT TESTED!).
func (Build) Darwin() error {
	mg.Deps(Build.preBuild)

	fmt.Println("Building Tecla for Darwin...")
	return build(darwinEnv, darwinExecutable)
}

// Builds Tecla for Linux.
func (Build) Linux() error {
	mg.Deps(Build.preBuild)

	fmt.Println("Building Tecla for Linux...")
	return build(linuxEnv, linuxExecutable)
}

func (Build) preBuild() {
	fmt.Println("Executing pre-build operations...")
	mg.SerialDeps(
		Tools.Install,
		Server.InstallDeps,
		Client.Build,
		Static.Generate,
		Build.chdir,
	)
}

func (Build) chdir() error {
	return os.Chdir(projectRootDir)
}

func build(env env, executable string, args ...string) error {
	// Insert args after build and before output flag
	args = append([]string{"build", "-v"}, args...)
	args = append(args, []string{"-o", executable, teclaCmd}...)
	return sh.RunWith(env, "go", args...)
}

// Server namespace
type Server mg.Namespace

// Runs the server tests.
func (Server) Test() error {
	mg.Deps(Server.InstallDeps)

	fmt.Println("Running server tests...")
	return goTest("./server/pkg/...")
}

// Lints the server.
// Requires golangci-lint.
func (Server) Lint() error {
	fmt.Println("Linting server...")
	return goLint("./server/...")
}

// Lists all server production dependencies.
func (Server) ProdDeps() error {
	fmt.Println("Server production dependencies:")
	return goListProdDeps(teclaCmd)
}

// Installs the server dependencies.
func (Server) InstallDeps() error {
	fmt.Println("Installing server dependencies...")
	return goModDownload()
}

// Client namespace
type Client mg.Namespace

// Serves the client for development.
// This command assumes that client dependencies are already installed.
func (Client) Serve() error {
	mg.Deps(Client.chdir)

	fmt.Println("Serving client...")
	return npmRunServe()
}

// Builds the client for deployment.
func (Client) Build() error {
	mg.Deps(Client.chdir, Client.InstallDeps)

	fmt.Println("Building client...")
	if clientDistDirExists() {
		fmt.Println("dist directory exists, skipping npm build")
		return nil
	}
	return npmRunBuild()
}

func clientDistDirExists() bool {
	info, _ := os.Stat(clientDistDir)
	return info != nil && info.Mode().IsDir()
}

// Lints the client.
func (Client) Lint() error {
	mg.Deps(Client.chdir)

	fmt.Println("Linting client...")
	return npmRunLint()
}

// Lists all client production dependencies.
func (Client) ProdDeps() error {
	mg.Deps(Client.chdir)

	fmt.Println("Client production dependencies:")
	return npmListProdDeps()
}

// Installs the client dependencies.
func (Client) InstallDeps() error {
	mg.Deps(Client.chdir)

	fmt.Println("Installing client dependencies...")
	if nodeModulesDirExists() {
		fmt.Println("node_modules directory exists, skipping npm install")
		return nil
	}
	return npmInstall()
}

func (Client) chdir() error {
	return os.Chdir(clientDir)
}

func nodeModulesDirExists() bool {
	info, _ := os.Stat(clientNodeModulesDir)
	return info != nil && info.Mode().IsDir()
}

// Static namespace
type Static mg.Namespace

// Generates static packages.
func (Static) Generate() {
	fmt.Println("Generating static packages...")
	mg.SerialDeps(
		Static.chdir,
		Static.generateInfo,
		Static.generateCredits,
		Static.generateClient,
		Static.test,
	)
}

func (Static) generateInfo() error {
	fmt.Println("Generating static data for info package...")
	return esc(
		"-o=./static/info/static.go",
		"-pkg=info",
		"-prefix=static/info/license",
		"-private",
		"./static/info/license",
	)
}

func (Static) generateCredits() error {
	fmt.Println("Generating static data for credits package...")
	return esc(
		"-o=./static/credits/static.go",
		"-pkg=credits",
		"-prefix=static/credits/licenses",
		"-private",
		"./static/credits/licenses",
	)
}

func (Static) generateClient() error {
	fmt.Println("Generating static client package...")
	return statik(
		"-src=./client/dist",
		"-dest=./static",
		"-p=client",
	)
}

func (Static) chdir() error {
	return os.Chdir(projectRootDir)
}

func (Static) test() error {
	fmt.Println("Testing static packages...")
	return goTest("./static/...")
}

// Tools namespace
type Tools mg.Namespace

// Installs all build tools.
func (Tools) Install() {
	fmt.Println("Installing build tools...")
	mg.SerialDeps(Tools.installStatik, Tools.installEsc)
}

// Installs the statik binary.
func (Tools) installStatik() error {
	fmt.Println("Installing statik...")
	return goGet("github.com/rakyll/statik")
}

// Installs the esc binary.
func (Tools) installEsc() error {
	fmt.Println("Installing esc...")
	return goGet("github.com/mjibson/esc")
}

// Clean namespace
type Clean mg.Namespace

// Cleans everything.
func (Clean) All() {
	fmt.Println("Cleaning...")
	mg.Deps(Clean.Artifacts, Clean.NodeModules)
}

// Removes build artifacts.
func (Clean) Artifacts() {
	fmt.Println("Removing all artifacts...")
	mg.Deps(Clean.build, Clean.client, Clean.static)
}
func (Clean) build() error {
	fmt.Println("Removing build artifacts...")
	return sh.Rm(buildDir)
}

func (Clean) client() error {
	fmt.Println("Removing client artifacts...")
	return sh.Rm(clientDistDir)
}

func (Clean) static() error {
	as := []string{
		filepath.Join(staticDir, "client", "statik.go"),
		filepath.Join(staticDir, "info", "static.go"),
		filepath.Join(staticDir, "credits", "static.go"),
	}

	fmt.Println("Removing static artifacts...")
	for _, a := range as {
		if err := sh.Rm(a); err != nil {
			return err
		}
	}

	return nil
}

// Removes node_modules directory.
func (Clean) NodeModules() error {
	fmt.Println("Removing node_modules directory...")
	return sh.Rm(clientNodeModulesDir)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
