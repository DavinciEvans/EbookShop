package ApiV1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

/**
 * @Author: DavinciEvans
 * @Author: zhouningsiyuan@foxmail.com
 * @Date: 2020/12/5 15:21
 * @Desc:
 */

type categoryForm struct {
	Name string `form:"Name" binding:"required"`
}

// /api/v1/categories GET 获取全部的分类信息以及分类条目下的书本数量
func GetCategoriesInfo(c *gin.Context) {
	var id uint
	var name string
	var count int64
	var returnJSON []map[string]interface{}
	rows, err := DB.Raw(
		"SELECT categories.id, categories.name, count(category_id)" +
			" FROM categories LEFT JOIN books b" +
			" on categories.id = b.category_id " +
			"GROUP BY categories.id, categories.name;").Rows()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {

		err := rows.Scan(&id, &name, &count)
		if err != nil {
			continue
		}
		returnJSON = append(returnJSON, map[string]interface{}{
			"id":    id,
			"name":  name,
			"count": count,
		})
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": returnJSON})
}

// /api/v1/categories/:id DELETE 删除某个分类
func DeleteCategory(c *gin.Context) {
	// 管理员身份验证
	var claims *userClaims
	claims = c.MustGet("claims").(*userClaims)
	loginId := claims.Uid
	if 1 != loginId {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "msg": "Unauthorized because is wrong user"})
		return
	}

	// 删除数据
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": err.Error()})
		return
	}

	if id == 1 {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": "Can not delete first category"})
	}

	result := DB.Delete(&Category{}, id)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "msg": "OK!"})
}

// /api/v1/categories/:id PUT 更新某个分类
func UpdateCategory(c *gin.Context) {
	var form categoryForm
	var category Category

	// 管理员身份验证
	var claims *userClaims
	claims = c.MustGet("claims").(*userClaims)
	loginId := claims.Uid
	if 1 != loginId {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "msg": "Unauthorized because is wrong user"})
		return
	}

	// validator 验证
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": err.Error()})
		return
	}
	result := DB.Find(&category, id)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": result.Error.Error()})
		return
	}
	category.Name = form.Name
	DB.Save(&category)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "msg": "OK!"})
}

// /api/v1/categories/ POST 创建分类
func PostNewCategory(c *gin.Context) {
	var form categoryForm

	// 管理员身份验证
	var claims *userClaims
	claims = c.MustGet("claims").(*userClaims)
	loginId := claims.Uid
	if 1 != loginId {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "msg": "Unauthorized because is wrong user"})
		return
	}

	// validator 验证
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": err.Error()})
		return
	}

	// 创建新目录
	newCategory := Category{
		Name: form.Name,
	}
	result := DB.Create(&newCategory)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "msg": "OK!"})
}
