package repository

import (
	"context"
	"errors"
	"fmt"
	"pheet-fiber-backend/constants"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/users"
	"pheet-fiber-backend/service/users/patterns"

	"github.com/Pheethy/psql/orm"
	"github.com/Pheethy/sqlx"
	"github.com/gofrs/uuid"
)

type usersRepository struct {
	psqlDB *sqlx.DB
}

func NewUsersRepository(psqlDB *sqlx.DB) users.IUsersRepository {
	return usersRepository{
		psqlDB: psqlDB,
	}
}

func (u usersRepository) InsertUser(ctx context.Context, userReq *models.User, isAdmin bool) (*models.UserPassport, error) {
	pattern := patterns.NewUsersPatterns(ctx, userReq, u.psqlDB)
	/* Get Insert User*/
	var err error
	/* Checking Admin */
	if isAdmin {
		pattern, err = pattern.Admin()
		if err != nil {
			return nil, err
		}
	} else {
		pattern, err = pattern.Customer()
		if err != nil {
			return nil, err
		}
	}
	/* Get User Passport */
	userPass, err := pattern.Result()
	if err != nil {
		return nil, err
	}

	return userPass, nil
}

func (u usersRepository) UpsertOAuth(ctx context.Context, req *models.OAuth) error {
	sql := `
		INSERT INTO oauth (id, user_id, access_token, refresh_token, created_at, updated_at)
		VALUES (
			$1::uuid,
			$2::uuid,
			$3::text,
			$4::text,
	    $5::timestamp,
			$6::timestamp
		)
		ON CONFLICT (id)
		DO UPDATE SET
			refresh_token=$7::text
	`
	stmt, err := u.psqlDB.PreparexContext(ctx, sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx,
		/* Create */
		req.Id,
		req.UserId,
		req.AccessToken,
		req.RefreshToken,
		req.CreatedAt,
		req.UpdatedAt,
		/* update */
		req.RefreshToken,
	); err != nil {
		return err
	}

	return nil
}

func (u usersRepository) FetchOneUserByEmail(ctx context.Context, email string) (*models.User, error) {
	sql := `
		SELECT
			users.id,
			users.username,
			users.password,
			users.role_id,
			users.email,
			users.created_at,
			users.updated_at
		FROM
			users
		WHERE
			LOWER(users.email) = LOWER($1::text)
	`

	stmt, err := u.psqlDB.PreparexContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryxContext(ctx, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user, err := u.ormOneUser(ctx, rows)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u usersRepository) FetchOneUserById(ctx context.Context, userId *uuid.UUID) (*models.User, error) {
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
			users.id=$1::uuid
	`,
		orm.GetSelector(models.User{}),
	)

	stmt, err := u.psqlDB.PreparexContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryxContext(ctx, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user, err := u.ormOneUser(ctx, rows)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u usersRepository) FetchOAuthByRefreshToken(ctx context.Context, refreshToken string) (*models.OAuth, error) {
	sql := fmt.Sprintf(`
		SELECT
			%s
		FROM
			oauth
		WHERE
			refresh_token=$1::text
	`,
		orm.GetSelector(models.OAuth{}))

	stmt, err := u.psqlDB.PreparexContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryxContext(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	oauth, err := u.ormOneOAuth(ctx, rows)
	if err != nil {
		return nil, err
	}

	return oauth, nil
}

func (u usersRepository) DeleteOAuth(ctx context.Context, oauthId *uuid.UUID) error {
	sql := `
		DELETE
		FROM
			oauth
		WHERE
			oauth.id=$1::uuid
	`
	stmt, err := u.psqlDB.PreparexContext(ctx, sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, oauthId); err != nil {
		return err
	}

	return nil
}

func (u usersRepository) ormOneOAuth(ctx context.Context, rows *sqlx.Rows) (*models.OAuth, error) {
	mapper, err := orm.OrmContext(ctx, new(models.OAuth), rows, orm.NewMapperOption())
	if err != nil {
		return nil, err
	}
	oauths := mapper.GetData().([]*models.OAuth)
	if len(oauths) == 0 {
		return nil, errors.New(constants.ERROR_OAUTH_NOT_FOUND)
	}
	return oauths[0], nil
}

func (u usersRepository) ormOneUser(ctx context.Context, rows *sqlx.Rows) (*models.User, error) {
	mapper, err := orm.OrmContext(ctx, new(models.User), rows, orm.MapperOption{})
	if err != nil {
		return nil, err
	}
	users := mapper.GetData().([]*models.User)
	if len(users) == 0 {
		return nil, errors.New(constants.ERROR_USER_NOT_FOUND)
	}

	return users[0], err
}
