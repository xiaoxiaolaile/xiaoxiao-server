package core

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
)

//go:embed static
var embedFiles embed.FS

func initStatic() {
	server.StaticFS("/admin/", getFileSystem(false))
}
func getFileSystem(useOS bool) http.FileSystem {
	if useOS {
		log.Print("using live mode")
		return http.FS(os.DirFS("static"))
	}

	log.Print("using embed mode")

	fsys, err := fs.Sub(embedFiles, "static")
	if err != nil {
		panic(err)
	}
	return http.FS(fsys)
}
