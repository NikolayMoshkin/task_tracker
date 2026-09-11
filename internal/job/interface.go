package job

import "io"

type Job interface {
	Execute(params []string) (io.Reader, error)
	Name() JobType
}
