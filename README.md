# caching-service
caching service in golang with microservices and REST, which is integrated with Kafka, Redis and MongoDB

# to run on docker swarm 
1. cd docker

2. Export env vars in deploy.sh

2. sh deploy.sh

# to deploy in minikube
1. cd k8s

2. sh deploy.sh

# Application is integrated with
1. Kafka for pub/sub model

2. Redis for caching

3. Mongo for data persistence

# On application startup swagger docs can be accessed at
1. http://<deployment ip>:port/api/v1/docs
  
# Employee CRUD APIs:
1. GET: api/v1/employee?pageNo=2&pageSize=3
2. GET: api/v1/employee
3. GET: api/v1/employee/{employee_name}
4. POST: api/v1/employee '{"name":"foo", "unit":"bar"}'

