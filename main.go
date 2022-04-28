package main

import (
	"fmt"
	"go-tutorial/connection"
	"os"
)

func main() {
	c := connection.Conn()

	if !connection.CheckIfTableExists("tickets") {
		fmt.Println("Creating databse for tickets")
		connection.CreateTickectTable()
	}
	if !connection.CheckIfTableExists("users") {
		fmt.Println("Creating databse for users")
		connection.CreateUser()
	}

	// Create ticket table with limit of 20 tickets
	connection.InsertIntoTickets("Majesty Conference", 20)

	intro()

	// Ask users for their details
	// takeUserDetails()

	// connection.ReturnUsers()
	// for _, x := range connection.ReturnUsers() {
	// 	fmt.Println("=====")
	// 	fmt.Println(x)
	// }

	// fmt.Println(connection.ReturnUsers()[0])

	defer c.Close()

}

func intro() {
	fmt.Println("Welcome to the booking platform")
	fmt.Println("Type 'c' to continue or 'q' to quit")

	var option string
	fmt.Scan(&option)
	if option == "c" || option == "C" {
		takeUserDetails()
	} else if option == "q" || option == "Q" {
		os.Exit(3)
	} else {
		fmt.Println("Invalid input")
		intro()
	}
}

func takeUserDetails() {
	var firstName string
	var lasttName string
	var email string
	var num_of_tickets int

	fmt.Println("Please enter your first name:")
	fmt.Scan(&firstName)
	fmt.Println("Please enter your last name:")
	fmt.Scan(&lasttName)
	fmt.Println("Please enter your email address:")
	fmt.Scan(&email)
	fmt.Println("Please enter number of tickets:")
	fmt.Scan(&num_of_tickets)

	connection.InsertIntoUsers(firstName, lasttName, num_of_tickets, email)

	connection.UpdateTicketTable(num_of_tickets)
}
