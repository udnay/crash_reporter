package sink

import "io"

type ObjectStore interface {
	Store(in io.Reader, fileName string) error
}
