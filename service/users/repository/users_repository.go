package repository

import (
	"context"
	"fmt"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/users"
	"pheet-fiber-backend/service/users/pattern"

	"github.com/jmoiron/sqlx"
)

type usersRepository struct {
	psqlDB *sqlx.DB
}

func NewUsersRepository(psqlDB *sqlx.DB) users.IUsersRepository {
	return usersRepository{
		psqlDB: psqlDB,
	}
}

func (u usersRepository) InsertUser(userReq *models.UserRegisterReq, isAdmin bool) (*models.UserPassport, error) {
	iUserPattern := pattern.InsertUser(u.psqlDB, userReq, isAdmin)

	var err error
	if isAdmin {
		iUserPattern, err = iUserPattern.Admin()
		if err != nil {
			return nil, fmt.Errorf("insert admin failed: %v", err)
		}
		
	} else {
		iUserPattern, err = iUserPattern.Customer()
		if err != nil {
			return nil, fmt.Errorf("insert customer failed: %v", err)
		}
	}

	user, err := iUserPattern.Result()
	if err != nil {
		return nil, fmt.Errorf("insert user failed: %v", err)
	}

	return user, err
}

func (u usersRepository) FindOneUserByEmail(ctx context.Context, email string) (*models.UserCredentialCheck, error) {
	sql := `
		SELECT
			id,
			username,
			password,
			email,
			role_id
		FROM 
			"users"
		WHERE
			"email" = $1
	`

	var uCredential = new(models.UserCredentialCheck)
	if err := u.psqlDB.GetContext(ctx, uCredential, sql, email); err != nil {
		return nil, err
	}

	return uCredential, nil
}

func (u usersRepository) InsertOauth(ctx context.Context, req *models.UserPassport) error {
	sql := `
		INSERT INTO "oauth" (
			"user_id",
			"access_token",
			"refresh_token"
		)
		VALUES (
			$1,
			$2,
			$3
		)
		RETURNING "id";
	`
	if err := u.psqlDB.QueryRowContext(
			ctx, 
			sql,
			req.User.Id,
			req.Token.AccessToken,
			req.Token.RefreshToken,
	).Scan(&req.Token.Id); err != nil {
		return fmt.Errorf("insert oauth failed: %v", err)
	}

	return nil
}