package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"simulator/entity"
	"simulator/queue"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

var active []string

func init() {// Executa antes do main, é nativa do go.
	err := godotenv.Load()// pesso para ler as variaveis de ambiente.
	if err != nil {
		panic("Error loading .env file")
	}
}

func main() {

	in := make(chan []byte)// uso um canal com rabbitmq de comunicação

	ch := queue.Connect()// Aqui é a conexão com rabbitmq
	queue.StartConsuming(in, ch)// essa funçao está no arquivo queue.go

	for msg := range in {// toda vez que ele consumir mensagens, elas vão para o canal IN
		var order entity.Order
		// Gerei uma nova ordem
		err := json.Unmarshal(msg, &order) // aqui eu to falando para colocar o json recebido da fila nessa mensagem dentro do objeto order tipo json.decode

		if err != nil {// se tiver erro lança
			fmt.Println(err.Error())
		}

		fmt.Println("New order Received: ", order.Uuid)// sucesso!

        // inicia o jogo
		start(order, ch)// aqui eu dou o start, pego a ordem, passo no start junto com a conecão do rabbitmq
	}
}

// esse start executa uma função SimulatorWorker passando  o order e a conexão.
// é a parte mais importante do programa, quando passamos por exemplo a order com ID = 1 ele passa ela para o simulador rodar!
// essa função será uma nova thread do go e ficará rodando em background e automaticamente ele cria 1 worker para cada order em paralelo tudo graças ao GO
// Esse SimulatorWorker é o cara que fica atualizando a posição do veiculo!
func start(order entity.Order, ch *amqp.Channel) {

// se o cara der um START e a ORDEM ja estiver em EXECUÇÂO ou FO ENTREGUE OU ESTA EM PROCESSAMENTO
// para garantir que ele não vai iniciar um worker se ele não tiver certeza que a entrega foi realizada
	if !stringInSlice(order.Uuid, active) {
		active = append(active, order.Uuid)
		go SimulatorWorker(order, ch)
	} else {
		fmt.Println("Order", order.Uuid, "was already completed or is on going...")
	}
}

func SimulatorWorker(order entity.Order, ch *amqp.Channel) {

	f, err := os.Open("destinations/" + order.Destination + ".txt")
	// Abre o arquivo de destino com as coordenadas do destino de entrega.

	if err != nil {
		panic(err.Error())
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)// lendo arquivo aberto

	for scanner.Scan() {// faço um FOR linha a linha do arquivo
		data := strings.Split(scanner.Text(), ",") // para separar virgulas para array
		json := destinationToJson(order, data[0], data[1]) // crio o json baseado no array

		time.Sleep(1 * time.Second) // Espero 1 segundo, posso colocar o tempo que eu quiser
		queue.Notify(string(json), ch)// Mando para a fila a mensagem com as posições vinda do json, com a conexão do rabbitmq
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	json := destinationToJson(order, "0", "0")
	queue.Notify(string(json), ch)// quando acaba o loop eu digo que acabou o loop ou seja, finalizou a entrega, eu mando 0 0 de coordenada.
}

func destinationToJson(order entity.Order, lat string, lng string) []byte {
	dest := entity.Destination{ // objeto vai virar json
		Order: order.Uuid,
		Lat:   lat,
		Lng:   lng,
	}
	json, _ := json.Marshal(dest)// Marshal faz o destino virar json
	return json
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
