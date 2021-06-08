package ApiV1

/**
 * @Author: DavinciEvans
 * @Author: zhouningsiyuan@foxmail.com
 * @Date: 2020/11/17 20:12
 * @Desc:
 */

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type userClaims struct {
	jwt.StandardClaims
	Uid uint `json:"uid"`
}

type userLoginForm struct {
	Username string `form:"Username" binding:"required"`
	Password string `form:"Password" binding:"required,min=6,max=20"`
}

type UserForm struct {
	Username string `form:"Username" binding:"required"`
	Password string `form:"Password" binding:"required"`
	Name     string `form:"Name" binding:"required"`
	Email    string `form:"Email" binding:"required"`
}

// 创建 JWT Token，传入密钥和 Uid，返回 JWT Token
func CreateToken(SecretKey []byte, Uid uint) (tokenString string, err error) {
	claims := &userClaims{
		jwt.StandardClaims{
			// 设置过期时间为 12 小时之后
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
			Issuer:    "Eshop",
			IssuedAt:  time.Now().Unix(),
		},
		Uid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(SecretKey)
	return
}

// ParseToken 解析JWT，传入 Token 和密钥，返回JWT结构
func ParseToken(tokenSrt string, secretKey []byte) (claims *userClaims, err error) {
	var token *jwt.Token
	var c = userClaims{}
	token, err = jwt.ParseWithClaims(tokenSrt, &c, func(*jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid {
		claims = &c
		return
	}
	return nil, errors.New("invalid token")
}

// 验证登录的中间件
func ValidateLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if len(token) < 8 {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "msg": "Unauthorized"})
			c.Abort()
			return
		}
		token = token[7:]
		claims, err := ParseToken(token, []byte(SecretKey))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "msg": err.Error()})
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}

// 管理员身份验证中间件
func ValidateAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var claims *userClaims
		claims = c.MustGet("claims").(*userClaims)
		loginId := claims.Uid
		if 1 != loginId {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "msg": "Unauthorized because is not admin."})
			c.Abort()
		}
		c.Set("claims", claims)
		c.Next()
	}
}

// /api/v1/auth POST 用户登录，返回 JWT
func Login(c *gin.Context) {
	var user User
	var form userLoginForm

	// 检查表单
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": err.Error()})
		return
	}

	// 查找用户
	result := DB.Where("username = ?", form.Username).Find(&user)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": "Not exist Users id"})
		return
	}

	// 验证密码
	pwdRight := user.ValidatePwdHash(form.Password)
	if pwdRight && user.Username == form.Username {
		// 生成 token
		token, err := CreateToken([]byte(SecretKey), user.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest, "msg": "Error", "Error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"msg":    "success",
			"token":  token,
			"userID": user.ID,
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "msg": "The Password or Username is not right"})
	}
}

// /api/v1/auth/newAuth POST 用户注册
func Register(c *gin.Context) {
	var form UserForm
	var newUser User
	// 检查表单
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": err.Error()})
		return
	}

	err := newUser.GeneratePwdHash(form.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": err.Error()})
		return
	}

	newUser.Username = form.Username
	newUser.Name = form.Name
	newUser.Mail = form.Email

	result := DB.Create(&newUser)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "msg": "success"})
}

// /api/v1/auth GET 获取用户数据，只有用户自己本人可以获取
func GetUserInfo(c *gin.Context) {
	var user User

	// 获取登录信息
	var claims *userClaims
	claims = c.MustGet("claims").(*userClaims)
	id := claims.Uid

	// 查找用户
	result := DB.Find(&user, id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": "Not exist Users id"})
		return
	}

	// 返回表单
	returnJSON := make(map[string]interface{})
	returnJSON["ID"] = user.ID
	returnJSON["Mail"] = user.Mail
	returnJSON["Username"] = user.Username
	returnJSON["Name"] = user.Name

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "msg": "success", "data": returnJSON})
}

// /api/v1/auth/:id PUT 用户数据更新，尽管表单需要提交Username，但是Username是无法修改的，做前端的时候要注意设置成disabled
func UpdateUserInfo(c *gin.Context) {
	var form UserForm
	var user User
	// 检查路由参数
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": "Illegal primary key. It must be a number."})
		return
	}

	// 检验当前登录是否为查询对象
	var claims *userClaims
	claims = c.MustGet("claims").(*userClaims)
	loginId := claims.Uid
	if uint(id) != loginId {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "msg": "Unauthorized because is wrong user"})
		return
	}

	// 验证表单
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": err.Error()})
		return
	}
	// 查找用户
	result := DB.Find(&user, id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": "Not exist Users id"})
		return
	}

	// 更新信息
	err = user.GeneratePwdHash(form.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": err.Error()})
		return
	}

	user.Name = form.Name
	user.Mail = form.Email

	DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "msg": "success"})
}
