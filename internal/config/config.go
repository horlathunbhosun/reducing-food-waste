package config

import "os"

func ConnectionStringAndDriver() (string, string) {
	connStr := os.Getenv("DB_CONNECTION_STRING")

	DBDriver := os.Getenv("DB_DRIVER")

	return connStr, DBDriver
}
