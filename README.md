# Home Recipes Server

## Goals

This app implements a recipes app server.

## Technical solution

We created this application using Golang to define REST services to manage a mongo database, allowing CRUD operations for recipes and ingredients. These Golang services use:

- A cli to allow to run our application as a system service, with its options and flags.
- Use of reflection over an interface (a bit different to use it over a struct directly). See mongoconnector, patch operation.
- Use of channeling to orchestate our application
- Use of mux library to expose our endpoints
- Use of mongo connector library (at this moment __mgo__, but [mongo-go-driver](https://github.com/mongodb/mongo-go-driver) comming soon).

Example of reflection:

```go
import (
	"fmt"
	"reflect"
		"strconv"
	"strings"
)


type T struct {
    A int
    B string
}

type MetadataObject interface {
	GetObjectInfo() string
}

// Recipe DTO
type Recipe struct {
	Code        string   `json:"code" bson:"_id,omitempty"`
	Name        string   `json:"name" bson:"name,omitempty"`
	Description string   `json:"description" bson:"description,omitempty"`
	Steps       []string `json:"steps" bson:"steps,omitempty"`
	Ingredients []string `json:"ingredients" bson:"ingredients,omitempty"`
}

// GetCode function
func (r *Recipe) GetCode() string {
	return r.Code
}

// SetCode function
func (r *Recipe) SetCode(code string) {
	r.Code = code
}

// GetName function
func (r *Recipe) GetName() string {
	return r.Name
}

// SetName function
func (r *Recipe) SetName(name string) {
	r.Name = name
}

// GetDescription function
func (r *Recipe) GetDescription() string {
	return r.Description
}

// SetDescription function
func (r *Recipe) SetDescription(description string) {
	r.Description = description
}

// GetSteps function
func (r *Recipe) GetSteps() []string {
	return r.Steps
}

// SetSteps function
func (r *Recipe) SetSteps(steps []string) {
	r.Steps = steps
}

// SetIngredients function
func (r *Recipe) SetIngredients(ingredients []string) {
	r.Ingredients = ingredients
}

// GetIngredients function
func (r *Recipe) GetIngredients() []string {
	return r.Ingredients
}

// GetObjectInfo - Interface DTOObject Implementation
func (r *Recipe) GetObjectInfo() string {
	info := []string{
		r.GetName(),
		r.GetDescription(),
	}
	resp := strings.Join(info, ": ")

	for i, step := range r.GetSteps() {
		resp += "\nStep " + strconv.Itoa(i) + ":" + step
	}

	return resp
}

// Function that receives a metadataObject (an Interface, not struct)
func ExecuteUpdate(obj MetadataObject) {
    //WITHOUT &
	s := reflect.ValueOf(obj).Elem()
	typeOfT := s.Type()
	
	for i := 0; i < s.NumField(); i++ {
    		f := s.Field(i)
    		fmt.Printf("%d: %s %s = %v\n", i,
        	typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
}

func main() {
    // REFLECTION ON A STRUCT
	recipe:= Recipe{"Code value","Name value", "Description value", nil, nil}
	// WITH &
	s = reflect.ValueOf(&recipe).Elem()
	typeOfT = s.Type()
	
		for i := 0; i < s.NumField(); i++ {
    		f := s.Field(i)
    		fmt.Printf("%d: %s %s = %v\n", i,
        	typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
    
    // REFLECTION ON A FUNCTION WITH A INTERFACE PARAMETER. USING & HERE, NOT INSIDE
	ExecuteUpdate(&recipe)
}

```
