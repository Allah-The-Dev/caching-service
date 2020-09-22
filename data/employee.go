package data

import (
	"encoding/json"
	"io"
)

//Employee ...
// swagger:model
type Employee struct {
	// the id of employee
	//
	// required: false
	// min: 1
	ID int `json:"id"`
	// the name of employee
	//
	// required: true
	// max length: 255
	Name string `json:"name"`
	// the unit of employee
	//
	// required: true
	// max length: 255
	Unit string `json:"unit"`
}

//ToJSON ... serialize data and write to a destination
func (emp *Employee) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(emp)
}

//FromJSON .... deserialize employee data into a destination pointer
func (emp *Employee) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(emp)
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

//AddEmployee ...
func AddEmployee(emp *Employee) int {
	emp.ID = getNextID()
	empList = append(empList, emp)
	return emp.ID
}

func getNextID() int {
	lastEmp := empList[len(empList)-1]
	return lastEmp.ID + 1
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
