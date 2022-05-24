package sqlite

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func Restore(filepath string) {
	log.Printf("go get file:%v\n", filepath)
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	defer f.Close()
	if err != nil {
		log.Panic("打开备份文件异常")
	}
	b := bufio.NewReader(f)
	for {
		str, err := b.ReadString('\n')
		if err == io.EOF {
			break
		}
		fmt.Println(str)
		//Exec(str)

	}

	log.Printf("restore finished")

}
