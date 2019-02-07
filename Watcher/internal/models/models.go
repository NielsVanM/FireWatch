package models

var SetupQueries = []string{
	`CREATE TABLE account (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE,
		password TEXT
	);`,
	`CREATE TABLE session (
		id SERIAL PRIMARY KEY,
		user_id INT REFERENCES account(id),
		token TEXT UNIQUE, 
		address TEXT UNIQUE,
		expiry_date TIMESTAMP
	)`,
	`CREATE TABLE device (
		id SERIAL PRIMARY KEY,
		uuid TEXT UNIQUE,
		longitude FLOAT,
		latitude FLOAT
	);`,
}
