package sqlFunctions

var PostCity = `INSERT INTO "public"."Cities" (city_name, country, lat, lon) 
	VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING;`

var PostWeather = `INSERT INTO "public"."Weather" (city_id, date, temperature, date_txt) 
	VALUES ($1, $2, $3, $4) ON CONFLICT (city_id, date) DO UPDATE 
	SET temperature = $3, date_txt = $4;`

var GetCityWeather = `SELECT w.temperature, w.date, w.date_txt FROM "Weather" AS w 
   INNER JOIN "Cities" AS c 
   ON w.city_id = c.city_id AND c.city_name = $1 
   AND c.country = $2 
   AND EXTRACT(epoch FROM now()) <= w.date`
