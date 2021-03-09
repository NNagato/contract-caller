package etherscan

import (
	"fmt"
	"net/http"

	ethereum "github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/contract-caller/common"
	libhttp "github.com/KyberNetwork/contract-caller/lib/http"
)

// Etherscan ...
type Etherscan struct {
	apiKey  string
	baseAPI string
	cli     *libhttp.RestClient
}

// NewEtherscan ...
func NewEtherscan(apiKey string) *Etherscan {
	return &Etherscan{
		apiKey:  apiKey,
		baseAPI: "https://api.etherscan.io",
		cli:     libhttp.NewRestClient(&http.Client{}),
	}
}

type etherscanResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

func (e *Etherscan) baseAPIURLFromNetwork(network string) string {
	switch network {
	case common.EthereumMainnet:
		return "https://api.etherscan.io"
	case common.EthereumRopsten:
		return "https://api-ropsten.etherscan.io"
	case common.EthereumKovan:
		return "https://api-kovan.etherscan.io"
	case common.BSCMainnet:
		return "https://api.bscscan.com"
	case common.BSCTestnet:
		return "https://api-testnet.bscscan.com"
	default:
		return e.baseAPI
	}
}

// GetContractABI ...
func (e *Etherscan) GetContractABI(contractAddress ethereum.Address, network string) (string, error) {
	url := fmt.Sprintf("%s/api?module=contract&action=getabi&address=%s&apikey=%s",
		e.baseAPIURLFromNetwork(network), contractAddress.Hex(), e.apiKey)
	fmt.Printf("\n\n\n %s", e.baseAPIURLFromNetwork(network))
	var resp etherscanResponse
	if err := e.cli.DoReq(url, http.MethodGet, nil, &resp); err != nil {
		return "", err
	}
	if resp.Status != "1" {
		return "", fmt.Errorf("error msg: %s", resp.Message)
	}
	return resp.Result, nil
}
