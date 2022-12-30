package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"io/fs"
	"log"
	"net/http"
	"os"
)

//go:embed static
var embedFiles embed.FS

//type Resource struct {
//	fs   embed.FS
//	path string
//}
//
//func NewResource() *Resource {
//	return &Resource{
//		fs:   embedFiles,
//		path: "html",
//	}
//}

func main() {
	server := gin.New()
	server.StaticFS("/", getFileSystem(false))
	server.Run(":8888")

	//http.Handle("/", http.FileServer(getFileSystem(false)))
	//http.ListenAndServe(":8888", nil)

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
