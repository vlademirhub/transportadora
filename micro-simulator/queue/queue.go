package queue

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
)

// Crio uma função connect que vai retornar um canal do RabbitMQ pois o rabbitmq precisa de um CHANNEL para estabelecer a conexão
func Connect() *amqp.Channel {
// Montando a dsn a partir das variaveis de ambiente
// ATENÇÂO. RABBITMQ_DEFAULT_VHOST é utilizado para CONTEXTOS diferentes na nosssa aplicação. Ex.: Um app A usa o VHOST A, um APP B usa o VHOST B
	dsn := "amqp://" + os.Getenv("RABBITMQ_DEFAULT_USER") + ":" + os.Getenv("RABBITMQ_DEFAULT_PASS") + "@" + os.Getenv("RABBITMQ_DEFAULT_HOST") + ":" + os.Getenv("RABBITMQ_DEFAULT_PORT") + os.Getenv("RABBITMQ_DEFAULT_VHOST")
	conn, err := amqp.Dial(dsn) // Dial retonar o resultado da tentativa de obter a conexão. https://godoc.org/github.com/streadway/amqp
	failOnError(err, "Failed to connect to RabbitMQ")
	//defer conn.Close()
	// se chegou aqui é porque se conectou corretamente

    // Baseado na conexão realizada com sucesso, eu estabeleço um canal.
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	//defer ch.Close()
	return ch
}

// Aqui é onde o RabbitMQ vai começar a escutar tudo de uma fila
func StartConsuming(in chan []byte, ch *amqp.Channel) {
// Crio um canal de slice de byte, tipo um array com tamanho infinito e passo a conexão com rabbitmq no ultimo params

    // Aqui eu digo que eu vou conectar numa fila, declaro ela para se conectar
	q, err := ch.QueueDeclare(
		os.Getenv("RABBITMQ_CONSUMER_QUEUE"), // name nome da fila
		true,                                 // durable se matar o processo do rabbitmq ela ainda fica, persistente
		false,                                // delete when usused para não deletar automatica quando fica na espera
		false,                                // exclusive
		false,                                // no-wait
		nil,                                  // arguments
	)
	failOnError(err, "Failed to declare a queue")// se der erro!

    // Baseado na declaração da fila que quero consumir, eu começo a consumir!
	msgs, err := ch.Consume(
		q.Name,      // queue
		"go-worker", // consumer
		true,        // auto-ack para garantir que recebemos a mensagem, aqui quando ele receber a mensagem e já tira ela da fila!
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	failOnError(err, "Failed to register a consumer")
	// Quando ele começar a consumir, as mensagens vão para variavel msgs

    // aqui a mensagens são processadas nessa função anonima
    // go func é uma GREEN THREAD ou seja GO ROUTING por que eu coloquei o GO na frente ele vira uma função ASINCRONA!
	go func() {
		for d := range msgs {
			in <- []byte(d.Body)// o conteudo da mensagem eu jogo para o meu canal IN passado como parametro StartConsuming
		}
		close(in)// quando termina o loop eu fecho o canal porque nao será mais utilizado
	}()
}

func Notify(payload string, ch *amqp.Channel) {

	err := ch.Publish(
		os.Getenv("RABBITMQ_DESTINATION_POSITIONS_EX"), // exchange
		os.Getenv("RABBITMQ_DESTINATION_ROUTING_KEY"),                 // routing key
		false,                             // mandatory
		false,                             // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(payload),// payload é o objeto json que ta vindo do simulator go linha 82
		})
	failOnError(err, "Failed to publish a message")
	fmt.Println("Message sent: ", payload)

}

func failOnError(err error, msg string) {// se der erro  na declaração
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
