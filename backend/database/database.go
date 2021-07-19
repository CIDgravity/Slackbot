package database

import (
	"context"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v4"
	"log"
	"twinQuasarAppV2/backend/customTypes"
)

func Connect() *pgx.Conn {

	// Read the configuration file
	var config customTypes.Config

	errReadConfig := cleanenv.ReadConfig("config/config.yaml", &config)

	if errReadConfig != nil {
		log.Fatal(errReadConfig)
	}

	var databaseUrl = config.Database.EndpointWithCredentials + config.Database.DatabaseName
	conn, err := pgx.Connect(context.Background(), databaseUrl)

	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n\n", err)
	}

	return conn
}
