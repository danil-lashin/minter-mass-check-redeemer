package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/MinterTeam/minter-go-node/core/check"
	tx "github.com/MinterTeam/minter-go-node/core/transaction"
	"github.com/MinterTeam/minter-go-node/core/types"
	"github.com/MinterTeam/minter-go-node/crypto"
	"github.com/MinterTeam/minter-go-node/crypto/sha3"
	"github.com/MinterTeam/minter-go-node/rlp"
	"github.com/btcsuite/btcd/btcec"
	"github.com/valyala/fasthttp"
	"math/big"
	"os"
	"strconv"
)

var (
	coin = types.GetBaseCoin()
	senderPkBytes, _ = hex.DecodeString(os.Args[1])
	senderPrivateKey, _ = crypto.ToECDSA(senderPkBytes)
	passphrasePkBytes = sha256.Sum256([]byte("password"))
	passphrasePk, _ = crypto.ToECDSA(passphrasePkBytes[:])

	checkValue = big.NewInt(1)
)

func main() {
	t, err := strconv.Atoi(os.Args[2])

	if err != nil {
		panic(err)
	}

	for j := t; j < t + 10000; j++ {
		createAndRedeemCheck(uint64(j))
	}

	println("All Done!")
}

func createAndRedeemCheck(i uint64) {
	receiverPrivateKey, _ := ecdsa.GenerateKey(btcec.S256(), rand.Reader)
	receiverAddr := crypto.PubkeyToAddress(receiverPrivateKey.PublicKey)

	c := check.Check{
		Nonce:    i,
		DueBlock: 1000000,
		Coin:     coin,
		Value:    checkValue,
	}

	lock, _ := crypto.Sign(c.HashWithoutLock().Bytes(), passphrasePk)
	c.Lock = big.NewInt(0).SetBytes(lock)

	c.Sign(senderPrivateKey)

	rawCheck, _ := rlp.EncodeToBytes(c)

	var senderAddressHash types.Hash
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, []interface{}{receiverAddr})
	hw.Sum(senderAddressHash[:0])

	sig, _ := crypto.Sign(senderAddressHash.Bytes(), passphrasePk)

	proof := [65]byte{}
	copy(proof[:], sig)

	data := tx.RedeemCheckData{
		RawCheck: rawCheck,
		Proof:    proof,
	}

	encodedData, _ := rlp.EncodeToBytes(data)

	transaction := tx.Transaction{
		Nonce:         1,
		GasPrice:      big.NewInt(1),
		GasCoin:       coin,
		Type:          tx.TypeRedeemCheck,
		Data:          encodedData,
		SignatureType: tx.SigTypeSingle,
	}

	transaction.Sign(receiverPrivateKey)

	encodedTx, _ := rlp.EncodeToBytes(transaction)
	req := "http://" + os.Args[3] + "/broadcast_tx_sync?tx=0x" + fmt.Sprintf("%x", encodedTx)
	_, body, _ := fasthttp.Get(nil, req)
	println(i, string(body))
}