package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/database"
	"github.com/urfave/cli/v2"
	"github.com/vineboneto/rest-api-golang/internal/infra/db"
	"github.com/vineboneto/rest-api-golang/migrations"
)

func printMigrationResults(results ...*goose.MigrationResult) {
	log.Println("\n=== migration results  ===")
	for _, r := range results {
		log.Printf("%-3s %-2v done: %v\n", r.Source.Type, r.Source.Version, r.Duration)
	}
}

func printMigrationStatus(results ...*goose.MigrationStatus) {
	log.Println("\n=== migration status  ===")
	for _, r := range results {
		log.Printf("%-3s %-2v applied: %v status:%v\n", r.Source.Type, r.Source.Version, r.AppliedAt, r.State)
	}
}

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}

	db, err := sql.Open("postgres", db.BuildDSN())

	if err != nil {
		log.Fatalf("sql: failed to open DB: %v\n", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("sql: failed to close DB: %v\n", err)
		}
	}()

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
							provider, err := goose.NewProvider(database.DialectPostgres, db, migrations.Embed)

							if err != nil {
								return err
							}

							ctx := context.Background()
							results, err := provider.Up(ctx)

							if err != nil {
								log.Fatal(err)
							}

							printMigrationResults(results...)

							return nil
						},
					},
					{
						Name:  "down",
						Usage: "Desfaz a última migração",
						Action: func(c *cli.Context) error {
							provider, err := goose.NewProvider(database.DialectPostgres, db, migrations.Embed)

							if err != nil {
								return err
							}

							ctx := context.Background()
							results, err := provider.Down(ctx)

							if err != nil {
								log.Fatal(err)
							}

							printMigrationResults(results)

							return nil
						},
					},
					{
						Name:  "status",
						Usage: "Verifica status do banco",
						Action: func(c *cli.Context) error {
							provider, err := goose.NewProvider(database.DialectPostgres, db, migrations.Embed)

							if err != nil {
								return err
							}

							ctx := context.Background()
							results, err := provider.Status(ctx)

							if err != nil {
								log.Fatal(err)
							}

							printMigrationStatus(results...)

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
