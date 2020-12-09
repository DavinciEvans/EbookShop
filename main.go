package main

import (
	"EbookShop/ApiV1"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strconv"
)

/**
 * @Author: DavinciEvans
 * @Author: zhouningsiyuan@foxmail.com
 * @Date: 2020/11/6 11:32
 * @Desc:
 */

func main() {
	// 设置日志
	gin.DisableConsoleColor()
	f, err := os.Create("./logs/run.log")
	if err != nil {
		fmt.Println("Could not open log.")
		panic(err.Error())
	}
	//gin.SetMode(gin.DebugMode)
	gin.DefaultWriter = io.MultiWriter(f)

	gin.SetMode(gin.ReleaseMode)
	// 创建实例
	r := gin.Default()
	r.LoadHTMLFiles("html/index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", "")
	})

	// ApiV1
	v1 := r.Group("/api/v1")
	v1Book := v1.Group("/book")
	v1Book.GET("/", ApiV1.GetAllBooks)
	v1Book.GET("/:id", ApiV1.GetSingleBook)
	v1Book.POST("/", ApiV1.ValidateLogin(), ApiV1.ValidateAdmin(), ApiV1.PostNewBook)
	v1Book.PUT("/:id", ApiV1.ValidateLogin(), ApiV1.ValidateAdmin(), ApiV1.UpdateBook)
	v1Book.DELETE("/:id", ApiV1.ValidateLogin(), ApiV1.ValidateAdmin(), ApiV1.DeleteBook)

	v1Auth := v1.Group("/auth")
	v1Auth.POST("/:id", ApiV1.Login)
	v1Auth.POST("/", ApiV1.Register)
	v1Auth.PUT("/:id", ApiV1.ValidateLogin(), ApiV1.UpdateUserInfo)

	v1Categories := v1.Group("/categories")
	v1Categories.GET("/", ApiV1.GetCategoriesInfo)
	v1Categories.DELETE("/:id", ApiV1.ValidateLogin(), ApiV1.ValidateAdmin(), ApiV1.DeleteCategory)
	v1Categories.PUT("/:id", ApiV1.ValidateLogin(), ApiV1.ValidateAdmin(), ApiV1.UpdateCategory)
	v1Categories.POST("/", ApiV1.ValidateLogin(), ApiV1.ValidateAdmin(), ApiV1.PostNewCategory)

	V1Comments := v1.Group("/comments")
	V1Comments.GET("/", ApiV1.GetAllComments)
	V1Comments.GET("/:id", ApiV1.GetSingleComment)
	V1Comments.POST("/", ApiV1.ValidateLogin(), ApiV1.PostComment)
	V1Comments.DELETE("/:id", ApiV1.ValidateLogin(), ApiV1.DeleteComment)

	V1Carts := v1.Group("/carts")
	V1Carts.GET("/", ApiV1.ValidateLogin(), ApiV1.GetUsersCarts)
	V1Carts.DELETE("/:id", ApiV1.ValidateLogin(), ApiV1.DeleteCart)
	V1Carts.POST("/", ApiV1.ValidateLogin(), ApiV1.AddNewCart)

	V1Purchased := v1.Group("/purchased")
	V1Purchased.POST("/:id", ApiV1.ValidateLogin(), ApiV1.BuySingleBook)
	V1Purchased.GET("/", ApiV1.ValidateLogin(), ApiV1.GetAllPurchasedBooks)
	V1Purchased.GET("/:id", ApiV1.ValidateLogin(), ApiV1.GetSinglePurchasedBook)

	// 错误处理
	r.NoRoute(func(context *gin.Context) {
		context.JSON(http.StatusNotFound, gin.H{"Status": 404, "msg": "Not Found"})
	})

	// 程序入口
	err = r.Run(":" + strconv.Itoa(ApiV1.Config.NetPort))

	if err != nil {
		fmt.Println(err.Error())
	}
}
