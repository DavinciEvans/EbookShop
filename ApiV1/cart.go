package ApiV1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

/**
 * @Author: DavinciEvans
 * @Author: zhouningsiyuan@foxmail.com
 * @Date: 2020/12/5 21:31
 * @Desc:
 */

type cartForm struct {
	BookID uint
}

// /api/v1/cart GET 获得当前登录用户的所有购物车
func GetUsersCarts(c *gin.Context) {
	var claims *userClaims
	var carts []Cart
	var returnJSON []map[string]interface{}
	claims = c.MustGet("claims").(*userClaims)
	loginId := claims.Uid

	DB.Where("user_id = ?", loginId).Find(&carts)

	for _, cart := range carts {
		book := Book{}
		DB.Find(&book, cart.BookID)
		temp := map[string]interface{}{
			"id":     book.ID,
			"name":   book.Name,
			"price":  book.Price,
			"author": book.Author,
		}
		returnJSON = append(returnJSON, temp)
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "msg": "success", "data": returnJSON})
}

// /api/v1/cart/:id DELETE 删除某一个购物车条目
func DeleteCart(c *gin.Context) {
	var claims *userClaims
	var cart Cart
	claims = c.MustGet("claims").(*userClaims)
	loginId := claims.Uid

	// 解析路由
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": err.Error()})
		return
	}

	// 获取购物车信息
	result := DB.Where("book_id = ?", id).Find(&cart)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": "This items is not exist."})
		return
	}

	if cart.UserID != loginId {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "msg": "this is not your cart."})
		return
	}

	DB.Delete(&cart)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "msg": "OK!"})
}

// /api/v1/cart/:id PUT 在当前用户下增加一条购物车条目
func AddNewCart(c *gin.Context) {
	var claims *userClaims
	var form cartForm
	claims = c.MustGet("claims").(*userClaims)
	loginId := claims.Uid

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": err.Error()})
		return
	}
	newCart := Cart{
		BookID: form.BookID,
		UserID: loginId,
	}

	result := DB.Create(&newCart)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": "This items is already exist!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "msg": "OK!"})
}
