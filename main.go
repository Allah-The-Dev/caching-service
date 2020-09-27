package main

import (
	"caching-service/config"
	"caching-service/handlers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"

	"github.com/go-openapi/runtime/middleware"

	"caching-service/data"
)

// Tweak configuration values here.
const (
	httpServerPort    = ":8090"
	readHeaderTimeout = 1 * time.Second
	writeTimeout      = 10 * time.Second
	idleTimeout       = 90 * time.Second
	maxHeaderBytes    = http.DefaultMaxHeaderBytes
)

var (
	cacheServiceDB = "cacheService"
)

func main() {

	//initialize app config
	config.InitializeAppConfig()

	//global mongo client
	mongoClient, err := data.InitializeMongoClient()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = mongoClient.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()
	data.MongoClient = mongoClient

	//global redis client pool
	data.RedisClientPool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", config.RedisURI)
		},
	}

	//http REST routes
	router := initializeHTTPRouter()

	// HTTP server configuration
	httpServer := &http.Server{
		Addr:              httpServerPort,
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}

	// start http server
	go func() {
		log.Printf("**************http server listening on port %s *************", httpServerPort)

		err := httpServer.ListenAndServe()
		if err != nil {
			log.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	//tap interrupt and kill signal and gracefully shutdown server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	//global kafka consumer
	go data.StartKafkaConsumer()

	//block until signal is received
	sig := <-c
	log.Println("Got os signal : ", sig)

	//gracefully shutdown server, waiting 30 second for shutting down server
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	httpServer.Shutdown(ctx)

	log.Println("Sutting down server")
	os.Exit(0)
}

func initializeHTTPRouter() *mux.Router {

	//employee handler
	empHandler := handlers.NewEmployee(config.EmpAPILogger)

	//gorilla mux router
	router := mux.NewRouter()

	//set router prefix
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	//get router
	getRouter := subRouter.Methods(http.MethodGet).Subrouter()

	//post router
	postRouter := subRouter.Methods(http.MethodPost).Subrouter()

	//employee get router
	getRouter.HandleFunc("/employee", empHandler.GetEmployees)
	getRouter.HandleFunc("/employee/{name:[a-zA-Z]+}", empHandler.GetEmployee)

	//employee post router
	postRouter.HandleFunc("/employee", empHandler.AddEmployee)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml", BasePath: "/api/v1"}
	sh := middleware.Redoc(opts, nil)
	getRouter.Handle("/docs", sh)
	router.Handle("/swagger.yaml", http.FileServer(http.Dir("./"))).Methods(http.MethodGet)

	return router
}
