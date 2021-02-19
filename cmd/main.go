package main

import (
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli"
	"go.uber.org/zap"

	"github.com/KyberNetwork/contract-caller/core"
	"github.com/KyberNetwork/contract-caller/lib/etherscan"
	"github.com/KyberNetwork/contract-caller/server"
	"github.com/KyberNetwork/contract-caller/storage"
)

var (
	sugar = zap.S()

	hostHTTPFlag        = "host"
	defaultHost         = "localhost:3001"
	etherscanAPIKeyFlag = "etherscan-apikey"
	nodeFlag            = "node"
	dbPathFlag          = "db-path"
	defaultDBPath       = "contract.db"
)

func main() {
	app := cli.NewApp()
	app.Name = "cex-assistant"
	app.Version = "0.0.1"
	app.Usage = "easy interface to call contract"
	app.Action = run

	app.Flags = append(app.Flags, cli.StringFlag{
		Name:   hostHTTPFlag,
		Usage:  "host",
		EnvVar: "HOST",
		Value:  defaultHost,
	}, cli.StringFlag{
		Name:   etherscanAPIKeyFlag,
		Usage:  "etherscan API key",
		EnvVar: "ETHERSCAN_APIKEY",
	}, cli.StringFlag{
		Name:   nodeFlag,
		Usage:  "ethereum node",
		EnvVar: "NODE",
	}, cli.StringFlag{
		Name:   dbPathFlag,
		Usage:  "db path",
		Value:  defaultDBPath,
		EnvVar: "DB_PATH",
	},
	)

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	sugar = logger.Sugar()
	defer func() {
		_ = sugar.Sync()
	}()
	if err := app.Run(os.Args); err != nil {
		sugar.Errorw("app error", "err", err)
	}
}

func run(c *cli.Context) error {
	esc := etherscan.NewEtherscan(c.String(etherscanAPIKeyFlag))
	ecli, err := ethclient.Dial(c.String(nodeFlag))
	if err != nil {
		return err
	}
	str, err := storage.NewStorage(c.String(dbPathFlag))
	if err != nil {
		return err
	}
	s := server.NewServer(c.String(hostHTTPFlag), core.NewCore(esc, ecli, str))
	return s.Run()
}
