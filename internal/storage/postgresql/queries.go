package postgresql

const (
	CreateTagStmt = "INSERT INTO tag(name) VALUES(@name) RETURNING id;"
	CheckTagNameUniqueStmt = "SELECT id FROM tag WHERE name=@name;"
	ListTagsStmt = "SELECT id, name from tag;"
	GetMultipleTagsStmt = "SELECT id FROM tag WHERE name=ANY(@tag_array)"
	DeleteTagStmt = "DELETE FROM tag WHERE id=@id"
	
	CreateCityStmt = "INSERT INTO city(name) VALUES(@name) RETURNING id;"
	CheckCityNameUniqueStmt = "SELECT id FROM city WHERE name=@name;"
	ListCitiesStmt = "SELECT id, name from city;"
	DeleteCityStmt = "DELETE FROM city WHERE id=@id"

	// CreateHotelStmt = "INSERT INTO hotel(name, description, city_id) VALUES(@name, @desc, @city_id) RETURNING id;"
	CreateTagHotelStmt = "INSERT INTO tag_hotel(hotel_id, tag_id) VALUES(@hotel_id, @tag_id);"
	CheckHotelNameUniqueStmt = "SELECT id FROM hotel WHERE name=@name"
	CreateHotelStmt = "INSERT INTO hotel(name, description, city_id) VALUES(@name, @desc, (SELECT id FROM city WHERE name=@city_name LIMIT 1)) RETURNING id;"
)