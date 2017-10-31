package frontend

import (
	"fmt"

	"github.com/micro/go-micro/transport"
	"github.com/micro/go-plugins/transport/rabbitmq"
)

// var amqpTransport transport.Transport
var amqpClient transport.Client

const amqpAddr = "amqp://localhost"

// func subscribe() {
// 	litener, err := transport.Listen(amqpAddr)
//
// 	fmt.Println("Listener:", litener)
//
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

// Run Frontend service
func Run() {
	fmt.Println("Frontend service is running")

	amqpTransport := rabbitmq.NewTransport()

	fmt.Println("Transport", amqpTransport)

	var err error

	amqpClient, err = amqpTransport.Dial(amqpAddr)

	if err != nil {
		fmt.Println("Transport Dial error:", err)
	}

	// if err := amqpTransport.Init(); err != nil {
	// 	log.Fatalf("Broker Init error: %v", err)
	// }
	//
	// if err := amqpBroker.Connect(); err != nil {
	// 	log.Fatalf("Broker Connect error: %v", err)
	// }

	// subscribe()

	pub()
}

func pub() {
	msg := &transport.Message{
		Header: map[string]string{
			// "Index": fmt.Sprintf("%d", postIndex),
			"Index": "12345",
		},
		Body: []byte(`test`),
	}
	err := amqpClient.Send(msg)

	if err != nil {
		fmt.Println("Tranport Send error:", err)
	} else {
		fmt.Println("Tranport Send")
	}
}
