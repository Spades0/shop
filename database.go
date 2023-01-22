package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func connectDB() {
	db, error := sql.Open(`mysql`, `root:mysqlpassword@tcp(127.0.0.1:3306)/spades_shop`)
	
	error = db.Ping()
	if error != nil {
		fmt.Println("A connection error occurred!")
		panic(error.Error())
	}

	fmt.Println("\nSuccessfully connected to database")
	// defer db.Close()

	DB = db
}

func createInventoryTable() {
	_,  error := DB.Exec(`CREATE TABLE inventory (
		inventory_id INT NOT NULL AUTO_INCREMENT,
		name VARCHAR(45) NOT NULL,
		price INT NOT NULL,
		qty INT NOT NULL,
		qty_sold INT NULL,
		removed TINYINT NOT NULL DEFAULT 0,
		UNIQUE INDEX inventory_id_UNIQUE (inventory_id ASC));
	  `)

	if error != nil {
		panic(error)
	}	

	fmt.Println("Inventory table created")
}

func insertDefaultInventory() {
	insert, error := DB.Query(`INSERT INTO inventory (name, price, qty, qty_sold, removed)  
	VALUES
	('Mercedes-AMG GT 53', '25000', '3', '0', '0'),
    ('Ferrari Roma', '15000', '1', '0', '0'),
    ('BMW X6', '19000', '5', '0', '0'),
    ('Honda Accord', '17000', '5', '0', '0');`)

	if error != nil {
		fmt.Println("Database insert error!")
		panic(error.Error())
	}
	fmt.Println("Inventory created")
    defer insert.Close()
}

func database () {
	connectDB()

	_, table_check := DB.Query(`SELECT * FROM inventory`)
	if table_check != nil {
		createInventoryTable()
		insertDefaultInventory()
	} else {
		return
	}
	
}