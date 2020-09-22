package handlers

import (
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

//GetEmployees ... http request handler to return all employees
func (emp *Employee) GetEmployees(w http.ResponseWriter, r *http.Request) {
	emp.l.Println("Handle Get all employees")

	empList := data.GetEmployees()

	err := empList.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to serialize employee list", http.StatusInternalServerError)
	}
}

//GetEmployees ... http request handler to return all employees
func (emp *Employee) GetEmployee(w http.ResponseWriter, r *http.Request) {
	emp.l.Println("Handle Get employee")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Unable to convert id", http.StatusBadRequest)
		return
	}
	empInfo := data.GetEmployee(id)

	err = empInfo.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to serialize employee info", http.StatusInternalServerError)
	}
}
