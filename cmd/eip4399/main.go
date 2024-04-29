package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	txfuzz "github.com/waterfall-network/tx-fuzz"
	"gitlab.waterfall.network/waterfall/protocol/gwat/common"
	"gitlab.waterfall.network/waterfall/protocol/gwat/core/types"
	"gitlab.waterfall.network/waterfall/protocol/gwat/crypto"
	"gitlab.waterfall.network/waterfall/protocol/gwat/ethclient"
	"gitlab.waterfall.network/waterfall/protocol/gwat/rpc"
)

var (
	address = "http://127.0.0.1:8545"
)

func main() {
	cl, sk := getRealBackend()
	backend := ethclient.NewClient(cl)
	sender := common.HexToAddress(txfuzz.ADDR)
	nonce, err := backend.PendingNonceAt(context.Background(), sender)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Nonce: %v\n", nonce)
	gp, _ := backend.SuggestGasPrice(context.Background())
	tx := types.NewContractCreation(nonce, common.Big1, 500000, gp, []byte{0x44, 0x44, 0x55})
	signedTx, _ := types.SignTx(tx, types.NewLondonSigner(big.NewInt(1337702)), sk)
	backend.SendTransaction(context.Background(), signedTx)
}

func getRealBackend() (*rpc.Client, *ecdsa.PrivateKey) {
	// eth.sendTransaction({from:personal.listAccounts[0], to:"0xb02A2EdA1b317FBd16760128836B0Ac59B560e9D", value: "100000000000000"})

	sk := crypto.ToECDSAUnsafe(common.FromHex(txfuzz.SK))
	if crypto.PubkeyToAddress(sk.PublicKey).Hex() != txfuzz.ADDR {
		panic(fmt.Sprintf("wrong address want %s got %s", crypto.PubkeyToAddress(sk.PublicKey).Hex(), txfuzz.ADDR))
	}
	cl, err := rpc.Dial(address)
	if err != nil {
		panic(err)
	}
	return cl, sk
}

func sendTx(sk *ecdsa.PrivateKey, backend *ethclient.Client, to common.Address, value *big.Int) {
	sender := common.HexToAddress(txfuzz.ADDR)
	nonce, err := backend.PendingNonceAt(context.Background(), sender)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Nonce: %v\n", nonce)
	gp, _ := backend.SuggestGasPrice(context.Background())
	tx := types.NewTransaction(nonce, to, value, 500000, gp, nil)
	signedTx, _ := types.SignTx(tx, types.HomesteadSigner{}, sk)
	backend.SendTransaction(context.Background(), signedTx)
}
