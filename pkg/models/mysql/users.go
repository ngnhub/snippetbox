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

func (model *UserModel) authencticate(email, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) GetBy(id int) (*models.User, error) {
	return nil, nil
}
