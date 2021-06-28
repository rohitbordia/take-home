package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	sqlxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/jmoiron/sqlx"
	"os"
)

func Connect() (*sqlx.DB, func() error) {
	sqltrace.Register("postgres", &pq.Driver{}, sqltrace.WithServiceName("postgresql-"+os.Getenv("APP_NAME")))
	db := sqlxtrace.MustConnect("postgres", os.Getenv("DATABASE_URL"))
	return db, db.Close
}
