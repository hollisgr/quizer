package main

import (
	"quizer_server/internal/app"
	"quizer_server/internal/config"
)

func main() {
	cfg := config.MustLoad("config.env")
	pool := app.ConnectToDB(cfg)
	defer pool.Close()

	services := app.SetupServices(cfg, pool)

	router := app.SetupRouter(services)

	srv := app.SetupServer(cfg, router)

	app.StartServer(srv)

	app.HandleQuit(srv)
}
