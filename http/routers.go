package http

import (
	"net/http"
	"os"
)

type IRoute interface {
	Get(pattern string, handlers ...HandlerFunc)
	Post(pattern string, handlers ...HandlerFunc)
	Put(pattern string, handlers ...HandlerFunc)
	Delete(pattern string, handlers ...HandlerFunc)
	Options(pattern string, handlers ...HandlerFunc)
	Head(pattern string, handlers ...HandlerFunc)
	Any(pattern string, handlers ...HandlerFunc)
}

type onlyFilesFS struct {
	fs http.FileSystem
}

type neuteredReaddirFile struct {
	http.File
}

// Dir returns a http.FileSystem that can be used by http.FileServer(). It is used internally
// in router.Static().
// if listDirectory == true, then it works the same as http.Dir() otherwise it returns
// a filesystem that prevents http.FileServer() to list the directory files.
func Dir(root string, listDirectory bool) http.FileSystem {
	fs := http.Dir(root)
	if listDirectory {
		return fs
	}
	return &onlyFilesFS{fs}
}

// Open conforms to http.Filesystem.
func (fs onlyFilesFS) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return neuteredReaddirFile{f}, nil
}

// Readdir overrides the http.File default implementation.
func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	// this disables directory listing
	return nil, nil
}

type FsRoute interface {
	// Static serves files from the given file system root.
	// Internally a http.FileServer is used, therefore http.NotFound is used instead
	// of the Router's NotFound handler.
	// To use the operating system's file system implementation,
	// use :
	//     router.Static("/static", "/var/www")
	Static(relativePath, root string)

	StaticFS(relativePath string, fs http.FileSystem)
}
