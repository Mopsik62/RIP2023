package dsn

import (
	"fmt"
	"os"
)

func FromEnv() string {
	// Set the USERNAME environment variable to "MattDaemon"
	os.Setenv("DB_HOST", "0.0.0.0")
	os.Setenv("DB_NAME", "One-pot")
	os.Setenv("DB_PASS", "1")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PORT", "5432")

	host := os.Getenv("DB_HOST")

	if host == "" {
		return ""
	}

	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)
}
