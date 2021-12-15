package users

//data access object

import (
	"fmt"
	"users-api/datasources/mysql/users_db"
	"users-api/utils/errors"
	"users-api/utils/mysql_utils"
)

const perPage = 9

const (
	queryInsertUser       = "INSERT INTO users(first_name,last_name,national_id,date_created,status,password) VALUES(?,?,?,?,?,?);"
	queryGetUser          = "SELECT id,first_name,last_name,national_id,date_created,status FROM users WHERE id=?;"
	queryGetAllUser       = "SELECT id,first_name,last_name,national_id,date_created,status FROM users WHERE status=? LIMIT ? OFFSET ?;"
	queryUpdateUser       = "UPDATE users SET first_name=?,last_name=?,national_id=? WHERE id=?;"
	queryDeleteUser       = "DELETE FROM users WHERE id=?;"
	queryFindUserByStatus = "SELECT id,first_name,last_name,national_id,date_created,status FROM users WHERE status=? LIMIT ? OFFSET ?;;"
	queryTotalSize        = "SELECT COUNT(*) FROM users;"
)

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id) //Return single row
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.NationalID, &user.DateCreated, &user.Status); getErr != nil {
		return mysql_utils.ParseError(getErr)
	}
	return nil
}

func (user *User) GetAll(page int, status string) ([]User, *errors.RestErr, int64) {
	stmt, err := users_db.Client.Prepare(queryGetAllUser)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error()), 0
	}
	defer stmt.Close()

	var total int64
	errQuery := users_db.Client.QueryRow(queryTotalSize).Scan(&total)
	if errQuery != nil {
		return nil, errors.NewInternalServerError(err.Error()), 0
	}

	rows, err := stmt.Query(status, perPage, (page-1)*perPage)

	if err != nil {
		return nil, errors.NewInternalServerError(err.Error()), 0
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id,
			&user.FirstName,
			&user.LastName,
			&user.NationalID,
			&user.DateCreated,
			&user.Status); err != nil {
			return nil, mysql_utils.ParseError(err), 0
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("users was not found")), 0
	}
	return results, nil, total

}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.NationalID, user.DateCreated, user.Status, user.Password)
	if saveErr != nil {
		return mysql_utils.ParseError(saveErr)
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		return mysql_utils.ParseError(saveErr)
	}

	user.Id = userId
	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	_, _err := stmt.Exec(user.FirstName, user.LastName, user.NationalID, user.Id)
	if _err != nil {
		return mysql_utils.ParseError(err)
	}
	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	if _, err := stmt.Exec(user.Id); err != nil {
		return mysql_utils.ParseError(err)
	}
	return nil
}

func (user *User) FindByStatus(page int, status string) ([]User, *errors.RestErr, int64) {
	stmt, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error()), 0
	}
	defer stmt.Close()

	var total int64
	errQuery := users_db.Client.QueryRow(queryTotalSize).Scan(&total)
	if errQuery != nil {
		return nil, errors.NewInternalServerError(err.Error()), 0
	}

	rows, err := stmt.Query(status, perPage, (page-1)*perPage)

	if err != nil {
		return nil, errors.NewInternalServerError(err.Error()), 0
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id,
			&user.FirstName,
			&user.LastName,
			&user.NationalID,
			&user.DateCreated,
			&user.Status); err != nil {
			return nil, mysql_utils.ParseError(err), 0
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status)), 0
	}
	return results, nil, total

}
