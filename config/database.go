package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	Client  *mongo.Client
	MongoDB *mongo.Database
)

func ConnectDatabase() {
    // 💡 MUDANÇA AQUI
	// Tenta carregar o .env, mas não falha se não existir
	err := godotenv.Load()
	if err != nil {
		// Em um ambiente de contêiner (Docker), o .env não existirá,
		// mas as variáveis serão injetadas pelo docker-compose.
		// Por isso, apenas logamos um aviso e continuamos.
		log.Println("Aviso: Arquivo .env não encontrado. Usando variáveis de ambiente do sistema.")
	}

    // A partir daqui, o código depende das variáveis de ambiente (seja do .env ou do Docker)
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI não está definido no ambiente")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	Client, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	if err := Client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	dbName := os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		log.Fatal("MONGO_DB_NAME não está definido no ambiente")
	}

	MongoDB = Client.Database(dbName)

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}

func GetMongoDB() *mongo.Database {
	return MongoDB
}