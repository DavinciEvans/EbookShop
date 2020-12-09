package ApiV1

/**
 * @Author: DavinciEvans
 * @Author: zhouningsiyuan@foxmail.com
 * @Date: 2020/11/6 11:31
 * @Desc:
 */

import (
	"github.com/brianvoe/gofakeit/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// 电子书数据库模型
type Book struct {
	ID             uint `gorm:"primaryKey"`
	Name           string
	Price          float32
	Content        string
	Author         string
	StarSum        int
	PayNumber      int
	CategoryID     uint
	Carts          []Cart
	PurchasedBooks []PurchasedBook
	Comments       []Comment
}

// 用户数据库模型
type User struct {
	ID             uint `gorm:"primaryKey"`
	Mail           string
	Username       string `gorm:"unique;not null"`
	PasswordHash   string `gorm:"not null"`
	Name           string `gorm:"not null"`
	Carts          []Cart
	PurchasedBooks []PurchasedBook
	Comments       []Comment
}

// 分类目录数据库模型
type Category struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Books []Book
}

// 购物车条目数据库类型
type Cart struct {
	UserID uint `gorm:"primaryKey"`
	BookID uint `gorm:"primaryKey"`
}

// 已购书籍数据库类型
type PurchasedBook struct {
	Time   time.Time
	UserID uint `gorm:"primaryKey"`
	BookID uint `gorm:"primaryKey"`
}

// 书评数据库类型
type Comment struct {
	ID      uint   `gorm:"primaryKey"`
	Content string `gorm:"not null"`
	Star    int
	UserID  uint `gorm:"not null"`
	BookID  uint `gorm:"not null"`
}

// 初始化模型函数
func modelsInit(development bool) {
	err := DB.AutoMigrate(&Category{}, &Book{}, &User{}, &Cart{}, &PurchasedBook{}, &Comment{})
	if err != nil {
		panic("Can't Migrate")
	}
}

// 载入仿真数据
func Forge() {
	err := DB.AutoMigrate(&Category{}, &Book{}, &User{}, &Cart{}, &PurchasedBook{}, &Comment{})
	if err != nil {
		panic("Can't Migrate")
	}
	var categories []Category
	var books []Book
	var users []User
	var comments []Comment
	var carts []Cart
	var purchasedBooks []PurchasedBook
	// 生成随机数据
	gofakeit.Seed(time.Now().Unix())
	// 生成目录分类
	for i := 0; i < 5; i++ {
		category := Category{Name: gofakeit.Word()}
		if i == 0 {
			category.Name = "uncategorized"
		}
		categories = append(categories, category)
	}

	// 生成书本
	for i := 0; i < 80; i++ {
		book := Book{
			Name:      gofakeit.Sentence(gofakeit.Number(3, 6)),
			Price:     gofakeit.Float32Range(10, 50),
			Content:   gofakeit.Paragraph(3, 20, 160, "\n"),
			Author:    gofakeit.Name(),
			PayNumber: gofakeit.Number(2, 4),
			StarSum:   gofakeit.Number(10, 20),
		}
		chooseCategory := gofakeit.Number(0, 4)
		categories[chooseCategory].Books = append(categories[chooseCategory].Books, book)
		books = append(books, book)
	}

	// 生成用户和他的相关购买信息
	for i := 0; i < 10; i++ {
		// 用户信息
		user := User{
			Username: gofakeit.Username(),
			Name:     gofakeit.Name(),
			Mail:     gofakeit.Email(),
		}
		if i == 0 {
			user.Username = "admin"
		}
		pwd := "123456"
		err := user.GeneratePwdHash(pwd)
		if err != nil {
			panic("Can't generate Password")
		}
		users = append(users, user)

		// 生成已购书籍
		for j := 0; j < gofakeit.Number(0, 3); j++ {
			chooseBook := 8*i + j
			purchasedBook := PurchasedBook{
				BookID: uint(chooseBook),
				UserID: uint(i + 1),
				Time:   gofakeit.Date(),
			}

			purchasedBooks = append(purchasedBooks, purchasedBook)

			// 随机生成书评
			if gofakeit.Bool() {
				star := gofakeit.Number(1, 10)
				comment := Comment{
					Content: gofakeit.Paragraph(1, gofakeit.Number(1, 3), gofakeit.Number(10, 20), "\n"),
					Star:    star,
					UserID:  uint(i + 1),
					BookID:  uint(chooseBook),
				}

				comments = append(comments, comment)
			}
		}

		// 生成购物车
		for j := 0; j < gofakeit.Number(0, 3); j++ {
			chooseBook := 8*i + j + 5
			cart := Cart{
				BookID: uint(chooseBook),
				UserID: uint(i + 1),
			}
			carts = append(carts, cart)
		}

	}

	DB.Create(&categories)
	DB.Create(&books)
	DB.Create(&users)
	DB.Create(&comments)
	DB.Create(&carts)
	DB.Create(&purchasedBooks)
}

// 生成用户加密密码
func (user *User) GeneratePwdHash(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hash)
	return nil
}

// 验证用户密码是否正确
func (user *User) ValidatePwdHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return false
	} else {
		return true
	}
}
