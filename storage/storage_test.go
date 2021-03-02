package storage

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestWriteAndRead(t *testing.T) {
	s, err := NewStorage("db_test.db")
	require.NoError(t, err)
	var (
		contract = common.HexToAddress("0xbc5b5c036eb41a1a85af0b4da13d56420e8a0a92")
		abi      = "abi"
		newABI   = "newABI"
	)

	err = s.StoreContractABI(contract, abi)
	require.NoError(t, err)
	sABI, err := s.GetContractABI(contract)
	require.NoError(t, err)
	require.Equal(t, abi, sABI)

	err = s.StoreContractABI(contract, newABI)
	require.NoError(t, err)
	sABI, err = s.GetContractABI(contract)
	require.NoError(t, err)
	require.Equal(t, newABI, sABI)
}
