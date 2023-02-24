package DatabaseAbstraction

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"runtime/debug"
	"strconv"
	"time"
)

type DBConnector struct {
	DB *pgxpool.Pool
}

type DBOrm interface {
	GetAllLicenseKeys() ([]LicenseKey, error)
	GetLicenseKeyByKeyID(keyID string) (LicenseKey, error)
	GetLicenseKeyByIndexID(indexID int) (LicenseKey, error)
	AddLicenseKey(licenseKey LicenseKey) (int, error)
	DeleteLicenseKey(indexID int) error

	GetAllProducts() ([]Product, error)
	GetProductByIndexID(indexID int) (Product, error)
	GetProductLicenseKeys(indexID int) ([]LicenseKey, error)
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
	GetOwnedLicenseKeys(indexID int) ([]LicenseKey, error)
	IncreaseUserBalance(indexID int, amount int) error
	DecreaseUserBalance(indexID int, amount int) error
	AddOwnedProduct(indexID int, productID int) error
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

	_, err = db.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, username VARCHAR(255) UNIQUE, password VARCHAR(255) NOT NULL, balance int DEFAULT 0 NOT NULL, created_at DATE DEFAULT CURRENT_DATE NOT NULL, updated_at DATE DEFAULT CURRENT_DATE NOT NULL, CONSTRAINT balance_never_negative CHECK ( balance >= 0 ));")
	if err != nil {
		logrus.Error(err)
		debug.PrintStack()
		return nil, err
	}
	_, err = db.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS products (id SERIAL PRIMARY KEY, name VARCHAR(255), description TEXT, price INTEGER, image VARCHAR(255), mpd_url VARCHAR(255), created_at DATE DEFAULT CURRENT_DATE, updated_at DATE DEFAULT CURRENT_DATE);")
	if err != nil {
		logrus.Error(err)
		debug.PrintStack()
		return nil, err
	}
	_, err = db.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS licenses (id SERIAL PRIMARY KEY, key_id VARCHAR(255), encryption_key VARCHAR(255), product_id INTEGER REFERENCES products(id), created_at DATE DEFAULT CURRENT_DATE, updated_at DATE DEFAULT CURRENT_DATE);")
	if err != nil {
		logrus.Error(err)
		debug.PrintStack()
		return nil, err
	}
	_, err = db.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS user_purchases (id SERIAL PRIMARY KEY, user_id INTEGER, product_id INTEGER, FOREIGN KEY(user_id) REFERENCES users(id), FOREIGN KEY(product_id) REFERENCES products(id));")
	if err != nil {
		logrus.Error(err)
		debug.PrintStack()
		return nil, err
	}
	_, err = db.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS user_tokens (id SERIAL PRIMARY KEY, user_id INTEGER, token VARCHAR(255), expiry DATE, FOREIGN KEY(user_id) REFERENCES users(id));")
	if err != nil {
		logrus.Error(err)
		debug.PrintStack()
		return nil, err
	}
	_, err = db.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS product_licenses (id SERIAL PRIMARY KEY, product_id INTEGER, license_id INTEGER, FOREIGN KEY(product_id) REFERENCES products(id), FOREIGN KEY(license_id) REFERENCES licenses(id));")
	if err != nil {
		logrus.Error(err)
		debug.PrintStack()
		return nil, err
	}
	log.Printf("Database connection pool has %d connections", stat.TotalConns())
	log.Printf("Database connection pool has %d idle connections", stat.MaxConns())

	return db, nil
}
