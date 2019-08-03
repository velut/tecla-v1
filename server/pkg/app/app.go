package app

import (
	"reflect"
	"strings"

	"github.com/velut/tecla/static/credits"

	"github.com/velut/tecla/static/info"

	"github.com/velut/tecla/server/pkg/core"
	"github.com/velut/tecla/server/pkg/gui"
)

// App represents the main application.
type App struct {
	options *Options
}

// Options represents the options available to configure the application.
type Options struct {
	ConfigValidator core.ConfigValidatorAPI
	Organizer       core.OrganizerAPI
	WindowWidth     int
	WindowHeight    int
	ProductionMode  bool
}

// DevOptions returns a set of options for a development application.
func DevOptions() *Options {
	o := DefaultOptions()
	o.ProductionMode = false
	return o
}

// DefaultOptions returns a set of default options for a production application.
func DefaultOptions() *Options {
	return &Options{
		ConfigValidator: core.NewConfigValidator(),
		Organizer:       core.NewOrganizer(),
		WindowWidth:     1280,
		WindowHeight:    720,
		ProductionMode:  true,
	}
}

// NewApp returns a new application.
// Options are required.
func NewApp(options *Options) *App {
	return &App{
		options: options,
	}
}

// Run starts the application and blocks until the GUI is closed.
func (a *App) Run() error {
	defer a.close()
	return a.run()
}

func (a *App) run() error {
	return a.startGUI()
}

func (a *App) startGUI() error {
	gui := gui.NewGUI(a.guiOptions())
	return gui.Start()
}

func (a *App) guiOptions() *gui.Options {
	return &gui.Options{
		Width:            a.options.WindowWidth,
		Height:           a.options.WindowHeight,
		BoundFuncs:       a.guiBoundFuncs(),
		ProductionClient: a.options.ProductionMode,
	}
}

func (a *App) guiBoundFuncs() []*gui.BoundFunc {
	var bfs []*gui.BoundFunc

	bfs = append(bfs, a.configValidatorMethods()...)
	bfs = append(bfs, a.organizerMethods()...)
	bfs = append(bfs, a.appInfoFunc())
	bfs = append(bfs, a.appCreditsFunc())

	return bfs
}

func (a *App) configValidatorMethods() []*gui.BoundFunc {
	api := reflect.TypeOf((*core.ConfigValidatorAPI)(nil)).Elem()
	return extractMethods(api, a.options.ConfigValidator)
}

func (a *App) organizerMethods() []*gui.BoundFunc {
	api := reflect.TypeOf((*core.OrganizerAPI)(nil)).Elem()
	return extractMethods(api, a.options.Organizer)
}

func extractMethods(api reflect.Type, impl interface{}) []*gui.BoundFunc {
	methodNames := extractAPIMethodNames(api)

	bfs := []*gui.BoundFunc{}
	implValue := reflect.ValueOf(impl)
	for _, name := range methodNames {
		// Binding name starts with a lower case character.
		bindingName := strings.ToLower(string(name[0])) + name[1:]
		bindingFunc := implValue.MethodByName(name).Interface()
		bf := &gui.BoundFunc{
			Name: bindingName,
			Func: bindingFunc,
		}
		bfs = append(bfs, bf)
	}

	return bfs
}

func extractAPIMethodNames(api reflect.Type) []string {
	names := []string{}
	for i := 0; i < api.NumMethod(); i++ {
		names = append(names, api.Method(i).Name)
	}
	return names
}

func (a *App) appInfoFunc() *gui.BoundFunc {
	return &gui.BoundFunc{
		Name: "appInfo",
		Func: info.AppInfo,
	}
}

func (a *App) appCreditsFunc() *gui.BoundFunc {
	return &gui.BoundFunc{
		Name: "appCredits",
		Func: credits.AppCredits,
	}
}

func (a *App) close() error {
	if err := a.closeOrganizer(); err != nil {
		return err
	}

	return nil
}

func (a *App) closeOrganizer() error {
	_, err := a.options.Organizer.DropConfig()
	return err
}
