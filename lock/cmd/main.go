package main

import (
	"context"
	"fmt"
	"github.com/shantanubansal/go-kit/lock"
	"github.com/shantanubansal/go-kit/lock/datastore"
	log "github.com/sirupsen/logrus"
	"time"
)

var lockManager *lock.LockManager

func initLockManager() {
	mongoStore, err := datastore.GetMongoDBDataStoreInstance("mongodb://localhost:27017", "mydb")
	if err != nil {
		panic(err)
	}
	lockManager = lock.NewLockManager(mongoStore)
}
func main() {
	initLockManager()

	if lockManager == nil {
		log.Fatalf("lock manager is not initailized")
	}
	lockID := "resource1"
	ctx := context.Background()

	requestID, err := lockManager.AcquireLock(ctx, lockID)
	if err != nil {
		fmt.Println("Failed to acquire lock:", err)
		return
	}

	if requestID != "" {
		fmt.Println("Lock acquired:", requestID)
		time.Sleep(10 * time.Second)
		if err := lockManager.ReleaseLock(ctx, lockID, requestID); err != nil {
			fmt.Println("Failed to release lock:", err)
			return
		}
		fmt.Println("Lock released")
	} else {
		fmt.Println("Lock was already held")
	}
}
