package contractcaller

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

// Caller ...
type Caller struct {
	l        *zap.SugaredLogger
	cli      *ethclient.Client
	cABI     abi.ABI
	cAddress common.Address
}

// NewContractCaller ...
func NewContractCaller(cABI abi.ABI, cli *ethclient.Client, contractAddress common.Address) *Caller {
	return &Caller{
		l:        zap.S(),
		cli:      cli,
		cABI:     cABI,
		cAddress: contractAddress,
	}
}

// ensureContext is a helper method to ensure a context is not nil, even if the
// user specified it as such.
func ensureContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.TODO()
	}
	return ctx
}

// Call ...
func (c *Caller) Call(opts *bind.CallOpts, methodName string, params ...interface{}) ([]interface{}, error) {
	input, err := c.cABI.Pack(methodName, params...)
	if err != nil {
		return nil, err
	}
	msg := ethereum.CallMsg{To: &c.cAddress, Data: input}
	ctx := ensureContext(opts.Context)
	outByte, err := c.cli.CallContract(ctx, msg, opts.BlockNumber)
	if err != nil {
		return nil, err
	}
	return c.cABI.Unpack(methodName, outByte)
}

// CallWithInput ...
func (c *Caller) CallWithInput(opts *bind.CallOpts, methodName string, input []byte) ([]interface{}, error) {
	msg := ethereum.CallMsg{To: &c.cAddress, Data: input}
	ctx := ensureContext(opts.Context)
	outByte, err := c.cli.CallContract(ctx, msg, opts.BlockNumber)
	if err != nil {
		return nil, err
	}
	return c.cABI.Unpack(methodName, outByte)
}
