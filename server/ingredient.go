package server

import (
	"strconv"
	"strings"
)

/** INGREDIENT DEFINITION **/

// Ingredient DTO
type Ingredient struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}

// Interface ResponseObject Implementations
func (i *Ingredient) getObjectInfo() string {
	info := []string{
		i.Name,
		i.Description,
	}
	resp := strings.Join(info, ": ")
	quantity := "\nQuantity: " + strconv.Itoa(i.Quantity)
	resp += quantity
	return resp
}

// GetID function
func (i *Ingredient) GetID() string {
	return i.ID
}

// SetID function
func (i *Ingredient) SetID(id string) {
	i.ID = id
}

func (i *Ingredient) getName() string {
	return i.Name
}

func (i *Ingredient) setName(name string) {
	i.Name = name
}

func (i *Ingredient) getDescription() string {
	return i.Description
}

func (i *Ingredient) setDescription(description string) {
	i.Description = description
}

func (i *Ingredient) getQuantity() int {
	return i.Quantity
}

func (i *Ingredient) setQuantity(quantity int) {
	i.Quantity = quantity
}
