import {
    Repository,
} from 'typeorm';
import {Order} from "../order.model";
import {RabbitSubscribe} from "@golevelup/nestjs-rabbitmq";
import {InjectRepository} from "@nestjs/typeorm";
import {Injectable} from "@nestjs/common";

// Este serviço fica escutando o rabbitmq
// Na hora que chegar o novo pedido, ele ja cadastra aqui no sistema
@Injectable()
export class NewOrderService {
    constructor(
        @InjectRepository(Order)
        private readonly orderRepo: Repository<Order>
    ) {

    }

    @RabbitSubscribe({
        exchange: 'amq.direct',
        routingKey: 'orders.new',
        queue: 'micro-mapping/orders-new'
    })
    public async rpcHandler(message) {
        const order = this.orderRepo.create({
            id: message.id,
            driver_name: message.driver_name,
            location_id: message.location_id,
            location_geo: message.location_geo,
        });
        await this.orderRepo.save(order);
    }
}

