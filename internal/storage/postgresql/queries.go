package postgresql

const (
	CreateTagStmt = "INSERT INTO tag(name) VALUES(@name);"
	ListTagsStmt = "SELECT id, name from tag;"
	CreateCityStmt = "INSERT INTO city(name) VALUES(@name);"
	ListCitiesStmt = "SELECT id, name from city;"
)