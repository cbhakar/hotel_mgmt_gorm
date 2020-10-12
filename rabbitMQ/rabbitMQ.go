package rabbitMQ

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"test_task/store"
)

var (
	url       = "amqp://ssvvtrpr:aW-jF853zRnQ8LHajrYfvAwi04bknIZn@lionfish.rmq.cloudamqp.com/ssvvtrpr"
	queueName = "test"
	exchange  = "hello"
)

var Ch *amqp.Channel
var Conn *amqp.Connection

func InitiateRabbitMQConn() {
	connection, err := amqp.Dial(url)
	if err != nil {
		panic("could not establish connection with RabbitMQ:" + err.Error())
	}
	Conn = connection

	channel, err := connection.Channel()
	if err != nil {
		panic("could not open RabbitMQ channel:" + err.Error())
	}
	log.Println("rabbitmq channel created")
	Ch = channel

	_, err = channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		panic("error declaring the queue: " + err.Error())
	}
	log.Println("queue created with name : ", queueName)

	Msgs, err := channel.Consume("test", "", false, false, false, false, nil)
	if err != nil {
		panic("error consuming the queue: " + err.Error())
	}

	go CheckMsg(Msgs)
}

func init() {
	InitiateRabbitMQConn()
}

func CheckMsg(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		log.Println("task received on queue")
		var task = store.Data{}
		err := json.Unmarshal(msg.Body, &task)
		if err != nil {
			log.Println("error : ", err.Error())
		}

		err = store.AddNewOffer(task)
		if err != nil {
			log.Println("error storing offers : ", err.Error())
		}

		err = msg.Ack(false)
		if err != nil {
			log.Println("error acknowledging rabbitmq msg : ", err.Error())
		}
		log.Println("task processed")
	}
}

func Push(data []byte) (err error) {
	err = Ch.Publish("", queueName, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         data,
	})
	return
}

func CloseRabbitMqConn() (err error) {
	err = Ch.Close()
	if err != nil {
		log.Println(err.Error())
	}
	err = Conn.Close()
	if err != nil {
		log.Println(err.Error())
	}
	return
}
