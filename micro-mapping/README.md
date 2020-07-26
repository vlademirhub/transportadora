## Descrição
Sistema responsável pelo rastreamento dos motoristas via geolocalização.
Todo pedido poderá ser monitorado via mapa mostrando o ponto de partida do motorista, o destino e todo trajeto feito pelo motorista.

Casos de comunicação

1º Filas

2º Websockets - Client <-> Server

- Esse microsserviço vai ter 2 vertentes:
    * O backend vai ficar consumindo o tempo todo a localização do entregador, através de um fila.
        * Quando receber a informação, ele vai abrir um canal de comunicação WEBSOCKET e REPLICARÁ os dados recebidos nesse canal.
    * No mesmo tempo que o backend replica a informação no canal a gente vai utilizar o frontend
    com React.js para acessar esse canal, pegar as informação e exibir em tempo real no mapa a localização do entregador.
     

Microsserviço de mapeamento de entregas construído com Nest.js Framework + Socket.io +React.js + RabbitMQ

## Tecnologias
- Google Maps
- Node.js
- Nest.js - Framework Backend em Node.js - (SEM HANDLEBLARS)
- React (Create React App) - Frontend (com SPA e React Router DOM)
- Socket.io - Para comunicação Websocket client e server
- Notistack - SnackbarProvider é que permite que trabalhamos com as nossas mensagens.
- Axios
- MySQL

## Rodar a aplicação

#### Antes de começar

O microsserviço de mapeamento necessita que os microsserviços de Drivers, Simulador e Orders já estejam rodando antes de inicia-lo. 

Microsserviço Drivers

Microsserviço Simulador

Microsserviço Order


#### Rodar o RabbitMQ
Entre em micro-tabbitmq

Rode ```docker-compose up```. 

#### Crie o .env e configure as variáveis de ambiente do projeto frontend

```bash
$ cd frontend
$ cp .env.example .env
```

#### Crie os containers com Docker

```bash
$ docker-compose up
```

#### Accesse no browser

```
http://localhost:3001 para a API Rest do microserviço
http://localhost:3002 para o front-end que contém a interface de mapeamento da entrega
```
