package data

import (
	"context"
	"encoding/json"
	"io"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var CLogger *log.Logger

//Employee ...
// swagger:model
type Employee struct {
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
func GetEmployees() (Employees, error) {

	empList := make([]*Employee, 0)

	cur, err := GetCollection().Find(context.TODO(), bson.D{{}})
	if err != nil {
		return empList, err
	}

	for cur.Next(context.TODO()) {
		var emp Employee
		if err := cur.Decode(&emp); err != nil {
			return empList, err
		}
		empList = append(empList, &emp)
	}
	return empList, nil
}

//GetEmployee ...
func GetEmployee(name string) (*Employee, error) {

	emp := &Employee{}

	if err := emp.GetEmployeeFromCache(name); err != nil {
		return emp, err
	}

	query := bson.M{"name": name}
	opts := options.FindOne()
	if err := GetCollection().FindOne(context.TODO(), query, opts).Decode(emp); err != nil {
		return emp, err
	}
	return emp, nil
}

//AddEmployee ...
func AddEmployee(emp *Employee) (interface{}, error) {
	var insertOneRes *mongo.InsertOneResult
	var err error
	if insertOneRes, err = GetCollection().InsertOne(context.TODO(), emp); err != nil {
		return 0, err
	}
	emp.UpdateEmployeeCache()
	return insertOneRes.InsertedID, nil
}
