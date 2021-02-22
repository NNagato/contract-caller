package server

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/contract-caller/core"
)

// Server ...
type Server struct {
	sugar *zap.SugaredLogger
	host  string
	r     *gin.Engine
	core  *core.Core
}

// NewServer ...
func NewServer(host string, core *core.Core) *Server {

	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.MaxAge = 5 * time.Minute
	r.Use(cors.New(corsConfig))

	return &Server{
		sugar: zap.S(),
		host:  host,
		r:     r,
		core:  core,
	}
}

type inputMethods struct {
	Contract    string `json:"contract" binding:"required"`
	ABI         string `json:"abi"`
	RememberABI bool   `json:"rememberABI"`
}

func (s *Server) methods(c *gin.Context) {
	var input inputMethods
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{
				"err": err.Error(),
			},
		)
		return
	}
	if !ethereum.IsHexAddress(input.Contract) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"err": "contract is not a valid ethereum address",
			},
		)
		return
	}
	result, err := s.core.ContractMethods(ethereum.HexToAddress(input.Contract), input.ABI, input.RememberABI)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{
				"err": err.Error(),
			},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"data": result,
		},
	)
}

// inputCall ...
type inputCall struct {
	Contract    string                 `json:"contract" binding:"required"`
	ABI         string                 `json:"abi"`
	Method      string                 `json:"method" binding:"required"`
	BlockNumber string                 `json:"blockNumber"`
	Params      map[string]interface{} `json:"params"`
}

func (s *Server) call(c *gin.Context) {
	var input inputCall
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{
				"err": err.Error(),
			},
		)
		return
	}
	if !ethereum.IsHexAddress(input.Contract) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"err": "contract is not a valid ethereum address",
			},
		)
		return
	}
	result, err := s.core.CallContract(ethereum.HexToAddress(input.Contract), input.ABI,
		input.Method, input.BlockNumber, input.Params)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{
				"err": err.Error(),
			},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"data": result,
		},
	)
}

func (s *Server) register() {
	g := s.r.Group("contract")
	g.POST("/methods", s.methods)
	g.POST("/call", s.call)
}

// Run ...
func (s *Server) Run() error {
	s.register()
	return s.r.Run(s.host)
}
