package models

import (
	"github.com/nielsvanm/firewatch/internal/database"
)

// Device represents a physical node in the network
type Device struct {
	ID        int
	UUID      string
	Longitude float32
	Latitude  float32
}

func (d *Device) Save() {
	database.Database.Exec(`
		INSERT INTO device (uuid, longitude, latitude)
		VALUES ($1, $2, $3);`, d.UUID, d.Longitude, d.Latitude)
}

func GetAllDevices() []*Device {
	rows := database.Database.Query(`
	SELECT * FROM device;`)

	devs := []*Device{}
	for rows.Next() {
		d := Device{}
		rows.Scan(
			&d.ID,
			&d.UUID,
			&d.Longitude,
			&d.Latitude,
		)

		devs = append(devs, &d)
	}

	return devs
}

func GetDeviceByID(id int) *Device {
	rows := database.Database.Query(`
	SELECT * FROM device
	WHERE id = $1;`, id)

	d := Device{}
	for rows.Next() {
		rows.Scan(
			&d.ID,
			&d.UUID,
			&d.Longitude,
			&d.Latitude,
		)
	}

	return &d
}
