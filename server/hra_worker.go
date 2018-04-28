package server

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Init - Starts the worker
func (w *Worker) Init(ctx context.Context, logger *log.Entry) {
	w.SetLogger(logger)
	w.Ctx = ctx
}

// CreateRecipe - Creates a new recipe
func (w *Worker) CreateRecipe(recipe Recipe) HRAResponse {
	w.logger.Debugf("Worker - CreateRecipe [IN]")
	rsp := HRAResponse{}
	rsp.Status = Status{
		Code:        "200",
		Description: "Recipe created successfully!",
	}
	rsp.RespObj = &Recipe{
		ID:   "1", // TODO: How to use UUID¿
		Name: "Dummy recipe",
	}
	rsp.SetError(nil)
	w.logger.Debugf("Worker - CreateRecipe [OUT]")
	return rsp
}

// GetRecipeByID - Given an id, returns a recipe
func (w *Worker) GetRecipeByID(id string) HRAResponse {
	w.logger.Debugf("Worker - GetRecipebyId [IN]")
	rsp := HRAResponse{}

	dummyRecipe := Recipe{}

	if dummyRecipe.Name == "" {
		funcError := FunctionalError{}
		funcError.SetError(fmt.Errorf("Bad initialization"))

		//  everything is ok if we try to assign a value of type *FunctionalError to HRSError
		rsp.Error = &funcError
	}

	rsp.Status = Status{
		Code:        "200",
		Description: "Query completed",
	}

	rsp.Status.getObjectInfo()
	w.logger.Debugf("Worker - GetRecipebyId [OUT]")
	return rsp
}
