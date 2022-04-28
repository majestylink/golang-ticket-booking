package connection

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "majesty"
	password = "majesty"
	dbname   = "godb"
)

func Conn() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	// defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// fmt.Println("Successfully connected!")
	return db
}

func CheckIfTableExists(arg string) bool {
	_, table_check := Conn().Query("select * from " + arg + ";")

	if table_check == nil {
		// fmt.Printf("Skipped creating database %v. Already exists", arg)
		return true
	} else {
		// fmt.Printf("Successfully created database %v.", arg)
		return false
	}
}

func CreateTickectTable() {
	query := `
	CREATE TABLE tickets (
		id SERIAL PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		total INT
	  );`
	fmt.Printf("%T\n", query)
	_, err := Conn().Exec(query)
	if err != nil {
		panic(err)
	}
}

func CreateUser() {
	query := `CREATE TABLE users (
			id SERIAL PRIMARY KEY,
			first_name TEXT,
			last_name TEXT,
			num_of_tickets INT,
			email TEXT UNIQUE NOT NULL
		);`
	_, err := Conn().Exec(query)
	if err != nil {
		panic(err)
	}
}

func InsertIntoTickets(ticket_name string, total int) {
	query := `
	INSERT INTO tickets (name, total)
	VALUES ($1, $2)`
	_, err := Conn().Exec(query, ticket_name, total)
	if err != nil {
		panic(err)
	}
}

func InsertIntoUsers(fname string, lname string, num_of_tickets int, email string) {
	query := `
	INSERT INTO users (first_name, last_name, num_of_tickets, email)
	VALUES ($1, $2, $3, $4)`
	_, err := Conn().Exec(query, fname, lname, num_of_tickets, email)
	if err != nil {
		panic(err)
	}
}

func UpdateTicketTable(num_of_tickets int) {
	query := "SELECT * FROM tickets WHERE name = 'Majesty Conference'"

	var total int
	var name string
	var id int

	row := Conn().QueryRow(query)
	switch err := row.Scan(&id, &name, &total); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		if total < 1 {
			fmt.Println("Tickects sold out, wait for the next conference")
		} else {
			if num_of_tickets > 0 && num_of_tickets <= total {
				fmt.Println(id, total)
				if total < num_of_tickets {
					fmt.Printf("We only have %v tickets available\n", total)
				} else {
					updateQuery := `UPDATE tickets SET total = $1 WHERE name = $2;`
					_, err := Conn().Exec(updateQuery, total-num_of_tickets, "Majesty Conference")
					if err != nil {
						panic(err)
					}
					fmt.Printf("%v number of tickets remaining", total-num_of_tickets)
				}
			} else {
				fmt.Println("Invalid input for number of tickets")
			}
		}
	default:
		panic(err)
	}
}

func ReturnUsers() []string {
	query := "SELECT * FROM users"

	var (
		id             int
		first_name     string
		last_name      string
		num_of_tickets int
		email          string
	)
	rows, err := Conn().Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	booking := []string{}
	for rows.Next() {
		err := rows.Scan(&id, &first_name, &last_name, &num_of_tickets, &email)
		if err != nil {
			panic(err)
		}
		// fmt.Println("\n", first_name, num_of_tickets)
		booking = append(booking, first_name+"\t"+strconv.FormatInt(int64(num_of_tickets), 10))
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return booking
}
