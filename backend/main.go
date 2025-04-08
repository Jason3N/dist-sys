package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/Jason3N/super-duper-high-dist-sys/userpb"
	"github.com/goombaio/namegenerator"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
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
type server struct {
	userpb.UnimplementedUserServiceServer
	db *pgxpool.Pool
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
	defer conn.Close()
	// create a GPRC server at port :6000
	lis, err := net.Listen("tcp", ":6000")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, &server{db: conn})

	fmt.Printf("gRPC server up at port :6000")

	go func() {
		time.Sleep(3 * time.Second)
		RunClient()
	}()

	if err := grpcServer.Serve(lis); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}

}

func RunClient() {
	conn, err := grpc.Dial("localhost:6000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := userpb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	req := &userpb.BatchRequest{Amount: 50}
	res, err := client.CreateRandomUsersBatch(ctx, req)
	if err != nil {
		log.Fatalf("Error calling CreateRandomUsersBatch: %v", err)
	}

	fmt.Printf("Successfully inserted %d users into DB\n", res.GetAmount())
}

// func (s *Server) CreateRandomUsersBatch(ctx context.Context, req *userpb.BatchRequest) (*userpb.BatchResponse, error) {
// 	numOfUsers := int(req.GetAmount())
// 	users := make([]Users, numOfUsers)
// 	for i := 0; i < numOfUsers; i++ {
// 		user, pass := generateRandomName()
// 		users[i] = Users{username: user, password: pass}
// 	}
// 	batches := partitionBatch(users, 10)
// 	handleConcurrency(s.db, batches)
// 	return &userpb.BatchResponse{Amount: int32(numOfUsers)}, nil
// }

func handleConcurrency(conn *pgxpool.Pool, batches [][]Users) {
	for _, batch := range batches {
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
		batch.Queue(`INSERT INTO public."gRPC_table" (username, password) VALUES ($1, $2)`, users.username, users.password)
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
