package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/goombaio/namegenerator"
	"github.com/jackc/pgx"
	"github.com/joho/godotenv"
)

var (
	global = 0
	mu     sync.Mutex
	wg     sync.WaitGroup
)

func main() {
	start := time.Now()
	handleConcurrency()
	handleConnection()
	fmt.Printf("Total time: %v\n", time.Since(start))
	fmt.Printf("%v\n", global)
}

func handleConnection() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("env file cannot be read")
		os.Exit(1)
	}

	CONNECTION_STRING := os.Getenv("DATABASE_URL")

	config, err := pgx.ParseConnectionString(CONNECTION_STRING)
	if err != nil {
		fmt.Println("Database URL is wrong, please check again")
		os.Exit(1)
	}

	conn, err := pgx.Connect(config)
	if err != nil {
		fmt.Println("Cannot connect to db")
		os.Exit(1)
	}

	fmt.Println("Connected to db")
	defer conn.Close()

	addRandomName(conn)

	var id int
	var username, password string
	err = conn.QueryRow(`SELECT id, username, password FROM "user_t" LIMIT 1`).Scan(&id, &username, &password)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(id, username, password)
}

func addRandomName(conn *pgx.Conn) {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)
	name := nameGenerator.Generate()
	password := name + "pw"
	fmt.Printf("%v, %v\n", name, password)
	var lastID int
	err := conn.QueryRow(`INSERT INTO "user_t" (username, password) values ($1, $2) RETURNING ID`, name, password).Scan(&lastID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("last id is %v\n", lastID)
}

func handleConcurrency() {
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go heavyHandler(&wg)
	}
	wg.Wait()
}

func heavyHandler(wg *sync.WaitGroup) {
	defer wg.Done()
	// some action must be done here, simulated with accessing atomic var
	// this will be placeholder for read/write on db
	// will stop race conditions by placing mutex lock
	mu.Lock()
	global++
	mu.Unlock()
	// assume action takes 1 second
	time.Sleep(1 * time.Second)
}

func mediumHandler(wg *sync.WaitGroup) {
	defer wg.Done()
	// 0.5 second
	mediumTime := new(decimal.Big).SetFloat64(0.5)
	mediumTimeFloat, _ := mediumTime.Float64()
	sleepDuration := time.Duration(mediumTimeFloat * float64(time.Second))
	time.Sleep(sleepDuration)
}

func lowHandler(wg *sync.WaitGroup) {
	defer wg.Done()
	// 0.25 second
	lowTime := new(decimal.Big).SetFloat64(0.5)
	lowFloat, _ := lowTime.Float64()
	sleepDuration := time.Duration(lowFloat * float64(time.Second))
	time.Sleep(sleepDuration)
}
