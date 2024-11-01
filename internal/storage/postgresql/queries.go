package postgresql

const (
	CreateTagStmt = "INSERT INTO tag(name) VALUES(@name) RETURNING id;"
	ListTagsStmt = "SELECT id, name from tag;"
	CreateCityStmt = "INSERT INTO city(name) VALUES(@name) RETURNING id;"
	ListCitiesStmt = "SELECT id, name from city;"
)