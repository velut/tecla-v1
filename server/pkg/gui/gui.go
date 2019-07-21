package gui

import (
	"net/url"
	"reflect"
	"runtime"
	"strings"

	"github.com/zserge/lorca"

	"github.com/velut/tecla/server/pkg/core"
)

// Open opens the GUI window with the given width and height,
// then binds the API methods, and blocks until the GUI is closed.
func Open(apiImpl core.API, width, height int) error {
	ui, err := lorca.New(getLoadingPage(), "", width, height, getArgs()...)
	if err != nil {
		return err
	}
	defer ui.Close()

	// Collect API methods' names
	interfaceType := reflect.TypeOf((*core.API)(nil)).Elem()
	var methodNames []string
	for i := 0; i < interfaceType.NumMethod(); i++ {
		methodNames = append(methodNames, interfaceType.Method(i).Name)
	}

	// Bind API methods
	implType := reflect.ValueOf(apiImpl)
	for _, name := range methodNames {
		method := implType.MethodByName(name).Interface()
		methodName := strings.ToLower(string(name[0])) + name[1:]
		err = ui.Bind(methodName, method)
		if err != nil {
			return err
		}
	}

	// TODO load client
	_ = ui.Load("http://localhost:8080")

	<-ui.Done()
	return nil
}

func getLoadingPage() string {
	return "data:text/html," + url.PathEscape(loadingPage)
}

func getArgs() []string {
	args := []string{}
	if runtime.GOOS == "linux" {
		args = append(args, "--class=Tecla")
	}
	return args
}
