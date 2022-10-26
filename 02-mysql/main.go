package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading environment")
	}

	c := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASS"),
		Addr:                 os.Getenv("DB_HOST"),
		Net:                  "tcp",
		DBName:               os.Getenv("DB_NAME"),
		ParseTime:            true,
		AllowNativePasswords: true,
	}

	fmt.Println(c.FormatDSN())

	db, err := sql.Open("mysql", c.FormatDSN())
	if err != nil {
		fmt.Println("sql.Open", err)
		return
	}
	defer func() {
		_ = db.Close()
		fmt.Println("closed")
	}()

	if err := db.PingContext(context.Background()); err != nil {
		fmt.Println("db.PingContext", err)
		return
	}

	row := db.QueryRowContext(context.Background(),
		"SELECT naziv FROM objekat WHERE objekat_id = ? or id = ?",
		"obj2", 2)
	if err := row.Err(); err != nil {
		fmt.Println("db.QueryRowContext", err)
		return
	}

	var naziv string
	if err := row.Scan(&naziv); err != nil {
		fmt.Println("row.Scan", err)
		return
	}

	fmt.Println("Objekat je", naziv)
}
