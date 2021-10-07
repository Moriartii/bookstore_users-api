package users

import (
	"fmt"
	"strings"

	"github.com/Moriartii/bookstore_users-api/datasources/postgres/users_db"
	"github.com/Moriartii/bookstore_users-api/logger"
	"github.com/Moriartii/bookstore_users-api/utils/errors"
	"github.com/Moriartii/bookstore_users-api/utils/postgres_utils"
)

const (
	//indexUniqueEmail = "duplicate key value violates unique constraint \"users_email_key"
	queryInsertUser   = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES($1, $2, $3, $4, $5, $6) returning id;"
	queryGetUser      = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id = $1;"
	queryUpdateUser   = "UPDATE users SET first_name=$1, last_name=$2, email=$3 WHERE id=$4;"
	queryDeleteUser   = "DELETE FROM users WHERE id=$1;"
	queryFindByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=$1;"

	queryFindByEmailAndPassword = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE email=$1 AND password=$2 AND status=$3;"
)

func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user statment", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()
	if _, err = stmt.Exec(user.Id); err != nil {
		logger.Error("error when trying to delete user", err)
		return errors.NewInternalServerError("database error")
		// return postgres_utils.ParseError(err)
	}
	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statment", err)
		return errors.NewInternalServerError(("database error"))
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		logger.Error("error when trying to update user", err)
		return errors.NewInternalServerError(("database error"))
		//return postgres_utils.ParseError(err)
	}
	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statment", err)
		return errors.NewInternalServerError(("database error"))
	}
	defer stmt.Close()

	//user.DateCreated = date_utils.GetNowString()

	saveErr := stmt.QueryRow(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password).Scan(&user.Id)
	if saveErr != nil {
		logger.Error("error when trying to save user", saveErr)
		return errors.NewInternalServerError(("database error"))
		//return postgres_utils.ParseError(saveErr)
	}
	return nil
}

func (user *User) Get() *errors.RestErr {
	//if err := users_db.Client.Ping(); err != nil {
	//	panic(err)
	//}
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get user statment", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()
	// results, err := stmt.QueryRow()
	// defer results.Close() 		// NEED TO CLOSE IF MANY!!
	result := stmt.QueryRow(user.Id)
	getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status)
	if getErr != nil {
		logger.Error("error when trying to get user by ID", getErr)
		return errors.NewInternalServerError("database error")
		// return postgres_utils.ParseError(getErr)
	}
	return nil
}

func (user *User) FindByStatus(status string) ([]User, *errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find users by status statment", err)
		return nil, errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error when trying to find users by status", err)
		return nil, errors.NewInternalServerError("database error")
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			logger.Error("error when trying to scan user row into struct", err)
			return nil, errors.NewInternalServerError("database error")
			// return nil, postgres_utils.ParseError(err)
		}
		results = append(results, user)
	}

	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}
	return results, nil
}

func (user *User) FindByEmailAndPassword() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryFindByEmailAndPassword)
	if err != nil {
		logger.Error("error when trying to prepare get user by email and pasword statment", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()
	result := stmt.QueryRow(user.Email, user.Password, StatusActive)
	getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status)
	if getErr != nil {
		if strings.Contains(getErr.Error(), postgres_utils.ErrorNoRows) {
			return errors.NewNotFoundError("invalid user credentials")
		}
		logger.Error("error when trying to get user by email and password", getErr)
		return errors.NewInternalServerError("database error")
	}
	return nil
}
