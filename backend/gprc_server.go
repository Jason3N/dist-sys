package main

import (
	"context"
	"fmt"

	"github.com/Jason3N/super-duper-high-dist-sys/userpb"
)

func (s *server) CreateRandomUsersBatch(ctx context.Context, req *userpb.BatchRequest) (*userpb.BatchResponse, error) {
	numOfUsers := int(req.GetAmount())
	users := make([]Users, numOfUsers)
	for i := 0; i < numOfUsers; i++ {
		user, pass := generateRandomName()
		users[i] = Users{username: user, password: pass}
	}
	batches := partitionBatch(users, 10)
	handleConcurrency(s.db, batches)
	fmt.Printf("done")
	return &userpb.BatchResponse{Amount: int32(numOfUsers)}, nil
}

func (s *server) GetGlobalStat(ctx context.Context, _ *userpb.Empty) (*userpb.GlobalResponse, error) {
	mu.Lock()
	defer mu.Unlock()
	return &userpb.GlobalResponse{Amount: int32(global)}, nil
}
