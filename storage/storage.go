package storage

import (
	"database/sql"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // sqlite driver
	"go.uber.org/zap"
)

// Storage ...
type Storage struct {
	l  *zap.SugaredLogger
	db *sqlx.DB
}

// NewStorage init storage
func NewStorage(dbPath string) (*Storage, error) {
	db, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	s := &Storage{
		db: db,
		l:  zap.S(),
	}
	return s, s.initDB()
}

// initDB ...
func (s *Storage) initDB() error {
	schema := `
		CREATE TABLE IF NOT EXISTS "abis" (
			contract TEXT PRIMARY KEY,
			abi      TEXT NOT NULL
		);
	`
	if _, err := s.db.Exec(schema); err != nil {
		return err
	}
	return nil
}

// GetContractABI return abi of given contract in db
func (s *Storage) GetContractABI(contract ethereum.Address) (string, error) {
	var (
		query = `SELECT abi FROM "abis" WHERE contract=$1;`
		abi   string
	)
	queryX, err := s.db.Preparex(query)
	if err != nil {
		return "", err
	}
	if err := queryX.Get(&abi, contract.Hex()); err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return abi, nil
}

// StoreContractABI ...
func (s *Storage) StoreContractABI(contract ethereum.Address, abi string) error {
	var (
		query = `REPLACE INTO "abis" (contract, abi) VALUES ($1, $2);`
	)
	queryX, err := s.db.Preparex(query)
	if err != nil {
		return err
	}
	if _, err := queryX.Exec(contract.Hex(), abi); err != nil {
		return err
	}
	return nil
}
