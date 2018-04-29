package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	// DECODE_ERROR Constant
	DECODE_ERROR = "Failed validation"
	FATAL_ERROR  = "Fatal error"
)

// Init the configuration needed to start the server
func (s *Server) Init() bool {
	s.router = mux.NewRouter()

	//s.router.PathPrefix("/hrs")

	//Reading configuration file
	dat, err := ioutil.ReadFile("jsons/mongoconfig.json")

	if err != nil {
		// TODO: Error handling
	}

	var result MongoConf
	err = json.Unmarshal(dat, &result)

	if err != nil {
		// TODO: Error handling
	}

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

	// Go routines and channel to orchestrate
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
		log.Printf("Listening on... %s", s.Addr)
		log.Fatal(http.ListenAndServe(addr, s.router))
	}()

	return exitChan
}

// addRoutes - Define API routes
func (s *Server) addRoutes() {

	/** RECIPES ENDPOINTS**/
	s.router.HandleFunc("/hrs/recipes", func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debugln("creating recipe...")

		var recipe Recipe
		var data []byte
		var err error

		status := http.StatusCreated
		hrsResp := initResponse()

		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&recipe)
		defer r.Body.Close()

		if err != nil {
			decodeError(&hrsResp, &data, &err)
			status = http.StatusConflict
		} else {
			hrsResp = s.worker.CreateRecipe(&recipe)
			data, err = json.Marshal(hrsResp)

			if err != nil {
				s.logger.Errorln("Json marshaling error")
				marshallError(&hrsResp, &data, &err)
			} else {
				s.logger.Infoln("Recipe created")
			}
		}

		w.WriteHeader(status)
		w.Write(data)
	}).Methods("POST")

	s.router.HandleFunc("/hrs/recipes/{id}", func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debugln("searching recipe...")
		status := http.StatusOK

		vars := mux.Vars(r)
		id := vars["id"]

		hrsResp := s.worker.GetRecipeByID(id)
		data, err := json.Marshal(hrsResp)

		if err != nil {
			s.logger.Errorln("Json marshaling error")
			marshallError(&hrsResp, &data, &err)
			status = http.StatusConflict
		} else {
			s.logger.Infoln("Recipe returned: ", id)
		}

		w.WriteHeader(status)
		w.Write(data)
	}).Methods("GET")

	s.router.HandleFunc("/hrs/recipes/{id}", func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debugln("patchting recipe...")
		var data []byte
		var recipe Recipe

		status := http.StatusOK
		vars := mux.Vars(r)
		id := vars["id"]
		hrsResp := initResponse()
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		err := decoder.Decode(&recipe)

		if err != nil {
			decodeError(&hrsResp, &data, &err)
			status = http.StatusConflict
		} else {
			hrsResp = s.worker.PatchRecipeByID(id, &recipe)
			data, err = json.Marshal(hrsResp)

			if err != nil {
				s.logger.Errorln("Json marshaling error")
				marshallError(&hrsResp, &data, &err)
				status = http.StatusConflict
			} else {
				s.logger.Infoln("Recipe patched: ", id)
			}

		}

		w.WriteHeader(status)
		w.Write(data)
	}).Methods("PATCH")

	s.router.HandleFunc("/hrs/recipes/{id}", func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debugln("deleting recipe...")
		status := http.StatusNoContent
		vars := mux.Vars(r)
		id := vars["id"]

		hrsResp := s.worker.DeleteRecipe(id)
		data, err := json.Marshal(hrsResp)

		if err != nil {
			s.logger.Errorln("Json marshaling error")
			marshallError(&hrsResp, &data, &err)
			status = http.StatusConflict
		} else {
			s.logger.Infoln("Recipe deleted: ", id)
		}

		w.WriteHeader(status)
		w.Write(data)
	}).Methods("DELETE")

	/** INGREDIENTS ENDPOINTS **/
	s.router.HandleFunc("/hrs/ingredients", func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debugln("creating ingredients...")

		var data []byte
		var err error

		status := http.StatusCreated
		hrsResp := initResponse()

		decoder := json.NewDecoder(r.Body)
		var ingredient Ingredient
		err = decoder.Decode(&ingredient)

		if err != nil {
			decodeError(&hrsResp, &data, &err)
			status = http.StatusConflict
		} else {
			hrsResp := s.worker.CreateIngredient(&ingredient)
			data, err = json.Marshal(hrsResp)

			if err != nil {
				s.logger.Errorln("Json marshaling error")
				marshallError(&hrsResp, &data, &err)
				status = http.StatusConflict
			} else {
				s.logger.Infoln("Ingredient created")
			}
		}
		defer r.Body.Close()

		w.WriteHeader(status)
		w.Write(data)
	}).Methods("POST")

	s.router.HandleFunc("/hrs/ingredients/{id}", func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debugln("searching ingredients...")
		status := http.StatusOK
		vars := mux.Vars(r)
		id := vars["id"]

		hrsResp := s.worker.GetIngredientByID(id)
		data, err := json.Marshal(hrsResp)

		if err != nil {
			s.logger.Errorln("Json marshaling error")
			marshallError(&hrsResp, &data, &err)
			status = http.StatusConflict
		} else {
			s.logger.Infoln("Ingredient returned: ", id)
		}

		w.WriteHeader(status)
		w.Write(data)
	}).Methods("GET")

	s.router.HandleFunc("/hrs/ingredients/{id}", func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debugln("patching ingredients...")
		var data []byte
		var ingredient Ingredient

		status := http.StatusOK
		hrsResp := initResponse()
		vars := mux.Vars(r)
		id := vars["id"]
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		err := decoder.Decode(&ingredient)

		if err != nil {
			decodeError(&hrsResp, &data, &err)
			status = http.StatusConflict
		} else {
			hrsResp = s.worker.PatchIngredientByID(id, &ingredient)
			data, err = json.Marshal(hrsResp)

			if err != nil {
				s.logger.Errorln("Json marshaling error")
				marshallError(&hrsResp, &data, &err)
				status = http.StatusConflict
			} else {
				s.logger.Infoln("Ingredient modified: ", id)
			}
		}

		w.WriteHeader(status)
		w.Write(data)
	}).Methods("PATCH")

	s.router.HandleFunc("/hrs/ingredients/{id}", func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debugln("deleting ingredient...")
		status := http.StatusNoContent
		vars := mux.Vars(r)
		id := vars["id"]

		hrsResp := s.worker.DeleteIngredient(id)
		data, err := json.Marshal(hrsResp)

		if err != nil {
			s.logger.Errorln("Json marshaling error")
			marshallError(&hrsResp, &data, &err)
			status = http.StatusConflict
		} else {
			s.logger.Infoln("Ingredient deleted: ", id)
		}

		w.WriteHeader(status)
		w.Write(data)
	}).Methods("DELETE")

	/** OTHER ENDPOINTS **/
	s.router.HandleFunc("/hrs/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "WTF\n")
	}).Methods("GET")
}

/** PRIVATE METHODS **/

func initResponse() HRAResponse {
	resp := HRAResponse{}
	return resp
}

func fatalResponse(err error) HRAResponse {
	status := Status{
		Code:        http.StatusConflict,
		Description: FATAL_ERROR,
	}
	hrsError := FatalError{}
	hrsError.SetError(err)
	resp := HRAResponse{
		Status: status,
		Error:  &hrsError,
	}

	return resp
}

func decodeError(hrsResp *HRAResponse, data *[]byte, err *error) {
	errRsp := initResponse()
	errRsp.Status = Status{
		Code:        http.StatusConflict,
		Description: DECODE_ERROR,
	}
	*data, *err = json.Marshal(errRsp)

	if err != nil {
		*hrsResp = fatalResponse(*err)
		*data, *err = json.Marshal(hrsResp)
	}
}

func marshallError(hrsResp *HRAResponse, data *[]byte, err *error) {
	*hrsResp = fatalResponse(*err)
	*data, *err = json.Marshal(hrsResp)

	if err != nil {
		*hrsResp = HRAResponse{}
	}
}
