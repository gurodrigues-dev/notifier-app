package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gurodrigues-dev/notifier-app/config"
	"github.com/gurodrigues-dev/notifier-app/internal/infra/setup"
	"github.com/spf13/viper"
)

func main() {
	setup.NewSetup()

	DBURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_HOST"),
		viper.GetString("DB_PORT"),
		viper.GetString("DB_NAME"),
	)

	m, err := migrate.New("file://database/migrations", DBURL)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Down(); err != nil {
		if fmt.Sprintf("%s", err) != "no change" {
			log.Fatal(err)
		}
	}

	if config.DevEnv() {
		cmd := exec.Command("pg_dump",
			"-s",
			"-h", viper.GetString("DB_HOST"),
			"-p", viper.GetString("DB_PORT"),
			"-U", viper.GetString("DB_USER"),
			"-d", viper.GetString("DB_NAME"),
		)

		cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", viper.GetString("DB_PASSWORD")))

		outFile, err := os.Create("database/structure.sql")
		if err != nil {
			log.Fatalf("erro ao criar structure.sql: %v", err)
		}
		defer outFile.Close()
		cmd.Stdout = outFile

		if err := cmd.Run(); err != nil {
			log.Fatalf("erro ao rodar pg_dump: %v", err)
		}
	}
}
