package service

import (
	"encoding/json"
	"log"
	"test_task/rabbitMQ"
	"test_task/store"
)

func AddData(req store.Data) (err error) {

	data, err := json.Marshal(req)
	if err != nil {
		log.Println("error in marshal : ", err.Error())
		panic(err.Error())
	}

	err = rabbitMQ.Push(data)
	return
}

func Stop() (err error) {
	err = rabbitMQ.CloseRabbitMqConn()
	if err != nil {
		log.Println("error : ", err.Error())
	}
	err = store.CloseDbConn()
	if err != nil {
		log.Println("error : ", err.Error())
	}
	return
}
