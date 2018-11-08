package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ninh0gauch0/hrstypes"
	mongo "github.com/ninh0gauch0/mongoconnector"
	mngtypes "github.com/ninh0gauch0/mongoconnector/types"
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
	// OPNOTCOMPLETED Constant
	OPNOTCOMPLETED = "Operation not completed"
	// INGREDIENTCOLL Constant
	INGREDIENTCOLL = "ingredients"
	// RECIPECOLL Constant
	RECIPECOLL = "recipes"
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

	manager := mongo.Manager{
		Ctx: w.Ctx,
	}

	if manager.Init() {
		recipeDTO := mapRecipeToMetadataObject(recipe)

		res, err := manager.ExecuteInsert(RECIPECOLL, recipeDTO)

		if err == nil {
			if res == 0 {
				rsp.Status = hrstypes.Status{
					Code:        http.StatusCreated,
					Description: CREATED,
				}
				recipe = mapRecipeToDTOObject(recipeDTO)
				rsp.RespObj = recipe
				rsp.SetError(nil)
			} else {
				techErr := hrstypes.TechnicalError{}
				return generateErrorResponse(OPNOTCOMPLETED, fmt.Sprintf("Insertion can't be accomplished"), techErr, http.StatusConflict)
			}
		} else {
			return generateErrorResponse(TECHNICAL, fmt.Sprintf("Fatal error trying to insert"), err, http.StatusInternalServerError)
		}
	} else {
		techErr := hrstypes.TechnicalError{}
		return generateErrorResponse(TECHNICAL, fmt.Sprintf("Connection problem"), techErr, http.StatusInternalServerError)
	}

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

	manager := mongo.Manager{
		Ctx: w.Ctx,
	}

	if manager.Init() {
		res, err := manager.ExecuteSearchByID(RECIPECOLL, id)

		if err == nil {
			if res != nil {
				rsp.Status = hrstypes.Status{
					Code:        http.StatusOK,
					Description: QUERIED,
				}

				// Casting MetadataObject to Recipe
				original, ok := res.(*mngtypes.Recipe)

				if ok {
					recipe := mapRecipeToDTOObject(original)
					rsp.RespObj = recipe
					rsp.SetError(nil)
				} else {
					techErr := hrstypes.TechnicalError{}
					return generateErrorResponse(OPNOTCOMPLETED, fmt.Sprintf("Error mapping the response"), techErr, http.StatusConflict)
				}
			} else {
				techErr := hrstypes.TechnicalError{}
				return generateErrorResponse(OPNOTCOMPLETED, fmt.Sprintf("Query can't be accomplished"), techErr, http.StatusConflict)
			}
		} else {
			return generateErrorResponse(TECHNICAL, fmt.Sprintf("Fatal error trying to query"), err, http.StatusInternalServerError)
		}
	} else {
		techErr := hrstypes.TechnicalError{}
		return generateErrorResponse(TECHNICAL, fmt.Sprintf("Connection problem"), techErr, http.StatusInternalServerError)
	}

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

	manager := mongo.Manager{
		Ctx: w.Ctx,
	}

	if manager.Init() {
		res, err := manager.ExecuteUpdate(RECIPECOLL, id, recipe)

		if err == nil {
			if res == 0 {
				rsp.Status = hrstypes.Status{
					Code:        http.StatusOK,
					Description: PATCHED,
				}
				rsp.RespObj = nil
				rsp.SetError(nil)
			} else {
				techErr := hrstypes.TechnicalError{}
				return generateErrorResponse(OPNOTCOMPLETED, fmt.Sprintf("Patch can't be accomplished"), techErr, http.StatusConflict)
			}
		} else {
			return generateErrorResponse(TECHNICAL, fmt.Sprintf("Fatal error trying to patch"), err, http.StatusInternalServerError)
		}
	} else {
		techErr := hrstypes.TechnicalError{}
		return generateErrorResponse(TECHNICAL, fmt.Sprintf("Connection problem"), techErr, http.StatusInternalServerError)
	}

	w.logger.Debugf(rsp.RespObj.GetObjectInfo())
	w.logger.Debugf("Worker - PatchRecipeByID [OUT]")
	return rsp
}

// DeleteRecipe - Deletes a recipe by id
func (w *Worker) DeleteRecipe(id string) hrstypes.HRAResponse {
	w.logger.Debugf("Worker - DeleteRecipe [IN]")
	rsp := hrstypes.HRAResponse{}

	if id == "" {
		err := hrstypes.FunctionalError{}
		rsp = generateErrorResponse(FAIL, fmt.Sprintf("Mandatory parameter %s", id), err, http.StatusConflict)
		return rsp
	}

	manager := mongo.Manager{
		Ctx: w.Ctx,
	}

	if manager.Init() {
		res, err := manager.ExecuteDelete(RECIPECOLL, id)

		if err == nil {
			if res == 0 {
				rsp.Status = hrstypes.Status{
					Code:        http.StatusNoContent,
					Description: REMOVED,
				}
				rsp.RespObj = nil
				rsp.SetError(nil)
			} else {
				techErr := hrstypes.TechnicalError{}
				return generateErrorResponse(OPNOTCOMPLETED, fmt.Sprintf("Remove can't be accomplished"), techErr, http.StatusConflict)
			}
		} else {
			return generateErrorResponse(TECHNICAL, fmt.Sprintf("Fatal error trying to remove"), err, http.StatusInternalServerError)
		}
	} else {
		techErr := hrstypes.TechnicalError{}
		return generateErrorResponse(TECHNICAL, fmt.Sprintf("Connection problem"), techErr, http.StatusInternalServerError)
	}

	w.logger.Debugf(rsp.Status.GetObjectInfo())
	w.logger.Debugf("Worker - DeleteRecipe [OUT]")
	return rsp
}

// CreateIngredient - creates an ingredient
func (w *Worker) CreateIngredient(ingredient *hrstypes.Ingredient) hrstypes.HRAResponse {
	w.logger.Debugf("Worker - CreateIngredient [IN]")
	rsp := hrstypes.HRAResponse{}

	manager := mongo.Manager{
		Ctx: w.Ctx,
	}

	if manager.Init() {
		ingredientDTO := mapIngredientToMetadataObject(ingredient)

		res, err := manager.ExecuteInsert(INGREDIENTCOLL, ingredientDTO)

		if err == nil {
			if res == 0 {
				rsp.Status = hrstypes.Status{
					Code:        http.StatusCreated,
					Description: CREATED,
				}
				ingredient = mapIngredientToDTOObject(ingredientDTO)
				rsp.RespObj = ingredient
				rsp.SetError(nil)
			} else {
				techErr := hrstypes.TechnicalError{}
				return generateErrorResponse(OPNOTCOMPLETED, fmt.Sprintf("Insertion can't be accomplished"), techErr, http.StatusConflict)
			}
		} else {
			return generateErrorResponse(TECHNICAL, fmt.Sprintf("Fatal error trying to insert"), err, http.StatusInternalServerError)
		}
	} else {
		techErr := hrstypes.TechnicalError{}
		return generateErrorResponse(TECHNICAL, fmt.Sprintf("Connection problem"), techErr, http.StatusInternalServerError)
	}

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

	manager := mongo.Manager{
		Ctx: w.Ctx,
	}

	if manager.Init() {
		res, err := manager.ExecuteSearchByID(INGREDIENTCOLL, id)

		if err == nil {
			if res != nil {
				rsp.Status = hrstypes.Status{
					Code:        http.StatusOK,
					Description: QUERIED,
				}

				// Casting MetadataObject to Ingredient
				original, ok := res.(*mngtypes.Ingredient)

				if ok {
					ingredient := mapIngredientToDTOObject(original)
					rsp.RespObj = ingredient
					rsp.SetError(nil)
				} else {
					techErr := hrstypes.TechnicalError{}
					return generateErrorResponse(OPNOTCOMPLETED, fmt.Sprintf("Error mapping the response"), techErr, http.StatusConflict)
				}
			} else {
				techErr := hrstypes.TechnicalError{}
				return generateErrorResponse(OPNOTCOMPLETED, fmt.Sprintf("Query can't be accomplished"), techErr, http.StatusConflict)
			}
		} else {
			return generateErrorResponse(TECHNICAL, fmt.Sprintf("Fatal error trying to query"), err, http.StatusInternalServerError)
		}

	} else {
		techErr := hrstypes.TechnicalError{}
		return generateErrorResponse(TECHNICAL, fmt.Sprintf("Connection problem"), techErr, http.StatusInternalServerError)
	}

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

	manager := mongo.Manager{
		Ctx: w.Ctx,
	}

	if manager.Init() {
		res, err := manager.ExecuteUpdate(INGREDIENTCOLL, id, ingredient)

		if err == nil {
			if res == 0 {
				rsp.Status = hrstypes.Status{
					Code:        http.StatusOK,
					Description: PATCHED,
				}
				rsp.RespObj = nil
				rsp.SetError(nil)
			} else {
				techErr := hrstypes.TechnicalError{}
				return generateErrorResponse(OPNOTCOMPLETED, fmt.Sprintf("Patch can't be accomplished"), techErr, http.StatusConflict)
			}
		} else {
			return generateErrorResponse(TECHNICAL, fmt.Sprintf("Fatal error trying to patch"), err, http.StatusInternalServerError)
		}
	} else {
		techErr := hrstypes.TechnicalError{}
		return generateErrorResponse(TECHNICAL, fmt.Sprintf("Connection problem"), techErr, http.StatusInternalServerError)
	}

	w.logger.Debugf(rsp.RespObj.GetObjectInfo())
	w.logger.Debugf("Worker - PatchIngredientByID [IN]")
	return rsp
}

// DeleteIngredient - Deletes an ingredient by id
func (w *Worker) DeleteIngredient(id string) hrstypes.HRAResponse {
	w.logger.Debugf("Worker - DeleteIngredient [IN]")
	rsp := hrstypes.HRAResponse{}

	if id == "" {
		err := hrstypes.FunctionalError{}
		rsp = generateErrorResponse(FAIL, fmt.Sprintf("Mandatory parameter %s", id), err, http.StatusConflict)
		return rsp
	}

	manager := mongo.Manager{
		Ctx: w.Ctx,
	}

	if manager.Init() {
		res, err := manager.ExecuteDelete(INGREDIENTCOLL, id)

		if err == nil {
			if res == 0 {
				rsp.Status = hrstypes.Status{
					Code:        http.StatusNoContent,
					Description: REMOVED,
				}
				rsp.RespObj = nil
				rsp.SetError(nil)
			} else {
				techErr := hrstypes.TechnicalError{}
				return generateErrorResponse(OPNOTCOMPLETED, fmt.Sprintf("Remove can't be accomplished"), techErr, http.StatusConflict)
			}
		} else {
			return generateErrorResponse(TECHNICAL, fmt.Sprintf("Fatal error trying to remove"), err, http.StatusInternalServerError)
		}
	} else {
		techErr := hrstypes.TechnicalError{}
		return generateErrorResponse(TECHNICAL, fmt.Sprintf("Connection problem"), techErr, http.StatusInternalServerError)
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
	rsp.RespObj = nil

	switch err.(type) { // this is an assert; asserting, "x is of this type"
	case hrstypes.TechnicalError:
		hrsError := hrstypes.TechnicalError{}
		hrsError.SetError(errors.New(errorMsg))
		//  everything is ok if we try to assign a value of type *technicalError to HRSError
		rsp.SetError(&hrsError)
	case hrstypes.FunctionalError:
		hrsError := hrstypes.FunctionalError{}
		hrsError.SetError(errors.New(errorMsg))
		rsp.SetError(&hrsError)
	default:
		hrsError := hrstypes.FatalError{}
		hrsError.SetError(errors.New(errorMsg))
		rsp.SetError(&hrsError)
	}

	rsp.Status = hrstypes.Status{
		Code:        status,
		Description: desc,
	}
	return rsp
}

// -- MAPPERS -- //

// maps a recipe DTOOBject to a recipe MetadataObject
func mapRecipeToMetadataObject(recipe *hrstypes.Recipe) *mngtypes.Recipe {

	rsp := &mngtypes.Recipe{}

	if recipe != nil {
		rsp.SetID(recipe.GetID())
		rsp.SetName(recipe.GetName())
		rsp.SetDescription(recipe.GetDescription())
		rsp.SetIngredients(recipe.GetIngredients())
		rsp.SetSteps(recipe.GetSteps())
	}
	return rsp
}

// maps a recipe MetadataObject to a recipe DTOOBject
func mapRecipeToDTOObject(recipe *mngtypes.Recipe) *hrstypes.Recipe {
	rsp := &hrstypes.Recipe{}

	if recipe != nil {
		rsp.SetID(recipe.GetID())
		rsp.SetName(recipe.GetName())
		rsp.SetDescription(recipe.GetDescription())
		rsp.SetIngredients(recipe.GetIngredients())
		rsp.SetSteps(recipe.GetSteps())
	}
	return rsp
}

// maps an ingredient DTOOBject to an ingredient MetadataObject
func mapIngredientToMetadataObject(ingredient *hrstypes.Ingredient) *mngtypes.Ingredient {
	rsp := &mngtypes.Ingredient{}

	if ingredient != nil {
		rsp.SetID(ingredient.GetID())
		rsp.SetName(ingredient.GetName())
		rsp.SetDescription(ingredient.GetDescription())
		rsp.SetQuantity(ingredient.GetQuantity())
	}
	return rsp
}

// maps an ingredient MetadataObject to an ingredient DTOOBject
func mapIngredientToDTOObject(ingredient *mngtypes.Ingredient) *hrstypes.Ingredient {
	rsp := &hrstypes.Ingredient{}

	if ingredient != nil {
		rsp.SetID(ingredient.GetID())
		rsp.SetName(ingredient.GetName())
		rsp.SetDescription(ingredient.GetDescription())
		rsp.SetQuantity(ingredient.GetQuantity())
	}
	return rsp
}
