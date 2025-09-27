package database
import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var Postdb *pgxpool.Pool

func ConnectDB() {
	godotenv.Load()

	dsn := os.Getenv("DB_DSN")

	db,err:= pgxpool.New(context.Background(), dsn)
	if err != nil{
		log.Fatal("cannot connect")
	}

	Postdb = db
	fmt.Println("db Connected")
}
