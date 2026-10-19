//go:build integration

package integration

import (
	"adweb3.local/platform/internal/campaigns"
	"adweb3.local/platform/internal/chain"
	"adweb3.local/platform/internal/store"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"os"
	"strings"
	"testing"
)

func artifact(t *testing.T, name string) (abi.ABI, []byte) {
	t.Helper()
	b, e := os.ReadFile("../../contracts/out/" + name + ".sol/" + name + ".json")
	if e != nil {
		t.Fatal(e)
	}
	var a struct {
		ABI      json.RawMessage `json:"abi"`
		Bytecode struct {
			Object string `json:"object"`
		} `json:"bytecode"`
	}
	if e = json.Unmarshal(b, &a); e != nil {
		t.Fatal(e)
	}
	parsed, e := abi.JSON(strings.NewReader(string(a.ABI)))
	if e != nil {
		t.Fatal(e)
	}
	return parsed, common.FromHex(a.Bytecode.Object)
}
func TestActualChainFundingAndRefund(t *testing.T) {
	ctx := context.Background()
	rpc := os.Getenv("EVM_RPC_URL")
	if rpc == "" {
		t.Fatal("EVM_RPC_URL required")
	}
	client, e := ethclient.Dial(rpc)
	if e != nil {
		t.Fatal(e)
	}
	defer client.Close()
	chainID, e := client.ChainID(ctx)
	if e != nil || chainID.Int64() != 31337 {
		t.Fatal("requires isolated local Anvil", e)
	}
	// Public Anvil development key, valid only for the explicitly checked local chain.
	key, e := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if e != nil {
		t.Fatal(e)
	}
	signer, e := bind.NewKeyedTransactorWithChainID(key, chainID)
	if e != nil {
		t.Fatal(e)
	}
	signer.Context = ctx
	tokenABI, tokenCode := artifact(t, "TestToken")
	asset, tx, tokenContract, e := bind.DeployContract(signer, tokenABI, tokenCode, client)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = bind.WaitMined(ctx, client, tx); e != nil {
		t.Fatal(e)
	}
	escrowABI, escrowCode := artifact(t, "CampaignEscrow")
	escrowAddress, tx, escrow, e := bind.DeployContract(signer, escrowABI, escrowCode, client, asset, common.HexToAddress("0x1111"), common.HexToAddress("0x2222"), common.HexToAddress("0x3333"))
	if e != nil {
		t.Fatal(e)
	}
	if _, e = bind.WaitMined(ctx, client, tx); e != nil {
		t.Fatal(e)
	}
	db := database(t)
	s := campaigns.Service{DB: db, ChainID: 31337, Token: asset.Hex()}
	header, e := client.HeaderByNumber(ctx, nil)
	if e != nil {
		t.Fatal(e)
	}
	o, e := s.CreateOrder(ctx, "chain-advertiser", store.UUID(), campaigns.Order{Payer: signer.From.Hex(), Amount: 25_000_000, Token: asset.Hex(), ChainID: 31337, Start: int64(header.Time), End: int64(header.Time) + 3600})
	if e != nil {
		t.Fatal(e)
	}
	id := crypto.Keccak256Hash([]byte(o.ID))
	terms := common.HexToHash(o.Terms)
	transact := func(c *bind.BoundContract, method string, args ...any) {
		t.Helper()
		tx, e := c.Transact(signer, method, args...)
		if e != nil {
			t.Fatal(method, e)
		}
		r, e := bind.WaitMined(ctx, client, tx)
		if e != nil || r.Status != 1 {
			t.Fatal(method, r, e)
		}
	}
	transact(escrow, "registerOrder", id, signer.From, big.NewInt(o.Amount), terms, uint64(o.End+86400))
	transact(tokenContract, "mint", signer.From, big.NewInt(o.Amount))
	transact(tokenContract, "approve", escrowAddress, big.NewInt(o.Amount))
	transact(escrow, "fundOrder", id, big.NewInt(o.Amount))
	unstable := chain.Indexer{Client: &changingRPC{Client: client, changeAt: 3}, DB: db, Campaigns: s, Escrow: escrowAddress}
	if e = unstable.Scan(ctx); e == nil {
		t.Fatal("mid-scan fork accepted")
	}
	before, e := s.Get(ctx, "chain-advertiser", o.ID)
	if e != nil || before.Funded != 0 {
		t.Fatal("orphan funding committed", before, e)
	}
	// Test-only recovery to canonical local chain after verifying the persistent quarantine.
	if _, e = db.Exec(ctx, "UPDATE chain_health SET quarantined=false,reason='' WHERE chain_id=31337"); e != nil {
		t.Fatal(e)
	}
	indexer := chain.Indexer{Client: client, DB: db, Campaigns: s, Escrow: escrowAddress, Confirmations: 0}
	if e = indexer.Scan(ctx); e != nil {
		t.Fatal("index actual funding", e)
	}
	got, e := s.Get(ctx, "chain-advertiser", o.ID)
	if e != nil || got.Funded != o.Amount {
		t.Fatal("real chain credit missing", got, e)
	}
	refundID, e := s.RequestRefund(ctx, "chain-advertiser", o.ID, store.UUID(), 15_000_000)
	if e != nil {
		t.Fatal(e)
	}
	_ = refundID
	worker := chain.Worker{DB: db, Client: client, Campaigns: s, Escrow: escrowAddress, Signer: signer, ABI: escrowABI, Confirmations: 0}
	// Simulate a crash after durable signing but before broadcast.
	prepared, e := worker.Prepare(ctx, "refund:"+refundID, "refundOrder", crypto.Keccak256Hash([]byte(refundID)), id, signer.From, big.NewInt(15_000_000))
	if e != nil {
		t.Fatal(e)
	}
	recovered, e := worker.Prepare(ctx, "refund:"+refundID, "refundOrder", crypto.Keccak256Hash([]byte(refundID)), id, signer.From, big.NewInt(15_000_000))
	if e != nil || prepared.Hash() != recovered.Hash() {
		t.Fatal("intent not durable", e)
	}
	if e = worker.ProcessOnce(ctx); e != nil {
		t.Fatal("worker refund", e)
	}
	if e = worker.ProcessOnce(ctx); e != nil {
		t.Fatal("worker restart refund", e)
	}
	// Registered-but-unfunded orders cannot starve later order pages.
	for n := 0; n < 101; n++ {
		if _, e = s.CreateOrder(ctx, "worker", store.UUID(), campaigns.Order{Payer: signer.From.Hex(), Amount: 100, Token: asset.Hex(), ChainID: 31337, Start: int64(header.Time), End: int64(header.Time) + 3600}); e != nil {
			t.Fatal(e)
		}
	}
	other := campaigns.Service{DB: db, ChainID: 43113, Token: asset.Hex()}
	if _, e = other.CreateOrder(ctx, "other", store.UUID(), campaigns.Order{Payer: signer.From.Hex(), Amount: 100, Token: asset.Hex(), ChainID: 43113, Start: 1, End: 2}); e != nil {
		t.Fatal(e)
	}
	for n := 0; n < 2; n++ {
		if e = worker.ProcessOnce(ctx); e != nil {
			t.Fatal("registration page", e)
		}
	}
	var count int
	if e = db.QueryRow(ctx, "SELECT count(*) FROM chain_registrations WHERE chain_id=31337").Scan(&count); e != nil || count != 101 {
		t.Fatal("registration starvation", count, e)
	}
	if e = db.QueryRow(ctx, "SELECT count(*) FROM chain_registrations WHERE chain_id=43113").Scan(&count); e != nil || count != 0 {
		t.Fatal("wrong-chain worker registration", count, e)
	}
	if e = indexer.Scan(ctx); e != nil {
		t.Fatal(e)
	}
	if e = indexer.Scan(ctx); e != nil {
		t.Fatal(e)
	}
	got, e = s.Get(ctx, "chain-advertiser", o.ID)
	if e != nil || got.Spent != 15_000_000 {
		t.Fatal("refund/replay mismatch", got, e)
	}
}

func TestIndexerRejectsDifferentEscrowToken(t *testing.T) {
	ctx := context.Background()
	client, e := ethclient.Dial(os.Getenv("EVM_RPC_URL"))
	if e != nil {
		t.Fatal(e)
	}
	defer client.Close()
	key, _ := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	signer, _ := bind.NewKeyedTransactorWithChainID(key, big.NewInt(31337))
	a, code := artifact(t, "TestToken")
	asset, tx, _, e := bind.DeployContract(signer, a, code, client)
	if e != nil {
		t.Fatal(e)
	}
	bind.WaitMined(ctx, client, tx)
	a, code = artifact(t, "CampaignEscrow")
	escrow, tx, _, e := bind.DeployContract(signer, a, code, client, asset, signer.From, signer.From, signer.From)
	if e != nil {
		t.Fatal(e)
	}
	bind.WaitMined(ctx, client, tx)
	db := database(t)
	s := campaigns.Service{DB: db, ChainID: 31337, Token: token}
	indexer := chain.Indexer{Client: client, DB: db, Campaigns: s, Escrow: escrow}
	if e = indexer.Scan(ctx); e == nil {
		t.Fatal("escrow token mismatch accepted")
	}
}

type changingRPC struct {
	*ethclient.Client
	reads    int
	changeAt int
}

func (r *changingRPC) HeaderByNumber(ctx context.Context, n *big.Int) (*types.Header, error) {
	h, e := r.Client.HeaderByNumber(ctx, n)
	if e != nil {
		return nil, e
	}
	if n != nil {
		r.reads++
		if r.reads >= r.changeAt {
			copy := *h
			copy.Extra = []byte("replacement-fork")
			return &copy, nil
		}
	}
	return h, nil
}
func TestCheckpointReorgQuarantinesRefunds(t *testing.T) {
	ctx := context.Background()
	client, e := ethclient.Dial(os.Getenv("EVM_RPC_URL"))
	if e != nil {
		t.Fatal(e)
	}
	defer client.Close()
	key, _ := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	signer, _ := bind.NewKeyedTransactorWithChainID(key, big.NewInt(31337))
	a, code := artifact(t, "TestToken")
	asset, tx, _, e := bind.DeployContract(signer, a, code, client)
	if e != nil {
		t.Fatal(e)
	}
	bind.WaitMined(ctx, client, tx)
	a, code = artifact(t, "CampaignEscrow")
	escrow, tx, _, e := bind.DeployContract(signer, a, code, client, asset, signer.From, signer.From, signer.From)
	if e != nil {
		t.Fatal(e)
	}
	bind.WaitMined(ctx, client, tx)
	db := database(t)
	s := campaigns.Service{DB: db, ChainID: 31337, Token: asset.Hex()}
	o, e := s.CreateOrder(ctx, "a", store.UUID(), campaigns.Order{Payer: payer, Amount: 100, Token: asset.Hex(), ChainID: 31337, Start: 1, End: 2})
	if e != nil {
		t.Fatal(e)
	}
	if e = s.ApplyFinalizedEvent(ctx, campaigns.ChainEvent{ID: store.UUID(), OrderID: o.ID, Kind: "fund", Payer: payer, Token: asset.Hex(), ChainID: 31337, Amount: 100, Final: true, TxHash: "fixture", Block: 1}); e != nil {
		t.Fatal(e)
	}
	head, _ := client.HeaderByNumber(ctx, nil)
	if _, e = db.Exec(ctx, "INSERT INTO chain_checkpoints VALUES($1,$2,$3,$4)", 31337, escrow.Hex(), head.Number.Int64(), head.Hash().Hex()); e != nil {
		t.Fatal(e)
	}
	indexer := chain.Indexer{Client: &changingRPC{Client: client, changeAt: 1}, DB: db, Campaigns: s, Escrow: escrow}
	if e = indexer.Scan(ctx); e == nil {
		t.Fatal("reorg not rejected")
	}
	if _, e = s.RequestRefund(ctx, "a", o.ID, store.UUID(), 50); e == nil {
		t.Fatal("quarantined refund accepted")
	}
	var quarantined bool
	if e = db.QueryRow(ctx, "SELECT quarantined FROM chain_health WHERE chain_id=31337").Scan(&quarantined); e != nil || !quarantined {
		t.Fatal("missing persistent quarantine", e)
	}
}

func workerFixture(t *testing.T, confirmations uint64) (chain.Worker, common.Address) {
	t.Helper()
	ctx := context.Background()
	client, e := ethclient.Dial(os.Getenv("EVM_RPC_URL"))
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(client.Close)
	key, _ := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	signer, _ := bind.NewKeyedTransactorWithChainID(key, big.NewInt(31337))
	a, code := artifact(t, "TestToken")
	asset, tx, _, e := bind.DeployContract(signer, a, code, client)
	if e != nil {
		t.Fatal(e)
	}
	bind.WaitMined(ctx, client, tx)
	a, code = artifact(t, "CampaignEscrow")
	escrow, tx, _, e := bind.DeployContract(signer, a, code, client, asset, signer.From, signer.From, signer.From)
	if e != nil {
		t.Fatal(e)
	}
	bind.WaitMined(ctx, client, tx)
	db := database(t)
	return chain.Worker{DB: db, Client: client, Campaigns: campaigns.Service{DB: db, ChainID: 31337, Token: asset.Hex()}, Escrow: escrow, Signer: signer, ABI: a, Confirmations: confirmations}, asset
}
func TestExpiredOrdersDoNotStarveRegistration(t *testing.T) {
	ctx := context.Background()
	w, asset := workerFixture(t, 0)
	for n := 0; n < 101; n++ {
		if _, e := w.Campaigns.CreateOrder(ctx, "expired", store.UUID(), campaigns.Order{Payer: w.Signer.From.Hex(), Amount: 100, Token: asset.Hex(), ChainID: 31337, Start: 0, End: 1}); e != nil {
			t.Fatal(e)
		}
	}
	head, _ := w.Client.HeaderByNumber(ctx, nil)
	o, e := w.Campaigns.CreateOrder(ctx, "live", store.UUID(), campaigns.Order{Payer: w.Signer.From.Hex(), Amount: 100, Token: asset.Hex(), ChainID: 31337, Start: int64(head.Time), End: int64(head.Time) + 3600})
	if e != nil {
		t.Fatal(e)
	}
	// Deterministically place the live fixture after all 101 expired UUIDs.
	for index := 0; index < 101; index++ {
		_, e = w.DB.Exec(ctx, "UPDATE orders SET id=$1 WHERE id=(SELECT id FROM orders WHERE advertiser='expired' AND id::text NOT LIKE '00000000-%' LIMIT 1)", fmt.Sprintf("00000000-0000-4000-8000-%012d", index))
		if e != nil {
			t.Fatal(e)
		}
	}
	oldID := o.ID
	o.ID = "ffffffff-ffff-4fff-8fff-ffffffffffff"
	if _, e = w.DB.Exec(ctx, "UPDATE orders SET id=$1,chain_order_id=$2 WHERE id=$3", o.ID, crypto.Keccak256Hash([]byte(o.ID)).Hex(), oldID); e != nil {
		t.Fatal(e)
	}
	if e = w.ProcessOnce(ctx); e != nil {
		t.Fatal(e)
	}
	var count int
	if e = w.DB.QueryRow(ctx, "SELECT count(*) FROM chain_registrations WHERE order_id=$1", o.ID).Scan(&count); e != nil || count != 1 {
		t.Fatal("expired rows starved live registration", count, e)
	}
}
func TestMinedRegistrationReorgAndRestart(t *testing.T) {
	ctx := context.Background()
	w, asset := workerFixture(t, 2)
	for n := 0; n < 3; n++ {
		if e := w.Client.Client().CallContext(ctx, nil, "evm_mine"); e != nil {
			t.Fatal(e)
		}
	}
	head, _ := w.Client.HeaderByNumber(ctx, nil)
	o, e := w.Campaigns.CreateOrder(ctx, "live", store.UUID(), campaigns.Order{Payer: w.Signer.From.Hex(), Amount: 100, Token: asset.Hex(), ChainID: 31337, Start: int64(head.Time), End: int64(head.Time) + 3600})
	if e != nil {
		t.Fatal(e)
	}
	var snapshot string
	if e = w.Client.Client().CallContext(ctx, &snapshot, "evm_snapshot"); e != nil {
		t.Fatal(e)
	}
	if e = w.ProcessOnce(ctx); e != nil {
		t.Fatal(e)
	}
	var count int
	if e = w.DB.QueryRow(ctx, "SELECT count(*) FROM chain_registrations WHERE order_id=$1", o.ID).Scan(&count); e != nil || count != 0 {
		t.Fatal("unfinal registration became permanent", count, e)
	}
	var reverted bool
	if e = w.Client.Client().CallContext(ctx, &reverted, "evm_revert", snapshot); e != nil || !reverted {
		t.Fatal(e)
	}
	if e = w.ProcessOnce(ctx); e != nil {
		t.Fatal("restart recovery", e)
	}
	for n := 0; n < 3; n++ {
		if e = w.Client.Client().CallContext(ctx, nil, "evm_mine"); e != nil {
			t.Fatal(e)
		}
	}
	if e = w.ProcessOnce(ctx); e != nil {
		t.Fatal("final recovery", e)
	}
	if e = w.DB.QueryRow(ctx, "SELECT count(*) FROM chain_registrations WHERE order_id=$1", o.ID).Scan(&count); e != nil || count != 1 {
		t.Fatal("final registration missing", count, e)
	}
}
