package database

func NewDatabase() (Database, error) {
	return NewSQLiteDatabase()
}
