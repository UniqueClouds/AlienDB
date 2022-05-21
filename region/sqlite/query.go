package sqlite

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Query(sqlstring string) (string, error) {
	//db, err := sql.Open("sqlite3", "./foo2.db")
	//defer db.Close()
	//if err != nil {
	//	return "database open failed", err
	//}
	rows, err := db.Query(sqlstring) // ignore_security_alert
	if err != nil {
		return "query error check the sql statement", err
	}

	input, err := rows.Columns()
	if err != nil {
		return "query error check the sql statement", err
	}

	l := len(input)
	fmt.Println(l)
	var resString string
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
			return "data find error, try another sql statement", err
		}
		resString = fmt.Sprintf("%s\n%v", resString, s)
	}
	//region test
	fmt.Println(resString)
	return resString, nil
}
