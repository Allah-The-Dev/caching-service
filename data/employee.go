package data

import (
	"encoding/json"
	"io"
)

//Employee ...
type Employee struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Unit string `json:"unit"`
}

//ToJSON ... serialize data and write to a destination
func (emp *Employee) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(emp)
}

//Employees ...
type Employees []*Employee

//ToJSON ... serialize data and write to a destination
func (emps *Employees) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(emps)
}

//GetEmployees ...
func GetEmployees() Employees {
	return empList
}

//GetEmployee ...
func GetEmployee(id int) *Employee {
	for _, emp := range empList {
		if emp.ID == id {
			return emp
		}
	}
	return nil
}

var empList = []*Employee{
	&Employee{
		ID:   1,
		Name: "foo",
		Unit: "ISRO",
	},
	&Employee{
		ID:   2,
		Name: "bar",
		Unit: "NASA",
	},
}
