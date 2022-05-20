package main

import (
	"fmt"
	"my/AlienDB/region/sqlite"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	res, err := sqlite.Exec("SELECT * FROM userinfo")
	if err != nil {
		fmt.Println(res)
	}

}
