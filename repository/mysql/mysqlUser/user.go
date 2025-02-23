package mysqlUser

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"mymodule/entity"
	"mymodule/pkg/richerr"
	"mymodule/repository/mysql"
)

func (d *DB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	//todo implemnet IsPhoneNumberUnique
	user := entity.User{}
	var createdAt []uint8
	var role string
	row := d.conn.NewConn().QueryRow(`select * from users where phone_number=?`, phoneNumber)
	err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &createdAt, &user.Password, &role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}
		return false, fmt.Errorf("mysql query row scan error: %w", err)
	}

	return false, richerr.New().
		SetKind(richerr.KindUnexpected).
		SetOperation("mysql.IsPhoneNumberUnique").
		SetMsg("\"phone number is not unique\"")
}

func (d *DB) RegisterUser(user entity.User) (entity.User, error) {
	hashedPassword, gErr := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if gErr != nil {
		return entity.User{}, richerr.New().
			SetMsg("failed to hash password").
			SetOperation("mysql.RegisterUser").
			SetKind(richerr.KindUnexpected).
			SetWrappedErr(gErr)
	}

	res, dErr := d.conn.NewConn().Exec(`insert into users(name,phone_number,password,role)values(?,?,?,?)`, user.Name, user.PhoneNumber, hashedPassword, user.Role.String())
	if dErr != nil {
		return entity.User{}, richerr.New().
			SetMsg("failed to insert user").
			SetOperation("mysql.RegisterUser").
			SetKind(richerr.KindUnexpected).
			SetWrappedErr(dErr)
	}

	id, lErr := res.LastInsertId()
	if lErr != nil {
		return entity.User{}, richerr.New().
			SetMsg("failed to give id to record").
			SetOperation("mysql.RegisterUser").
			SetKind(richerr.KindUnexpected).
			SetWrappedErr(lErr)
	}
	user.ID = uint(id)

	return user, nil
}

func (d *DB) GetUserByPhoneNumber(phoneNumber string) (entity.User, error) {
	user := entity.User{}
	var role string

	row := d.conn.NewConn().QueryRow(`select id,name,phone_number,password,role from users where phone_number=?`, phoneNumber)
	if sErr := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Password, &role); sErr != nil {
		if errors.Is(sErr, sql.ErrNoRows) {
			return entity.User{}, richerr.New().
				SetKind(richerr.KindUnexpected).
				SetMsg("user not found").
				SetOperation("mysql.GetUserByPhoneNumber").
				SetWrappedErr(sErr)
		}
		return entity.User{}, richerr.New().
			SetKind(richerr.KindUnexpected).
			SetMsg("failed to scan query").
			SetOperation("mysql.GetUserByPhoneNumber").
			SetWrappedErr(sErr)
	}
	user.Role.RoleId(role)
	//user.Role = tt(role)
	return user, nil

}

// get user information acording to claim.id
func (d *DB) GetUserById(id uint) (entity.User, error) {
	row := d.conn.NewConn().QueryRow(`select * from users where id=?`, id)
	user, sErr := ScanUser(row)
	if sErr != nil {
		return entity.User{}, richerr.New().
			SetKind(richerr.KindUnexpected).
			SetMsg("failed to scan query").
			SetOperation("mysql.GetUserById").
			SetWrappedErr(sErr)
	}
	return user, nil
}

func ScanUser(row mysql.Scanner) (entity.User, error) {
	user := entity.User{}
	var role string
	var createdAt []uint8
	sErr := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &createdAt, &user.Password, &role)
	user.Role.RoleId(role)
	return user, sErr
}

func tt(role string) entity.Role {
	switch role {
	case "user":
		return entity.UserRole
	case "admin":
		return entity.Admin
	default:
		return 0
	}
}
