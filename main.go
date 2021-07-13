package main

import (
	"embed"
	"io/fs"

	"github.com/tylfin/gin-swagger-files/internal"
	"golang.org/x/net/webdav"
)

var (
	// Handler is used to server files through an http.Handler
	Handler *webdav.Handler

	//go:embed dist
	dist embed.FS

	static fs.FS
)

func init() {
	// Static will store the embedded swagger-UI files for use by the Handler.
	static, _ = fs.Sub(dist, "dist")

	Handler = &webdav.Handler{
		FileSystem: internal.NewWebDAVFileSystemFromFS(static),
		LockSystem: webdav.NewMemLS(),
	}
}
