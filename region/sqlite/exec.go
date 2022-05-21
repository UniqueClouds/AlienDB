package sqlite

func Exec(sqlString string) (string, error) {

	_, err := db.Exec(sqlString) // ignore_security_alert

	if err != nil {
		return "database open fail check if there is a database on the region", err
	}
	return "success", nil
}
