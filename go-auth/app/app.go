package app

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"../database"
	"../user"
	"github.com/gorilla/mux"
	"github.com/op/go-logging"
)

type GoAuth interface {
	Run()
}
type GoAuthS struct {
	port      int
	userData  database.UserData
	tokenData database.TokenData
}

var log = logging.MustGetLogger("app")

func NewApplication(dbUser string, dbPassword string, databaseName string,
	dbHost string, dbPort int, port int) (GoAuth, error) {

	log.Debug("Connection to database %s:%d\n", dbHost, dbPort)
	userDB, tokenDB, err := database.Connect("mysql", dbUser, dbPassword, dbHost, dbPort, databaseName)
	if err != nil {
		return nil, err
	} else {
		return &GoAuthS{port, userDB, tokenDB}, nil
	}
}

func (s *GoAuthS) Run() {
	r := mux.NewRouter()
	r.HandleFunc("/user", s.putUser).Methods("PUT")
	r.HandleFunc("/user", s.deleteUser).Methods("DELETE")
	r.HandleFunc("/token", s.getToken).Methods("GET")
	r.HandleFunc("/token/{username}", s.verifyToken).Methods("POST")
	http.Handle("/", r)
	log.Debug("Listening on port %d", s.port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (s *GoAuthS) putUser(w http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("userid")
	password := r.FormValue("password")
	status := user.RegisterUser(s.userData, userId, password)
	if status != 200 {
		http.Error(w, "Unable to save user.", status)
	}
}

func (s *GoAuthS) deleteUser(w http.ResponseWriter, r *http.Request) {
	log.Info("Delete User Called")
	userId := r.FormValue("userid")
	password := r.FormValue("password")
	status := user.DeleteUser(s.userData, userId, password)
	if status != 200 {
		http.Error(w, "Unable to delete user.", status)
	}
}

func (s *GoAuthS) getToken(w http.ResponseWriter, r *http.Request) {
	log.Info("Get Token Called")
	userId := r.FormValue("userid")
	password := r.FormValue("password")

	passwordHash, err := s.userData.GetUser(userId)
	if err != nil {
		http.Error(w, "Unable to find user.", http.StatusNotFound)
		return
	}

	err = user.CompareHashAndPassword(passwordHash, password)
	if err != nil {
		http.Error(w, "Password not valid.", http.StatusBadRequest)
		return
	}

	token, err := s.tokenData.CreateToken(userId)
	if err != nil {
		http.Error(w, "Internal error.", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(token))
}

func (s *GoAuthS) verifyToken(w http.ResponseWriter, r *http.Request) {
	log.Info("Verify Token Called")
	vars := mux.Vars(r)
	username := vars["username"]
	token, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read token.", http.StatusBadRequest)
		return
	}

	status, _ := s.tokenData.ValidateToken(string(token[:]), username)
	if status != 200 {
		http.Error(w, "Unable to verify token.", status)
	}

}