package v1

import (
	"NFTmarket/internal/nft"
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	shell "github.com/ipfs/go-ipfs-api"
)

// IpfsGateway 是您的 IPFS 网关地址，用于生成前端可访问的 URL
const IpfsGateway = "https://ipfs.io/ipfs/"

// IpfsApiUrl 是您本地 IPFS 节点的 API 地址
const IpfsApiUrl = "localhost:5001"

// initIPFSClient 初始化 IPFS 客户端
var sh = shell.NewShell(IpfsApiUrl)

func Upload(c *gin.Context) {

	name := c.PostForm("name")
	description := c.PostForm("description")
	if name == "" || description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and description cannot be empty"})
		return
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未找到名为 'image' 的文件"})
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法打开上传的文件"})
		return
	}
	defer file.Close() // 确保文件流在函数退出时关闭
	//上传文件到 IPFS
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
	imageUrl := IpfsGateway + cid

	metadata := nft.Metadata{
		Name:        name,
		Description: description,
		Image:       imageUrl,
	}
	metadataJson, _ := json.Marshal(metadata)

	metaCid, err := sh.Add(bytes.NewReader(metadataJson))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "metadata 上传失败"})
		return
	}
	metadataUrl := IpfsGateway + metaCid

	c.JSON(http.StatusOK, gin.H{
		"image_url":    imageUrl,
		"metadata_url": metadataUrl,
		"token_uri":    "ipfs://" + metaCid,
	})
}
