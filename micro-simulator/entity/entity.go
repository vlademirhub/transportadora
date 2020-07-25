package entity

// Struct = Class só que do Golang
// Objeto Order, que vai ter Uuid e Destination do tipo STRING  que serão transformados em order e destination quando virar json.
type Order struct {
	Uuid        string `json:"order"`
	Destination string `json:"destination"`
}

// Essa classe aqui a gente vai enviar do nossso simulador para uma exchange no rabbitmq
type Destination struct {
	Order string `json:"order"`
	Lat   string `json:"lat"`
	Lng   string `json:"lng"`
}
