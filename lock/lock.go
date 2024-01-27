package lock

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type DistributedLock struct {
	LockID    string    `bson:"lock_id"`
	RequestID string    `bson:"request_id"`
	Timestamp time.Time `bson:"timestamp"`
}

type LockManager struct {
	store DataStore
}

func NewLockManager(store DataStore) *LockManager {
	return &LockManager{store: store}
}

func (lm *LockManager) AcquireLock(ctx context.Context, lockID string) (string, error) {
	requestID := uuid.New().String()
	lock := DistributedLock{
		LockID:    lockID,
		RequestID: requestID,
		Timestamp: time.Now(),
	}

	err := lm.store.Create(ctx, "locks", lock)
	if err != nil {
		var existingLock DistributedLock
		if readErr := lm.store.Read(ctx, "locks", map[string]string{"lock_id": lockID}, &existingLock); readErr != nil {
			return "", readErr
		}
		if existingLock.RequestID > requestID {
			return "", nil
		}
	}
	return requestID, nil
}

func (lm *LockManager) ReleaseLock(ctx context.Context, lockID, requestID string) error {
	return lm.store.Delete(ctx, "locks", map[string]string{"lock_id": lockID, "request_id": requestID})
}
