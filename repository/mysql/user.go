package mysql

//
//import (
//	"database/sql"
//	"errors"
//	"fmt"
//	"golang.org/x/crypto/bcrypt"
//	"mymodule/entity"
//	"mymodule/pkg/richerr"
//)
//
//func (d *MysqlDB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
//	//todo implemnet IsPhoneNumberUnique
//	user := entity.User{}
//	var role string
//	var createdAt []uint8
//	row := d.db.QueryRow(`select * from users where phone_number=?`, phoneNumber)
//	err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &createdAt, &user.Password, &role)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return true, nil
//		}
//		return false, fmt.Errorf("mysql query row scan error: %w", err)
//	}
//
//	return false, richerr.New().
//		SetKind(richerr.KindUnexpected).
//		SetOperation("mysql.IsPhoneNumberUnique").
//		SetMsg("\"phone number is not unique\"")
//}
//
//func (d *MysqlDB) RegisterUser(user entity.User) (entity.User, error) {
//	hashedPassword, gErr := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
//	if gErr != nil {
//		return entity.User{}, richerr.New().
//			SetMsg("failed to hash password").
//			SetOperation("mysql.RegisterUser").
//			SetKind(richerr.KindUnexpected).
//			SetWrappedErr(gErr)
//	}
//
//	res, dErr := d.db.Exec(`insert into users(name,phone_number,password)values(?,?,?)`, user.Name, user.PhoneNumber, hashedPassword)
//	if dErr != nil {
//		return entity.User{}, richerr.New().
//			SetMsg("failed to insert user").
//			SetOperation("mysql.RegisterUser").
//			SetKind(richerr.KindUnexpected).
//			SetWrappedErr(dErr)
//	}
//
//	id, lErr := res.LastInsertId()
//	if lErr != nil {
//		return entity.User{}, richerr.New().
//			SetMsg("failed to give id to record").
//			SetOperation("mysql.RegisterUser").
//			SetKind(richerr.KindUnexpected).
//			SetWrappedErr(lErr)
//	}
//	user.ID = uint(id)
//
//	return user, nil
//}
//
//func (d *MysqlDB) GetUserByPhoneNumber(phoneNumber string) (entity.User, error) {
//	user := entity.User{}
//	row := d.db.QueryRow(`select id,name,phone_number,password from users where phone_number=?`, phoneNumber)
//	if sErr := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Password); sErr != nil {
//		if errors.Is(sErr, sql.ErrNoRows) {
//			return entity.User{}, richerr.New().
//				SetKind(richerr.KindUnexpected).
//				SetMsg("user not found").
//				SetOperation("mysql.GetUserByPhoneNumber").
//				SetWrappedErr(sErr)
//		}
//		return entity.User{}, richerr.New().
//			SetKind(richerr.KindUnexpected).
//			SetMsg("failed to scan query").
//			SetOperation("mysql.GetUserByPhoneNumber").
//			SetWrappedErr(sErr)
//	}
//
//	return user, nil
//
//}
//
//func (d *MysqlDB) GetUserById(id uint) (entity.User, error) {
//	row := d.db.QueryRow(`select * from users where id=?`, id)
//	user, sErr := ScanUser(row)
//	if sErr != nil {
//		return entity.User{}, richerr.New().
//			SetKind(richerr.KindUnexpected).
//			SetMsg("failed to scan query").
//			SetOperation("mysql.GetUserById").
//			SetWrappedErr(sErr)
//	}
//	return user, nil
//}
//
//func ScanUser(row Scanner) (entity.User, error) {
//	user := entity.User{}
//	var createdAt []uint8
//	sErr := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &createdAt, &user.Password)
//	return user, sErr
//}
