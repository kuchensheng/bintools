package cacher

import "io"

type optionString interface {
	CacheInterface

	//Save Write the cache's item (using Gob) to an io.writer
	Save(w io.Writer) error
}
