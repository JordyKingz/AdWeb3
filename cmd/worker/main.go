package main

import (
	"adweb3.local/platform/internal/campaigns"
	"adweb3.local/platform/internal/chain"
	"adweb3.local/platform/internal/store"
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	db, e := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if e != nil {
		log.Fatal(e)
	}
	defer db.Close()
	if e = store.Migrate(ctx, db); e != nil {
		log.Fatal(e)
	}
	client, e := ethclient.DialContext(ctx, os.Getenv("EVM_RPC_URL"))
	if e != nil {
		log.Fatal(e)
	}
	defer client.Close()
	id, e := client.ChainID(ctx)
	if e != nil {
		log.Fatal(e)
	}
	configured, e := strconv.ParseInt(os.Getenv("CHAIN_ID"), 10, 64)
	if e != nil || configured != id.Int64() || (configured != 31337 && configured != 43113) {
		log.Fatal("only configured local/Fuji testnet enabled")
	}
	key, e := crypto.HexToECDSA(os.Getenv("RELAYER_PRIVATE_KEY"))
	if e != nil {
		log.Fatal("invalid relayer key")
	}
	signer, e := bind.NewKeyedTransactorWithChainID(key, id)
	if e != nil {
		log.Fatal(e)
	}
	token, escrow := os.Getenv("TOKEN_ADDRESS"), os.Getenv("ESCROW_ADDRESS")
	if !common.IsHexAddress(token) || !common.IsHexAddress(escrow) {
		log.Fatal("invalid token/escrow")
	}
	abi, e := chain.ContractABI("contracts/out/CampaignEscrow.sol/CampaignEscrow.json")
	if e != nil {
		log.Fatal(e)
	}
	confirmations := uint64(2)
	if configured == 31337 {
		confirmations = 0
	}
	w := chain.Worker{DB: db, Client: client, Campaigns: campaigns.Service{DB: db, ChainID: configured, Token: token}, Escrow: common.HexToAddress(escrow), Signer: signer, ABI: abi, Confirmations: confirmations}
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		if e = w.ProcessOnce(ctx); e != nil && ctx.Err() == nil {
			log.Printf("worker requires attention: %v", e)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}
