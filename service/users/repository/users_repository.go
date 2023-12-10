package repository

import (
	"context"
	"fmt"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/users"
	"pheet-fiber-backend/service/users/pattern"

	"github.com/Pheethy/sqlx"
	"github.com/gofrs/uuid"
)

type usersRepository struct {
	psqlDB *sqlx.DB
}

func NewUsersRepository(psqlDB *sqlx.DB) users.IUsersRepository {
	return &usersRepository{
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

func (u usersRepository) FetchUserProfile(ctx context.Context, id string) (*models.Users, error) {
	sql := `
		SELECT
			id,
			username,
			email,
			role_id
		FROM 
			"users"
		WHERE
			"id" = $1
	`

	var uCredential = new(models.Users)
	if err := u.psqlDB.GetContext(ctx, uCredential, sql, id); err != nil {
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

func (u usersRepository) FetchOneOauth(ctx context.Context, reToken string) (*models.Oauth, error) {
	sql := `
		SELECT 
			id,
			user_id
		FROM 
			"oauth"
		WHERE
			"refresh_token" = $1;
	`
	var oauth = new(models.Oauth)
	if err := u.psqlDB.GetContext(ctx, oauth, sql, reToken); err != nil {
		return nil, fmt.Errorf("fetch oauth failed: %v", err)
	}

	return oauth, nil
}

func (u usersRepository) UpdateOauth(ctx context.Context, req *models.UserToken) error {
	tx, err := u.psqlDB.Beginx()
	if err != nil {
		panic(err)
	}

	var id uuid.UUID
	if req != nil {
		id = uuid.FromStringOrNil(req.Id)
	}

	sql := `
		UPDATE 
			oauth 
		SET
			access_token=$1::text,
			refresh_token=$2::text
		WHERE
			id=$3::uuid
		
	`
	stmt, err := tx.PreparexContext(ctx, sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx,
		req.AccessToken,
		req.RefreshToken,
		id,
	); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (u usersRepository) DeleteOauth(ctx context.Context, oId string) error {
	sql := `
		DELETE FROM
			"oauth"
		WHERE
			"id" = $1;
	`

	if _, err := u.psqlDB.ExecContext(ctx, sql, oId); err != nil {
		return fmt.Errorf("oauth not found")
	}

	return nil
}
