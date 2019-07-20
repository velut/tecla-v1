package app

import (
	"github.com/velut/tecla/server/pkg/core"
	"github.com/velut/tecla/server/pkg/gui"
)

// App represents the main application.
type App struct {
	ConfigValidator core.ConfigValidatorAPI
	Organizer       core.OrganizerAPI
	WindowWidth     int
	WindowHeight    int
}

// NewDefaultApp creates a new App with default implementations.
func NewDefaultApp() *App {
	return &App{
		ConfigValidator: core.NewConfigValidator(),
		Organizer:       core.NewOrganizer(),
		WindowWidth:     1280,
		WindowHeight:    720,
	}
}

// Run runs the application and blocks.
func (a *App) Run() error {
	return a.run()
}

func (a *App) run() error {
	apiImpl := struct {
		core.ConfigValidatorAPI
		core.OrganizerAPI
	}{
		a.ConfigValidator,
		a.Organizer,
	}
	return gui.Open(apiImpl, a.WindowWidth, a.WindowHeight)
}

// Close closes the application.
func (a *App) Close() error {
	return a.close()
}

func (a *App) close() error {
	_, err := a.Organizer.DropConfig()
	return err
}
