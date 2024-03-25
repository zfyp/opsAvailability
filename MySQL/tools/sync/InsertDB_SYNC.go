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
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	defaultHost        = "127.0.0.1"
	defaultPort        = "3306"
	defaultUser        = "root"
	defaultPassword    = "123456"
	defaultDatabase    = "web"
	defaultTable       = "t01"
	defaultConcurrency = 10
)

var (
	host        string
	port        string
	user        string
	password    string
	database    string
	table       string
	concurrency int
)

func init() {
	flag.StringVar(&host, "h", defaultHost, "Database host")
	flag.StringVar(&port, "P", defaultPort, "Database port")
	flag.StringVar(&user, "u", defaultUser, "Database user")
	flag.StringVar(&password, "p", defaultPassword, "Database password")
	flag.StringVar(&database, "D", defaultDatabase, "Database name")
	flag.StringVar(&table, "t", defaultTable, "Table name")
	flag.IntVar(&concurrency, "c", defaultConcurrency, "Concurrency level")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample: %s -h %s -P %s -u %s -p %s -D %s -t %s -c %d\n",
			os.Args[0], defaultHost, defaultPort, defaultUser, defaultPassword, defaultDatabase, defaultTable, defaultConcurrency)
	}
}

func main() {
	flag.Parse()

	numRecords := promptNumRecords()

	db := openDBConnection()
	defer db.Close()

	insertRecords(db, numRecords)

	fmt.Println("Successfully inserted records into the table.")
}

func promptNumRecords() int {
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
	return numRecords
}

func openDBConnection() *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	return db
}

func insertRecords(db *sql.DB, numRecords int) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrency)

	for i := 0; i < numRecords; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			insertRecord(db)
		}()
	}

	wg.Wait()
}

func insertRecord(db *sql.DB) {
	stmtText := fmt.Sprintf("INSERT INTO %s(name, age, homeaddr, schooladdr) VALUES(?,?,?,?)", table)
	stmt, err := db.Prepare(stmtText)
	if err != nil {
		log.Fatalf("Error preparing statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		randomString(10),
		rand.Intn(100),
		randomString(40),
		randomString(40),
	)
	if err != nil {
		log.Fatalf("Error executing statement: %v", err)
	}
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
