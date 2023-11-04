package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"playground/configs"
	"playground/routes"
	"playground/services"
	"playground/utils"
	"time"
)

var logFile *os.File

func logMessage(format string, a ...interface{}) {
	logTime := time.Now().Format("2006-01-02 15:04:05")
	logLine := fmt.Sprintf(format, a...)
	logEntry := fmt.Sprintf("[%s] %s", logTime, logLine)

	// Print to the terminal
	fmt.Println(logEntry)

	// Write to the log file (if it's opened)
	if logFile != nil {
		log.Print(logEntry)
	}
}

func main() {
	err := configs.LoadEnv()
	if err != nil {
		return
	}

	// Open the log file
	logFile, err := os.OpenFile("logfile.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()

	// Set the log output to use UTF-8 encoding
	logFile.WriteString("\xEF\xBB\xBF") // Byte Order Mark (BOM) for UTF-8
	log.SetOutput(logFile)

	address := os.Getenv("HOST_PORT")
	logMessage("Server is starting at http://%s", address)

	minValue := 2000
	maxValue := 60000

	go func() {
		for {
			licensePlate := utils.GenerateThaiLicensePlate()
			logMessage("Generated license plate: %s", licensePlate)
			services.ConnectToRabbitMQ(licensePlate)

			randomDelay := rand.Intn(maxValue-minValue+1) + minValue
			logMessage("Next data in %ds", randomDelay/1000)
			time.Sleep(time.Duration(randomDelay) * time.Millisecond)
		}
	}()

	r := routes.SetupRouter()

	err = http.ListenAndServe(address, r)
	if err != nil {
		return
	}
}
