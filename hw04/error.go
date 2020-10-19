package tinyini

// New returns an error that formats as the given text.
// Each call to New returns a distinct error value even if the text is identical.
func New(text string) error {
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

var (
	ErrNoFile = New("no such file")
	ErrOpenFile = New("Open file failed!")
	ErrReadFile = New("Can not read the file!")
)
