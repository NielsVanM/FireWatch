package models

// SetupQueries is a list of strings that setup the database tables
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
		expiry_date TIMESTAMP
	);`,
	`CREATE TABLE apitoken (
		id SERIAL PRIMARY KEY,
		user_id INT REFERENCES account(id),
		token TEXT UNIQUE
	)`,
	`CREATE TABLE device (
		id SERIAL PRIMARY KEY,
		uuid TEXT UNIQUE,
		longitude FLOAT,
		latitude FLOAT
	);`,
}
