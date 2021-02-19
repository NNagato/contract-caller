package storage

import (
	"github.com/boltdb/bolt"
	ethereum "github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

const (
	contractBucket = "contract"
)

// Storage ...
type Storage struct {
	sugar *zap.SugaredLogger
	db    *bolt.DB
}

// NewStorage ...
func NewStorage(path string) (*Storage, error) {
	// init instance
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	// init buckets
	if err = db.Update(func(tx *bolt.Tx) error {
		if _, cErr := tx.CreateBucketIfNotExists([]byte(contractBucket)); cErr != nil {
			return cErr
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &Storage{
		sugar: zap.S(),
		db:    db,
	}, nil
}

// StoreContractABI ...
func (s *Storage) StoreContractABI(contract ethereum.Address, abi string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(contractBucket))
		return b.Put(contract.Bytes(), []byte(abi))
	})
}

// GetContractABI ...
func (s *Storage) GetContractABI(contract ethereum.Address) (string, error) {
	var result string
	if err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(contractBucket))
		v := b.Get(contract.Bytes())
		result = string(v)
		return nil
	}); err != nil {
		return "", err
	}
	return result, nil
}
