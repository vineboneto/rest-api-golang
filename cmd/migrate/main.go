package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/urfave/cli/v2"
	"github.com/vineboneto/rest-api-golang/internal/infra/db"
)

func openDBMigration() (*migrate.Migrate, error) {

	absPath, err := filepath.Abs("migrations")
	if err != nil {
		log.Printf("[E] abs migrations path failed. err:%v", err)
		return nil, err
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s", absPath),
		db.BuildDSN(),
	)

	if err != nil {
		return nil, err
	}

	return m, nil
}

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}

	app := &cli.App{
		Name:  "ponkan-cli",
		Usage: "Ferramenta CLI para executar tarefas administrativas",
		Commands: []*cli.Command{
			{
				Name:  "migrate",
				Usage: "Executa migrações do banco de dados",
				Subcommands: []*cli.Command{
					{
						Name:  "up",
						Usage: "Aplica todas as migrações pendentes",
						Action: func(c *cli.Context) error {
							m, err := openDBMigration()

							if err != nil {
								return err
							}

							defer func() error {
								sourceErr, dbErr := m.Close()
								if sourceErr != nil {
									return sourceErr
								}
								if dbErr != nil {
									return dbErr
								}

								return nil
							}()

							m.Up()
							return nil
						},
					},
					{
						Name:  "down",
						Usage: "Desfaz a última migração",
						Action: func(c *cli.Context) error {

							return nil
						},
					},
				},
			},
		},
	}

	err = app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
