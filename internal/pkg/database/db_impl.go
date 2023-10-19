package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pjmessi/golang-practice/config"
	"github.com/pjmessi/golang-practice/internal/model"
)

type RawDbImpl struct {
	db *sql.DB
}

func NewDb(appConf *config.AppConfig) (Db, error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", appConf.DB_USER, appConf.DB_PASSWORD, appConf.DB_HOST, appConf.DB_PORT, appConf.DB_DATABASE)
	db, err := sql.Open("mysql", dns)
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

func (r *RawDbImpl) CheckHealth() error {
	var total int
	res, err := r.db.Query("SELECT 2 + 2;")
	if err != nil {
		return fmt.Errorf("database.CheckHealth(): %w", err)
	}

	defer res.Close()

	if res.Next() {
		err := res.Scan(&total)
		if err != nil {
			return fmt.Errorf("database.CheckHealth(): %w", err)
		}
	}

	if total != 4 {
		return fmt.Errorf("database.CheckHealth(): expected result 4 but received %d", total)
	}

	return nil
}

func (r *RawDbImpl) SaveUser(ctx context.Context, user *model.User) error {
	stmt, err := r.db.Prepare("INSERT INTO users (id, email, password, first_name, last_name, created_at, updated_at) VALUE (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("database.SaveUser(): %w", err)
	}

	defer func() {
		if cerr := stmt.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("database.SaveUser(): %w", cerr)
		}
	}()

	_, err = stmt.Exec(user.Id, user.Email, user.Password, user.FirstName, user.LastName, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("database.SaveUser(): %w", err)
	}

	return nil
}

func (r *RawDbImpl) IsUserEmailTaken(ctx context.Context, email string) (bool, error) {
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

func (r *RawDbImpl) GetUserByEmail(ctx context.Context, email string) (bool, model.User, error) {
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
