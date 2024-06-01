package graph

import (
	"fmt"
	"os"
	"time"

	chain "github.com/wetee-dao/go-sdk"
	"github.com/wetee-dao/go-sdk/gen/system"
	"github.com/wetee-dao/indexer/store"
	"github.com/wetee-dao/indexer/util"
)

var (
	ChainClient     *chain.ChainClient
	DefaultChainUrl string = "ws://wetee-node.worker-addon.svc.cluster.local:9944"
)

func SubEvent() error {
	var url = DefaultChainUrl
	if u, exists := os.LookupEnv("CHAIN_URI"); exists {
		url = u
	}

	// 挖矿开始
mintStart:
	client, err := chain.ClientInit(url, true)
	if err != nil {
		return err
	}
	ChainClient = client
	chainAPI := client.Api

	for {
		header, err := chainAPI.RPC.Chain.GetHeaderLatest()
		if err != nil {
			util.LogWithRed("GetHeaderLatest", err)
			// 失败后等待10秒重新尝试
			// Wait 10 seconds to try again
			time.Sleep(time.Second * 10)
			continue
		}

		storeBlock, err := store.GetChainBlock()
		if err != nil {
			util.LogWithRed("GetChainBlock", err)
			// 失败后等待10秒重新尝试
			// Wait 10 seconds to try again
			time.Sleep(time.Second * 10)
			continue
		}

		if uint32(storeBlock) > uint32(header.Number) {
			store.SetChainBlock(uint64(header.Number))
			storeBlock = uint64(header.Number)
		}

		if uint32(storeBlock) == uint32(header.Number) {
			// 已经到达最新的区块
			fmt.Println("已经到达最新的区块")
			break
		}

		err = ExpEvent(storeBlock)
		if err != nil {
			util.LogWithRed("ExpEvent", err)
			// 失败后等待10秒重新尝试
			// Wait 10 seconds to try again
			time.Sleep(time.Second * 10)
			continue
		}

		err = store.SetChainBlock(storeBlock + 1)
		if err != nil {
			util.LogWithRed("SetChainBlock", err)
			// 失败后等待10秒重新尝试
			// Wait 10 seconds to try again
			time.Sleep(time.Second * 10)
			continue
		}
	}

	sub, err := chainAPI.RPC.Chain.SubscribeNewHeads()
	if err != nil {
		util.LogWithRed("SubscribeNewHeads", err)
		// 失败后等待10秒重新尝试
		// Wait 10 seconds to try again
		time.Sleep(time.Second * 10)
		goto mintStart
	}

	defer sub.Unsubscribe()
	for {
		head := <-sub.Chan()
		util.LogWithRed("Chain is at block: #", fmt.Sprint(head.Number))

		err = ExpEvent(uint64(head.Number))
		if err != nil {
			util.LogWithRed("ExpEvent", err)
			// 失败后等待10秒重新尝试
			// Wait 10 seconds to try again
			time.Sleep(time.Second * 10)
			continue
		}

		err = store.SetChainBlock(uint64(head.Number))
		if err != nil {
			util.LogWithRed("SetChainBlock", err)
			// 失败后等待10秒重新尝试
			// Wait 10 seconds to try again
			time.Sleep(time.Second * 10)
			continue
		}
	}
	return nil
}

func ExpEvent(number uint64) error {
	chainAPI := ChainClient.Api

	hash, err := chainAPI.RPC.Chain.GetBlockHash(number)
	if err != nil {
		util.LogWithRed("GetBlockHash", err)
		return err
	}

	events, err := system.GetEvents(chainAPI.RPC.State, hash)
	if err != nil {
		util.LogWithRed("GetEventsLatest", err)
		return err
	}

	for _, event := range events {
		e := event.Event
		if e.IsWeteeWorker {
			if e.AsWeteeWorkerField0.IsWorkRuning {
				fmt.Println("程序启动")
			}
			if e.AsWeteeWorkerField0.IsWorkStoped {
				fmt.Println("程序停止")
			}
			if e.AsWeteeWorkerField0.IsWorkContractUpdated {
				fmt.Println("程序上传工作量证明")
			}
		}
		if e.IsContracts {
			if e.AsContractsField0.IsInstantiated {
				fmt.Println("合约部署成功")
			}
			if e.AsContractsField0.IsCodeStored {
				fmt.Println("合约代码上传成功")
			}
			if e.AsContractsField0.IsContractCodeUpdated {
				fmt.Println("合约代码更新成功")
			}
		}
	}

	return nil
}
