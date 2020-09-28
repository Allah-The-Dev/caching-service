export KAFKA_HOSTNAME_EXT=
export KAFKA_PORT_EXT=9094
#mongo db server and port like localhost:27017
export MONGODB_SERVER=
export REDIS_SERVER=
export REDIS_PORT=6379
export KAFKA_SERVER=

docker stack deploy -c stack.yaml caching-service