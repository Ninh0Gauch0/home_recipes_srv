package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

/** MAIN TYPES **/

// Server struct
type Server struct {
	LoggerTrait
	Server      *http.Server
	Addr        string
	router      *mux.Router
	Ctx         context.Context
	worker      *Worker
	initialized bool
}

// Worker struct
type Worker struct {
	LoggerTrait
	Ctx context.Context
}

/* Logger */

// LoggerTrait - a logger trait that let's you configure a log
type LoggerTrait struct {
	logger *log.Entry
}

// SetLogger - let's you configure a logger
func (lt *LoggerTrait) SetLogger(l *log.Entry) {
	if l != nil {
		lt.logger = l
	}
}

// GetLogger - returns the logger
func (lt *LoggerTrait) GetLogger() *log.Entry {
	return lt.logger
}

/** General DTOs **/

// // ResponseObject interface
// type ResponseObject interface {
// 	getObjectInfo() string
// }

// /* Status Definition */

// // Status DTO
// type Status struct {
// 	Code        int    `json:"code"`
// 	Description string `json:"description"`
// }

// // Interface ResponseObject Implementation
// func (r *Status) getObjectInfo() string {
// 	info := []string{
// 		strconv.Itoa(r.Code),
// 		r.Description,
// 	}
// 	resp := strings.Join(info, ": ")
// 	return resp
// }
