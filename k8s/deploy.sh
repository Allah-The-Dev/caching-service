#create secrets
kubectl apply -f mongo-secrets.yaml

#create config map
kubectl apply -f configmap.yaml

#create mongo deployment and service
kubectl apply -f mongo.yaml

#create zookeeper deployment and service
kubectl apply -f zookeeper.yaml

#create kafka deployment and service
kubectl apply -f kafka.yaml

#create redis deployment and service
kubectl apply -f redis.yaml

#create caching-service deployment and service
kubectl apply -f caching-service.yaml

#get external ip for service
minikube service caching-service