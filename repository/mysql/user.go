package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"mymodule/entity"
)

func (d *MysqlDB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	//todo implemnet IsPhoneNumberUnique
	user := entity.User{}
	var createdAt []uint8
	row := d.db.QueryRow(`select * from users where phone_number=?`, phoneNumber)
	err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Password, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}
		return false, fmt.Errorf("mysql query row scan error: %w", err)
	}
	return false, fmt.Errorf("phone number %s already exists", phoneNumber)
}

func (d *MysqlDB) RegisterUser(user entity.User) (entity.User, error) {
	hashedPassword, gErr := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if gErr != nil {
		return entity.User{}, fmt.Errorf("failed to hash password: %v\n", gErr)
	}

	res, dErr := d.db.Exec(`insert into users(name,phone_number,password)values(?,?,?)`, user.Name, user.PhoneNumber, hashedPassword)
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

func (d *MysqlDB) GetUserByPhoneNumber(phoneNumber string) (entity.User, error) {
	user := entity.User{}
	row := d.db.QueryRow(`select id,name,phone_number,password from users where phone_number=?`, phoneNumber)
	if sErr := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Password); sErr != nil {
		if errors.Is(sErr, sql.ErrNoRows) {
			return entity.User{}, sErr
		}
		return entity.User{}, fmt.Errorf("failed to scan query:%v\n", sErr)
	}

	return user, nil

}
