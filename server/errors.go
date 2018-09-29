package server

import (
	"fmt"

	"github.com/Ninh0Gauch0/home_recipes_srv/hola"
)

// HRSError interface
type HRSError interface {
	ShowError() string
}

/** FatalError Error **/

// FatalError Error
type FatalError struct {
	Err error
}

// ShowError - Shows
func (fte *FatalError) ShowError() string {
	return holamundos.Hola()
	//return fmt.Sprintf("[ERROR] A fatal error occured: %v", fte.Err)
}

// SetError function
func (fte *FatalError) SetError(err error) {
	fte.Err = err
}

// GetError function
func (fte *FatalError) GetError() error {
	return fte.Err
}

/** FunctionalError Error **/

// FunctionalError Error
type FunctionalError struct {
	Err error
}

// ShowError - Shows
func (fe *FunctionalError) ShowError() string {
	return fmt.Sprintf("[ERROR] A functional error occured: %v", fe.Err)
}

// SetError function
func (fe *FunctionalError) SetError(err error) {
	fe.Err = err
}

// GetError function
func (fe *FunctionalError) GetError() error {
	return fe.Err
}

/** TechnicalError Error **/

// TechnicalError - Struct
type TechnicalError struct {
	Err error
}

// ShowError - Shows
func (te *TechnicalError) ShowError() string {
	return fmt.Sprintf("[ERROR] A functional error occured: %v", te.Err)
}

// SetError function
func (te *TechnicalError) SetError(err error) {
	te.Err = err
}

// GetError function
func (te *TechnicalError) GetError() error {
	return te.Err
}
