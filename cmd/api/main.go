package main

import (
	"adweb3.local/platform/internal/campaigns"
	"adweb3.local/platform/internal/store"
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"os/signal"
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
	var tokens map[string]string
	if e = json.Unmarshal([]byte(os.Getenv("API_TOKENS")), &tokens); e != nil || len(tokens) == 0 {
		log.Fatal("API_TOKENS must map private tokens to advertisers")
	}
	for key, actor := range tokens {
		if len(key) < 32 || actor == "" {
			log.Fatal("invalid bootstrap credential")
		}
	}
	token := os.Getenv("TOKEN_ADDRESS")
	if !common.IsHexAddress(token) || common.HexToAddress(token) == (common.Address{}) {
		log.Fatal("invalid test token address")
	}
	chain := int64(31337)
	if os.Getenv("CHAIN_ID") == "43113" {
		chain = 43113
	} else if v := os.Getenv("CHAIN_ID"); v != "" && v != "31337" {
		log.Fatal("only local/Fuji testnet enabled")
	}
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	server := campaigns.Server(addr, campaigns.Handler(campaigns.Service{DB: db, ChainID: chain, Token: token}, tokens))
	go func() {
		<-ctx.Done()
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(c)
	}()
	log.Printf("testnet API listening on %s", addr)
	if e = server.ListenAndServe(); e != nil && ctx.Err() == nil {
		log.Fatal(e)
	}
}
