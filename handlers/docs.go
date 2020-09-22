// Package classification of Employee API
//
// Documentation for Employee API
//
//	Schemes: http
//	BasePath: /api/v1
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta

package handlers

import "caching-service/data"

// A list of employee
// swagger:response employeesResponse
type employeesResponseWrapper struct {
	// All current products
	// in: body
	Body []data.Employee
}
