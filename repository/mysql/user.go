package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"mymodule/entity"
)

func (d *MysqlDB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	//todo implemnet IsPhoneNumberUnique
	user := entity.User{}
	var createdAt []uint8
	row := d.db.QueryRow(`select * from users where phone_number=?`, phoneNumber)
	err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}
		return false, fmt.Errorf("mysql query row scan error: %w", err)
	}
	return false, fmt.Errorf("phone number %s already exists", phoneNumber)
}

func (d *MysqlDB) RegisterUser(user entity.User) (entity.User, error) {
	res, dErr := d.db.Exec(`insert into users(name,phone_number)values(?,?)`, user.Name, user.PhoneNumber)
	if dErr != nil {
		return entity.User{}, fmt.Errorf("failed to insert user to db :%w\n", dErr)
	}

	id, lErr := res.LastInsertId()
	if lErr != nil {
		return entity.User{}, lErr
	}
	user.ID = uint(id)

	return user, nil
}
