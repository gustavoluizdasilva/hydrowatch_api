package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hydrowatch-api/src/config"
	"net/http"
)

func GetConnection() (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI("mongodb+srv://" + config.USER + ":" + config.PASS +
			"@cluster0.deee8be.mongodb.net/?retryWrites=true&w=majority").
		SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)

	return client, err
}

func CloseConnection(client *mongo.Client, ctx context.Context, w http.ResponseWriter) {
	err := client.Disconnect(ctx)
	if err != nil {
		http.Error(w, "Erro ao desconectar o banco de dados", http.StatusBadRequest)
		return
	}
}
