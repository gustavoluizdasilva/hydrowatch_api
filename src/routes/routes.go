package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"hydrowatch-api/src/config"
	"hydrowatch-api/src/controller"
	"net/http"
)

func InitRoutes() {
	r := mux.NewRouter()

	greenhouseRoutes(r)
	webSocketRoute(r)

	http.Handle("/", r)
	fmt.Println("Server is up")

	err := http.ListenAndServe(":"+config.PORT, nil)

	if err != nil {
		return
	}
}

func greenhouseRoutes(r *mux.Router) {
	r.HandleFunc("/getDataSensor", controller.GetDataSensor).Methods("GET")
	r.HandleFunc("/getByIDGreenHouseConfig", controller.GetGreenHouseConfig).Methods("GET")
	r.HandleFunc("/getGreenHouses", controller.GetGreenHouses).Methods("GET")
	r.HandleFunc("/save-dataSensor", controller.SaveDataGreenhouse).Methods("POST")
	r.HandleFunc("/save-greenhouse", controller.SaveGreenHouseConfig).Methods("POST")
	r.HandleFunc("/delete-greenhouse", controller.DeleteGreenHouse).Methods("DELETE")
	r.HandleFunc("/save-greenhouse/{id}", controller.UpdateGreenHouse).Methods("PUT")
}

func webSocketRoute(r *mux.Router) {
	r.HandleFunc("/ws", controller.HandleWebSocket)
}
