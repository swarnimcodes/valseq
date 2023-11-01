package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/alexbrainman/odbc"
	"os"
)

func main() {
	var server, database, username, password, connStr string
	// Accepting user input for server credentials
	fmt.Print("Enter Server IP: ")
	fmt.Scanln(&server)
	fmt.Print("Enter Database name: ")
	fmt.Scanln(&database)
	fmt.Print("Enter Username: ")
	fmt.Scanln(&username)
	fmt.Print("Enter Password: ")
	fmt.Scanln(&password)

	connStr = fmt.Sprintf(
		"driver=SQL Server;"+
			"Server=%s;"+
			"Database=%s;"+
			"UID=%s;"+
			"PWD=%s;",
		server, database, username, password,
	)

	fmt.Printf("Connection String: %s\n", connStr)

	db, err := sql.Open("odbc", connStr)
	if err != nil {
		fmt.Printf("Connection failed: %v\n", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Printf("Connection failed: %v\n", err)
		return
	}
	fmt.Println("Connection Successful!")
	var query string
	fmt.Printf("Enter your SQL Query: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		query = scanner.Text()
	} else {
		fmt.Printf("No input provided!\n")
	}
	fmt.Println("Entered SQL query:", query)
	// query = "SELECT @@VERSION"
	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("Query execution failed: %v\n", err)
		return
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	count := len(columns)
	scanArgs := make([]interface{}, count)
	values := make([]interface{}, count)

	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			fmt.Printf("Error scanning rows: %v\n", err)
			return
		}
		for i, col := range values {
			switch v := col.(type) {
			case nil:
				fmt.Printf("%s: NULL\n", columns[i])
			case []byte:
				fmt.Printf("%s: %s\n", columns[i], string(v))
			case int64:
				fmt.Printf("%s: %d\n", columns[i], v)
				// Add more cases if there are other expected types
			default:
				fmt.Printf("%s: Unexpected type %T\n", columns[i], v)
			}
		}

	}
}
