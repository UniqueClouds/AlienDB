package client

import (
	"encoding/json"
	"fmt"
)

func Test() {
	table := map[string]string{"col1": "xxx", "col2": "yyy"}
	s, _ := json.Marshal(table)
	sNew := make([]byte, 255)
	copy(sNew, s)
	m := make(map[string]string)
	json.Unmarshal(sNew, &m)
	fmt.Println(table)
	fmt.Println(s)
	fmt.Println(m["col1"])
}

func main() {
	RunClient()
}
