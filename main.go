package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	config := Config{}
	initialiseEnv()
	populateConfig(&config)

	router := gin.Default()
	router.POST("/queue/:queueName", func(c *gin.Context) {
		msg, err := sendMessageToQueue(c, &config)

		if err != nil {
			c.String(http.StatusBadGateway, msg)
		}
		c.Status(http.StatusOK)
	})

	err := router.Run(":" + config.ProxyPort)
	if err != nil {
		os.Exit(1)
	}
}

type Config struct {
	ProxyPort        string
	RabbitMQHost     string
	RabbitMQPort     string
	RabbitMQUsername string
	RabbitMQPassword string
}

func handleFatalError(msg string, err error) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func initialiseEnv() {
	if os.Getenv("RABBITMQ_HOST") == "" {
		err := os.Setenv("RABBITMQ_HOST", "localhost")
		if err != nil {
			handleFatalError("Error setting RABBITMQ_HOST", err)
		}
	}

	if os.Getenv("RABBITMQ_PORT") == "" {
		err := os.Setenv("RABBITMQ_PORT", "5672")
		if err != nil {
			handleFatalError("Error setting RABBITMQ_PORT", err)
		}
	}

	if os.Getenv("RABBITMQ_USERNAME") == "" {
		err := os.Setenv("RABBITMQ_USERNAME", "guest")
		if err != nil {
			handleFatalError("Error setting RABBITMQ_USERNAME", err)
		}
	}

	if os.Getenv("RABBITMQ_PASSWORD") == "" {
		err := os.Setenv("RABBITMQ_PASSWORD", "guest")
		if err != nil {
			handleFatalError("Error setting RABBITMQ_PASSWORD", err)
		}
	}

	if os.Getenv("PROXY_PORT") == "" {
		err := os.Setenv("PROXY_PORT", "5555")
		if err != nil {
			handleFatalError("Error setting PROXY_PORT", err)
		}
	}
}

func populateConfig(config *Config) {
	config.ProxyPort = os.Getenv("PROXY_PORT")
	config.RabbitMQHost = os.Getenv("RABBITMQ_HOST")
	config.RabbitMQPort = os.Getenv("RABBITMQ_PORT")
	config.RabbitMQUsername = os.Getenv("RABBITMQ_USERNAME")
	config.RabbitMQPassword = os.Getenv("RABBITMQ_PASSWORD")
}

func formatError(err error, msg string) (string, error) {
	if err != nil {
		return fmt.Sprintf("%s: %s", msg, err), err
	}
	return "", nil
}

func parseDefaultBool(queryParam string, defaultBool bool) bool {
	value, err := strconv.ParseBool(queryParam)
	if err != nil {
		return defaultBool
	}
	return value
}

func sendMessageToQueue(c *gin.Context, config *Config) (string, error) {
	queueName := c.Param("queueName")
	durable := parseDefaultBool(c.Query("durable"), false)
	autoDelete := parseDefaultBool(c.Query("autoDelete"), false)
	exclusive := parseDefaultBool(c.Query("exclusive"), false)
	noWait := parseDefaultBool(c.Query("noWait"), false)

	connectionString := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		config.RabbitMQUsername,
		config.RabbitMQPassword,
		config.RabbitMQHost,
		config.RabbitMQPort,
	)
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		return formatError(err, "Failed to open a channel")
	}
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
			log.Printf("Error closing connection: %s", err)
		}
	}(conn)

	ch, err := conn.Channel()
	if err != nil {
		return formatError(err, "Failed to open a channel")
	}

	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			log.Printf("Error closing channel: %s", err)
		}
	}(ch)

	q, err := ch.QueueDeclare(
		queueName,
		durable,
		autoDelete,
		exclusive,
		noWait,
		nil,
	)
	if err != nil {
		return formatError(err, "Failed to open a channel")
	}

	bodyBuffer := make([]byte, c.Request.ContentLength)

	body, err := c.Request.Body.Read(bodyBuffer)
	bodyBuffer = bodyBuffer[:body]

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			Body: bodyBuffer,
		})
	if err != nil {
		return formatError(err, "Failed to open a channel")
	}
	log.Printf(" [x] Sent %x", bodyBuffer)

	return "", nil
}
