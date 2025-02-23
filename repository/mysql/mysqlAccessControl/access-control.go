package mysqlAccessControl

import (
	"database/sql"
	"fmt"
	"mymodule/entity"
	"mymodule/pkg/richerr"
	"mymodule/repository/mysql"
	"strings"
)

func (m *DB) GetUserAcl(userid uint, role entity.Role) ([]entity.PermissionTitle, error) {
	//get acls for role
	//get acls for user
	//merge role & user acls
	rows, qErr := m.conn.NewConn().Query(`select * from access_controls where actor_id = ? and actor_type=?`, role, entity.RoleActorType)
	if qErr != nil {
		return nil, richerr.New().
			SetOperation("mysqlAccessControl.GetUserAcl").
			SetMsg(qErr.Error()).
			SetKind(richerr.KindUnexpected)
	}

	defer rows.Close()

	roleAcls := []entity.AccessControl{}

	for rows.Next() {
		acl, sErr := ScanRoleAcl(rows)
		if sErr != nil {
			return nil, sErr
		}
		roleAcls = append(roleAcls, acl)
	}

	if err := rows.Err(); err != nil {
		return nil, richerr.New().
			SetOperation("mysqlAccessControl.GetUserAcl").
			SetMsg(err.Error()).
			SetKind(richerr.KindUnexpected)
	}

	userRows, qErr := m.conn.NewConn().Query(`select * from access_controls where actor_id = ? and actor_type=?`, userid, entity.UserActorType)
	if qErr != nil {
		return nil, richerr.New().
			SetOperation("mysqlAccessControl.GetUserAcl").
			SetMsg(qErr.Error()).
			SetKind(richerr.KindUnexpected)
	}

	defer rows.Close()

	if err := userRows.Err(); err != nil {
		return nil, richerr.New().
			SetOperation("mysqlAccessControl.GetUserAcl").
			SetMsg(err.Error()).
			SetKind(richerr.KindUnexpected)
	}

	userAcls := []entity.AccessControl{}

	for rows.Next() {
		acl, sErr := ScanRoleAcl(userRows)
		if sErr != nil {
			return nil, sErr
		}
		userAcls = append(userAcls, acl)
	}

	allAcls := append(roleAcls, userAcls...)

	var permissionIds []interface{}

	//for _, acl := range allAcls {
	//if !doesExist(permissionIds,acl.PermissionId) {
	//	permissionIds = append(permissionIds, acl.PermissionId)
	//}
	//}

	for _, acl := range allAcls {
		isExist := false
		for _, permissionId := range permissionIds {
			if acl.PermissionId == permissionId {
				isExist = true
				break
			}
		}
		if !isExist {
			permissionIds = append(permissionIds, acl.PermissionId)
		}
	}

	args := make([]string, len(permissionIds))
	for i := range args {
		args[i] = "?"
	}

	fmt.Printf("userAcls: %+v\n,roleAcls:%+v,allAcls:%+v\n", userAcls, roleAcls, allAcls)

	queryParams := strings.Join(args, ",")
	//fmt.Println("log check", fmt.Sprintf("select * from permissions where id in(%s)", queryParams))

	pRows, pErr := m.conn.NewConn().Query(fmt.Sprintf("select * from permissions where id in(%s)", queryParams), permissionIds...)
	if pErr != nil {
		return nil, richerr.New().
			SetOperation("mysqlAccessControl.GetUserAcl").
			SetMsg(pErr.Error()).
			SetKind(richerr.KindUnexpected)
	}

	permissionTitles := make([]entity.PermissionTitle, len(permissionIds))

	for pRows.Next() {
		permission, sErr := ScanPermission(pRows)
		if sErr != nil {
			return nil, sErr
		}
		permissionTitles = append(permissionTitles, permission.Title)
	}

	if err := pRows.Err(); err != nil {
		return nil, richerr.New().
			SetOperation("mysqlAccessControl.GetUserAcl").
			SetMsg(err.Error()).
			SetKind(richerr.KindUnexpected)
	}

	return permissionTitles, nil
}

func ScanRoleAcl(rows mysql.Scanner) (entity.AccessControl, error) {
	var acl entity.AccessControl
	var createdAt []uint8
	err := rows.Scan(&acl.Id, &acl.ActorId, &acl.ActorType, &acl.PermissionId, &createdAt)
	if err != nil {
		return entity.AccessControl{}, richerr.New().
			SetOperation("mysqlAccessControl.ScanRoleAcl").
			SetMsg(err.Error()).
			SetKind(richerr.KindUnexpected)
	}
	return acl, nil
}

func ScanPermission(rows *sql.Rows) (entity.Permission, error) {
	var permission entity.Permission
	var createdAt []uint8
	err := rows.Scan(&permission.Id, &permission.Title, &createdAt)
	if err != nil {
		return entity.Permission{}, richerr.New().
			SetOperation("mysqlAccessControl.ScanPermission").
			SetMsg(err.Error()).
			SetKind(richerr.KindUnexpected)
	}
	return permission, nil
}

func doesExist(permissionIds []uint, value uint) bool {
	for _, p := range permissionIds {
		if value == p {
			return true
		}
	}
	return false
}
