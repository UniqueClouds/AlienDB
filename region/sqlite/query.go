package sqlite

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Query(sqlstring string) ([]map[string]interface{}, error) {
	//db, err := sql.Open("sqlite3", "./foo2.db")
	//defer db.Close()
	//if err != nil {
	//	return "database open failed", err
	//}
	rows, err := db.Query(sqlstring) // ignore_security_alert
	var m []map[string]interface{}
	if err != nil {
		return nil, err
	}

	input, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	l := len(input)
	//var resString string
	for rows.Next() {
		s := make([]interface{}, l)
		ps := make([]*interface{}, l)
		c := make([]interface{}, l)
		for i := 0; i < l; i++ {
			ps[i] = &s[i]
			c[i] = ps[i]
		}
		err = rows.Scan(c...)
		if err != nil {
			return nil, err
		}
		temp := make(map[string]interface{})
		for i, Col := range input {
			temp[Col] = s[i]
		}
		m = append(m, temp)
		//resString = fmt.Sprintf("%s\n%v", resString, s)
	}
	//region test
	fmt.Println(m)
	return m, nil
}
