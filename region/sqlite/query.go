package sqlite

import (
	"fmt"
)

func Query(sqlstring string) (string, error) {

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
		s := make([]string, l)
		ps := make([]*string, l)
		c := make([]interface{}, l)
		for i := 0; i < l; i++ {
			ps[i] = &s[i]
			c[i] = ps[i]
		}
		err = rows.Scan(c...)
		if err != nil {
			return "data find error, try another sql statement", err
		}
		resString = fmt.Sprintf("%s\n%s", resString, s)
	}
	return resString, nil
}
