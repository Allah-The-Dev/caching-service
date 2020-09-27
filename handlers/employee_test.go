package handlers

import (
	"bytes"
	"caching-service/config"
	"caching-service/data"
	"encoding/json"
	"log"
	"net/http/httptest"
	"os"
	"testing"
)

var testLogger = log.New(os.Stdout, "emp : ", log.LstdFlags)

func TestMain(m *testing.M) {
	//set env vars
	os.Setenv("MONGODB_SERVER", "")
	os.Setenv("MONGODB_ADMINUSERNAME", "")
	os.Setenv("MONGODB_ADMINPASSWORD", "")
	os.Setenv("REDIS_SERVER", "")
	os.Setenv("REDIS_PORT", "6379")
	os.Setenv("KAFKA_SERVER", "")

	config.InitializeAppConfig()

	if err := data.InitializeMongoClient(); err != nil {
		panic(err)
	}

	data.InitializeRedisClientPool()

	os.Exit(m.Run())
}

func BenchmarkGetEmployees(b *testing.B) {

	r := httptest.NewRequest("GET", "/api/v1/employee", nil)

	emp := &Employee{testLogger}

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		emp.GetEmployees(w, r)
	}
}

func BenchmarkGetEmployee(b *testing.B) {

	r := httptest.NewRequest("GET", "/api/v1/employee/foo", nil)

	emp := &Employee{testLogger}

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		emp.GetEmployee(w, r)
	}
}

func BenchmarkPostEmployee(b *testing.B) {

	reqBody := map[string]string{"name": "foo", "unit": "bar"}

	reqBodyJSON, _ := json.Marshal(reqBody)

	r := httptest.NewRequest("POST", "/api/v1/employee", bytes.NewBuffer(reqBodyJSON))

	emp := &Employee{testLogger}

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		emp.AddEmployee(w, r)
	}
}
