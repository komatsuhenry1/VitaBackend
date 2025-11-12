package main

import (
	"fmt"
	"medassist/config"
	"medassist/router"
	"net"
	"os"
)

// @title Sua API Incrível
// @version 1.0
// ...
// @host localhost:8081
// @BasePath /api/v1
//
// 💡 ADICIONE ISSO PARA DEFINIR O HEADER DE AUTORIZAÇÃO
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	// Inicializa o banco de dados
	fmt.Println("Iniciando o servidor...")
	config.ConnectDatabase()
	fmt.Println("Banco de dados conectado com sucesso!")
	fmt.Printf("IP local: %s\n", getLocalIPv4())
	// Inicializar o router
	r := router.InitializeRoutes()

	// Inicializar o servidor
	r.Run(":" + os.Getenv("SERVER_PORT"))
}

func getLocalIPv4() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "desconhecido"
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ip := ipnet.IP.To4(); ip != nil {
				return ip.String()
			}
		}
	}
	return "desconhecido"
}
