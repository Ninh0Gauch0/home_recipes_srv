package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Init the configuration needed to start the server
func (s *Server) Init() bool {
	s.router = mux.NewRouter()

	//s.router.PathPrefix("/hrs")

	//Reading configuration file
	dat, err := ioutil.ReadFile("jsons/mongoconfig.json")
	check(err)
	var result MongoConf
	check(json.Unmarshal(dat, &result))

	s.initialized = true
	return true
}

// Start the server
func (s *Server) Start(config map[string]string) chan bool {

	// Recovering config server
	addr, ok := config["addr"]
	if !ok {
		addr = ":8080"
	}
	s.Addr = addr

	if s.initialized != true {
		err := s.Init()
		if err {
			return nil
		}
	}

	s.logger.Infof("Starting server....")
	s.worker = &Worker{}
	s.worker.Init(s.Ctx, s.GetLogger())

	s.addRoutes()

	exitChan := make(chan bool)

	// Go routines and channel to orchestate
	go func() {
		<-exitChan
		s.logger.Infoln("Stopping server")
		// Server shutdown
		err := s.Server.Shutdown(s.Ctx)

		if err != nil {
			s.logger.Errorln("Error shutdowning server - error: ", err.Error())
		}
	}()
	go func() {
		log.Printf("Listening on.... %s", s.Addr)
		log.Fatal(http.ListenAndServe(addr, s.router))
	}()

	return exitChan
}

// addRoutes - Define API routes
func (s *Server) addRoutes() {

	/** RECIPES ENDPOINTS**/
	s.router.HandleFunc("/hrs/recipes", func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debugln("creating recipe...")
		decoder := json.NewDecoder(r.Body)
		var recipe Recipe
		err := decoder.Decode(&recipe)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()

		hrsResp := s.worker.CreateRecipe(recipe)
		data, _ := json.Marshal(hrsResp)
		/*
			if hrsResp.GetError() != nil {
				s.logger.Errorln("[POST] - ERROR: %s", hrsResp.GetError())
			}
		*/
		s.logger.Infoln("Recipe created")
		w.WriteHeader(http.StatusCreated)
		w.Write(data)
	}).Methods("POST")

	s.router.HandleFunc("/hrs/recipes/{id}", func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debugln("searching recipe...")
		vars := mux.Vars(r)
		id := vars["id"]

		hrsResp := s.worker.GetRecipeByID(id)
		data, _ := json.Marshal(hrsResp)
		s.logger.Infoln("Recipe returned: ", id)

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}).Methods("GET")

	s.router.HandleFunc("/hrs/recipes/{id}", func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debugln("patchting recipe...")
		vars := mux.Vars(r)
		id := vars["id"]

		decoder := json.NewDecoder(r.Body)
		var recipe Recipe
		err := decoder.Decode(&recipe)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()

		hrsResp := s.worker.PatchRecipeByID(id, recipe)
		data, _ := json.Marshal(hrsResp)
		s.logger.Infoln("Recipe modified: ", id)

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}).Methods("PATCH")

	s.router.HandleFunc("/hrs/recipes/{id}", func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debugln("deleting recipe...")
		vars := mux.Vars(r)
		id := vars["id"]

		hrsResp := s.worker.DeleteRecipe(id)
		data, _ := json.Marshal(hrsResp)
		s.logger.Infoln("Recipe deleted: ", id)

		w.WriteHeader(http.StatusNoContent)
		w.Write(data)
	}).Methods("DELETE")

	/** INGREDIENTS ENDPOINTS **/
	s.router.HandleFunc("/hrs/ingredients", func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debugln("creating ingredients...")
		decoder := json.NewDecoder(r.Body)
		var ingredient Ingredient
		err := decoder.Decode(&ingredient)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()

		hrsResp := s.worker.CreateIngredient(ingredient)
		data, _ := json.Marshal(hrsResp)
		/*
			if hrsResp.GetError() != nil {
				s.logger.Errorln("[POST] - ERROR: %s", hrsResp.GetError())
			}
		*/
		s.logger.Infoln("Ingredient created")
		w.WriteHeader(http.StatusCreated)
		w.Write(data)
	}).Methods("POST")

	s.router.HandleFunc("/hrs/ingredients/{id}", func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debugln("searching ingredients...")
		vars := mux.Vars(r)
		id := vars["id"]

		hrsResp := s.worker.GetIngredientByID(id)
		data, _ := json.Marshal(hrsResp)
		s.logger.Infoln("Ingredient returned: ", id)

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}).Methods("GET")

	s.router.HandleFunc("/hrs/ingredients/{id}", func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debugln("patching ingredients...")
		vars := mux.Vars(r)
		id := vars["id"]

		decoder := json.NewDecoder(r.Body)
		var ingredient Ingredient
		err := decoder.Decode(&ingredient)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()

		hrsResp := s.worker.PatchIngredientByID(id, ingredient)
		data, _ := json.Marshal(hrsResp)
		s.logger.Infoln("Ingredient modified: ", id)

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}).Methods("PATCH")

	s.router.HandleFunc("/hrs/ingredients/{id}", func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debugln("deleting ingredient...")

		vars := mux.Vars(r)
		id := vars["id"]

		hrsResp := s.worker.DeleteIngredient(id)
		data, _ := json.Marshal(hrsResp)
		s.logger.Infoln("Ingredient deleted: ", id)

		w.WriteHeader(http.StatusNoContent)
		w.Write(data)
	}).Methods("DELETE")

	/** OTHER ENDPOINTS **/
	s.router.HandleFunc("/hrs/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "WTF\n")
	}).Methods("GET")
}

//TODO: Change error handling
func check(e error) {
	if e != nil {
		panic(e)
	}
}
