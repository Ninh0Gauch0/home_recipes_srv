package server

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

var (
	dummyRecipe = Recipe{
		Name: "Dummy",
	}
	dummyIngredient = Ingredient{
		Name: "Dummy",
	}
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
	id, _ := newUUID()
	/*	rsp.RespObj = &Recipe{
			ID:   uuid,
			Name: "Dummy recipe",
		}
	*/
	dummyRecipe.SetID(id)
	rsp.RespObj = &dummyRecipe
	rsp.SetError(nil)
	w.logger.Debugf(rsp.RespObj.getObjectInfo())
	w.logger.Debugf("Worker - CreateRecipe [OUT]")
	return rsp
}

// GetRecipeByID - Given an id, returns a recipe
func (w *Worker) GetRecipeByID(id string) HRAResponse {
	w.logger.Debugf("Worker - GetRecipebyId [IN]")
	rsp := HRAResponse{}

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
	rsp.RespObj = &dummyRecipe

	w.logger.Debugf(rsp.RespObj.getObjectInfo())
	w.logger.Debugf("Worker - PatchRecipeByID [OUT]")
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
	w.logger.Debugf(rsp.Status.getObjectInfo())
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
	id, _ := newUUID()
	dummyIngredient.SetID(id)
	rsp.RespObj = &dummyIngredient
	rsp.SetError(nil)

	w.logger.Debugf(rsp.RespObj.getObjectInfo())
	w.logger.Debugf("Worker - CreateIngredient [OUT]")
	return rsp
}

// GetIngredientByID - Given an id, returns an ingredient
func (w *Worker) GetIngredientByID(id string) HRAResponse {
	w.logger.Debugf("Worker - GetIngredientByID [IN]")
	rsp := HRAResponse{}

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
	rsp.RespObj = &dummyIngredient

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

	w.logger.Debugf(rsp.Status.getObjectInfo())
	w.logger.Debugf("Worker - DeleteIngredient [OUT]")
	return rsp
}

// newUUID generates a random UUID according to RFC 4122
func newUUID() (string, error) {
	uuid, err := uuid.NewV4()

	if err != nil {
		return "", err
	}

	return uuid.String(), nil
}
