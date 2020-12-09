package ApiV1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

/**
 * @Author: DavinciEvans
 * @Author: zhouningsiyuan@foxmail.com
 * @Date: 2020/12/6 0:08
 * @Desc:
 */

type commentForm struct {
	Content string
	Star    int
	BookID  uint
}

// /api/v1/comments GET 获取全部书评
func GetAllComments(c *gin.Context) {
	var comments []Comment
	const pageSize = 8 // 设定每页显示八个书评
	var count int64

	// 获取参数
	// page 参数，默认为 1
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}
	offset := (page - 1) * pageSize // 查询偏移量
	// book 参数，默认为 无
	book := c.DefaultQuery("book", "")
	bookPayload := ""
	if book == "" {
		book = "1" // 这个为任意的都可以，因为之后会被 payload 给覆盖掉
		bookPayload = " OR 1=1"
	}

	// 获得数量方便前端制作分页
	DB.Where("book_id = ?"+bookPayload, book).Find(&[]Comment{}).Count(&count)

	DB.Where("book_id = ?"+bookPayload, book).Offset(offset).Limit(pageSize).
		Find(&comments)

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "all comments count": count, "data": comments})
}

// /api/v1/comments/:id GET 获取单个书评
func GetSingleComment(c *gin.Context) {
	var comment Comment
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": "Illegal primary key. It must be a number."})
		return
	}
	result := DB.Find(&comment, id)
	if result.RowsAffected != 0 {
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": comment})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "msg": "This book is not exist."})
	}
}

// /api/v1/comments POST 发布书评
func PostComment(c *gin.Context) {

	var claims *userClaims
	claims = c.MustGet("claims").(*userClaims)
	loginId := claims.Uid

	var form commentForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": err.Error()})
		return
	}

	newComment := Comment{
		Content: form.Content,
		Star:    form.Star,
		BookID:  form.BookID,
		UserID:  loginId,
	}

	var commentBook Book
	result := DB.Find(&commentBook, form.BookID)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "data": "The book is not exist."})
		return
	} else {
		commentBook.StarSum += form.Star
	}

	result = DB.Create(&newComment)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "msg": "OK!"})
}

// /api/v1/comments/:id DELETE 删除书评（需要管理员权限或者用户本身）
func DeleteComment(c *gin.Context) {
	var comment Comment
	var commentBook Book
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": err.Error()})
		return
	}
	// 读取原评论相关信息
	result := DB.Find(&comment, id)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": "This items is not exist."})
		return
	}

	DB.Find(&commentBook, id)
	commentBook.StarSum -= comment.Star
	DB.Save(&commentBook)

	// 检查是否为管理员或者为该评论发布者
	var claims *userClaims
	claims = c.MustGet("claims").(*userClaims)
	loginId := claims.Uid
	if 1 != loginId {
		if loginId != comment.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "msg": "Unauthorized because is not admin."})
			return
		}
	}

	DB.Delete(&comment)

	c.JSON(http.StatusOK, gin.H{"status": "OK!"})
}
