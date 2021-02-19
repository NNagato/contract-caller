package core

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/KyberNetwork/contract-caller/common"
	cc "github.com/KyberNetwork/contract-caller/lib/contract-caller"
	"github.com/KyberNetwork/contract-caller/lib/etherscan"
	"github.com/KyberNetwork/contract-caller/storage"
	"go.uber.org/zap"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

// Core ...
type Core struct {
	l    *zap.SugaredLogger
	esc  *etherscan.Etherscan
	ecli *ethclient.Client
	s    *storage.Storage
}

// NewCore ...
func NewCore(esc *etherscan.Etherscan, ecli *ethclient.Client, s *storage.Storage) *Core {
	return &Core{
		l:    zap.S(),
		esc:  esc,
		ecli: ecli,
		s:    s,
	}
}

// verifyContract ...
func (c *Core) verifyContract(contract ethereum.Address) error {
	code, err := c.ecli.CodeAt(context.Background(), contract, nil)
	if err != nil {
		return err
	}
	if len(code) == 0 {
		return errors.New("no code at given contract")
	}
	return nil
}

// getContractABIFromEtherscan ...
func (c *Core) getContractABIFromEtherscan(contract ethereum.Address) (string, error) {
	return c.esc.GetContractABI(contract)
}

// ContractMethods ...
func (c *Core) ContractMethods(contract ethereum.Address, contractABI string) ([]common.Method, error) {
	l := c.l.With("func", "core/ContractMethods", "contract", contract.Hex())
	if len(contractABI) == 0 {
		var (
			rawABI string
			err    error
		)
		rawABI, err = c.s.GetContractABI(contract)
		if err != nil {
			l.Errorw("cannnot get abi from storage", "err", err)
		}
		if len(rawABI) == 0 {
			rawABI, err = c.getContractABIFromEtherscan(contract)
			if err != nil {
				return nil, fmt.Errorf("cannot get contract ABI, err: %s", err.Error())
			}
		}
		contractABI = rawABI
	} else {
		if err := c.verifyContract(contract); err != nil {
			return nil, fmt.Errorf("cannot verify contract, err: %s", err.Error())
		}
	}
	cABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return nil, fmt.Errorf("cannot read abi, err: %s", err.Error())
	}
	if err := c.s.StoreContractABI(contract, contractABI); err != nil {
		l.Errorw("cannot store contract abi", "err", err)
	}
	var result []common.Method
	for name, detail := range cABI.Methods {
		if !detail.IsConstant() {
			continue // skip write function
		}
		var args []common.Argument
		for _, i := range detail.Inputs {
			args = append(args, common.Argument{
				Name: i.Name,
				Type: i.Type.String(),
			})
		}
		result = append(result, common.Method{
			Name:      name,
			Arguments: args,
		})
	}
	return result, nil
}

// CallContract ...
func (c *Core) CallContract(
	contract ethereum.Address,
	contractABI, methodName string,
	blockNumber string,
	params map[string]interface{}) (interface{}, error) {

	l := c.l.With("func", "core/CallContract", "contract", contract.Hex())
	if contractABI == "" {
		storedABI, err := c.s.GetContractABI(contract)
		if err != nil {
			l.Errorw("cannnot get abi from storage", "err", err)
		}
		contractABI = storedABI
	}
	cABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		l.Errorw("cannot read abi", "err", err)
		return nil, fmt.Errorf("cannot read abi, err: %s", err.Error())
	}
	method, ok := cABI.Methods[methodName]
	if !ok {
		l.Errorw("method is not available in this contract", "method", methodName)
		return nil, fmt.Errorf("method is not available in this contract, method = %s", methodName)
	}

	var input []interface{}
	for _, arg := range method.Inputs {
		var (
			ps string
		)
		p, ok := params[arg.Name]
		if ok {
			ps, ok = p.(string)
			if !ok {
				l.Errorw("wrong data type", "method", methodName, "arg name", arg.Name)
				return nil, fmt.Errorf("wrong data type, method = %s, arg name = %s", methodName, arg.Name)
			}
		}
		i, err := handleData(arg, ps)
		if err != nil {
			l.Errorw("cannot handle data", "method", methodName, "err", err)
			return nil, err
		}
		input = append(input, i)
	}
	var (
		bn   *big.Int
		errB error
	)
	if blockNumber != "" {
		if strings.Contains(blockNumber, "0x") {
			bn, errB = hexutil.DecodeBig(blockNumber)
			if errB != nil {
				l.Errorw("cannot handle block number", "err", errB)
				return nil, errB
			}
		} else {
			var ok bool
			bn, ok = big.NewInt(0).SetString(blockNumber, 10)
			if !ok {
				return nil, fmt.Errorf("wrong data type block number, input=%s", blockNumber)
			}
		}
	}
	caller := cc.NewContractCaller(cABI, c.ecli, contract)
	result, err := caller.Call(&bind.CallOpts{
		BlockNumber: bn,
	}, methodName, input...)
	if err != nil {
		l.Errorw("cannot get contract data", "err", err)
		return nil, fmt.Errorf("cannot get data from contract, err=%s", err)
	}
	return result, nil
}

func handleData(arg abi.Argument, ps string) (interface{}, error) {
	typeName := arg.Type.String()
	switch typeName {
	case "uint256", "int256":
		b, ok := big.NewInt(0).SetString(ps, 10)
		if !ok {
			return nil, fmt.Errorf("wrong data type, arg=%s, expected type=%s, actual value=%s",
				arg.Name, typeName, ps)
		}
		return b, nil
	case "uint256[]", "int256[]":
		ps = strings.ReplaceAll(ps, " ", "")
		nums := strings.Split(ps, ",")
		var bs []*big.Int
		for _, n := range nums {
			b, ok := big.NewInt(0).SetString(n, 10)
			if !ok {
				return nil, fmt.Errorf("wrong data type, arg=%s, expected type=%s, actual value=%s",
					arg.Name, typeName, ps)
			}
			bs = append(bs, b)
		}
		return bs, nil
	case "address[]":
		ps = strings.ReplaceAll(ps, " ", "")
		addresses := strings.Split(ps, ",")
		var as []ethereum.Address
		for _, a := range addresses {
			if !ethereum.IsHexAddress(a) {
				return nil, fmt.Errorf("wrong data type, arg=%s, expected type=%s, actual value=%s",
					arg.Name, typeName, ps)
			}
			as = append(as, ethereum.HexToAddress(a))
		}
		return as, nil
	case "address":
		if !ethereum.IsHexAddress(ps) {
			return nil, fmt.Errorf("wrong data type, arg=%s, expected type=%s, actual value=%s",
				arg.Name, typeName, ps)
		}
		return ethereum.HexToAddress(ps), nil
	case "bool":
		var b bool
		switch ps {
		case "false":
		case "true":
			b = true
		default:
			return nil, fmt.Errorf("wrong data type, arg=%s, expected type=%s, actual value=%s",
				arg.Name, typeName, ps)
		}
		return b, nil
	case "bytes", "int8", "bytes32":
		if strings.Contains(ps, "0x") {
			b, err := hexutil.Decode(ps)
			if err != nil {
				return nil, fmt.Errorf("wrong data type, arg=%s, expected type=%s, actual value=%s",
					arg.Name, typeName, ps)
			}
			return b, err
		}
		return []byte(ps), nil
	default:
		return nil, fmt.Errorf("wrong data type, arg=%s, expected type=%s, actual value=%s",
			arg.Name, typeName, ps)
	}
}
