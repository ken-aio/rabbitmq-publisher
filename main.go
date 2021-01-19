package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/NeowayLabs/wabbit"
	"github.com/NeowayLabs/wabbit/amqp"
	"github.com/pkg/errors"
)

var (
	uri        string
	exchange   string
	routingKey string
	filePath   string
	conn       *amqp.Conn
)

func validateFlags() error {
	if uri == "" {
		return errors.New("uri cannot be blank")
	}

	if exchange == "" && routingKey == "" {
		return errors.New("exchange and routing-key cannot both be blank")
	}

	if filePath == "" {
		return errors.New("file-path cannot be blank")
	}

	return nil
}

func init() {
	flag.StringVar(&uri, "uri", "", "AMQP URI amqp://<user>:<password>@<host>:<port>/[vhost]")
	flag.StringVar(&exchange, "exchange", "", "Exchange name")
	flag.StringVar(&routingKey, "routing-key", "", `Routing key. Use queue name with blank exchange to publish directly to queue.`)
	flag.StringVar(&filePath, "file-path", "", "message file path")

	flag.Parse()

	if err := validateFlags(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := createConnection(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	messages, err := readFile()
	if err != nil {
		panic(err)
	}

	for _, m := range messages {
		if err := publish(ch, m); err != nil {
			log.Println("publish error: ", m, err)
		}
	}
}

func createConnection() error {
	var err error
	conn, err = amqp.Dial(uri)
	if err != nil {
		return errors.Wrapf(err, "connect to RabbitMQ error. uri: %s", uri)
	}

	return nil
}

func readFile() ([]string, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "error read file")
	}
	return strings.Split(string(bytes), "\n"), nil
}

func publish(ch wabbit.Channel, message string) error {
	if message == "" {
		return nil
	}
	err := ch.Publish(exchange, routingKey, []byte(message), wabbit.Option{
		"contentType":  "application/json",
		"deliveryMode": 2, // persistent
	})
	if err != nil {
		return err
	}
	logMsg := message
	if len([]rune(message)) > 100 {
		logMsg = string([]rune(message)[0:99])
	}
	log.Println("success to publish: ", logMsg)
	return nil
}
