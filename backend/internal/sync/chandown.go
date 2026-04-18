package sync

import (
	"NFTmarket/internal/nft"
	"NFTmarket/internal/nft/contract"
	"context"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
)

type Task struct {
	TokenID uint64
	URI     string
}
type Manager struct {
	DB        *gorm.DB
	BaseDir   string
	TaskQueue chan Task // 改为结构体字段
}

func NewManager(db *gorm.DB) *Manager {
	return &Manager{
		DB:        db,
		BaseDir:   "./assets/nfts",
		TaskQueue: make(chan Task, 100), // 在这里初始化！
	}
}

var TaskQueue = make(chan Task, 100)

func (s *Manager) StartWorkerPool(ctx context.Context, workerCount int) {
	for i := 0; i < workerCount; i++ {
		go func(workerID int) {
			log.Println("start worker", workerID)
			for {
				select {
				case task := <-TaskQueue:
					s.processTask(ctx, task)
				case <-ctx.Done():
					return
				}
			}
		}(i)
	}
}
func (s *Manager) processTask(ctx context.Context, task Task) {
	var n nft.NFT
	if err := s.DB.WithContext(ctx).First(&n, "token_id = ?", task.TokenID).Error; err != nil {
		log.Printf("查询错误: %v", err)
		return
	}

	// 尝试下载
	err := s.downloadFile(ctx, task.URI, n.LocalPath)

	if err != nil {
		// 增加重试次数
		n.RetryCount++

		if n.RetryCount >= 5 {
			// 彻底失败，不再重试
			n.SyncStatus = 3 // 失败状态
			log.Printf("TokenID %d 同步彻底失败", task.TokenID)
		} else {
			// 依然标记为待处理，并利用时间间隔以后再入队
			n.SyncStatus = 0
			// 可以在这里写一个延时入队逻辑
			time.AfterFunc(time.Duration(n.RetryCount*10)*time.Second, func() {
				TaskQueue <- task
			})
		}
	} else {
		// 同步成功
		n.SyncStatus = 2
		fileName := filepath.Base(n.LocalPath)
		n.Image = "http://localhost:8080/static/" + fileName
	}

	s.DB.Save(&n)
}
func (s *Manager) downloadFile(ctx context.Context, url string, destPath string) error {
	// 确保父目录存在
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	// 创建带 Context 的请求
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// 创建文件
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {

		}
	}(out)

	// 使用 io.Copy 实现流式写入 (处理几百兆大图也不会撑爆内存)
	_, err = io.Copy(out, resp.Body)
	return err
}
func (s *Manager) WatchMintEvents(ctx context.Context, client *ethclient.Client, contractAddress common.Address) {
	// 1. 设置过滤器：只关心 from 是 0x00... 的 Transfer 事件
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics: [][]common.Hash{
			{crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))},
			{common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")},
		},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(ctx, query, logs)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("监视器已就绪，正在监听 Mint...")

	for {
		select {
		case vLog := <-logs:
			// 解析出 TokenID (vLog.Topics[3] 是 tokenId)
			tokenId := vLog.Topics[3].Big().Uint64()
			log.Printf("发现新 NFT: %d, 准备入队同步...", tokenId)

			owner := common.HexToAddress(vLog.Topics[2].Hex()).Hex()
			fileName := fmt.Sprintf("%d.png", tokenId)
			localPath := filepath.Join(s.BaseDir, fileName)

			newNFT := nft.NFT{
				TokenID:    tokenId,
				Owner:      owner,
				Name:       fmt.Sprintf("NFT #%d", tokenId),
				LocalPath:  localPath,
				SyncStatus: 0, // 待同步
				RetryCount: 0,
			}

			// 5. 入库 (使用 FirstOrCreate 防止重复)
			if err := s.DB.Where("token_id = ?", tokenId).FirstOrCreate(&newNFT).Error; err != nil {
				log.Printf("入库失败: %v", err)
				continue
			}

			log.Printf("NFT #%d 已入库，Owner: %s", tokenId, owner)
			//丢进任务队列
			uri := s.getChainURI(client, contractAddress, tokenId)
			s.TaskQueue <- Task{TokenID: tokenId, URI: uri}

		case err := <-sub.Err():
			log.Printf("监视连接中断，正在尝试重连...: %v", err)
			// 这里加入简单的重连逻辑
			time.Sleep(time.Second * 5)
		case <-ctx.Done():
			sub.Unsubscribe()
			return
		}
	}
}
func (s *Manager) getChainURI(client *ethclient.Client, addr common.Address, tid uint64) string {
	instance, err := contract.NewMyNft(addr, client)
	if err != nil {
		return ""
	}
	uri, err := instance.TokenURI(nil, big.NewInt(int64(tid)))
	if err != nil {
		return ""
	}
	return uri
}
