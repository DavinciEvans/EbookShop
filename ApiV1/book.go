package ApiV1

/**
 * @Author: DavinciEvans
 * @Author: zhouningsiyuan@foxmail.com
 * @Date: 2020/11/6 11:31
 * @Desc:
 */

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type bookForm struct {
	Name     string  `form:"Name" binding:"required"`
	Price    float32 `form:"Price" binding:"required"`
	Author   string  `form:"Author" binding:"required"`
	Content  string  `form:"Content" binding:"required"`
	Category uint    `form:"Category" binding:"required"`
}

// /api/v1/book GET 大量获取书籍信息（不包含书本内容）
func GetAllBooks(c *gin.Context) {
	var books []Book
	var returnJSON []map[string]interface{}
	const pageSize = 8 // 设定每页显示八本书
	var count int64

	// 获取参数
	// page 参数，默认为 1
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}
	offset := (page - 1) * pageSize // 查询偏移量
	// category 参数，默认为 无
	category := c.DefaultQuery("category", "")
	categoryPayload := ""
	if category == "" {
		category = "1" // 这个为任意的都可以，因为之后会被 payload 给覆盖掉
		categoryPayload = " OR 1=1"
	}
	// hot 参数，默认为 1 升序，-1 降序
	hot, err := strconv.Atoi(c.DefaultQuery("hot", ""))
	hotPayload := ""
	if err != nil {
		hotPayload = ""
	}
	if hot > 0 {
		hotPayload = "star_sum/pay_number desc"
	} else if hot == 0 {
		hotPayload = "id" // 即默认按照 id 来排序
	} else {
		hotPayload = "star_sum/pay_number asc"
	}

	// 获得数量方便前端制作分页
	DB.Where("category_id = ?"+categoryPayload, category).Find(&[]Book{}).Count(&count)

	// 查询
	condition := []string{"ID", "Name", "Price", "Author", "StarSum", "PayNumber", "CategoryID"}
	DB.Select(condition).
		Where("category_id = ?"+categoryPayload, category).
		Order(hotPayload).
		Offset(offset).Limit(pageSize).
		Find(&books)

	for _, book := range books {
		temp := map[string]interface{}{
			"ID":        book.ID,
			"Name":      book.Name,
			"Price":     book.Price,
			"Author":    book.Author,
			"Star":      book.StarSum / book.PayNumber,
			"PayNumber": book.PayNumber,
			"Category":  book.CategoryID,
		}
		returnJSON = append(returnJSON, temp)
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "all books count": count, "data": returnJSON})
}

// /api/v1/book/:id GET 获取对应主键的书籍信息（不包含书本内容）
func GetSingleBook(c *gin.Context) {
	var book Book
	returnJSON := make(map[string]interface{})
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": "Illegal primary key. It must be a number."})
		return
	}
	condition := []string{"ID", "Name", "Price", "Author", "Star", "PayNumber", "CategoryID"}
	result := DB.Select(condition).Find(&book, id)
	returnJSON["ID"] = book.ID
	returnJSON["name"] = book.Name
	returnJSON["Price"] = book.Price
	returnJSON["Author"] = book.Author
	returnJSON["Star"] = book.StarSum / book.PayNumber
	returnJSON["PayNumber"] = book.PayNumber
	returnJSON["Category"] = book.CategoryID
	if result.RowsAffected != 0 {
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": returnJSON})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "msg": "This book is not exist."})
	}
}

// /api/v1/book POST 创建新书 （需要管理员权限）
func PostNewBook(c *gin.Context) {
	var form bookForm
	// validator 验证
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": err.Error()})
		return
	}
	newBook := Book{
		Name:       form.Name,
		Author:     form.Author,
		Price:      form.Price,
		Content:    form.Content,
		StarSum:    0,
		PayNumber:  0,
		CategoryID: form.Category,
	}
	result := DB.Create(&newBook)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "msg": "OK!"})
}

// /api/v1/book/:id PUT 更新书信息 （需要管理员权限）
func UpdateBook(c *gin.Context) {
	var form bookForm
	var book Book

	// validator 验证
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": err.Error()})
		return
	}
	if id, err := strconv.Atoi(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": err.Error()})
		return
	} else {
		result := DB.Find(&book, id)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": result.Error.Error()})
			return
		}
		book.Name = form.Name
		book.Author = form.Author
		book.Price = form.Price
		book.Content = form.Content
		book.CategoryID = form.Category
		DB.Save(&book)
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "msg": "OK!"})
	}
}

// /api/v1/book/:id DELETE 删除书本 （需要管理员权限）
func DeleteBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": err.Error()})
		return
	}
	result := DB.Delete(&Book{}, id)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK!"})
}
