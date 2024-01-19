package main

import (
	"github.com/ThiagoDonadel/loan-management/internal/infra"
	"github.com/ThiagoDonadel/loan-management/internal/metrics"
	"github.com/ThiagoDonadel/loan-management/internal/registry"
)

func main() {

	infra.InitEnv()
	infra.InitlializeLogger()

	infra.Logger.Info("INICIALIZANDO")

	infra.Logger.Info("initializing configurations.")
	if err := infra.LoadConfigurationFromFile(); err != nil {
		infra.Logger.Fatal(err)
	}
	infra.Logger.Info("configurations initialized.")

	metrics.RegisterMetrics()

	infra.Logger.Info("initializing database connection.")
	if err := infra.ConnectToDatabase(); err != nil {
		infra.Logger.Fatal(err)
	}
	infra.Logger.Info("database connection initialized.")

	infra.Logger.Info("initializing dependency injection.")
	registry.Initialialize(infra.DBConnection)
	infra.Logger.Info("dependency injection initialized.")

	infra.Logger.Info("starting web server.")
	infra.StartGinServer(registry.GetControllers())
}
