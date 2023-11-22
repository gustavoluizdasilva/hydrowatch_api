package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hydrowatch-api/src/config"
	"hydrowatch-api/src/db"
	"hydrowatch-api/src/models"
	"net/http"
	"strconv"
	"time"
)

func SaveDataGreenhouse(w http.ResponseWriter, r *http.Request) {
	var data models.GreenhouseData

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, "Erro ao decodificar o JSON", http.StatusBadRequest)
		return
	}
	idGreenhouse := data.IdGreenhouse
	temperature := data.Temperature
	humidity := data.Humidity
	Flowrate := data.Flowrate

	client, err := db.GetConnection()
	if err != nil {
		fmt.Fprintln(w, "Erro ao conectar ao banco de dados!")
		return
	}

	collection := client.Database(config.DBNAME).Collection("greenHouseData")
	defer db.CloseConnection(client, context.TODO(), w)

	_, err = collection.InsertOne(context.TODO(), bson.D{
		{Key: "idgreenhouse", Value: idGreenhouse},
		{Key: "temperature", Value: temperature},
		{Key: "humidity", Value: humidity},
		{Key: "flowrate", Value: Flowrate},
		{Key: "timestamp", Value: time.Now()},
	})

	if err != nil {
		http.Error(w, "Erro ap gravar dados 'greenHouseData.InsertOne'", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Dados recebidos com sucesso")
}

func SaveGreenHouseConfig(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST	")
	var data models.GreenHouseConfig

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, "Erro ao decodificar o JSON", http.StatusBadRequest)
		return
	}

	client, err := db.GetConnection()
	if err != nil {
		fmt.Fprintln(w, "Erro ao conectar ao banco de dados!")
		return
	}
	defer db.CloseConnection(client, context.TODO(), w)

	collection := client.Database(config.DBNAME).Collection("greenHouseConfig")

	temperature, humidity, flowrate := createModelsConfig(data)

	_, err = collection.InsertOne(context.TODO(), bson.D{
		{Key: "greenhouse", Value: data.Greenhouse},
		{Key: "cultivar", Value: data.Cultivar},
		{Key: "idsensor", Value: data.IdSensor},
		{Key: "temperature", Value: temperature},
		{Key: "humidity", Value: humidity},
		{Key: "flowrate", Value: flowrate},
	})

	if err != nil {
		http.Error(w, "Erro ao decodificar o JSON", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
}

func createModelsConfig(data models.GreenHouseConfig) (models.Temperature, models.Humidity, models.Flowrate) {
	temperature := models.Temperature{
		Min: data.Temperature.Min,
		Max: data.Temperature.Max,
	}

	humidity := models.Humidity{
		Min: data.Humidity.Min,
		Max: data.Humidity.Max,
	}

	flowrate := models.Flowrate{
		Min: data.Flowrate.Min,
		Max: data.Flowrate.Max,
	}
	return temperature, humidity, flowrate
}

func GetDataSensor(w http.ResponseWriter, r *http.Request) {
	//startDateStr := r.FormValue("start_date")
	//endDateStr := r.FormValue("end_date")
	idGreenhouse, _ := strconv.Atoi(r.FormValue("idgreenhouse"))

	client, err := db.GetConnection()
	if err != nil {
		fmt.Fprintln(w, "Erro ao conectar ao banco de dados!")
		return
	}
	defer db.CloseConnection(client, context.TODO(), w)

	collection := client.Database(config.DBNAME).Collection("greenHouseData")

	//startDate, err := time.Parse(time.DateOnly, startDateStr)
	//if err != nil {
	//	http.Error(w, "Erro ao converter start_date", http.StatusBadRequest)
	//	return
	//}
	//
	//endDate, err := time.Parse(time.DateOnly, endDateStr)
	//if err != nil {
	//	http.Error(w, "Erro ao converter end_date", http.StatusBadRequest)
	//	return
	//}

	pipeline := bson.A{
		bson.M{
			"$match": bson.M{
				"idgreenhouse": idGreenhouse,
				//"timestamp": bson.M{
				//	"$gte": startDate,
				//	"$lte": endDate,
				//},
			},
		},
		bson.M{
			"$project": bson.M{
				"yearMonthDay": bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$timestamp"}},
				"temperature":  1,
				"humidity":     1,
				"flowrate":     1,
			},
		},
		bson.M{
			"$group": bson.M{
				"_id":            "$yearMonthDay",
				"avgTemperature": bson.M{"$avg": "$temperature"},
				"avgHumidity":    bson.M{"$avg": "$humidity"},
				"avgFlowrate":    bson.M{"$avg": "$flowrate"},
			},
		},
		bson.M{"$sort": bson.M{"_id": 1}},
	}

	// Executar a agregação
	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		http.Error(w, "Erro na agregação", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	// Iterar sobre os resultados e enviar a resposta
	var result []bson.M
	if err := cursor.All(context.Background(), &result); err != nil {
		http.Error(w, "Erro ao iterar sobre os resultados", http.StatusInternalServerError)
		return
	}

	// Converter os resultados para JSON e enviar a resposta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func GetGreenHouseConfig(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	client, err := db.GetConnection()
	if err != nil {
		fmt.Fprintln(w, "Erro ao conectar ao banco de dados!")
		return
	}

	defer db.CloseConnection(client, context.TODO(), w)

	collection := client.Database(config.DBNAME).Collection("greenHouseConfig")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": objectID}

	var result models.GreenHouseConfig
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Documento não encontrado", http.StatusNotFound)
		return
	}

	fmt.Println(result)
	json.NewEncoder(w).Encode(result)
}

func DeleteGreenHouse(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	client, err := db.GetConnection()
	if err != nil {
		fmt.Fprintln(w, "Erro ao conectar ao banco de dados!")
		return
	}

	collection := client.Database(config.DBNAME).Collection("greenHouseConfig")
	defer db.CloseConnection(client, context.TODO(), w)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	filter := bson.D{{"_id", objID}}

	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		http.Error(w, "Erro ao excluir o registro", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Dados deletados com sucesso!")
}

func UpdateGreenHouse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PUT")
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	// Parse o corpo da solicitação para obter os dados a serem atualizados
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Erro ao analisar o corpo da solicitação", http.StatusBadRequest)
		return
	}

	client, err := db.GetConnection()
	if err != nil {
		fmt.Fprintln(w, "Erro ao conectar ao banco de dados!")
		return
	}
	defer db.CloseConnection(client, context.TODO(), w)

	collection := client.Database(config.DBNAME).Collection("greenHouseConfig")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	filter := bson.D{{"_id", objID}}
	update := bson.D{{"$set", data}}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, "Erro ao atualizar o registro", http.StatusInternalServerError)
		return
	}

	if result.ModifiedCount == 0 {
		http.Error(w, "Nenhum registro foi atualizado", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Registro atualizado com sucesso!")
}
func GetGreenHouses(w http.ResponseWriter, r *http.Request) {
	client, err := db.GetConnection()
	if err != nil {
		fmt.Fprintln(w, "Erro ao conectar ao banco de dados!")
		return
	}

	defer db.CloseConnection(client, context.TODO(), w)

	collection := client.Database(config.DBNAME).Collection("greenHouseConfig")

	var results []models.GreenHouseConfig
	cursor, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Erro na consulta", http.StatusInternalServerError)
		return
	}

	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var result models.GreenHouseConfig
		if err := cursor.Decode(&result); err != nil {
			fmt.Println(err)
			http.Error(w, "Erro na decodificação", http.StatusInternalServerError)
			return
		}
		results = append(results, result)
	}

	if err := cursor.Err(); err != nil {
		fmt.Println(err)
		http.Error(w, "Erro no cursor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
