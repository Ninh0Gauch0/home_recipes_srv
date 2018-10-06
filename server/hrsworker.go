package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ninh0gauch0/hrstypes"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

const (
	// CREATED Constant
	CREATED = "Created successfully"
	// QUERIED Constant
	QUERIED = "Query completed"
	// PATCHED Constant
	PATCHED = "Element patched successfully"
	// REMOVED Constant
	REMOVED = "Element removed"
	// FAIL Constant
	FAIL = "Failed validation"
	// TECHNICAL Constant
	TECHNICAL = "Technical error"
)

var (
	dummyRecipe = hrstypes.Recipe{
		Name:        "Dummy Recipe",
		Description: "This recipe sucks!",
		Steps:       make([]string, 0, 4),
	}
	dummyIngredient = hrstypes.Ingredient{
		Name:        "Dummy Ingredient",
		Description: "This ingredient sucks!",
		Quantity:    2,
	}
)

// Init - Starts the worker
func (w *Worker) Init(ctx context.Context, logger *log.Entry) {
	w.SetLogger(logger)
	w.Ctx = ctx
}

// CreateRecipe - Creates a new recipe
func (w *Worker) CreateRecipe(recipe *hrstypes.Recipe) hrstypes.HRAResponse {
	w.logger.Debugf("Worker - CreateRecipe [IN]")

	rsp := hrstypes.HRAResponse{}
	rsp.Status = hrstypes.Status{
		Code:        http.StatusCreated,
		Description: CREATED,
	}
	id, err := newUUID()

	if err != nil {
		rsp = generateErrorResponse(TECHNICAL, fmt.Sprintf("uuid generation error"), err, http.StatusConflict)
		return rsp
	}

	dummyRecipe.SetID(id)
	dummyRecipe.Steps = append(dummyRecipe.Steps, "Llamar al Foster", "Llamar al chino")
	rsp.RespObj = &dummyRecipe
	rsp.SetError(nil)

	w.logger.Debugf(rsp.RespObj.GetObjectInfo())
	w.logger.Debugf("Worker - CreateRecipe [OUT]")
	return rsp
}

// GetRecipeByID - Given an id, returns a recipe
func (w *Worker) GetRecipeByID(id string) hrstypes.HRAResponse {
	w.logger.Debugf("Worker - GetRecipebyId [IN]")
	rsp := hrstypes.HRAResponse{}
	if id == "" {
		err := hrstypes.FunctionalError{}
		rsp = generateErrorResponse(FAIL, fmt.Sprintf("Mandatory parameter %s", id), err, http.StatusConflict)
		return rsp
	}

	rsp.Status = hrstypes.Status{
		Code:        http.StatusOK,
		Description: QUERIED,
	}
	rsp.RespObj = &dummyRecipe

	w.logger.Debugf(rsp.RespObj.GetObjectInfo())
	w.logger.Debugf("Worker - GetRecipebyId [OUT]")
	return rsp
}

// PatchRecipeByID - Given a id, a recipe is patched
func (w *Worker) PatchRecipeByID(id string, recipe *hrstypes.Recipe) hrstypes.HRAResponse {
	w.logger.Debugf("Worker - PatchRecipeByID [IN]")
	rsp := hrstypes.HRAResponse{}
	if id == "" {
		err := hrstypes.FunctionalError{}
		rsp = generateErrorResponse(FAIL, fmt.Sprintf("Mandatory parameter %s", id), err, http.StatusConflict)
		return rsp
	}

	rsp.Status = hrstypes.Status{
		Code:        http.StatusOK,
		Description: PATCHED,
	}
	rsp.RespObj = &dummyRecipe

	w.logger.Debugf(rsp.RespObj.GetObjectInfo())
	w.logger.Debugf("Worker - PatchRecipeByID [OUT]")
	return rsp
}

// DeleteRecipe - Deletes a recipe by id
func (w *Worker) DeleteRecipe(id string) hrstypes.HRAResponse {
	w.logger.Debugf("Worker - DeleteRecipe [IN]")
	rsp := hrstypes.HRAResponse{}

	rsp.Status = hrstypes.Status{
		Code:        http.StatusNoContent,
		Description: REMOVED,
	}
	w.logger.Debugf(rsp.Status.GetObjectInfo())
	w.logger.Debugf("Worker - DeleteRecipe [OUT]")
	return rsp
}

// CreateIngredient - creates an ingredient
func (w *Worker) CreateIngredient(ingredient *hrstypes.Ingredient) hrstypes.HRAResponse {
	w.logger.Debugf("Worker - CreateIngredient [IN]")
	rsp := hrstypes.HRAResponse{}
	rsp.Status = hrstypes.Status{
		Code:        http.StatusCreated,
		Description: CREATED,
	}
	id, err := newUUID()

	if err != nil {
		rsp = generateErrorResponse(TECHNICAL, fmt.Sprintf("uuid generation error"), err, http.StatusConflict)
		return rsp
	}

	dummyIngredient.SetID(id)
	rsp.RespObj = &dummyIngredient
	rsp.SetError(nil)

	w.logger.Debugf(rsp.RespObj.GetObjectInfo())
	w.logger.Debugf("Worker - CreateIngredient [OUT]")
	return rsp
}

// GetIngredientByID - Given an id, returns an ingredient
func (w *Worker) GetIngredientByID(id string) hrstypes.HRAResponse {
	w.logger.Debugf("Worker - GetIngredientByID [IN]")

	rsp := hrstypes.HRAResponse{}
	if id == "" {
		err := hrstypes.FunctionalError{}
		rsp = generateErrorResponse(FAIL, fmt.Sprintf("Mandatory parameter %s", id), err, http.StatusConflict)
		return rsp
	}

	rsp.Status = hrstypes.Status{
		Code:        http.StatusOK,
		Description: QUERIED,
	}
	rsp.RespObj = &dummyIngredient

	w.logger.Debugf(rsp.RespObj.GetObjectInfo())
	w.logger.Debugf("Worker - GetIngredientByID [OUT]")
	return rsp
}

// PatchIngredientByID - Given a id, an ingredient is patched
func (w *Worker) PatchIngredientByID(id string, ingredient *hrstypes.Ingredient) hrstypes.HRAResponse {
	w.logger.Debugf("Worker - PatchIngredientByID [IN]")
	rsp := hrstypes.HRAResponse{}

	if id == "" {
		err := hrstypes.FunctionalError{}
		rsp = generateErrorResponse(FAIL, fmt.Sprintf("Mandatory parameter %s", id), err, http.StatusConflict)
		return rsp
	}

	rsp.Status = hrstypes.Status{
		Code:        http.StatusOK,
		Description: PATCHED,
	}
	rsp.RespObj = &dummyIngredient

	w.logger.Debugf(rsp.RespObj.GetObjectInfo())
	w.logger.Debugf("Worker - PatchIngredientByID [IN]")
	return rsp
}

// DeleteIngredient - Deletes an ingredient by id
func (w *Worker) DeleteIngredient(id string) hrstypes.HRAResponse {
	w.logger.Debugf("Worker - DeleteIngredient [IN]")
	rsp := hrstypes.HRAResponse{}

	rsp.Status = hrstypes.Status{
		Code:        http.StatusNoContent,
		Description: "Ingredient patched successfully",
	}

	w.logger.Debugf(rsp.Status.GetObjectInfo())
	w.logger.Debugf("Worker - DeleteIngredient [OUT]")
	return rsp
}

/** PRIVATE METHODS **/

// newUUID generates a random UUID according to RFC 4122
func newUUID() (string, error) {
	newUUID, err := uuid.NewV4()

	if err != nil {
		return "", err
	}

	return newUUID.String(), nil
}

// GenerateErrorResponse - generates a error response
func generateErrorResponse(errorMsg string, desc string, err interface{}, status int) hrstypes.HRAResponse {
	rsp := hrstypes.HRAResponse{}

	switch err.(type) { // this is an assert; asserting, "x is of this type"
	case hrstypes.FunctionalError:
		hrsError := hrstypes.TechnicalError{}
		hrsError.SetError(errors.New(errorMsg))
		//  everything is ok if we try to assign a value of type *technicalError to HRSError
		rsp.SetError(&hrsError)
	case hrstypes.TechnicalError:
		hrsError := hrstypes.FunctionalError{}
		hrsError.SetError(errors.New(errorMsg))
		rsp.SetError(&hrsError)
	default:
		hrsError := hrstypes.FunctionalError{}
		hrsError.SetError(errors.New(errorMsg))
		rsp.SetError(&hrsError)
	}

	rsp.Status = hrstypes.Status{
		Code:        status,
		Description: desc,
	}
	return rsp
}
