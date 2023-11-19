package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/ngnhub/snippetbox/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (model *UserModel) Insert(email, name, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (email, name, hashed_password, created)
	VALUES (?, ?, ?, UTC_TIMESTAMP())`

	_, err = model.DB.Exec(stmt, email, name, string(hash))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrorDuplicatedEmail
			}
		}
		return err
	}
	return nil
}

func (model *UserModel) Authencticate(email, password string) (int, error) {
	stmt := `SELECT id, hashed_password FROM users
	WHERE email = ? AND active`

	var userId int
	var userPassword []byte

	row := model.DB.QueryRow(stmt, email)
	err := row.Scan(&userId, &userPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrorInvalidCredentials
		}
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword(userPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrorInvalidCredentials
		}
		return 0, err
	}
	return userId, nil
}

func (m *UserModel) GetBy(id int) (*models.User, error) {
	stmt := `SELECT id, name, email, created, active FROM users
	WHERE id = ?`

	user := models.User{}

	result := m.DB.QueryRow(stmt, id)

	err := result.Scan(&user.ID, &user.Name, &user.Email, &user.Created, &user.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErorNoRecord
		} else {
			return nil, err
		}
	}
	return &user, nil
}
