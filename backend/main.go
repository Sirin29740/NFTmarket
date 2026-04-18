package main

import (
	"NFTmarket/api"
	"NFTmarket/internal/database"
	"NFTmarket/internal/nft"
	"NFTmarket/internal/sync"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

func main() {
	// 1. 初始化数据库
	database.InitDB()

	// 2. 初始化客户端 (wss 模式)
	client, err := nft.GetClient()
	if err != nil {
		log.Fatal("无法连接至区块链:", err)
	}

	// 3. 初始化同步管理器
	manager := sync.NewManager(database.GetDB())

	// 4. 创建全局控制信号
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 5. 启动后台逻辑 (Worker + Watcher)
	manager.StartWorkerPool(ctx, 5)
	go manager.WatchMintEvents(ctx, client, common.HexToAddress("0x5368c8a780B05f4A5aa267c54503b84Ed688EbE9"))

	// 6. 启动 API 服务 (放在协程里)
	// 这样它就不会阻塞后面信号监听的运行
	go func() {
		log.Println("API 服务启动中...")
		api.Router()
	}()

	// 7. 阻塞主线程，监听退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Println("系统已就绪。按 Ctrl+C 退出。")
	<-quit

	log.Println("正在关闭服务...")
	// 执行 cancel 会通知所有使用了 ctx 的协程（Worker 和 Watcher）安全退出
	cancel()

	// 给一点时间让协程处理完手头的工作
	time.Sleep(2 * time.Second)
	log.Println("服务已优雅退出")
}
