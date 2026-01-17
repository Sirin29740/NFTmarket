package v1

import (
	"github.com/gin-gonic/gin"
	shell "github.com/ipfs/go-ipfs-api"
	"log"
	"net/http"
)

// IpfsGateway 是您的 IPFS 网关地址，用于生成前端可访问的 URL
const IpfsGateway = "https://dweb.link/ipfs/"

// IpfsApiUrl 是您本地 IPFS 节点的 API 地址
const IpfsApiUrl = "localhost:5001"

// initIPFSClient 初始化 IPFS 客户端
var sh = shell.NewShell(IpfsApiUrl)

func Upload(c *gin.Context) {
	// 1. 获取上传的文件
	// "image" 必须与前端 FormData.append('image', file) 中的字段名匹配
	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未找到名为 'image' 的文件"})
		return
	}

	// 2. 打开文件流
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法打开上传的文件"})
		return
	}
	defer file.Close() // 确保文件流在函数退出时关闭

	// 3. 上传文件到 IPFS
	// sh.Add 方法接受一个 io.Reader 接口，可以直接传入文件流
	cid, err := sh.Add(file)
	if err != nil {
		log.Printf("IPFS 上传失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "上传到 IPFS 失败",
			"details": err.Error(),
		})
		return
	}

	log.Printf("文件上传成功，CID: %s", cid)

	// 4. 构造可访问的 URL
	imageUrl := IpfsGateway + cid

	// 5. 返回前端所需的 JSON 格式
	c.JSON(http.StatusOK, gin.H{
		"image_url": imageUrl,
	})
}
