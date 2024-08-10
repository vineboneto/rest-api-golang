package db

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	Conn *gorm.DB
}

func BuildDSN() string {

	port, err := strconv.ParseUint(os.Getenv("POSTGRES_PORT"), 10, 16)

	if err != nil {
		panic("POSTGRES_PORT is not int")

	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASS"),
		os.Getenv("POSTGRES_HOST"),
		uint16(port),
		os.Getenv("POSTGRES_DATABASE"),
	)

	return dsn
}

func NewPostgresDB() (*DB, error) {

	config, err := gorm.Open(postgres.Open(BuildDSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %v", err)
	}

	db, err := config.DB()

	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %v", err)
	}

	db.SetMaxOpenConns(30)

	return &DB{Conn: config}, nil
}

func (db *DB) Close() {
	gorm, err := db.Conn.DB()

	if err != nil {
		log.Fatalln(err)
	}

	err = gorm.Close()

	if err != nil {
		log.Fatalln(err)
	}

}
