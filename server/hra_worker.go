package server

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"

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
		Code:        "recipeCreated",
		Description: "Recipe created successfully!",
	}
	uuid, _ := newUUID()
	rsp.RespObj = &Recipe{
		ID:   uuid,
		Name: "Dummy recipe",
	}
	rsp.SetError(nil)
	w.logger.Debugf(rsp.RespObj.getObjectInfo())
	w.logger.Debugf("Worker - CreateRecipe [OUT]")
	return rsp
}

// GetRecipeByID - Given an id, returns a recipe
func (w *Worker) GetRecipeByID(id string) HRAResponse {
	w.logger.Debugf("Worker - GetRecipebyId [IN]")
	rsp := HRAResponse{}

	dummyRecipe := Recipe{
		Name: "Dummy",
	}

	if dummyRecipe.Name == "" {
		funcError := FunctionalError{}
		funcError.SetError(fmt.Errorf("Bad initialization"))
		//  everything is ok if we try to assign a value of type *FunctionalError to HRSError
		rsp.Error = &funcError
	}

	rsp.Status = Status{
		Code:        "queryCompleted",
		Description: "Query completed",
	}
	rsp.RespObj = &dummyRecipe

	w.logger.Debugf(rsp.RespObj.getObjectInfo())
	w.logger.Debugf("Worker - GetRecipebyId [OUT]")
	return rsp
}

// PatchRecipeByID - Given a id, a recipe is patched
func (w *Worker) PatchRecipeByID(id string, recipe Recipe) HRAResponse {
	w.logger.Debugf("Worker - PatchRecipeByID [IN]")
	rsp := HRAResponse{}

	rsp.Status = Status{
		Code:        "updateCompleted",
		Description: "Recipe patched successfully",
	}
	w.logger.Debugf(rsp.RespObj.getObjectInfo())
	w.logger.Debugf("Worker - PatchRecipeByID [IN]")
	return rsp
}

// DeleteRecipe - Deletes a recipe by id
func (w *Worker) DeleteRecipe(id string) HRAResponse {
	w.logger.Debugf("Worker - DeleteRecipe [IN]")
	rsp := HRAResponse{}

	rsp.Status = Status{
		Code:        "deleteCompleted",
		Description: "Recipe patched successfully",
	}
	w.logger.Debugf(rsp.RespObj.getObjectInfo())
	w.logger.Debugf("Worker - DeleteRecipe [OUT]")
	return rsp
}

// CreateIngredient - creates an ingredient
func (w *Worker) CreateIngredient(ingredient Ingredient) HRAResponse {
	w.logger.Debugf("Worker - CreateIngredient [IN]")
	rsp := HRAResponse{}
	rsp.Status = Status{
		Code:        "ingredientCreated",
		Description: "Ingredient created successfully!",
	}
	uuid, _ := newUUID()
	rsp.RespObj = &Ingredient{
		ID:   uuid,
		Name: "Dummy ingredient",
	}
	rsp.SetError(nil)
	w.logger.Debugf(rsp.RespObj.getObjectInfo())
	w.logger.Debugf("Worker - CreateIngredient [OUT]")
	return rsp
}

// GetIngredientByID - Given an id, returns an ingredient
func (w *Worker) GetIngredientByID(id string) HRAResponse {
	w.logger.Debugf("Worker - GetIngredientByID [IN]")
	rsp := HRAResponse{}

	dummyIngredient := Ingredient{
		Name: "Dummy",
	}

	//TODO: Manage errors
	if dummyIngredient.Name == "" {
		funcError := FunctionalError{}
		funcError.SetError(fmt.Errorf("Bad initialization"))
		//  everything is ok if we try to assign a value of type *FunctionalError to HRSError
		rsp.Error = &funcError
	}

	rsp.Status = Status{
		Code:        "queryCompleted",
		Description: "Query completed",
	}
	rsp.RespObj = &dummyIngredient

	w.logger.Debugf(rsp.RespObj.getObjectInfo())
	w.logger.Debugf("Worker - GetIngredientByID [OUT]")
	return rsp
}

// PatchIngredientByID - Given a id, an ingredient is patched
func (w *Worker) PatchIngredientByID(id string, ingredient Ingredient) HRAResponse {
	w.logger.Debugf("Worker - PatchIngredientByID [IN]")
	rsp := HRAResponse{}

	rsp.Status = Status{
		Code:        "updateCompleted",
		Description: "Ingredient patched successfully",
	}

	w.logger.Debugf(rsp.RespObj.getObjectInfo())
	w.logger.Debugf("Worker - PatchIngredientByID [IN]")
	return rsp
}

// DeleteIngredient - Deletes an ingredient by id
func (w *Worker) DeleteIngredient(id string) HRAResponse {
	w.logger.Debugf("Worker - DeleteIngredient [IN]")
	rsp := HRAResponse{}

	rsp.Status = Status{
		Code:        "deleteCompleted",
		Description: "Ingredient patched successfully",
	}

	w.logger.Debugf(rsp.RespObj.getObjectInfo())
	w.logger.Debugf("Worker - DeleteIngredient [OUT]")
	return rsp
}

// newUUID generates a random UUID according to RFC 4122
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	resp := string(uuid[0:4]) + "-" + string(uuid[4:6]) + "-" + string(uuid[6:8]) + "-" + string(uuid[8:10]) + "-" + string(uuid[10:])
	return resp, nil
}
