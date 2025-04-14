package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"social-network/internal/config"
	model "social-network/internal/models"
	"social-network/internal/storage/postgres"
)

const defaultFilePath = "data/people.v2.csv"
const mockHashPW = "$2a$10$At88yngCSUEznEwyjkEt9.Litxb8ZRqixcN81CGPqmKJXn.pY.caG" // pw: qwerty2

func main() {
	csvPath := flag.String("csv", defaultFilePath, "Путь к CSV файлу с данными пользователей")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки конфига", err)
	}

	storage, err := postgres.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Ошибка подключения к бд", err)
	}
	defer storage.Close()

	records, err := loadCSVFile(csvPath)
	if err != nil {
		log.Fatal("Ошибка чтения csv", err)
	}

	for i, row := range records {
		if len(row) < 3 {
			log.Printf("Строка %d пропущена: недостаточно столбцов: %v", i, row)
			continue
		}

		var user model.User

		user.ID = uuid.NewString()
		user.PasswordHash = mockHashPW

		secondName, firstName, errFullName := parseFullName(row[0])
		if errFullName != nil {
			log.Printf("Ошибка парсинга полного имени: %v", errFullName)
			continue
		}

		user.FirstName = firstName
		user.SecondName = &secondName
		user.Birthday = &row[1]
		user.City = &row[2]

		if err = storage.CreateUser(&user); err != nil {
			log.Printf("Ошибка записи пользователя %v в бд: %v", user, err)
			continue
		}
	}

	log.Println("Импорт данных завершён успешно")
}

func parseFullName(fullName string) (lastName, firstName string, err error) {
	parts := strings.SplitN(strings.TrimSpace(fullName), " ", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("неверный формат полного имени: %q", fullName)
	}
	return parts[0], parts[1], nil
}

func loadCSVFile(csvPath *string) (records [][]string, err error) {
	var file *os.File

	file, err = os.Open(*csvPath)
	if err != nil {
		log.Printf("Ошибка открытия CSV файла: %v", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err = reader.ReadAll()
	if err != nil {
		log.Printf("Ошибка чтения CSV файла: %v", err)
		return
	}

	if len(records) == 0 {
		log.Print("CSV файл пуст или содержит только заголовок")
		return
	}

	return
}
