package main

import (
	"EbookShop/ApiV1"
	"fmt"
	"testing"
)

/**
 * @Author: DavinciEvans
 * @Author: zhouningsiyuan@foxmail.com
 * @Date: 2020/11/17 22:21
 * @Desc:
 */

func TestJWT(t *testing.T) {
	token, err := ApiV1.CreateToken([]byte("dev"), 25)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(token)
	}

	claims, err2 := ApiV1.ParseToken(token, []byte("dev"))
	if err2 != nil {
		fmt.Println("Error2: " + err2.Error())
		return
	}
	fmt.Println(claims.Uid)
}

func TestChangeDB(t *testing.T) {
	rows, err := ApiV1.DB.Select("book_id").Table("purchased_books").Joins("left outer join books on purchased_books.book_id = books.id and user_id = ?", 1).Rows()
	if err != nil {
		fmt.Println(err.Error())
	}
	for rows.Next() {
		var book uint
		fmt.Println(rows)
		rows.Scan("%d", &book)
		fmt.Println(book)
	}

}

func TestForgeDB(t *testing.T) {
	ApiV1.Forge()
}
