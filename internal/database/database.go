package database

import (
	"database/sql"
	"github.com/horlathunbhosun/wastewarrior-api/internal/config"
	"github.com/joho/godotenv"
	"log"
)

var DB *sql.DB

func InitDB() {

	err := godotenv.Load()
	if err != nil {
		return
	}
	connStr, dbServer := config.ConnectionStringAndDriver()

	DB, err = sql.Open(dbServer, connStr)

	if err != nil {
		panic("Could not connect to the database")
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	createTables()
}

func createTables() {
	createUsersTable()
	createPartnersTable()
	createProductsTable()
	createMagicBagsTable()
	createTransactionsTable()
	createMagicBagProductsTable()
	createFeedbackTable()
}

func createUsersTable() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
      id INTEGER PRIMARY KEY AUTO_INCREMENT,
    fullname VARCHAR(30) NOT NULL,
    email VARCHAR(50) UNIQUE,
    phone_number VARCHAR(40) UNIQUE,
    user_type ENUM('waste_warrior', 'partners', 'admin') NOT NULL,
    date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
	date_updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalln(err)
		panic("Can not users table")
	}
}

func createPartnersTable() {
	query := `
	CREATE TABLE IF NOT EXISTS partners (
	  id INTEGER PRIMARY KEY AUTO_INCREMENT,
	business_number VARCHAR(30) NOT NULL,
	user_id INTEGER NOT NULL,   
	logo VARCHAR(50) UNIQUE,
	address VARCHAR(40) UNIQUE,
	date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
	date_updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
	)`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalln(err)
		panic("Can not partners table")
	}
}

func createProductsTable() {
	query := `
	CREATE TABLE IF NOT EXISTS products (
	  id INTEGER PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR(30) NOT NULL,
	date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
	date_updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalln(err)
		panic("Can not products table")
	}
}

func createMagicBagsTable() {
	query := `
	CREATE TABLE IF NOT EXISTS magic_bags (
	  id INTEGER PRIMARY KEY AUTO_INCREMENT,
	  bag_price FLOAT,
	  partner_id INTEGER NOT NULL,
	date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
	date_updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	FOREIGN KEY (partner_id) REFERENCES partners(partner_id) ON DELETE CASCADE

	)`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalln(err)
		panic("Can not magic_bags table")
	}
}

func createTransactionsTable() {
	query := `
	CREATE TABLE IF NOT EXISTS transactions (
	  id INTEGER PRIMARY KEY AUTO_INCREMENT,
	  amount FLOAT,
	magic_bag_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
	date_updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	FOREIGN KEY (magic_bag_id) REFERENCES magic_bags(magic_bag_id) ON DELETE CASCADE,
	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
	UNIQUE KEY waste_warrior_purchase_unique (user_id, magic_bag_id, date_created),

	)`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalln(err)
		panic("Can not transactions table")
	}
}

func createMagicBagProductsTable() {
	query := `CREATE TABLE IF NOT EXISTS magic_bag_products (
    	id INTEGER PRIMARY KEY AUTO_INCREMENT,
    	quantity INTEGER,
    	magic_bag_id INTEGER NOT NULL,
    	product_id INTEGER NOT NULL,
    	date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
    	date_updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    	FOREIGN KEY (magic_bag_id) REFERENCES magic_bags(magic_bag_id) ON DELETE CASCADE,
    	FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE   
)`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalln(err)
		panic("Can not transactions table")
	}
}

func createFeedbackTable() {
	query := `CREATE TABLE IF NOT EXISTS feedback (
		id INTEGER PRIMARY KEY AUTO_INCREMENT,
		rating  INTEGER DEFAULT 0,
    	comment TEXT
		transaction_id INTEGER NOT NULL,
		date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
		date_updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (transaction_id) REFERENCES transaction(transaction_id) ON DELETE CASCADE
)`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalln(err)
		panic("Can not feedback table")
	}
}

//
