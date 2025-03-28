// services/tcp_server.go
package services

import (
	"bufio"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/ricardomussett/gotest/models"
)

type TCPServer struct {
	Port     string
	MongoSvc *MongoService
}

func NewTCPServer(port string, mongoSvc *MongoService) *TCPServer {
	return &TCPServer{
		Port:     port,
		MongoSvc: mongoSvc,
	}
}

func (s *TCPServer) Start() {
	ln, err := net.Listen("tcp", s.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	log.Printf("TCP Server listening on %s", s.Port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		nmeaSentence := strings.TrimSpace(scanner.Text())
		if nmeaSentence == "" {
			continue
		}

		log.Printf("xsReceived NMEA: %s", nmeaSentence)

		if strings.HasPrefix(nmeaSentence, "GPRMC") {
			log.Printf("pa ve: %s", nmeaSentence)
			data := parseGPRMC(nmeaSentence)
			if data != nil {
				if err := s.MongoSvc.SaveGPSData(data); err != nil {
					log.Printf("Error saving to MongoDB: %v", err)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from connection: %v", err)
	}
}

func parseGPRMC(nmea string) *models.GPSData {
	parts := strings.Split(nmea, ",")
	if len(parts) < 10 {
		return nil
	}

	// Validate checksum if present
	// if idx := strings.Index(nmea, "*"); idx != -1 {
	// 	if !validateChecksum(nmea) {
	// 		return nil
	// 	}
	// }

	// Skip if data is invalid (status != 'A')
	if parts[2] != "A" {
		return nil
	}
	lat := parseLatitude(parts[3], parts[4])
	lon := parseLongitude(parts[5], parts[6])
	speed := parseSpeed(parts[7])
	course := parseCourse(parts[8])

	return &models.GPSData{
		Timestamp: time.Now(),
		Latitude:  lat,
		Longitude: lon,
		Speed:     speed,
		Course:    course,
		RawData:   nmea,
	}
}

func parseLatitude(latStr, dir string) float64 {
	if latStr == "" || dir == "" || len(latStr) < 2 {
		return 0
	}

	degrees, err := strconv.ParseFloat(latStr[:2], 64)
	if err != nil {
		return 0
	}
	minutes, err := strconv.ParseFloat(latStr[2:], 64)
	if err != nil {
		return 0
	}

	decimal := degrees + (minutes / 60)

	if dir == "S" {
		decimal = -decimal
	}

	return decimal
}

func validateChecksum(nmea string) bool {
	idx := strings.Index(nmea, "*")
	if idx == -1 {
		return false
	}

	calculated := 0
	for i := 1; i < idx; i++ {
		calculated ^= int(nmea[i])
	}

	expected, err := strconv.ParseInt(nmea[idx+1:], 16, 64)
	if err != nil {
		return false
	}

	return calculated == int(expected)
}

func parseLongitude(lonStr, dir string) float64 {
	if lonStr == "" || dir == "" || len(lonStr) < 3 {
		return 0
	}

	degrees, err := strconv.ParseFloat(lonStr[:3], 64)
	if err != nil {
		return 0
	}
	minutes, err := strconv.ParseFloat(lonStr[3:], 64)
	if err != nil {
		return 0
	}

	decimal := degrees + (minutes / 60)

	if dir == "W" {
		decimal = -decimal
	}

	return decimal
}
func parseSpeed(speedStr string) float64 {
	if speedStr == "" {
		return 0
	}

	speed, err := strconv.ParseFloat(speedStr, 64)
	if err != nil {
		return 0
	}

	return speed
}
func parseCourse(courseStr string) float64 {
	if courseStr == "" {
		return 0
	}

	course, err := strconv.ParseFloat(courseStr, 64)
	if err != nil {
		return 0
	}

	return course
}
