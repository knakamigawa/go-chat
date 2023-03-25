package repository

import (
	"context"
	"database/sql"
	"go-chat-ai-server/domain/model"
	"go-chat-ai-server/infra/database/chat_db"
	"golang.org/x/crypto/bcrypt"
)

func ProvideDbUserMailPassword(db *sql.DB, queries *chat_db.Queries) *DbUserMailPassword {
	return &DbUserMailPassword{
		db:      db,
		queries: queries,
	}
}

type DbUserMailPassword struct {
	db      *sql.DB
	queries *chat_db.Queries
}

func (r DbUserMailPassword) Of(ctx context.Context, email model.Email, password model.Password) (model.User, error) {
	// メールアドレスだけで取得
	findUser, err := r.queries.FindUserByEmail(ctx, email.String())
	if err != nil {
		return model.User{}, err
	}
	// パスワードチェック
	err = bcrypt.CompareHashAndPassword([]byte(findUser.PasswordHash), []byte(password.String()))
	if err != nil {
		return model.User{}, err
	}

	return model.MakeUser(findUser.ID, findUser.LoginName), nil
}

func (r DbUserMailPassword) Save(ctx context.Context, user model.UserEmailPassword) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := r.queries.WithTx(tx)
	param := chat_db.CreateUserParams{
		ID:        user.User().ID(),
		LoginName: user.User().Name(),
	}
	_, err = qtx.CreateUser(ctx, param)
	if err != nil {
		return err
	}

	authParam := chat_db.CreateUserEmailPasswordParams{
		UserID:       user.User().ID(),
		Email:        user.Email(),
		PasswordHash: user.HashPassword(),
	}
	_, err = qtx.CreateUserEmailPassword(ctx, authParam)
	if err != nil {
		return err
	}

	return tx.Commit()
}
