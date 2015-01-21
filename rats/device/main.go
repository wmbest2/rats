package device

type Device interface {
	Manufacturer() string
	Name() string
	OS() string
	IsTablet() bool
	ApiVersion() int
	Version() string
	Reserve() bool
	Release() bool
	InUse() bool
	Identity() string

	Push(filename string, f io.Reader)
	Pull(filename string, f io.Writer)

	Install(filename string, f io.Reader)
	Uninstall(filename string)
}
