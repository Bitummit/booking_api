package postgresql

const (
	CreateTagStmt = "INSERT INTO tag(name) VALUES(@name) RETURNING id;"
	CheckTagNameUniqueStmt = "SELECT id FROM tag WHERE name=@name;"
	ListTagsStmt = "SELECT id, name from tag;"
	DeleteTagStmt = "DELETE FROM tag WHERE id=@id"
	
	CreateCityStmt = "INSERT INTO city(name) VALUES(@name) RETURNING id;"
	CheckCityNameUniqueStmt = "SELECT id FROM city WHERE name=@name;"
	ListCitiesStmt = "SELECT id, name from city;"
	DeleteCityStmt = "DELETE FROM city WHERE id=@id"
)