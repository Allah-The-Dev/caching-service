package data

import (
	"caching-service/config"
	"errors"
	"time"

	"github.com/gomodule/redigo/redis"
)

//RedisClientPool ...
var RedisClientPool *redis.Pool

//InitializeRedisClientPool ...
func InitializeRedisClientPool() {
	RedisClientPool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", config.RedisURI)
		},
	}
}

//UpdateEmployeeCache ...
func (emp *Employee) UpdateEmployeeCache() {

	config.EmpAPILogger.Printf("cache updated for emp %s\n", emp.Name)
	conn := RedisClientPool.Get()
	if _, err := conn.Do("HMSET", redis.Args{}.Add(emp.Name).AddFlat(emp)...); err != nil {
		config.EmpAPILogger.Println(err)
	}
}

//GetEmployeeFromCache ...
func (emp *Employee) GetEmployeeFromCache(name string) error {

	conn := RedisClientPool.Get()

	//check if this name exist in cache
	exists, err := redis.Int(conn.Do("EXISTS", name))
	if err != nil {
		config.EmpAPILogger.Println(err)
		return err
	} else if exists == 0 {
		return errors.New("Data for this emp do not exist")
	}
	config.EmpAPILogger.Printf("cache hit successful for %s.\n", name)

	v, err := redis.Values(conn.Do("HGETALL", name))
	if err != nil {
		config.EmpAPILogger.Println(err)
		return err
	}

	if err := redis.ScanStruct(v, emp); err != nil {
		config.EmpAPILogger.Println(err)
		return err
	}
	config.EmpAPILogger.Printf("cache hit successful for %s. Value is %v\n", name, *emp)
	return nil
}
