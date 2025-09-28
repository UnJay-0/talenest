package main

import (
	"fmt"
	"talenest/backend/internal/utils"
)

func main() {
	cfg := utils.LoadConfig()
	fmt.Println(cfg.SQLitePath)
}
