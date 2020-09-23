package handlers

import (
	"fmt"
	"log"
	"net/http"

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

	var empList data.Employees
	var err error
	if empList, err = data.GetEmployees(); err != nil {
		emp.l.Println(err)
		http.Error(w, "Unable to fetch employee list from DB", http.StatusInternalServerError)
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

//GetEmployees ... http request handler to return all employees
func (emp *Employee) GetEmployee(w http.ResponseWriter, r *http.Request) {
	emp.l.Println("Handle Get employee")

	vars := mux.Vars(r)
	name := vars["name"]

	var empInfo *data.Employee
	var err error
	if empInfo, err = data.GetEmployee(name); err != nil {
		emp.l.Println(err)
		http.Error(w, "Unable to get employee info from DB", http.StatusInternalServerError)
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
		http.Error(w, "Error reading employee info", http.StatusBadRequest)
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
