package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

var (
	host     = os.Getenv("DB_HOST")
	password = os.Getenv("DB_PASSWORD")
	portStr  = os.Getenv("DB_PORT")
	user     = os.Getenv("DB_USERNAME")
	database = os.Getenv("DB_NAME")
)

func GetConnection() (*sql.DB, error) {
	port, err := strconv.ParseUint(portStr, 10, 32)
	if err != nil {
		return nil, err
	}
	psqlInfo := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=require", user, password, host, port, database)
	return sql.Open("postgres", psqlInfo)
}

const createEmployeesTableSQL = `
CREATE SEQUENCE IF NOT EXISTS employees_id_seq;

CREATE TABLE IF NOT EXISTS employees (
	id integer DEFAULT nextval('employees_id_seq'),
	email text,
	first_name varchar,
	last_name varchar,
	PRIMARY KEY (id)
  );
`

/* const deleteEmployeesTableSQL = `
DROP TABLE employees;
` */

func CreateEmployeesTable(ctx context.Context, db *sql.DB) error {
	//db.ExecContext(ctx, deleteEmployeesTableSQL)
	_, err := db.ExecContext(ctx, createEmployeesTableSQL)
	return err
}
