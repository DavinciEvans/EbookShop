package ApiV1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

/**
 * @Author: DavinciEvans
 * @Author: zhouningsiyuan@foxmail.com
 * @Date: 2020/12/6 13:35
 * @Desc:
 */

// /api/v1/purchased/:id POST 购买某一本书
func BuySingleBook(c *gin.Context) {
	var claims *userClaims
	claims = c.MustGet("claims").(*userClaims)
	loginId := claims.Uid

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": "Illegal primary key. It must be a number."})
		return
	}

	puchasedBook := PurchasedBook{
		UserID: loginId,
		BookID: uint(id),
		Time:   time.Now(),
	}
	result := DB.Create(&puchasedBook)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": "You have already had this book!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "msg": "OK!"})

}

// /api/v1/purchased GET 获得所有已购书籍
func GetAllPurchasedBooks(c *gin.Context) {
	var claims *userClaims
	var purchasedBook []PurchasedBook
	var returnJSON []map[string]interface{}
	claims = c.MustGet("claims").(*userClaims)
	loginId := claims.Uid

	DB.Where("user_id = ?", loginId).Find(&purchasedBook)

	for _, info := range purchasedBook {
		var book Book
		bookID := info.BookID
		DB.Find(&book, bookID)
		temp := map[string]interface{}{
			"id":     bookID,
			"name":   book.Name,
			"author": book.Author,
		}
		returnJSON = append(returnJSON, temp)
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": returnJSON})
}

// /api/v1/purchased/:id GET 获得一本已购书籍
func GetSinglePurchasedBook(c *gin.Context) {
	var claims *userClaims
	var purchasedBook PurchasedBook
	var book Book
	claims = c.MustGet("claims").(*userClaims)
	loginId := claims.Uid

	bookId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": "Illegal primary key. It must be a number."})
		return
	}

	DB.Where("user_id = ? and book_id = ?", loginId, bookId).Find(&purchasedBook)

	if purchasedBook.BookID == 0 || purchasedBook.UserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "msg": "You have not had this book yet."})
		return
	}

	DB.Find(&book, bookId)
	returnBook := map[string]interface{}{
		"id":      bookId,
		"name":    book.Name,
		"author":  book.Author,
		"content": book.Content,
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": returnBook})
}
