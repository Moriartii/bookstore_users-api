package users_db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

const (
	postgres_users_username = "postgres_users_username"
	postgres_users_password = "postgres_users_password"
	postgres_users_host     = "postgres_users_host"
	postgres_users_database = "postgres_users_database"
	postgres_users_sslmode  = "postgres_users_sslmode"
)

var (
	Client *sql.DB

	username = os.Getenv(postgres_users_username)
	password = os.Getenv(postgres_users_password)
	host     = os.Getenv(postgres_users_host)
	database = os.Getenv(postgres_users_database)
	sslmode  = os.Getenv(postgres_users_sslmode)
)

func init() {
	//datasourceName := ("user=udemy password=udemy dbname=udemy host=127.0.0.1 sslmode=disable")
	datasourceName := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=%s",
		username,
		password,
		host, //
		database,
		sslmode,
	)

	var err error
	Client, err = sql.Open("postgres", datasourceName)
	if err != nil {
		panic(err)
	}

	if err = Client.Ping(); err != nil {
		panic(err)
	}
	//pq.SetLogger()
	log.Println("database successfully configured")
}
