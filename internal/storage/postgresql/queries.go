package postgresql

const (
	CreateTagStmt = "INSERT INTO tag(name) VALUES(@name);"
	CreateCityStmt = "INSERT INTO city(name) VALUES(@name);"
)