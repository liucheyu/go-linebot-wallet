2022-08-01  
1. 基本上使用狀態模式完成  
2. 因為欄位取值很少，  
將database/sql的Scan,Next封裝成sqlRowSets
方便做GetString

markdown語法教學[網站](https://markdown.tw)

# go專案結構及lint bot sdk介紹
- [介紹網站](https://github.com/golang-standards/project-layout/blob/master/README_zh-TW.md)  
- [youtube](https://www.youtube.com/watch?v=oL6JBUk6tj0)
- [line bot go sdk & example](https://github.com/line/line-bot-sdk-go)

# git
git初始化
~~~
git init
git add .
git commit -m "first commit"
git remote add origin https://github.com/liucheyu/go-linebot-wallet.git
git push -u origin master
~~~

# GO指令
使用go module
~~~
go mod init github.com/liucheyu/go-linebot-wallet
~~~

go module其他指令
~~~
go env -w GO111MODULE=on
#下載全部
go mod download
#取得套件
go get [library]
#刪除沒用套件和補足遺失套件
go mod tidy 
#列出使用套件
go list -m all
~~~

go get指令
```
-d 只下载不安装
-f 只有在你包含了 -u 参数的时候才有效，不让 -u 去验证 import 中的每一个都已经获取了，这对于本地 fork 的包特别有用
-fix 在获取源码之后先运行 fix，然后再去做其他的事情
-t 同时也下载需要为运行测试所需要的包
-u 强制使用网络去更新包和它的依赖包
-v 显示执行的命令
```

# Gin
go gin install
```
go get -u github.com/gin-gonic/gin
import "github.com/gin-gonic/gin"
#引http用於回應status code
import "net/http"
```

go gin usage
from [gin](https://github.com/gin-gonic/gin)
```
func main() {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	router.GET("/someGet", getting)
	router.POST("/somePost", posting)
	router.PUT("/somePut", putting)
	router.DELETE("/someDelete", deleting)
	router.PATCH("/somePatch", patching)
	router.HEAD("/someHead", head)
	router.OPTIONS("/someOptions", options)


	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
	// router.Run(":3000") for a hard coded port
}
```

帶url參數
```
// This handler will match /user/john but will not match /user/ or /user
	router.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})
```

GET查詢參數
```
router.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")
		lastname := c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})
```

POST form
```
router.POST("/form_post", func(c *gin.Context) {
		message := c.PostForm("message")
		nick := c.DefaultPostForm("nick", "anonymous")

		c.JSON(200, gin.H{
			"status":  "posted",
			"message": message,
			"nick":    nick,
		})
	})
```

#line bot
line bot sdk [url](https://github.com/line/line-bot-sdk-go)
```
go get -u github.com/line/line-bot-sdk-go/v7/linebot
```

ling bot的ReplyMessage或PushMessage會返回Client，  
執行Do()才會實際返回，也就是說可以做其他事再返回或蒐集好幾個client一次或分次返回


# Redis
- git: https://github.com/go-redis/redis
- doc: https://redis.uptrace.dev/


# DB
github.com/go-sql-driver/mysql

