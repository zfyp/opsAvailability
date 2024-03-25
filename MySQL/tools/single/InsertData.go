package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	// 自定义帮助信息
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample: %s -h 127.0.0.1 -P 3306 -u root -p'123456' -D 'web' -t t01\n", os.Args[0])
	}
}

// randomString generates a random string of a fixed length.
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func main() {
	// Parse command-line flags
	host := flag.String("h", "127.0.0.1", "Database host")
	port := flag.String("P", "3306", "Database port")
	user := flag.String("u", "root", "Database user")
	password := flag.String("p", "123456", "Database password")
	database := flag.String("D", "web", "Database name")
	table := flag.String("t", "t01", "Table name")
	flag.Parse()

	// Prompt user for the number of records to insert
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the number of records to insert: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error reading number of records to insert:", err)
	}
	input = strings.TrimSpace(input)
	numRecords, err := strconv.Atoi(input)
	if err != nil || numRecords <= 0 {
		log.Fatal("Invalid number of records to insert:", err)
	}

	// Build the DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", *user, *password, *host, *port, *database)

	// Open database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Prepare insert statement
	stmtText := fmt.Sprintf("INSERT INTO %s(name, age, homeaddr, schooladdr) VALUES(?,?,?,?)", *table)
	stmt, err := db.Prepare(stmtText)
	if err != nil {
		log.Fatalf("Error preparing statement: %v", err)
	}
	defer stmt.Close()

	// Insert records and display progress
	for i := 0; i < numRecords; i++ {
		_, err := stmt.Exec(
			randomString(10),
			rand.Intn(100),
			randomString(40),
			randomString(40),
		)
		if err != nil {
			log.Fatalf("Error executing statement: %v", err)
		}

		// Display progress
		fmt.Printf("\rInserted %d out of %d records (%.2f%% complete)", i+1, numRecords, float64(i+1)/float64(numRecords)*100)
	}

	fmt.Println("\nSuccessfully inserted records into the table.")
}
