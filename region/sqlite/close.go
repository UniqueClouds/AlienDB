package sqlite

func Close() {
	db.Close()
	logFile.Close()
}
