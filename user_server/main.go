package main

import (
	"fmt"
	"shopping-sys/user_server/utils"
)

func main() {
	pwd, salt := utils.CryptoPasswordWithSalt("123")
	fmt.Println("pwd, salt:", pwd, salt)
}
