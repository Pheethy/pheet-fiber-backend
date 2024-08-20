package patterns

import (
	"context"
	"errors"
	"fmt"
	"pheet-fiber-backend/constants"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/users"

	"github.com/Pheethy/psql/orm"
	"github.com/Pheethy/sqlx"
)

type usersPattern struct {
	ctx    context.Context
	req    *models.User
	psqlDB *sqlx.DB
}

func NewUsersPatterns(ctx context.Context, req *models.User, psqlDB *sqlx.DB) users.IUsersPattern {
	return usersPattern{
		ctx:    ctx,
		req:    req,
		psqlDB: psqlDB,
	}
}

func (u usersPattern) Customer() (users.IUsersPattern, error) {
	sql := `
		INSERT INTO users (id, username, password, email, role_id, created_at, updated_at)
		VALUES (
			$1::uuid,
			$2::text,
			$3::text,
			$4::text,
			$5::integer,
			$6::timestamp,
			$7::timestamp
		)
	`

	stmt, err := u.psqlDB.PreparexContext(u.ctx, sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(u.ctx,
		u.req.Id,
		u.req.Username,
		u.req.Password,
		u.req.Email,
		constants.USER_ROLE_CUSTOMER,
		u.req.CreatedAt,
		u.req.UpdatedAt,
	); err != nil {
		switch err.Error() {
		case constants.ERROR_PQ_UNIQUE_USERNAME:
			return nil, errors.New(constants.ERROR_USERNAME_HAS_BEEN_USED)
		case constants.ERROR_PQ_UNIQUE_EMAIL:
			return nil, errors.New(constants.ERROR_EMAIL_HAS_BEEN_USED)
		}
		return nil, err
	}

	return u, nil
}

func (u usersPattern) Admin() (users.IUsersPattern, error) {
	sql := `
		INSERT INTO users (id, username, password, email, role_id, created_at, updated_at)
		VALUES (
			$1::uuid,
			$2::text,
			$3::text,
			$4::text,
			$5::integer,
			$6::timestamp,
			$7::timestamp
		)
	`

	stmt, err := u.psqlDB.PreparexContext(u.ctx, sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(u.ctx,
		u.req.Id,
		u.req.Username,
		u.req.Password,
		u.req.Email,
		constants.USER_ROLE_ADMIN,
		u.req.CreatedAt,
		u.req.UpdatedAt,
	); err != nil {
		switch err.Error() {
		case constants.ERROR_PQ_UNIQUE_USERNAME:
			return nil, errors.New(constants.ERROR_USERNAME_HAS_BEEN_USED)
		case constants.ERROR_PQ_UNIQUE_EMAIL:
			return nil, errors.New(constants.ERROR_EMAIL_HAS_BEEN_USED)
		}
		return nil, err
	}

	return u, nil
}

func (u usersPattern) Result() (*models.UserPassport, error) {
	sql := fmt.Sprintf(`
		SELECT
			%s
		FROM
		(
			SELECT
				users.*,
				roles.title "role_title"
			FROM
				users
			JOIN
				roles
			ON
				users.role_id = roles.id
		) AS users
		WHERE
			id=$1::uuid
	`, orm.GetSelector(models.User{}))

	stmt, err := u.psqlDB.PreparexContext(u.ctx, sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryxContext(u.ctx, u.req.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user, err := u.ormOneUser(u.ctx, rows)
	if err != nil {
		return nil, err
	}

	userPass := &models.UserPassport{
		User: &models.User{
			Id:        user.Id,
			Username:  user.Username,
			Password:  "",
			Email:     user.Email,
			RoleId:    user.RoleId,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Token: nil,
	}

	return userPass, nil
}

func (u usersPattern) ormOneUser(ctx context.Context, rows *sqlx.Rows) (*models.User, error) {
	mapper, err := orm.OrmContext(ctx, new(models.User), rows, orm.NewMapperOption())
	if err != nil {
		return nil, err
	}
	users := mapper.GetData().([]*models.User)
	if len(users) == 0 {
		return nil, fmt.Errorf(constants.ERROR_USER_NOT_FOUND)
	}

	return users[0], nil
}
