package main

import (
	"fmt"

	"github.com/ccnj/go-utils/passhash"
)

func main() {
	fmt.Println("hello world")

	// salt, err := passhash.GenerateSalt(32)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(salt)
	// fmt.Println(string(salt))

	// salt := "zswZFdDf14pso5rTi6KmXLcV9SZLYqbu"
	password := "123456"

	// hash := passhash.HashPassword(password, []byte(salt))
	// fmt.Println(hash)

	salt, pwdHash, _ := passhash.EasyHash(password, 32)
	fmt.Println(salt)
	fmt.Println(pwdHash)

	fmt.Println(passhash.Verify(password, pwdHash, salt))

}
