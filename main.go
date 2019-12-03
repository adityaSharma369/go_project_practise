package main

import (
	fn "tericai/common/helpers"
	"fmt"
	"os"

	"github.com/globalsign/mgo"
	"github.com/joho/godotenv"
)

var mgoSession *mgo.Session

var collections = make(map[string]*mgo.Collection)

func main() {

	err := godotenv.Load()
	if err != nil {
		fn.Log("Error loading .env file", 0)
	}

	userSession, err := mgo.Dial("mongodb://" + os.Getenv("PROJECT_DB_HOST"))

	if err != nil {

		fmt.Println(err)

	}

	collections["Project"] = userSession.DB("ms_project").C("project")

	fmt.Println("Connected to MongoDB!")

	//rabbitMQ initialize
	exchange := os.Getenv("RABBIT_MQ_EXCHANGE_NAME")
	pattern := os.Getenv("RABBIT_MQ_PATTERN")
	queueName := os.Getenv("RABBIT_MQ_QUEUE_NAME")
	rabbitConn := "amqp://guest:guest@192.168.0.32:5673"
	fn.InitMsgBroker(rabbitConn)
	patternRequests, err := fn.NewQueue(exchange, pattern)
	fn.GetReplies(exchange, queueName)

	if err != nil {
		fmt.Println("Log -- cannot start Microservice", err)
		os.Exit(3)
	}

	fmt.Println("------------------------------------------------ \n Database + RabbitMQ initialized : ", "\n ------------------------------------------------ ")

	forever := make(chan bool)

	// listen to requests -- patterns
	go func() {
		for req := range patternRequests {
			go callMethod(req.Module+"_"+req.Action, req)
		}
	}()

	<-forever

}
