package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Driver struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}

type Drivers struct {
	Drivers []Driver
}

func loadDrivers() []byte {
// função que retorna o arquivo do slice de byte

	jsonFile, err := os.Open("drivers.json")
	// declaro a veriavel do arquivo e de erro, se der falha ao abrir o arquivo, err fica com valor.
	if err != nil {// se o erro for diferente de nullo
		panic(err.Error())// eu crash, dou um panic equivalente ao exception
	}

	defer jsonFile.Close()// se não tiver erro, fecha o arquivo.

	data, err := ioutil.ReadAll(jsonFile)
	// ler o arquivo com read
	if err != nil {
		panic(err.Error())/// se deu ruim eu lanço uma exception/panic
	}
	// se não tiver erro eu retorno o conteudo do arquivo
	return data
}

func ListDrivers(w http.ResponseWriter, r *http.Request) {
	drivers := loadDrivers()// chama a loadDrivers que trabalha com arquivo JSON(raiz do microserviço), para não ter trabalho de configurar um db para isso.
	w.Write([]byte(drivers))// esse Write é equivalente os Response, aqui retorno o resultado dos motoristas.
}

func GetDriverById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	data := loadDrivers()

	var drivers Drivers
	json.Unmarshal(data, &drivers)

	for _, v := range drivers.Drivers {
		if v.Uuid == vars["id"] {
			driver, _ := json.Marshal(v)
			w.Write([]byte(driver))
		}
	}
}

func main() {// Inicia aqui
	r := mux.NewRouter()// Esse pacote mux é para levantar um server http e o newRouter cria um sistema de roteamento
	r.HandleFunc("/drivers", ListDrivers)// quando bater no endpoint /drivers vai chamar a função  listdrivers para listar os motoristas
	r.HandleFunc("/drivers/{id}", GetDriverById)// // quando bater nesse endpoint vai pegar um motorista especifico
	http.ListenAndServe(":8081", r)// http://localhost:8081/drivers vai ser onde vai rodar nossa app
	// go run driver.go
	// go build driver.go para gerar o arquivo EXECUTAVEL do go, se quiser rodar esse comando que gerar o arquivo "driver" e executa com "./driver" RODA EM QUALQUER SO.
	// para windows precisa de adicionar uma variaveis de ambiente: "GOOS=windows go build driver.go"
	// bom porque na hora de fazer deploy é só 1 unico arquivo!
	// precisa do go instalado na sua maquina.
}
