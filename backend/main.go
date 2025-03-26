package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/goombaio/namegenerator"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

var (
	global = 0
	mu     sync.Mutex
	wg     sync.WaitGroup
)

type Users struct {
	username string
	password string
}

func main() {
	start := time.Now()
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

	conn, err := pgxpool.Connect(context.Background(), CONNECTION_STRING)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Println("Connected to db")

	// create X amount of users to create
	numOfUsers := 100
	users := make([]Users, numOfUsers)
	for i := 0; i < numOfUsers; i++ {
		user, pass := generateRandomName()
		users[i] = Users{username: user, password: pass}
	}

	// partition those users into amounts of X batches
	batchesSize := 10
	batches := partitionBatch(users, batchesSize)
	fmt.Printf("Total batches: %d\n", len(batches))
	handleConcurrency(conn, batches)

	defer conn.Close()
}

func handleConcurrency(conn *pgxpool.Pool, batches [][]Users) {
	for _, batch := range batches {
		// batch is a bunch of [][]Users
		wg.Add(1)
		go addRandomNameConBatch(conn, batch)
	}
	wg.Wait()
}

func partitionBatch(users []Users, batchesSize int) [][]Users {
	var batches [][]Users
	for i := 0; i < len(users); i += batchesSize {
		end := i + batchesSize
		if end > len(users) {
			end = len(users)
		}
		batches = append(batches, users[i:end])
	}
	return batches
}

func generateRandomName() (string, string) {
	// generate random_seed
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)
	// generate username / password
	name := nameGenerator.Generate()
	password := name + "pw"
	// return user, pw
	return name, password
}

// adding random users through single insertions iteratively
func addRandomName(conn *pgx.Conn) {
	user, pass := generateRandomName()
	var lastID int
	err := conn.QueryRow(context.Background(), `INSERT INTO "user_t" (username, password) values ($1, $2) RETURNING ID`, user, pass).Scan(&lastID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// adding random users through concurrency, single insertions
func addRandomNameCon(conn *pgxpool.Pool, wg *sync.WaitGroup) {
	defer wg.Done()
	mu.Lock()
	user, pass := generateRandomName()
	mu.Unlock()
	var lastID int
	err := conn.QueryRow(context.Background(), `INSERT INTO "user_t" (username, password) values ($1, $2) RETURNING ID`, user, pass).Scan(&lastID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// adding random users through concurrency, batch insertions
func addRandomNameConBatch(conn *pgxpool.Pool, user []Users) error {
	defer wg.Done()
	batch := &pgx.Batch{}
	for _, users := range user {
		batch.Queue(`INSERT INTO "user_t" (username, password) VALUES ($1, $2)`, users.username, users.password)
	}

	br := conn.SendBatch(context.Background(), batch)
	for range user {
		_, err := br.Exec()
		if err != nil {
			br.Close()
			os.Exit(1)
		}
	}
	return br.Close()
}

// TODO: MORE FUNCTIONS HERE
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
