module github.com/velut/tecla

go 1.12

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gammazero/workerpool v0.0.0-20190608213748-0ed5e40ec55e
	github.com/gen2brain/dlgs v0.0.0-20190708095831-3854608588f7
	github.com/magefile/mage v1.8.0
	github.com/mjibson/esc v0.2.0
	github.com/pkg/errors v0.8.1 // indirect
	github.com/rakyll/statik v0.1.6
	github.com/stretchr/testify v1.3.0
	github.com/velut/fsutils-go v0.1.4
	github.com/zserge/lorca v0.1.8
	golang.org/x/net v0.0.0-20190724013045-ca1201d0de80 // indirect
	golang.org/x/tools v0.0.0-20190809145639-6d4652c779c4 // indirect
)

// See https://github.com/velut/dlgs/tree/add-window-owner
replace github.com/gen2brain/dlgs => github.com/velut/dlgs v0.0.0-20190810153543-5240659e20bc
