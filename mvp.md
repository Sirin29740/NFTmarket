# 🟦 **模块 1：Solidity 智能合约（4 个功能）**

要写一个最简单的 **ERC721 NFT 合约**：

你必须有的内容：

### 1. `mint()`：铸造 NFT

接收：

- 接收者地址
- tokenURI（图片链接 + metadata）

### 2. `tokenURI()`

让前端能查到 NFT 对应的图片 metadata。

### 3. `ownerOf()`

后端/前端查 NFT 归属权。

### 4. `balanceOf()` / `tokenOfOwnerByIndex()`

用于展示用户的 NFT 列表。

📌 这个合约部署好就是整个项目的核心链上部分。



## ✅ **最简 ERC721 NFT 合约（完整可部署版本）**

> 使用 OpenZeppelin 标准库，安全、规范、最适合新手 + 产品上链。



```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract SimpleNFT is ERC721URIStorage, Ownable {
    uint256 private _tokenIds;

    constructor() ERC721("SimpleNFT", "SNFT") {}

    /// @notice 只允许合约拥有者铸造 NFT
    /// @param to 接收 NFT 的地址
    /// @param tokenURI NFT 对应的 metadata URI
    function mint(address to, string memory tokenURI) public onlyOwner returns (uint256) {
        _tokenIds += 1;
        uint256 newId = _tokenIds;

        _safeMint(to, newId);
        _setTokenURI(newId, tokenURI);

        return newId;
    }
}

```






------

# 🟩 **模块 2：Go 后端（Gin）需要实现的核心功能**

MVP only，需要写 5 个接口。

## **① 用户登录（简单版，不用数据库）**

实现一个最简单的登录方式，例如：

- 直接生成 JWT
- 或者 Web3 钱包登录（前端用 MetaMask）

MVP 建议：**随便一个假登录即可**。
 你重点是 NFT，而不是做复杂用户系统。

##  JWT 实现流程概览



1. **定义密钥：** 设置一个私密的签名密钥。
2. **定义 Claims：** 定义 JWT 中存储的用户信息（例如 `UserID`, `Username`）。
3. **生成 Token：** 登录/注册成功时，使用密钥签名生成 JWT。
4. **验证 Token：** 创建一个 Gin 中间件，用于解析和验证每个请求头中的 Token。



## 🛠️ 步骤 1：安装 JWT 库



我们需要使用 Go 社区最流行的 JWT 库：

Bash

```
go get github.com/golang-jwt/jwt/v5
```



## 步骤 2：创建 `auth` 包 (生成和解析)



我们创建一个 `internal/auth` 包来处理所有 JWT 逻辑。



### `internal/auth/jwt.go`



Go

```go
package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 假设这是你的私钥，生产环境中应从配置文件或环境变量中加载
var jwtSecret = []byte("YOUR_ULTRA_SECURE_AND_LONG_SECRET_KEY")

// UserClaims 结构体：定义了 JWT 中要存储的信息（Payload）
// 必须嵌入 jwt.RegisteredClaims
type UserClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 用于在用户登录/注册成功时生成 JWT
// 参数 userID 和 username 应该来自你的数据库模型
func GenerateToken(userID uint, username string) (string, error) {
	// 1. 设置 Token 的过期时间
	expirationTime := time.Now().Add(24 * time.Hour) // Token 24小时后过期

	// 2. 构造 Claims
	claims := &UserClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			// 可以添加 Issuer (iss), Subject (sub) 等
		},
	}

	// 3. 使用 HS256 签名方法创建 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 4. 使用密钥签名并获取完整的 Token 字符串
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken 用于解析和验证 JWT
func ValidateToken(tokenString string) (*UserClaims, error) {
	claims := &UserClaims{}

	// 解析 Token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 确保签名方法是预期的 (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenInvalid
	}

	// 返回解析出的 Claims
	return claims, nil
}
```



## 步骤 3：修改 Handler（集成生成 Token）



在你的 `Register` 和 `Login` Handler 中，在用户成功创建/验证后，调用 `GenerateToken`。



### `internal/api/user.go` (伪代码 - 注册/登录)



Go

```
package api

import (
	"net/http"
	"NFTmarket/internal/auth"      // 导入我们创建的 auth 包
	"NFTmarket/internal/model"     // 导入模型
	"github.com/gin-gonic/gin"
)

// 假设我们有一个安全的响应结构体
type AuthResponse struct {
	Token string             `json:"token"`
	User  model.UserResponse `json:"user"` // 使用安全的响应 DTO
}

func Register(c *gin.Context) {
	var user model.User // 数据库模型
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误"})
		return
	}
	
	// 1. **核心业务逻辑**：创建用户并写入数据库 (假设成功后 dbUser 获得 ID)
	// dbUser := service.CreateUser(&user)
	
	// 2. **生成 JWT**：使用数据库返回的 ID 和 Username
	// 假设 dbUser.ID 和 dbUser.Username 是从数据库返回的
    dbUser := model.User{ID: 101, Username: user.Username, Mail: user.Mail} // 模拟数据
	token, err := auth.GenerateToken(dbUser.ID, dbUser.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Token 生成失败"})
		return
	}

	// 3. 返回 JWT 和 User 信息
	safeUser := model.UserResponse{ID: dbUser.ID, Username: dbUser.Username, Mail: dbUser.Mail}
	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  safeUser,
	})
}

// Login 示例 (逻辑类似)
func Login(c *gin.Context) {
    // ... 1. 验证用户名/密码
    // ... 2. 从数据库加载用户模型 (dbUser)
    
	// 3. 生成 JWT
    dbUser := model.User{ID: 101, Username: "testUser"} // 模拟加载
	token, err := auth.GenerateToken(dbUser.ID, dbUser.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Token 生成失败"})
		return
	}
    
    // 4. 返回 Token
	safeUser := model.UserResponse{ID: dbUser.ID, Username: dbUser.Username}
	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  safeUser,
	})
}
```



## 步骤 4：创建 JWT 认证中间件 (验证 Token)



这个中间件会保护你的 Mint NFT 接口和其他需要登录的接口。



### `internal/api/middleware.go`



Go

```go
package api

import (
	"net/http"
	"strings"

	"NFTmarket/internal/auth"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware 是用于验证 JWT 的 Gin 中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从 Header 中获取 Token 字符串
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "请求头缺少 Authorization"})
			c.Abort() // 终止后续处理
			return
		}

		// 2. 检查格式：通常是 "Bearer [Token]"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Token 格式错误，应为 Bearer Token"})
			c.Abort()
			return
		}
		tokenString := parts[1]

		// 3. 验证 Token
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Token 无效或已过期"})
			c.Abort()
			return
		}

		// 4. 将用户信息（Claims）存储在 Gin Context 中，供后续 Handler 使用
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		
		// 继续处理请求
		c.Next()
	}
}
```



## 步骤 5：在路由中使用中间件



现在，你可以在你的 Mint NFT 路由上应用这个中间件。



### `cmd/main.go` (路由配置伪代码)



Go

```
package main

import (
	"NFTmarket/internal/api"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	
	// --- 公开路由 (无需认证) ---
	r.POST("/register", api.Register)
	r.POST("/login", api.Login)

	// --- 认证路由组 (需要 JWT) ---
	// 使用 AuthMiddleware() 保护 /api/v1 下的所有路由
	authGroup := r.Group("/api/v1")
	authGroup.Use(api.AuthMiddleware())
	{
		// 只有通过 JWT 认证的用户才能访问 Mint NFT 接口
		authGroup.POST("/nft/mint", api.MintNFT) 
		authGroup.GET("/user/profile", api.GetUserProfile)
	}

	return r
}

// MintNFT Handler 可以安全地从 Context 中获取用户 ID
func MintNFT(c *gin.Context) {
	// 从 Context 中取出中间件设置的 user_id
	userID, exists := c.Get("user_id")
	if !exists {
		// 不应该发生，因为中间件已经处理了，但出于安全仍需检查
		c.JSON(http.StatusInternalServerError, gin.H{"message": "无法获取用户ID"})
		return
	}
	
	// 现在你可以使用这个 userID 来执行铸造业务逻辑
	c.JSON(http.StatusOK, gin.H{"message": "铸造请求已接收", "user_id": userID})
}
```

这就是一个完整的、符合现代 API 规范的 JWT 认证实现。它将认证逻辑与业务逻辑清晰地分离，并使得 Go 后端保持**无状态**，非常适合你的 NFT 交易市场。

------

## **② 图片上传（IPFS 或本地存储）**

后端提供一个 API：

- 接受用户上传图片
- 保存到 `./uploads/xx.png` 或上传到 IPFS
- 生成一个 URL 返回给前端

结果要返回：

```
{ "image_url": "https://xxxx..." }
```

------

## **③ 生成 metadata（JSON）并上传**

NFT 需要 metadata，例如：

```
{
  "name": "MyNFT",
  "description": "...",
  "image": "https://image_url_here"
}
```

后端需要做：

- 根据上传的图片生成 metadata JSON
- 保存 JSON 到本地或 IPFS
- 返回 metadata_url（给 mint 用）

------

## **④ mint NFT（调用智能合约）**

Go 后端需要写一个接口：

- 接受前端传来的 metadata_url
- 使用 go-ethereum 调用合约的 `mint()`
- 返回 transaction hash

------

## **⑤ 查询用户 NFT 列表**

后端提供：

- 输入：用户钱包地址
- 输出：该用户所有 NFT 的 metadata_url + image_url

用于前端展示。

------

# 🟨 **模块 3：前端（最少页面即可）**

### 1. 首页（上传 + 预览）

用户在这里：

- 上传图片
- 点击 mint
- 等待 mint 成功

### 2. 我的 NFT 页面

功能：

- 查询当前用户钱包地址的 NFT 列表
- 显示图片

前端不需要复杂框架：
 Vue / React / Next.js / 甚至原生 HTML 都可以。

------

# 📦 **最终你需要做到的目录结构（MVP）**

```
nft-market/
│
├── contract/           # Solidity NFT 合约
│   └── MyNFT.sol
│
├── backend/            # Go 后端（Gin）
│   ├── main.go
│   ├── api/
│   │   ├── auth.go
│   │   ├── upload.go
│   │   ├── mint.go
│   │   └── nfts.go
│   ├── service/
│   ├── uploads/        # 存放图片
│   └── metadata/       # 存放 metadata JSON
│
└── frontend/           # 简单页面即可
    ├── index.html
    └── my-nfts.html
```

------

# 🚀 **这个就是一个 NFT 项目的 MVP 必须要做的事**

### 区块链部分（你要写）：

- 一个最简单的 ERC721 mint 合约（5 函数）

### 后端 Go（你要写）：

- 图片上传接口
- metadata 生成
- mint 的 RPC 调用
- 查询 NFT 的接口 -> 合约调用
- 用户简单登录（可选）

### 前端（你要写）：

- 上传图片按钮
- 预览
- 「Mint NFT」按钮
- 展示 NFT 列表的页面

------

# ❗ 你不需要写的（因为不是 MVP 范围）：

- 搜索过滤
- 排行榜
- 交易市场 buy/sell
- WebSocket 实时刷新
- IPFS pin 管理
- 后台系统
- 用户资料
- 评论
- 多链支持

这些以后可以加。