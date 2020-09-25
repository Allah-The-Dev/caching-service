package data

import (
	"errors"

	"github.com/gomodule/redigo/redis"
)

//RedisClientPool ...
var RedisClientPool *redis.Pool

//UpdateEmployeeCache ...
func (emp *Employee) UpdateEmployeeCache() {

	CLogger.Printf("cache updated for emp %s\n", emp.Name)
	conn := RedisClientPool.Get()
	if _, err := conn.Do("HMSET", redis.Args{}.Add(emp.Name).AddFlat(emp)...); err != nil {
		CLogger.Println(err)
	}
}

//GetEmployeeFromCache ...
func (emp *Employee) GetEmployeeFromCache(name string) error {

	conn := RedisClientPool.Get()

	//check if this name exist in cache
	exists, err := redis.Int(conn.Do("EXISTS", name))
	if err != nil {
		CLogger.Println(err)
		return err
	} else if exists == 0 {
		return errors.New("Data for this emp do not exist")
	}
	CLogger.Printf("cache hit successful for %s.\n", name)

	v, err := redis.Values(conn.Do("HGETALL", name))
	if err != nil {
		CLogger.Println(err)
		return err
	}

	if err := redis.ScanStruct(v, emp); err != nil {
		CLogger.Println(err)
		return err
	}
	CLogger.Printf("cache hit successful for %s. Value is %v\n", name, *emp)
	return nil
}
