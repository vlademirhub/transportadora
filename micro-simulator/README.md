# Go Routines

- Channels
- Queues
- Threads

_____

Arquivos:
- destinations/
    - Arquivo onde tem a latitude e longitude para todo o caminho que o carro vai andar, esse é o nosso simulador da rota!
    - Por padrão coloquei 1.txt para quando dizer que o destino é 1 ele abra o 1.txt se dizer destino CASA ele abre casa1.txt
- entity/
    - Entidades.
    - Ele que vai mandar o destino para a nossa exchange do RabbitMQ
- queue/
    - Queue
    - Esse é o arquivo que vai se comunicar com nosso rabbitmq
- simulador
    - To ordem processada, para saber se tem entrega em andamento e etc...


Deve ter queue positions no rabbitmq:

- go run simulator.go