// models/gps_data.go
package models

import "time"

type GPSData struct {
    Timestamp time.Time `bson:"timestamp"`
    Latitude  float64   `bson:"latitude"`
    Longitude float64   `bson:"longitude"`
    Speed     float64   `bson:"speed"`     // En nudos
    Course    float64   `bson:"course"`    // Dirección en grados
    RawData   string    `bson:"raw_data"`  // Mensaje NMEA original
}