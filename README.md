# caching-service
caching service in golang with microservices and REST

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
