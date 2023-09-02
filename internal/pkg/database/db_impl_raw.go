package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pjmessi/go-database-usage/config"
	"github.com/pjmessi/go-database-usage/internal/model"
)

type RawDbImpl struct {
	db *sql.DB
}

func NewDbImpl(appConf *config.AppConfig) (Db, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", appConf.DB_USER, appConf.DB_PASSWORD, appConf.DB_HOST, appConf.DB_PORT, appConf.DB_DATABASE))
	if err != nil {
		return nil, fmt.Errorf("database.NewDbImpl: %w", err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	impl := &RawDbImpl{
		db: db,
	}

	return impl, nil
}

func (r *RawDbImpl) IsHealthy() bool {
	var total int
	res, err := r.db.Query("SELECT 2 + 2;")
	if err != nil {
		log.Println(fmt.Errorf("database.IsHealthy: %w", err))
		return false
	}

	defer res.Close()

	if res.Next() {
		log.Println(res.Columns())
		err := res.Scan(&total)
		if err != nil {
			log.Println(fmt.Errorf("database.IsHealthy: %w", err))
			return false
		}
	}

	log.Printf("database.IsHealthy result: %d", total)

	return total == 4
}

func (r *RawDbImpl) CreateUser(user *model.User) error {
	stmt, err := r.db.Prepare("INSERT INTO users (id, email, password, first_name, last_name, created_at, updated_at) VALUE (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("database.CreateUser: %w", err)
	}

	defer func() {
		if cerr := stmt.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("database.CreateUser: %w", cerr)
		}
	}()

	_, err = stmt.Exec(user.Id, user.Email, user.Password, user.FirstName, user.LastName, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("database.CreateUser: %w", err)
	}

	return nil
}

func (r *RawDbImpl) IsUserEmailTaken(email string) (bool, error) {
	var isTaken bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT * FROM users WHERE email = \"%s\");", email)
	res, err := r.db.Query(query)
	if err != nil {
		return false, fmt.Errorf("database.IsUserEmailTaken: %w", err)
	}

	defer res.Close()

	if res.Next() {
		err := res.Scan(&isTaken)
		if err != nil {
			return false, fmt.Errorf("database.IsUserEmailTaken: %w", err)
		}
	}

	return isTaken, nil
}

func (r *RawDbImpl) GetUserByEmail(email string) (bool, model.User, error) {
	rows, err := r.db.Query("SELECT * FROM users WHERE email = ?;", email)
	if err != nil {
		return false, model.User{}, fmt.Errorf("database.GetUserByEmail: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var user model.User
		err := rows.Scan(&user.Id, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return false, model.User{}, fmt.Errorf("database.GetUserByEmail: %w", err)
		}
		return true, user, nil
	} else {
		return false, model.User{}, nil
	}
}

func (r *RawDbImpl) CloseConnection() {
	r.db.Close()
}
