package repository

import (
	"context"
	"errors"
	"fmt"
	"pheet-fiber-backend/constants"
	"pheet-fiber-backend/middleware"
	"pheet-fiber-backend/models"

	"github.com/Pheethy/psql/orm"
	"github.com/Pheethy/sqlx"
	"github.com/gofrs/uuid"
)

type middlewareRepository struct {
	psqlDB *sqlx.DB
}

func NewMiddlewareRepository(psqlDB *sqlx.DB) middleware.IMiddlewareRepository {
	return middlewareRepository{
		psqlDB: psqlDB,
	}
}

func (m middlewareRepository) FindAccessToken(ctx context.Context, userId *uuid.UUID, accessToken string) bool {
	var ok bool
	sql := `
		SELECT
			(CASE WHEN count(*) = 1 THEN TRUE ELSE FALSE END)
		FROM
			oauth
		WHERE
			oauth.user_id = $1::uuid
		AND
			oauth.access_token = $2::text 
	`
	stmt, err := m.psqlDB.PreparexContext(ctx, sql)
	if err != nil {
		return false
	}
	defer stmt.Close()

	if err := stmt.GetContext(ctx, &ok, userId, accessToken); err != nil {
		return false
	}

	return ok
}

func (m middlewareRepository) FetchRoles(ctx context.Context) ([]*models.Roles, error) {
	sql := fmt.Sprintf(`
		SELECT
			%s
		FROM
			roles
		ORDER BY roles.id DESC;
	`,
		orm.GetSelector(models.Roles{}),
	)

	stmt, err := m.psqlDB.PreparexContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryxContext(ctx)
	if err != nil {
		return nil, err
	}

	roles, err := m.ormRoles(ctx, rows)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (m middlewareRepository) ormRoles(ctx context.Context, rows *sqlx.Rows) ([]*models.Roles, error) {
	mapper, err := orm.OrmContext(ctx, new(models.Roles), rows, orm.NewMapperOption())
	if err != nil {
		return nil, err
	}
	roles := mapper.GetData().([]*models.Roles)
	if len(roles) == 0 {
		return nil, errors.New(constants.ERROR_ROLES_NOT_FOUND)
	}

	return roles, nil
}
