package main

import (
	"clean_architecture_fiber/config"
	"fmt"
)

func main() {
	cfg := config.LoadAppConfig()
	dsn := cfg.GetDatabaseURL()
	fmt.Print(dsn)
}
