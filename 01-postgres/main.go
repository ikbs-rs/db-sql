package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading environment")
	}

	dsn := url.URL{
		Scheme: os.Getenv("DB_SCHEMA"),
		Host:   os.Getenv("DB_HOST"),
		User:   url.UserPassword(os.Getenv("DB_USER"), os.Getenv("DB_PASS")),
		Path:   os.Getenv("DB_NAME"),
	}

	q := dsn.Query()
	q.Add("sslmode", "disable")

	dsn.RawQuery = q.Encode()

	conn, err := sql.Open("pgx", dsn.String())
	if err != nil {
		fmt.Println("sql.Open", err)
		return
	}
	defer func() {
		_ = conn.Close()
		fmt.Println("closed")
	}()

	if err := conn.PingContext(context.Background()); err != nil {
		fmt.Println("db.PingContext", err)
		return
	}
	// Jedan slog
	// Ova sintaksa prolazi sa navodnicima, mozda zbog slova u bazi.
	//row := conn.QueryRowContext(context.Background(), `SELECT "USER"."NAZIV" FROM "USER" where "USER"."USER_ID" = 'gost'`)
	row := conn.QueryRowContext(context.Background(), `select naziv, obj_id from objekat where obj_id = 'obj'`)
	if err := row.Err(); err != nil {
		fmt.Println("db.QueryRowContext", err)
		return
	}

	var naziv string

	if err := row.Scan(&naziv); err != nil {
		fmt.Println("row.Scan", err)
		return
	}

	fmt.Println("Naziv", naziv)

	//Svi slogovi
	rows, err := conn.QueryContext(context.Background(), `SELECT "USER"."NAZIV", "USER"."USER_ID" FROM "USER"`)
	if err != nil {
		fmt.Println("row.Scan", err)
		return
	}
	defer func() {
		_ = rows.Close()
	}()

	if rows.Err() != nil {
		fmt.Println("row.Err()", err)
		return
	}

	for rows.Next() {
		var naziv string
		var user_id string

		if err := rows.Scan(&naziv, &user_id); err != nil {
			fmt.Println("rows.Scan", err)
			return
		}

		fmt.Println("Naziv: ", naziv, "User ID: ", user_id)
	}
}
