import {AmqpConnection, RabbitSubscribe} from "@golevelup/nestjs-rabbitmq";
import {WebSocketGateway, WebSocketServer} from "@nestjs/websockets";
import {Server} from "socket.io";
import {InjectRepository} from "@nestjs/typeorm";
import {Order, OrderStatus} from "../order.model";
import {Repository} from "typeorm";

// Inicia a biblioteca do socket.io ja pronto, graças ao decorators nest, eu não configuro nada, só uso!
@WebSocketGateway() //Socket.io
export class MappingService {
    @WebSocketServer() server: Server;// Graças ao decorator @WebSocketGateway(), ele joga a instancia do server socket.io

    constructor(
        @InjectRepository(Order)
        private readonly orderRepo: Repository<Order>,
        private amqpConnection: AmqpConnection,
    ) {

    }


    @RabbitSubscribe({
        exchange: 'amq.direct',
        routingKey: 'mapping.new-position',
        queue: 'micro-mapping/new-position'
    })
    public async rpcHandler(message) { //lat, lng, order
        const lat = parseFloat(message.lat); // converto só por questões de segurança pois recebemos como string
        const lng = parseFloat(message.lng);
        this.server.emit(`order.${message.order}.new-position`,{lat, lng});// usamos o server do socket.io para emitir essa mensagem para o canal que estou criando com id da ordem
        if (lat === 0 && lng === 0) { // o simulador quando chega ao destino envia o sinal que é 0,0 e se for o caso cai aqui nesse if
            // pego o pedido, mudo o status para entregue e salvo.
            const order = await this.orderRepo.findOne(message.order);
            order.status = OrderStatus.DONE;
            await this.orderRepo.save(order);
            // depois disso eu publico essa mensagem para exchange direct no rabbitmq
            // por isso que o serviço de pedido / orders fica ouvindo essa fila de mudança de status, para mudar o status lá também!
            await this.amqpConnection.publish(
                'amq.direct',
                'orders.change-status',
                {
                    id: order.id,
                    status: OrderStatus.DONE
                }
            )
        }
    }
}
