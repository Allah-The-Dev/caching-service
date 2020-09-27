package data

import (
	"caching-service/config"
	"context"
	"encoding/json"
	"io"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
func GetEmployees(pageNo, pageSize int) (Employees, error) {

	empList := make([]*Employee, 0)

	//options
	options := options.FindOptions{}
	if pageNo != 0 && pageSize != 0 {

		skips := int64(pageSize * (pageNo - 1))
		int64PageSize := int64(pageSize)

		options.SetSkip(skips)
		options.SetLimit(int64PageSize)
	}

	cur, err := getCollection().Find(context.TODO(), bson.D{{}}, &options)
	if err != nil {
		config.EmpAPILogger.Println(err)
		return empList, err
	}

	for cur.Next(context.TODO()) {
		var emp Employee
		if err := cur.Decode(&emp); err != nil {
			config.EmpAPILogger.Println(err)
			return empList, err
		}
		empList = append(empList, &emp)
	}
	return empList, nil
}

//GetEmployee ...
func GetEmployee(name string) (*Employee, error) {

	emp := &Employee{}

	if err := emp.GetEmployeeFromCache(name); err == nil {
		return emp, nil
	}

	query := bson.M{"name": name}
	opts := options.FindOne()
	if err := getCollection().FindOne(context.TODO(), query, opts).Decode(emp); err != nil {
		config.EmpAPILogger.Println(err)
		return emp, err
	}
	return emp, nil
}

//AddEmployee ...
func AddEmployee(emp *Employee) (interface{}, error) {
	var insertOneRes *mongo.InsertOneResult
	var err error
	if insertOneRes, err = getCollection().InsertOne(context.TODO(), emp); err != nil {
		config.EmpAPILogger.Println(err)
		return 0, err
	}
	defer emp.PublishToKafka()
	return insertOneRes.InsertedID, nil
}
