package createfile

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type LogDir struct {
	LogDirectory string
}

func New(filename string) *LogDir {
	year, month, day := time.Now().Date()
	yearInt := int(year)
	monthInt := int(month)
	dayInt := int(day)

	py := filepath.Join("logs", fmt.Sprintf("%04d", yearInt)+"/")
	if _, err := os.Stat(py); os.IsNotExist(err) {
		err := os.Mkdir(py, 0755)
		if err != nil {
			fmt.Println("[x] Failed to Mkdir:", err)
			return nil
		}
	}

	pm := filepath.Join(py, fmt.Sprintf("%02d", monthInt)+"/")

	if _, err := os.Stat(pm); os.IsNotExist(err) {
		err := os.Mkdir(pm, 0755)
		if err != nil {
			fmt.Println("[x] Failed to Mkdir:", err)
			return nil
		}
	}

	pd := filepath.Join(pm, fmt.Sprintf("%02d", dayInt)+"/")

	if _, err := os.Stat(pd); os.IsNotExist(err) {
		err := os.Mkdir(pd, 0755)
		if err != nil {
			fmt.Println("[x] Failed to Mkdir:", err)
			return nil
		}
	}

	return &LogDir{
		LogDirectory: pd,
	}
}

func SetLogFile(filename string) *os.File {
	year, month, day := time.Now().Date()
	yearInt := int(year)
	monthInt := int(month)
	dayInt := int(day)

	fileName := fmt.Sprintf("%04d/%02d/%02d/%s.log", yearInt, monthInt, dayInt, filename)
	filePath, err := os.OpenFile(filepath.Join("logs", fileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("[x] Failed to OpenFile:", err)
		return nil
	}
	return filePath
}

func (l *LogDir) Info(filename string) *log.Logger {
	getFilePath := SetLogFile(filename)
	return log.New(getFilePath, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func (l *LogDir) Warning(filename string) *log.Logger {
	getFilePath := SetLogFile(filename)
	return log.New(getFilePath, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func (l *LogDir) Error(filename string) *log.Logger {
	getFilePath := SetLogFile(filename)
	return log.New(getFilePath, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func (l *LogDir) Fatal(filename string) *log.Logger {
	getFilePath := SetLogFile(filename)
	return log.New(getFilePath, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func CreateTextFile(path string, records string, appLogger *LogDir, filename string) error {

	// Open the file for writing. Create if it doesn't exist, truncate if it does.
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		appLogger.Error(filename).Println("Error opening file:", err)
		return err
	}
	defer file.Close()

	// Write JSON data to the file
	_, err = file.WriteString(records)
	if err != nil {
		appLogger.Error(filename).Println("Error writing to file:", err)
		return err
	}

	fmt.Println("Records inserted successfully.")
	return nil
}
