package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/plate/random", func(ctx *gin.Context) {

		licensePlate := GenerateThaiLicensePlate()
		log.Println("ป้ายทะเบียน -> " + licensePlate)

		err := ConnectToRabbitMQ(licensePlate)
		if err != nil {
			log.Printf("Error connecting to RabbitMQ: %v\n", err)
		}

		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/mqtt", func(ctx *gin.Context) {
		licensePlate := GenerateThaiLicensePlate()
		log.Println("ป้ายทะเบียน -> " + licensePlate)

		err := sentJsonMqtt(licensePlate)
		if err != nil {
			log.Printf("Error connecting to RabbitMQ: %v\n", err)
		}
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

func ConnectToRabbitMQ(message string) error {

	brokerURL := "amqp://" + "testoffice" + ":" + "sm2O0itJGrwP2NBz" + "@" + "mqtt.letmein.asia/"
	queueName := "testoffice1"
	conn, err := amqp.Dial(brokerURL)
	if err != nil {
		return err
	}
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			return
		}
	}(ch)

	_, err = ch.QueueDeclare(
		queueName, // Queue name
		true,      // Durable
		false,     // Delete when unused
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",        // Exchange
		queueName, // Routing key
		false,     // Mandatory
		false,     // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return err
	}

	log.Printf("Connected to RabbitMQ\n")
	log.Printf("Sent Message to RabbitMQ: %s %s\n", message, time.Now().Format(time.RFC3339))

	return nil
}

type Message struct {
	Data   Data `json:"data"`
	Status int  `json:"status"`
}

type Data struct {
	LicensePlate string `json:"license_plate"`
}

func sentJsonMqtt(plate string) error {
	brokerURL := "amqp://" + "testoffice" + ":" + "sm2O0itJGrwP2NBz" + "@" + "mqtt.letmein.asia/"
	queueName := "testoffice1"
	conn, err := amqp.Dial(brokerURL)
	if err != nil {
		return err
	}
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {

		}
	}(ch)

	message := Message{
		Data: Data{
			LicensePlate: plate,
		},
		Status: 1,
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}

	err = ch.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonMessage,
		},
	)
	if err != nil {
		return err
	}

	log.Printf("Connected to RabbitMQ\n")
	log.Printf("Sent Message to RabbitMQ: %s %s\n", message.Data.LicensePlate, time.Now().Format(time.RFC3339))

	return err
}

func GenerateThaiLicensePlate() string {
	thaiChars := []rune{
		'ก', 'ข', 'ค', 'ฆ', 'ง', 'จ', 'ฉ', 'ช', 'ฌ', 'ญ', 'ฎ', 'ฏ', 'ฐ', 'ฑ', 'ฒ', 'ณ',
		'ด', 'ต', 'ถ', 'ท', 'ธ', 'น', 'บ', 'ป', 'ผ', 'ฝ', 'พ', 'ฟ', 'ภ', 'ม', 'ย', 'ร',
		'ล', 'ว', 'ศ', 'ษ', 'ส', 'ห', 'ฬ', 'อ', 'ฮ',
	}

	rand.NewSource(time.Now().UnixNano())

	numericPart := fmt.Sprintf("%04d", rand.Intn(10000))

	prefix := ""
	for i := 0; i < 2; i++ {
		prefix += fmt.Sprintf("%c", thaiChars[rand.Intn(len(thaiChars))])
	}

	var numericPrefix string
	if rand.Float64() < 0.5 {
		numericPrefix = ""
	} else {
		numericPrefix = fmt.Sprintf("%d", rand.Intn(10))
	}

	licensePlate := numericPrefix + prefix + numericPart

	return licensePlate
}
