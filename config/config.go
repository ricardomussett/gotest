// config/config.go
package config

type Config struct {
	TCPPort    string
	MongoDBURI string
	DBName     string
	Collection string
}

func LoadConfig() *Config {
	return &Config{
		TCPPort:    ":8080", // Puerto TCP para recibir datos NMEA
		MongoDBURI: "mongodb://localhost:27017",
		DBName:     "gps_db",
		Collection: "gps_data",
	}
}
