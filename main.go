package main

import (
	"fmt"
	"log"
	"os"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jedib0t/go-pretty/v6/table"
)

func welcomePage() {
	newline(1)
	fmt.Println("Welcome to ")
	fmt.Println("\n███████╗██████╗  █████╗ ██████╗ ███████╗███████╗    ███████╗██╗  ██╗ ██████╗ ██████╗ ")
	fmt.Println("██╔════╝██╔══██╗██╔══██╗██╔══██╗██╔════╝██╔════╝    ██╔════╝██║  ██║██╔═══██╗██╔══██╗")
	fmt.Println("███████╗██████╔╝███████║██║  ██║█████╗  ███████╗    ███████╗███████║██║   ██║██████╔╝")
	fmt.Println("╚════██║██╔═══╝ ██╔══██║██║  ██║██╔══╝  ╚════██║    ╚════██║██╔══██║██║   ██║██╔═══╝")
	fmt.Println("███████║██║     ██║  ██║██████╔╝███████╗███████║    ███████║██║  ██║╚██████╔╝██║")
	fmt.Println("╚══════╝╚═╝     ╚═╝  ╚═╝╚═════╝ ╚══════╝╚══════╝    ╚══════╝╚═╝  ╚═╝ ╚═════╝ ╚═╝\n")
	fmt.Println("We sell Luxury cars and offer free deliveries worldwide.")
}

var (
	inventory Product
)

type Product struct {
	inventoryID uint
	productName     string
	productPrice    uint
	productQty uint
	qtySold uint
}

func (p Product) showPrice() {
	fmt.Println("The price of this product is: ", p.productPrice)
}

func main() {
	database()
	welcomePage()
	showInventory()
	displayMenu()
}

// Display different menu options
func displayMenu() {
	newline(1)
	fmt.Println("Select an operation:")
	fmt.Println("1. Show inventory.\t\t2. Buy product.")
	fmt.Println("3. Remove product\t\t4. Add product")
	fmt.Println("5. View sales\t\t\t6. Exit ")

	var menuNumber int
	_, err := fmt.Scan(&menuNumber)

	if err != nil {
		fmt.Println("Error: Please enter a number between 1 and 6")
	}

	//Switch for menu options
	switch menuNumber {
	case 1:
		showInventory()
	case 2:
		buyProduct()
	case 3:
		removeProduct()
	case 4:
		addProduct()
	case 5:
		showSales()
	case 6:
		exit()
	default:
		fmt.Println("Error: Please enter a number between 1 and 6")
		displayMenu()
	}
}

// Add new n line(s)
func newline(numberOfLines int) {
	i := 0
	for i < numberOfLines {
		fmt.Println("\n")
		i++
	}
}

// Display all information
func showInventory(){
	fmt.Println("This is our current inventory")
	showInventoryDatabase, error := DB.Query(`SELECT * FROM inventory;`)
	if error != nil{
		log.Fatal("Error occured when fetching inventory data:", error)
	}


	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"inventoryID", "Product name", "Price (USD)", "Quantity", "Quantity sold"})
	for showInventoryDatabase.Next() {
		error = showInventoryDatabase.Scan(&inventory.inventoryID, &inventory.productName, &inventory.productPrice, &inventory.productQty, &inventory.qtySold)
		if error != nil{
			panic(error.Error())
		}
        
		t.AppendRows([]table.Row{
			{inventory.inventoryID, inventory.productName, inventory.productPrice, inventory.productQty, inventory.qtySold},
		})
		t.AppendSeparator()
	}
	// newline(1)
	t.Render()
	displayMenu()
}

func addProduct() {
	newProduct := new(Product)
	fmt.Println("Enter the name of the Car: ")
	_, err_name := fmt.Scan(&newProduct.productName)
	if err_name != nil{
		log.Fatal(err_name)
	}

	fmt.Println("Enter the price of the product: ")
	_, err_price := fmt.Scan(&newProduct.productPrice)
	if err_price != nil{
		log.Fatal(err_price)
	} else if newProduct.productPrice * 1 != newProduct.productPrice || newProduct.productPrice * 0 != 0 {
		fmt.Println("Enter a number!")
		return
	}

	fmt.Println("Enter the quantity of the product: ")
	_, err_qty := fmt.Scan(&newProduct.productQty)
	if err_qty != nil{
		log.Fatal(err_qty)
	} else if newProduct.productQty * 1 != newProduct.productQty || newProduct.productQty * 0 != 0 {
		fmt.Println("Enter a number!")
		return
	}

	newProduct.qtySold = 0

	insertQuery := "INSERT INTO inventory (name, price, qty, qty_sold) VALUES (?, ?, ?, ?);"
	stmt, error := DB.Prepare(insertQuery)
	if error != nil {
		log.Fatal("Unable to prepare statement:", error)
	}

	_, err := stmt.Exec( newProduct.productName, newProduct.productPrice, newProduct.productQty, newProduct.qtySold)
	if err != nil {
		log.Fatal("Unable to execute statement:", error)
	}
	fmt.Printf("%s was succesfully added to the store.", newProduct.productName)
    
	// newline(1)
	displayMenu()
}

func buyProduct() {
	var id uint
	fmt.Println("Please select the inventoryID of the product you want to purchase:")
	_, error := fmt.Scan(&id)
	if error != nil {
		panic(error)
	}
	
	selector, error := DB.Query(`SELECT * FROM inventory WHERE inventory_id = ?;`, id)
	if error != nil{
		log.Fatal("Error occured when preparing the inventory data:", error)
	}

	for selector.Next() {
		error = selector.Scan(&inventory.inventoryID, &inventory.productName, &inventory.productPrice, &inventory.productQty, &inventory.qtySold)
		if error != nil{
			panic(error.Error())
		}
	}

	if error != nil{
		log.Fatal("Error occured when fetching inventory data:", error)
	}

	if id != inventory.inventoryID {
		fmt.Println("Please enter a valid product ID")
		return
	}

	if inventory.productQty == 0 {
		fmt.Println("This product is unavailable for sale")
		return
	} else {
		inventory.productQty = inventory.productQty - 1
		inventory.qtySold = inventory.qtySold + 1
		_, error := DB.Exec(`UPDATE inventory SET qty = ?, qty_sold = ? WHERE inventory_id = ?;`, inventory.productQty, inventory.qtySold, id)
		if error != nil {
			panic(error)
		}
	}

	fmt.Printf("You successfuly purchased %s costing $%d from Spades shop.", inventory.productName, inventory.productPrice)
	newline(1)

	displayMenu()
}

func showSales () {
	fmt.Println("Below shows the current sales of Spades shop")
	currentSales, error := DB.Query(`SELECT name, price, qty_sold FROM inventory;`)
	if error != nil{
		log.Fatal("Error occured when fetching sales data:", error)
	}

	for currentSales.Next() {
		error = currentSales.Scan(&inventory.productName, &inventory.productPrice, &inventory.qtySold)
		if error != nil{
			panic(error.Error())
		}
        
		fmt.Printf("We sold %d %s car(s) worth $%d\n", inventory.qtySold, inventory.productName, inventory.productPrice)
	}
	displayMenu()
}

func removeProduct() {
	var id uint
	fmt.Println("Please select the inventoryID of the product you want to remove from the store:")
	_, error := fmt.Scan(&id)
	if error != nil {
		panic(error)
	}
	
	remove, error := DB.Query(`DELETE FROM inventory WHERE inventory_id = ?;`, id)
	if error != nil{
		log.Fatal("Error occured when deleting product from the inventory:", error)
	}
    
	for remove.Next() {
		error = remove.Scan(&inventory.inventoryID, &inventory.productName, &inventory.productPrice, &inventory.productQty, &inventory.qtySold)
		if error != nil{
			panic(error.Error())
		}
	}

	if error != nil{
		log.Fatal("Error occured when fetching inventory data:", error)
	}

    defer remove.Close()
	fmt.Printf("You have successfuly removed %s from Spades shop.", inventory.productName)
	showInventory()
	displayMenu()
}

func exit() {
	fmt.Println("Thanks for using our shop.")
	newline(1)
	fmt.Println("\t██████╗  ██████╗  ██████╗ ██████╗     ██████╗ ██╗   ██╗███████╗")
	fmt.Println("\t██╔════╝ ██╔═══██╗██╔═══██╗██╔══██╗    ██╔══██╗╚██╗ ██╔╝██╔════╝")
	fmt.Println("\t██║  ███╗██║   ██║██║   ██║██║  ██║    ██████╔╝ ╚████╔╝ █████╗")
	fmt.Println("\t██║   ██║██║   ██║██║   ██║██║  ██║    ██╔══██╗  ╚██╔╝  ██╔══╝")
	fmt.Println("\t╚██████╔╝╚██████╔╝╚██████╔╝██████╔╝    ██████╔╝   ██║   ███████╗")
	fmt.Println("\t ╚═════╝  ╚═════╝  ╚═════╝ ╚═════╝     ╚═════╝    ╚═╝   ╚══════╝")
	os.Exit(0)
}