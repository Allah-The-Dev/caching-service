package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"caching-service/data"

	"github.com/gorilla/mux"
)

//Employee ... http.Handler
type Employee struct {
	l *log.Logger
}

//NewEmployee ... constructor for new Employee
func NewEmployee(l *log.Logger) *Employee {
	return &Employee{l}
}

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

// swagger:route GET /employees Employee listEmployee
// Return list of employees available in cache
//
// responses:
//	200: employeesResponse

//GetEmployees ... http request handler to return all employees
func (emp *Employee) GetEmployees(w http.ResponseWriter, r *http.Request) {
	emp.l.Println("Handle Get all employees")

	var pageNo, pageSize int
	var err error
	var pageNoStr, pageSizeStr string
	var pageNoPresent, pageSizePresent bool
	//query params
	query := r.URL.Query()
	if pageNoStr = query.Get("pageNo"); pageNoStr != "" {
		pageNoPresent = true
		if pageNo, err = validatePageNoAndSize(pageNoStr); err != nil {
			emp.l.Println("[Error] " + err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	if pageSizeStr = query.Get("pageSize"); pageSizeStr != "" {
		pageSizePresent = true
		if pageSize, err = validatePageNoAndSize(pageSizeStr); err != nil {
			emp.l.Println("[Error] " + err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if pageNoPresent && !pageSizePresent {
		errStr := "[Error] if page no is provided, pageSize must be provided in request"
		emp.l.Println(errStr)
		http.Error(w, errStr, http.StatusBadRequest)
		return
	}

	if !pageNoPresent && pageSizePresent {
		errStr := "[Error] if page size is provided, pageNo must be provided in request"
		emp.l.Println(errStr)
		http.Error(w, errStr, http.StatusBadRequest)
		return
	}

	//db call
	var empList data.Employees
	if empList, err = data.GetEmployees(pageNo, pageSize); err != nil {
		emp.l.Println(err)
		http.Error(w, "Unable to fetch employee list from DB", http.StatusInternalServerError)
		return
	}

	err = empList.ToJSON(w)
	if err != nil {
		emp.l.Println(err)
		http.Error(w, "Unable to serialize employee list", http.StatusInternalServerError)
	}
}

// swagger:route GET /employee/{id} Employee employeeInfo
// Return one employee from the cache based on id
//
// responses:
//	200: employeeResponse
//	400: badRequestResponse
// 	500: internalServerErrorResponse

//GetEmployee ... http request handler to return all employees
func (emp *Employee) GetEmployee(w http.ResponseWriter, r *http.Request) {
	emp.l.Println("Handle Get employee")

	vars := mux.Vars(r)
	name := vars["name"]

	var empInfo *data.Employee
	var err error
	if empInfo, err = data.GetEmployee(name); err != nil {
		emp.l.Println(err)
		http.Error(w, "Unable to get employee info from DB", http.StatusInternalServerError)
		return
	}

	err = empInfo.ToJSON(w)
	if err != nil {
		emp.l.Println(err)
		http.Error(w, "Unable to serialize employee info", http.StatusInternalServerError)
		return
	}
}

// swagger:route POST /employee Employee addEmployee
// Return new employee id for posted employee data
//
// responses:
//	200: employeesResponse
//	400: badRequestResponse

//AddEmployee ...
func (emp *Employee) AddEmployee(w http.ResponseWriter, r *http.Request) {
	emp.l.Println("Handle Post employee")

	empInfo := data.Employee{}
	if err := empInfo.FromJSON(r.Body); err != nil {
		emp.l.Println("[Error] unable to deserialize employee data")
		http.Error(w, "request body is not in proper format", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	newEmpID, err := data.AddEmployee(&empInfo)
	if err != nil {
		emp.l.Println("[Error] unable to insert data in database")
		http.Error(w, "Error inserting employee info into DB", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	msg := fmt.Sprintf("Employee added successfully with id %v", newEmpID)
	w.Write([]byte(msg))
}

func validatePageNoAndSize(queryParamValue string) (int, error) {
	var intQPValue int
	var err error
	if intQPValue, err = strconv.Atoi(queryParamValue); err != nil {
		return intQPValue, errors.New("pageNo and pageSize must be number")
	}
	if intQPValue <= 0 {
		return intQPValue, errors.New("pageNo and pageSize must be greater than zero")
	}
	return intQPValue, nil
}
