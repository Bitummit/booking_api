package postgresql

const (
	CreateTagStmt = "INSERT INTO tag(name) VALUES(@name) RETURNING id;"
	GetTagByName = "SELECT id FROM tag WHERE name=@name;"
	ListTagsStmt = "SELECT id, name from tag;"
	GetMultipleTagsStmt = "SELECT id FROM tag WHERE name=ANY(@tag_array)"
	DeleteTagStmt = "DELETE FROM tag WHERE id=@id"
	
	CreateCityStmt = "INSERT INTO city(name) VALUES(@name) RETURNING id;"
	GetCityByName = "SELECT id FROM city WHERE name=@name;"
	ListCitiesStmt = "SELECT id, name from city;"
	DeleteCityStmt = "DELETE FROM city WHERE id=@id"

	CreateTagHotelStmt = "INSERT INTO tag_hotel(hotel_id, tag_id) VALUES(@hotel_id, (SELECT id FROM tag WHERE name=@tag_name LIMIT 1));"
	CheckHotelNameUniqueStmt = "SELECT id FROM hotel WHERE name=@name"
	CreateHotelStmt = "INSERT INTO hotel(name, description, city_id) VALUES(@name, @desc, (SELECT id FROM city WHERE name=@city_name LIMIT 1)) RETURNING id;"
	GetOwnedHotelsStmt = `
		SELECT h.id, h.name, h.description, c.name, t.name 
		FROM hotel AS h 
		LEFT JOIN city AS c ON h.city_id=c.id 
		LEFT JOIN tag_hotel AS th ON th.hotel_id=h.id 
		LEFT JOIN tag AS t ON th.tag_id=t.id 
		WHERE h.manager_id=@user_id;
	`
	GetAllHotelsStmt = `
		SELECT h.id, h.name, h.description, c.name, t.name 
		FROM hotel AS h 
		LEFT JOIN city AS c ON h.city_id=c.id 
		LEFT JOIN tag_hotel AS th ON th.hotel_id=h.id
		LEFT JOIN tag AS t ON th.tag_id=t.id;
	`
	GetHotelStmt = `
		SELECT h.id, h.name, h.description, c.name, t.name
	`
)