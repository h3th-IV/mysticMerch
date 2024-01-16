package server

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/h3th-IV/mysticMerch/internal/api"
)

type neuteredFileSystem struct {
	fs http.FileSystem
}

// check if request is a directory
func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	httpFile, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	//get file info
	neuter, _ := httpFile.Stat()
	//check if file is a directory
	if neuter.IsDir() {
		//if path is directory join index.html to path
		index := filepath.Join(path, "index.html") //index.html is just prevent user from accesing dev pages
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := httpFile.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}
	return httpFile, nil
}

func SetUserRoutes(router *mux.Router) {
	UserRouter := router.PathPrefix("/user").Subrouter()

	//routes for the user
	UserRouter.HandleFunc("/signup", api.SignUp)
	UserRouter.HandleFunc("/login", api.LogIn)
	UserRouter.HandleFunc("/cart", api.UserCart)
}
