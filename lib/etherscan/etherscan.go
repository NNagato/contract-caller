package etherscan

import (
	"fmt"
	"net/http"

	ethereum "github.com/ethereum/go-ethereum/common"

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
		baseAPI: "https://api.etherscan.io/",
		cli:     libhttp.NewRestClient(&http.Client{}),
	}
}

type etherscanResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

// GetContractABI ...
func (e *Etherscan) GetContractABI(contractAddress ethereum.Address) (string, error) {
	url := fmt.Sprintf("%s/api?module=contract&action=getabi&address=%s&apikey=%s",
		e.baseAPI, contractAddress.Hex(), e.apiKey)
	var resp etherscanResponse
	if err := e.cli.DoReq(url, http.MethodGet, nil, &resp); err != nil {
		return "", err
	}
	if resp.Status != "1" {
		return "", fmt.Errorf("error msg: %s", resp.Message)
	}
	return resp.Result, nil
}
