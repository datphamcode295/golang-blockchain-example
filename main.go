package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	solana "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/k0kubun/pp/v3"
)

func main() {
	endpoint := rpc.MainNetBeta_RPC
	client := rpc.New(endpoint)
	txSig := solana.MustSignatureFromBase58("5xQU3AXuA7qUCf9zLWaTA7TwQPw7Y4weqCNWjuffKLoKYP4QsXbPXUWoABWuFCrbJsPriRdzpZh2c9BpwfJd9w93")

	//get change Balance of an account
	out, err := client.GetTransaction(
		context.TODO(),
		txSig,
		nil,
	)
	if err != nil {
		log.Fatalln("get balance error", err)
	}

	preToken, err := strconv.ParseInt(out.Meta.PreTokenBalances[1].UiTokenAmount.Amount, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	postToken, err := strconv.ParseInt(out.Meta.PostTokenBalances[1].UiTokenAmount.Amount, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	changeBalance := postToken - preToken
	pp.Println("Change of balance: ", changeBalance)

}
