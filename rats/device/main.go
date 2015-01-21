package device

type Device interface {
	Identity() string

	Reserve() bool
	Release() bool
	InUse() bool

	Push(filename string, f io.Reader)
	Pull(filename string, f io.Writer)

	Install(filename string, f io.Reader)
	Uninstall(filename string)
}
