package server

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/h3th-IV/mysticMerch/internal/api"
	"github.com/h3th-IV/mysticMerch/internal/database"
	"github.com/h3th-IV/mysticMerch/internal/utils"
	"github.com/justinas/alice"
)

var UserDB *database.DBModel

// to curb directory access to non-adminstrative user
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

	userMWchain := alice.New(utils.AuthRoute)

	//routes for the user
	UserRouter.HandleFunc("/signup", api.SignUp).Methods(http.MethodPost)
	UserRouter.HandleFunc("/login", api.LogIn).Methods(http.MethodPost)
	UserRouter.Handle("/addaddress", userMWchain.ThenFunc(api.AddNewAddr)).Methods(http.MethodPost)
	UserRouter.Handle("/removeaddress/{id:[0-9]+}", userMWchain.ThenFunc(api.RemoveAddress)).Methods(http.MethodDelete)
}
