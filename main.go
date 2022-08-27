package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	solana "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/k0kubun/pp/v3"
)

type UniSwapQuoteResponse struct {
	BlockNumber                 string `json:"blockNumber"`
	Amount                      string `json:"amount"`
	AmountDecimals              string `json:"amountDecimals"`
	Quote                       string `json:"quote"`
	QuoteDecimals               string `json:"quoteDecimals"`
	QuoteGasAdjusted            string `json:"quoteGasAdjusted"`
	QuoteGasAdjustedDecimals    string `json:"quoteGasAdjustedDecimals"`
	GasUseEstimateQuote         string `json:"gasUseEstimateQuote"`
	GasUseEstimateQuoteDecimals string `json:"gasUseEstimateQuoteDecimals"`
	GasUseEstimate              string `json:"gasUseEstimate"`
	GasUseEstimateUSD           string `json:"gasUseEstimateUSD"`
	GasPriceWei                 string `json:"gasPriceWei"`
	Route                       [][]struct {
		Type    string `json:"type"`
		Address string `json:"address"`
		TokenIn struct {
			ChainID  uint32 `json:"chainId"`
			Decimals string `json:"decimals"`
			Address  string `json:"address"`
			Symbol   string `json:"symbol"`
		} `json:"tokenIn"`
		TokenOut struct {
			ChainID  uint32 `json:"chainId"`
			Decimals string `json:"decimals"`
			Address  string `json:"address"`
			Symbol   string `json:"symbol"`
		} `json:"tokenOut"`
		Fee          string `json:"fee"`
		Liquidity    string `json:"liquidity"`
		SqrtRatioX96 string `json:"sqrtRatioX96"`
		TickCurrent  string `json:"tickCurrent"`
		AmountIn     string `json:"amountIn,omitempty"`
		AmountOut    string `json:"amountOut,omitempty"`
	} `json:"route"`
	RouteString string `json:"routeString"`
	QuoteID     string `json:"quoteId"`
}

func getChangeBalanceSolana(client *rpc.Client, txSig solana.Signature) {
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

func getLatestTimestampBlockSolana(client *rpc.Client) {
	slot, err := client.GetSlot(
		context.TODO(),
		rpc.CommitmentFinalized,
	)
	if err != nil {
		log.Fatalln("get balance error", err)
	}
	block, err := client.GetBlock(
		context.Background(),
		uint64(slot),
	)
	pp.Println("slot of the latest block :", slot)
	pp.Println("timestamp : ", uint64(*block.BlockTime))
}

func callUniswapAPI() {
	uniSwapQuoteAPIURL := "https://api.uniswap.org/v1/quote"
	queries := []string{
		fmt.Sprintf("protocols=%s", "v3"),
		fmt.Sprintf("tokenInAddress=%s", "0x6B175474E89094C44Da98b954EedeAC495271d0F"),
		fmt.Sprintf("tokenOutAddress=%s", "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
		fmt.Sprintf("amount=%s", "5000000000000000000"),
		fmt.Sprintf("type=%s", "exactIn"),
		fmt.Sprintf("tokenInChainId=%s", "1"),
		fmt.Sprintf("tokenOutChainId=%s", "1"),
	}
	url := fmt.Sprintf("%s?%s", uniSwapQuoteAPIURL, strings.Join(queries, "&"))
	uniReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error :", err)
	}
	uniReq.Header.Set("authority", "api.uniswap.org")
	uniReq.Header.Set("origin", "https://api.uniswap.org")
	uniReq.Header.Set("referer", "https://api.uniswap.org")

	resp, err := http.DefaultClient.Do(uniReq)
	if err != nil {
		fmt.Println("Error :", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error :", err)
	}
	var uniSwapQuoteResponse UniSwapQuoteResponse
	if err = json.Unmarshal(body, &uniSwapQuoteResponse); err != nil {
		fmt.Println("Error :", err)
	}

	pp.Println(uniSwapQuoteResponse)
}

func main() {
	endpoint := rpc.MainNetBeta_RPC
	//Solana client
	client := rpc.New(endpoint)
	txSig := solana.MustSignatureFromBase58("5xQU3AXuA7qUCf9zLWaTA7TwQPw7Y4weqCNWjuffKLoKYP4QsXbPXUWoABWuFCrbJsPriRdzpZh2c9BpwfJd9w93")

	//get change Balance of an account in Solana
	getChangeBalanceSolana(client, txSig)
	//get the the timestamp of the latest block
	getLatestTimestampBlockSolana(client)
	//call uniswap api
	callUniswapAPI()
}
