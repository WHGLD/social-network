package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	model "social-network/internal/models"
	"social-network/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(databaseURL string) (*Storage, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", storage.ErrURLOpening)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", storage.ErrDBConnection)
	}

	return &Storage{db}, nil
}

func (s *Storage) Close() error {
	s.db.Close()
	return nil
}

func (s *Storage) CreateUser(user *model.User) error {
	_, err := s.db.Exec(
		`INSERT INTO users 
        	(user_id, password_hash, first_name, second_name, birthday, sex, biography, city) 
        VALUES 
        	($1, $2, $3, $4, $5, $6, $7, $8)`,
		user.ID,
		user.PasswordHash,
		user.FirstName,
		user.SecondName,
		user.Birthday,
		user.Sex,
		user.Biography,
		user.City,
	)
	return err
}

func (s *Storage) GetUserByID(userID string) (*model.User, error) {
	var user model.User
	err := s.db.QueryRow(
		"SELECT user_id, password_hash, first_name, second_name, birthday, sex, biography, city FROM users WHERE user_id = $1",
		userID,
	).Scan(&user.ID, &user.PasswordHash, &user.FirstName, &user.SecondName, &user.Birthday, &user.Sex, &user.Biography, &user.City)
	return &user, err
}

func (s *Storage) UsersSearch(FirstName, SecondName string) ([]model.User, error) {
	var users []model.User

	query := `
		SELECT
			user_id,
			password_hash,
			first_name,
			second_name,
			birthday,
			sex,
			biography,
			city
		FROM users
		WHERE first_name ILIKE $1 AND second_name ILIKE $2
	`
	firstParam := FirstName + "%"
	secondParam := SecondName + "%"

	rows, err := s.db.Query(query, firstParam, secondParam)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User
		if err = rows.Scan(
			&user.ID,
			&user.PasswordHash,
			&user.FirstName,
			&user.SecondName,
			&user.Birthday,
			&user.Sex,
			&user.Biography,
			&user.City,
		); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return users, err
	}

	return users, nil
}
