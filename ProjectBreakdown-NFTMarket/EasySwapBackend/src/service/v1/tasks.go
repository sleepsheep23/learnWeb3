package service

import (
	"context"
	"fmt"
	"github.com/ProjectsTask/EasySwapBackend/src/service/svc"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"time"
)

func StartPriceMonitor(svcCtx *svc.ServerCtx) {
	ticker := time.NewTicker(5 * time.Minute)
	ctx := context.Background()
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			CheckUserCoinsPriceAlert(ctx, svcCtx)
		case <-ctx.Done():
			return
		}
	}
}

func CheckUserCoinsPriceAlert(ctx context.Context, svcCtx *svc.ServerCtx) error {
	coinAlerts, err := svcCtx.Dao.ListAllCoinsAlerts()
	if err != nil {
		return err
	}

	coinIdDedup := make(map[string]interface{})
	coinIdList := make([]string, 0)
	for _, alert := range coinAlerts {
		if _, ok := coinIdDedup[alert.CoinID]; ok {
			continue
		} else {
			coinIdList = append(coinIdList, alert.CoinID)
			coinIdDedup[alert.CoinID] = struct{}{}
		}
	}

	coinPrice, err := GetCoinsPrice(ctx, svcCtx, coinIdList, "usd")
	for _, alert := range coinAlerts {
		price, ok := coinPrice[alert.CoinID]
		if !ok {
			continue
		}
		if alert.AlertType == "below" { // 价格低于
			if price.Usd < alert.TargetPrice {
				fmt.Printf("alert user %d coin %s price %.4f below %.4f\n", alert.UserID, alert.CoinID, price.Usd, alert.TargetPrice)
			}
		} else if alert.AlertType == "above" { // 价格高于
			if price.Usd > alert.TargetPrice {
				fmt.Printf("alert user %d coin %s price %.4f above %.4f\n", alert.UserID, alert.CoinID, price.Usd, alert.TargetPrice)
			}
		} else if alert.AlertType == "between" {
			if price.Usd > alert.PriceRangeMin && price.Usd < alert.PriceRangeMax {
				fmt.Printf("alert user %d coin %s price %.4f between %.4f and %.4f\n", alert.UserID, alert.CoinID, price.Usd, alert.PriceRangeMin, alert.PriceRangeMax)
			}
		}
	}

	return nil
}

func StartBlockMonitor(svcCtx *svc.ServerCtx) {
	// 连接以太坊 WebSocket 节点
	client, err := ethclient.Dial(svcCtx.C.Alchemy.EthereumWss)
	if err != nil {
		fmt.Printf("Failed to connect to Ethereum WS: %v", err)
	}
	defer client.Close()

	// 创建订阅通道
	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		fmt.Printf("Failed to subscribe to new head: %v", err)
		return
	}

	fmt.Println("✅ Listening for new Ethereum blocks...")

	for {
		select {
		case err := <-sub.Err():
			fmt.Printf("Subscription error: %v", err)
			return
		case header := <-headers:
			// 收到新区块头
			fmt.Printf("⛓️  New Block: #%v | Hash: %s\n", header.Number.String(), header.Hash().Hex())

			err := svcCtx.Dao.AddEthereumBlock(header)
			if err != nil {
				fmt.Printf("Failed to save block header: %v", err)
			}

			block, err := client.BlockByNumber(context.Background(), header.Number)
			if err != nil {
				fmt.Printf("Failed to get full block: %v", err)
				continue
			}

			addressList, err := svcCtx.Dao.GetWatchAddressList()
			if err != nil {
				fmt.Printf("Failed to get watch address list: %v", err)
				continue
			}
			if len(addressList) == 0 {
				continue
			}
			for _, transaction := range block.Transactions() {
				chainID := big.NewInt(svcCtx.C.Alchemy.ChainID)
				signer := types.LatestSignerForChainID(chainID)
				addrSender, err := types.Sender(signer, transaction)
				if err != nil {
					fmt.Printf("Failed to get transaction sender: %v", err)
					continue
				}
				if _, ok := addressList[addrSender.String()]; ok {
					fmt.Printf("🚨 Alert! Transaction from watched address %s in block #%v | Tx Hash: %s\n", addrSender.Hex(), header.Number.String(), transaction.Hash().Hex())
				}

				addrTo := transaction.To()
				if addrTo != nil {
					if _, ok := addressList[addrTo.String()]; ok {
						fmt.Printf("🚨 Alert! Transaction to watched address %s in block #%v | Tx Hash: %s\n", addrTo.Hex(), header.Number.String(), transaction.Hash().Hex())
					}
				}
			}
		}
	}
}
