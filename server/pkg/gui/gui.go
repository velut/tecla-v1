package gui

import (
	"net/http"
	"net/url"
	"runtime"
	"time"

	"github.com/velut/tecla/static/client"

	"github.com/zserge/lorca"
)

const defaultClientDevServerAddr = "http://localhost:8080"

const defaultClientProdServerAddr = "http://localhost:5920"
const defaultClientProdServerPort = ":5920"

// GUI represents the GUI interface.
type GUI struct {
	ui      lorca.UI
	server  *http.Server
	options *Options
}

// Options represents the options available to configure the GUI.
type Options struct {
	Width            int          // GUI width
	Height           int          // GUI height
	BoundFuncs       []*BoundFunc // Go functions to bind to the GUI
	ProductionClient bool         // Use production client for stand-alone release
}

// BoundFunc represents a Go function or method that will be callable from the GUI.
type BoundFunc struct {
	// Binding name
	Name string

	// Bound function.
	// Function signatures must be the ones supported by Lorca.
	Func interface{}
}

// NewGUI returns a new GUI.
// Options are required.
func NewGUI(options *Options) *GUI {
	return &GUI{
		options: options,
	}
}

// Start opens the GUI, binds the required functions,
// loads the client, and then blocks until the GUI is closed.
func (g *GUI) Start() error {
	if err := g.start(); err != nil {
		return err
	}
	defer g.close()

	g.wait()
	return nil
}

func (g *GUI) start() error {
	if err := g.open(); err != nil {
		return err
	}

	if err := g.bindFuncs(); err != nil {
		return err
	}

	if err := g.bindGUIFuncs(); err != nil {
		return err
	}

	if err := g.loadClient(); err != nil {
		return err
	}

	return nil
}

func (g *GUI) open() error {
	ui, err := lorca.New(
		g.loadingPage(),
		g.profileDir(),
		g.width(),
		g.height(),
		g.chromeArgs()...,
	)
	if err != nil {
		return err
	}

	g.ui = ui
	return nil
}

func (g *GUI) loadingPage() string {
	return "data:text/html," + url.PathEscape(loadingPage)
}

func (g *GUI) profileDir() string {
	// Create a temporary profile directory.
	return ""
}

func (g *GUI) width() int {
	return g.options.Width
}

func (g *GUI) height() int {
	return g.options.Height
}

func (g *GUI) chromeArgs() []string {
	var args []string
	if runtime.GOOS == "linux" {
		// Differentiate Tecla from other X applications.
		args = append(args, "--class=Tecla")
	}
	return args
}

func (g *GUI) bindFuncs() error {
	funcs := g.options.BoundFuncs
	for _, f := range funcs {
		if err := g.ui.Bind(f.Name, f.Func); err != nil {
			return err
		}
	}
	return nil
}

func (g *GUI) bindGUIFuncs() error {
	funcs := []BoundFunc{
		{
			Name: "selectDirectory",
			Func: SelectDirectory,
		},
	}
	for _, f := range funcs {
		if err := g.ui.Bind(f.Name, f.Func); err != nil {
			return err
		}
	}
	return nil
}

func (g *GUI) loadClient() error {
	var clientAddr string
	prodMode := g.options.ProductionClient

	if prodMode {
		if err := g.serveClient(defaultClientProdServerPort); err != nil {
			return err
		}
		clientAddr = defaultClientProdServerAddr
	} else {
		clientAddr = defaultClientDevServerAddr
	}

	// Wait for server to start.
	// In production mode, loading the client immediately sometimes fails in Chrome.
	time.Sleep(100 * time.Millisecond)

	if err := g.ui.Load(clientAddr); err != nil {
		return err
	}

	return nil
}

func (g *GUI) serveClient(addr string) error {
	// Load client from binary.
	clientFS, err := client.Assets()
	if err != nil {
		return err
	}

	s := &http.Server{
		Addr:    addr,
		Handler: http.FileServer(clientFS),
	}

	// Serve client.
	go func() {
		_ = s.ListenAndServe()
	}()

	return nil
}

func (g *GUI) wait() {
	<-g.ui.Done()
}

func (g *GUI) close() error {
	if err := g.closeUI(); err != nil {
		return err
	}

	if err := g.closeServer(); err != nil {
		return err
	}

	return nil
}

func (g *GUI) closeUI() error {
	return g.ui.Close()
}

func (g *GUI) closeServer() error {
	if g.server != nil {
		return g.server.Close()
	}
	return nil
}
