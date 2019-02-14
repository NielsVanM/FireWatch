package models

import "github.com/nielsvanm/firewatch/internal/database"

// Device represents a physical node in the network
type Device struct {
	ID        int
	UUID      string
	Longitude float32
	Latitude  float32
}

// Save stores the device in the database
func (d *Device) Save() {
	database.DB.Exec(`
		INSERT INTO device (uuid, longitude, latitude)
		VALUES ($1, $2, $3);`, d.UUID, d.Longitude, d.Latitude)
}

// GetAllDevices returns a list of all devices from the database
func GetAllDevices() []*Device {
	rows, _ := database.DB.Query(`
	SELECT * FROM device;`)

	defer rows.Close()

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

// GetDeviceByID retrievs a device from the database marked by the provided ID
func GetDeviceByID(id int) *Device {
	rows, _ := database.DB.Query(`
	SELECT * FROM device
	WHERE id = $1;`, id)

	defer rows.Close()

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
