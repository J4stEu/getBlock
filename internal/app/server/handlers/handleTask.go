package handlers

import (
	"bytes"
	"encoding/json"
	internalError "github.com/J4stEu/getBlock/internal/app/errors"
	apiError "github.com/J4stEu/getBlock/internal/app/errors/api_errors"
	"io"
	"math/big"
	"net/http"
	"strings"
)

type GetLastBlockDataRaw struct {
	JsonRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      string        `json:"id"`
}
type GetLastBlockResponse struct {
	ID      string `json:"id"`
	JsonRPC string `json:"jsonrpc"`
	Result  string `json:"result"`
}

func GetLastBlock(APIkey string) (*GetLastBlockResponse, error) {

	//	curl --location --request POST 'https://eth.getblock.io/mainnet/' \
	//	--header 'x-api-key: YOUR-API-KEY' \
	//	--header 'Content-Type: application/json' \
	//	--data-raw '{"jsonrpc": "2.0",
	//	"method": "eth_blockNumber",
	//		"params": [],
	//	"id": "getblock.io"}'

	rawData := &GetLastBlockDataRaw{
		JsonRPC: "2.0",
		Method:  "eth_blockNumber",
		Params:  make([]interface{}, 0, 0),
		Id:      "getblock.io",
	}
	jsonPayload, err := json.Marshal(rawData)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "https://eth.getblock.io/mainnet/", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", APIkey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	respStruct := &GetLastBlockResponse{}
	if err = json.Unmarshal(bodyBytes, &respStruct); err != nil {
		return nil, err
	}
	return respStruct, nil
}

type GetBlockDataRaw struct {
	JsonRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      string        `json:"id"`
}

type Transactions struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	PublicKey        string `json:"publicKey"`
	R                string `json:"r"`
	Raw              string `json:"raw"`
	S                string `json:"s"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	V                string `json:"v"`
	Value            string `json:"value"`
}

type GetBlockResponse struct {
	ID      string `json:"id"`
	JsonRPC string `json:"jsonrpc"`
	Result  struct {
		Difficulty       string         `json:"difficulty"`
		ExtraData        string         `json:"extraData"`
		GasLimit         string         `json:"gasLimit"`
		GasUsed          string         `json:"gasUsed"`
		Hash             string         `json:"hash"`
		LogsBloom        string         `json:"logsBloom"`
		Miner            string         `json:"miner"`
		MixHash          string         `json:"mixHash"`
		Nonce            string         `json:"nonce"`
		Number           string         `json:"number"`
		ParentHash       string         `json:"parentHash"`
		ReceiptsRoot     string         `json:"receiptsRoot"`
		Sha3Uncles       string         `json:"sha3Uncles"`
		Size             string         `json:"size"`
		StateRoot        string         `json:"stateRoot"`
		Timestamp        string         `json:"timestamp"`
		TotalDifficulty  string         `json:"totalDifficulty"`
		Transactions     []Transactions `json:"transactions"`
		TransactionsRoot string         `json:"transactionsRoot"`
		Uncles           []interface{}  `json:"uncles"`
	} `json:"result"`
}

func GetBlock(APIkey, blockID string, byNumber bool) (*GetBlockResponse, error) {

	//	curl --location --request POST 'https://eth.getblock.io/mainnet/' \
	//	--header 'x-api-key: YOUR-API-KEY' \
	//	--header 'Content-Type: application/json' \
	//	--data-raw '{"jsonrpc": "2.0",
	//	"method": "eth_getBlockByNumber",
	//		"params": ["0x68B3", true],
	//	"id": "getblock.io"}'

	params := make([]interface{}, 0, 2)
	params = append(params, blockID)
	params = append(params, true)

	method := "eth_getBlockByHash"
	if byNumber {
		method = "eth_getBlockByNumber"
	}

	rawData := &GetBlockDataRaw{
		JsonRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      "getblock.io",
	}
	jsonPayload, err := json.Marshal(rawData)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "https://eth.getblock.io/mainnet/", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", APIkey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	respStruct := &GetBlockResponse{}
	if err = json.Unmarshal(bodyBytes, &respStruct); err != nil {
		return nil, err
	}
	return respStruct, nil
}

type PosNeg struct {
	positive uint64
	negative uint64
}

type BestDiff struct {
	address    string
	difference uint64
}

func BestTxDiff(blockTransactions []Transactions) (*BestDiff, error) {
	npTransactionsDiff := make(map[string]*PosNeg)
	for _, transaction := range blockTransactions {
		cleaned := strings.Replace(transaction.Gas, "0x", "", -1)
		gasValue, err := new(big.Int).SetString(cleaned, 16)
		if !err {
			return nil, internalError.SetError(internalError.ApiErrorLevel, apiError.BadValue)
		}
		cleaned = strings.Replace(transaction.Value, "0x", "", -1)
		value, err := new(big.Int).SetString(cleaned, 16)
		if !err {
			return nil, internalError.SetError(internalError.ApiErrorLevel, apiError.BadValue)
		}
		if _, ok := npTransactionsDiff[transaction.From]; !ok {
			npTransactionsDiff[transaction.From] = &PosNeg{
				positive: 0,
				negative: 0,
			}
		}
		if _, ok := npTransactionsDiff[transaction.To]; !ok {
			npTransactionsDiff[transaction.To] = &PosNeg{
				positive: 0,
				negative: 0,
			}
		}
		npTransactionsDiff[transaction.From].negative += gasValue.Uint64() + value.Uint64()
		npTransactionsDiff[transaction.To].positive += value.Uint64()
	}
	transactionsDiff := make(map[string]uint64)
	for key, transaction := range npTransactionsDiff {
		var actualDiff uint64
		if transaction.positive > transaction.negative {
			actualDiff = transaction.positive - transaction.negative
		} else {
			actualDiff = transaction.negative - transaction.positive
		}
		transactionsDiff[key] = actualDiff
	}
	maxDiff := &BestDiff{
		address:    "",
		difference: 0,
	}
	for key, value := range transactionsDiff {
		if value > maxDiff.difference {
			maxDiff.address = key
			maxDiff.difference = value
		}
	}
	return maxDiff, nil
}

func HandleTask(APIkey string) (string, error) {
	lastBlockInfo, err := GetLastBlock(APIkey)
	if err != nil {
		return "", internalError.SetError(internalError.ApiErrorLevel, apiError.LastBlockGetError)
	}
	blockID := lastBlockInfo.Result
	lastBlock, err := GetBlock(APIkey, blockID, true)
	if err != nil {
		return "", internalError.SetError(internalError.ApiErrorLevel, apiError.BlockGetError)
	}
	currentBestTxDiff, err := BestTxDiff(lastBlock.Result.Transactions)
	if err != nil {
		return "", err
	}
	blockID = lastBlock.Result.ParentHash
	for i := 0; i < 99; i++ {
		prevBlock, err := GetBlock(APIkey, blockID, false)
		if err != nil {
			return "", internalError.SetError(internalError.ApiErrorLevel, apiError.BlockGetError)
		}
		isBestTxDiff, err := BestTxDiff(prevBlock.Result.Transactions)
		if err != nil {
			return "", internalError.SetError(internalError.ApiErrorLevel, apiError.BestTxDiffError)
		}
		if isBestTxDiff.difference > currentBestTxDiff.difference {
			currentBestTxDiff = isBestTxDiff
		}
		blockID = prevBlock.Result.ParentHash
	}
	return currentBestTxDiff.address, nil
}
