package main

import (
	"database/sql"
	"log"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

func main() {
	// Connect to the PostgreSQL database
	serviceURI := "postgres://avnadmin:AVNS_h2hsxPs7zF_ljiAf1v4@deploy-postgresql-deploy-postgresql.aivencloud.com:12446/defaultdb?sslmode=require"

	conn, _ := url.Parse(serviceURI)
	conn.RawQuery = "sslmode=verify-ca;sslrootcert=ca.pem"
	db, err := sql.Open("postgres", conn.String())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create a new Echo instance
	e := echo.New()

	// Define routes
	e.GET("/users", getUsers(db))
	e.POST("/users", createUser(db))

	// Start the server
	e.Start(":8080")
}

func getUsers(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		rows, err := db.Query("SELECT id, name, message FROM users")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		users := []User{}
		for rows.Next() {
			var user User
			err := rows.Scan(&user.ID, &user.Name, &user.Message)
			if err != nil {
				log.Fatal(err)
			}
			users = append(users, user)
		}

		return c.JSON(http.StatusOK, users)
	}
}

func createUser(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := new(User)
		if err := c.Bind(user); err != nil {
			return err
		}

		_, err := db.Exec("INSERT INTO users (name, message) VALUES ($1, $2)", user.Name, user.Message)
		if err != nil {
			log.Fatal(err)
		}

		return c.JSON(http.StatusCreated, user)
	}
}
