package bootstrap

import (
	"hydrowatch-api/src/config"
	"hydrowatch-api/src/routes"
)

func Bootstrap() {
	config.CarrConfig()
	routes.InitRoutes()
}
