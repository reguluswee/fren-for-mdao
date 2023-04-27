package event

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
)

const (
	chainId  = "513100"
	chainRpc = "https://rpc.etherfair.org"
)

const mdaoStartBlock = uint64(15782000)
const mdaoEndBlock = uint64(16639289)
const mdaoContractAddress = "0xcDfd138a8E59916E687F869f5c9D6B6f4334aE73"
const bigChainId = uint64(513100)
const frenAddress = "0x7127deeff734cE589beaD9C4edEFFc39C9128771"

var frenEventRewardMethod = []byte("MintClaimed(address,uint256)")
var frenEventRewardHash = crypto.Keccak256Hash(frenEventRewardMethod)

var diffAmount, _ = decimal.NewFromString("99999")
var dec, _ = decimal.NewFromString("1000000000000000000")

var multiple = new(big.Int).SetUint64(uint64(99))

type ClaimReward struct {
	Proxy  string
	Reward *big.Int
}

var startBlock = mdaoStartBlock

func cutLeftZeroToHex(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	r := hexutil.Encode(data)[2:]

	rs := strings.Split(r, "")
	var index int = -1
	for i, v := range rs {
		if v != "0" {
			index = i
			break
		}
	}
	if index > -1 {
		return "0x" + strings.Join(rs[index:], "")
	}
	return "0x0"
}

func BatchMdaoIssue() ([]MdaoData, []MdaoBlockError) {
	// startBlock = uint64(15827523 + 1)
	fixedObj := new(big.Int).Mul(diffAmount.BigInt(), dec.BigInt())

	rpcClient, rpcError := ethclient.Dial(chainRpc)
	if rpcError != nil {
		log.Fatal(rpcError)
	}

	contractAbi, _ := abi.JSON(strings.NewReader(FrenAbi))

	var resultList []MdaoData
	var errorList []MdaoBlockError

	for iBlock := startBlock; iBlock <= mdaoEndBlock; iBlock++ {
		block, err := rpcClient.BlockByNumber(context.Background(), new(big.Int).SetUint64(iBlock))
		if err != nil {
			errorList = append(errorList, createErr(iBlock, ""))
			log.Println(iBlock, err)
			continue
		}
		txs := block.Transactions()

		var wg sync.WaitGroup
		mutex := sync.RWMutex{}
		for _, tx := range txs {
			wg.Add(1)

			go func(tx *types.Transaction) {
				defer wg.Done()
				tx, _, err2 := rpcClient.TransactionByHash(context.Background(), tx.Hash())
				txReceipt, err1 := rpcClient.TransactionReceipt(context.Background(), tx.Hash())
				if err1 != nil || err2 != nil {
					errorList = append(errorList, createErr(iBlock, tx.Hash().Hex()))
					return
				}
				if tx.To() == nil { // contract create tx
					return
				}
				msg, _ := tx.AsMessage(types.LatestSignerForChainID(new(big.Int).SetUint64(bigChainId)), nil)

				if txReceipt.Status == 1 && tx.To().Hex() == mdaoContractAddress {
					var minterArray []ClaimReward
					for _, log := range txReceipt.Logs {
						eventHash := log.Topics[0]
						evss := log.Address.Hex()
						if eventHash.Hex() == frenEventRewardHash.Hex() && evss == frenAddress {
							//MintClaim event
							vb := log.Topics[1].Bytes()

							minterHash := hexutil.Encode(vb[len(vb)-20:])

							event := struct {
								RewardAmount *big.Int
							}{}
							err = contractAbi.UnpackIntoInterface(&event, "MintClaimed", log.Data)
							fmt.Println(minterHash, event, err)

							minterArray = append(minterArray, ClaimReward{
								Proxy:  minterHash,
								Reward: event.RewardAmount,
							})
						}
					}
					if len(minterArray) > 0 {
						bd := tx.Data()
						mintCount := hexutil.MustDecodeUint64(cutLeftZeroToHex(bd[4:36]))
						termSel := hexutil.MustDecodeUint64(cutLeftZeroToHex(bd[36:]))
						timeNow := time.Now()
						for i, _ := range minterArray {
							loss := new(big.Int).SetUint64(0)
							if termSel > 10 && minterArray[i].Reward.Cmp(fixedObj) == -1 {
								loss = new(big.Int).Mul(minterArray[i].Reward, multiple)
							}
							mutex.Lock()
							resultList = append(resultList, MdaoData{
								AddTime: timeNow,
								Block:   iBlock,
								Txhash:  tx.Hash().Hex(),
								Wallet:  msg.From().Hex(),
								Minter:  minterArray[i].Proxy,
								Round:   mintCount,
								Term:    termSel,
								Rewards: decimal.NewFromBigInt(minterArray[i].Reward, 0),
								Loss:    decimal.NewFromBigInt(loss, 0),
								Ts:      time.Unix(int64(block.Time()), 0),
							})
							mutex.Unlock()
						}
					}
				}
			}(tx)
			wg.Wait()
		}
	}
	return resultList, errorList
}

func createErr(block uint64, txhash string) MdaoBlockError {
	return MdaoBlockError{
		AddTime: time.Now(),
		Block:   block,
		Txhash:  txhash,
	}
}
