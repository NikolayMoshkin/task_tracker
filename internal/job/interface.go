package job

import "io"

type Executable interface {
	Execute(params []string) (io.Reader, error)
}
