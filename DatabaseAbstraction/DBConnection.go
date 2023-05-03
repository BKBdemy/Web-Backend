package DatabaseAbstraction

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"strconv"
	"time"
)

type DBConnector struct {
	DB *pgxpool.Pool
}

type DBOrm interface {
	GetAllProducts() ([]Product, error)
	GetProductByIndexID(indexID int) (Product, error)
	AddProduct(NewProduct Product) (int, error)

	GetTokenByTokenID(tokenID string) (Token, error)
	GetTokenByHash(token string) (Token, error)
	AddToken(userID int, token string, expiry time.Time) error
	DeleteToken(tokenID int) error
	DeleteTokenByHash(token string) error

	GetAllUsers() ([]User, error)
	GetUserByUsername(username string) (User, error)
	GetUserByIndexID(indexID int) (User, error)
	AddUser(username string, password string) error
	DeleteUser(indexID int) error
	UpdateUserPassword(indexID int, newPassword string) error
	UpdateUserUsername(indexID int, newUsername string) error
	GetOwnedProducts(indexID int) ([]Product, error)
	IncreaseUserBalance(indexID int, amount int) error
	DecreaseUserBalance(indexID int, amount int) error
	AddOwnedProduct(indexID int, productID int) error

	MarkVideoAsWatched(indexID int, user User) error
	GetWatchedVideosByUser(user User) ([]Video, error)

	GetAllVideos() ([]Video, error)
	GetVideosByProductIndexID(productID int) ([]Video, error)
	GetVideoByIndexID(indexID int) (Video, error)
	GetProductByVideoIndexID(indexID int) (Product, error)
}

const (
	default_postgres_host     = "localhost"
	default_postgres_port     = 5432
	default_postgres_user     = "license"
	default_postgres_password = "license"
	default_postgres_dbname   = "license"
)

func getenvWithFallback(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func Connect() (*pgxpool.Pool, error) {
	host := getenvWithFallback("DB_HOST", default_postgres_host)
	port, err := strconv.Atoi(getenvWithFallback("DB_PORT", fmt.Sprintf("%d", default_postgres_port)))
	if err != nil {
		log.Fatal("DB_PORT is not a valid integer " + err.Error())
	}
	user := getenvWithFallback("DB_USER", default_postgres_user)
	password := getenvWithFallback("DB_PASSWORD", default_postgres_password)
	dbname := getenvWithFallback("DB_NAME", default_postgres_dbname)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	censoredPsqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, "********", dbname)

	log.Printf("Connecting to database with connection string: %s", censoredPsqlInfo)

	var db *pgxpool.Pool

	db, err = pgxpool.New(context.Background(), psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	stat := db.Stat()

	log.Printf("Database connection pool has %d connections", stat.TotalConns())
	log.Printf("Database connection pool has %d idle connections", stat.MaxConns())

	return db, nil
}
