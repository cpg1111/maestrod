package datastore

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	redis "gopkg.in/redis.v3"
)

type RedisStore struct {
	Datastore
	store *redis.Client
}

type redisData struct {
	Data interface{}
}

func NewRedis(host, port, password string) *RedisStore {
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
	}
	return &RedisStore{
		store: redis.NewClient(options),
	}
}

func (r RedisStore) Save(key string, data interface{}, callback NoResultCallback) {
	go func() {
		rData := &redisData{
			Data: data,
		}
		dataStr, marshalErr := json.Marshal(rData)
		if marshalErr != nil {
			callback(marshalErr)
			return
		}
		cmd := r.store.Set(key, dataStr, 0)
		res, resErr := cmd.Result()
		if resErr != nil {
			callback(resErr)
			return
		}
		log.Println("redis message: ", res)
		callback(nil)
	}()
}

func (r RedisStore) Find(queryStr string, callback ResultCallback) {
	go func() {
		cmd := r.store.Get(queryStr)
		callback(cmd.Result())
	}()
}

func (r RedisStore) Remove(queryStr string, callback NoResultCallback) {
	go func() {
		cmd := r.store.Set(queryStr, nil, 0)
		res, resErr := cmd.Result()
		if resErr != nil {
			callback(resErr)
			return
		}
		log.Println(res)
		callback(nil)
	}()
}

func (r RedisStore) Update(queryStr string, update interface{}, callback NoResultCallback) {
	resChan := make(chan string)
	errChan := make(chan error)
	go func() {
		gCMD := r.store.Get(queryStr)
		gRes, gResErr := gCMD.Result()
		if gResErr != nil {
			errChan <- gResErr
			return
		}
		resChan <- gRes
	}()
	go func() {
		for {
			select {
			case errMsg := <-errChan:
				if errMsg != nil {
					log.Fatal(errMsg)
					return
				}
			case resMsg := <-resChan:
				if resMsg == "" {
					callback(errors.New("no object found with that key"))
					return
				}
				rData := &redisData{
					Data: update,
				}
				dataStr, marshalErr := json.Marshal(rData)
				if marshalErr != nil {
					callback(marshalErr)
					return
				}
				sCMD := r.store.Set(queryStr, dataStr, 0)
				res, resErr := sCMD.Result()
				if resErr != nil {
					callback(resErr)
					return
				}
				log.Println(res)
				callback(nil)
			}
		}
	}()
}

func (r RedisStore) FindAndUpdate(queryStr string, update interface{}, callback ResultCallback) {
	r.Update(queryStr, update, func(err error) {
		if err != nil {
			callback(nil, err)
			return
		}
		r.Find(queryStr, callback)
	})
}
