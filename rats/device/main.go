package device

type Metadata map[string]string

type Device interface {
	Metadata() *Metadata

	Reserve() bool
	Release() bool

	Push(filename string, f io.Reader)
	Pull(filename string, f io.Writer)

	Install(filename string, f io.Reader)
	Uninstall(filename string)

	RunTest(app io.Reader, test io.Reader) *test.TestRunner
}
