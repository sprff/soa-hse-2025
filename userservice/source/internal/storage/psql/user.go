package psql

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"social/shared/models"
	"social/userservice/internal/utils"
)

func (p *PsqlStorage) userWithLoginExists(ctx context.Context, user models.User) (bool, error) {

	res, err := p.db.Query("SELECT * FROM users WHERE id!=$1 AND login=$2 ", user.ID, user.Login)
	if err != nil {
		return false, fmt.Errorf("can't select: %w", err)
	}

	return res.Next(), nil
}

func (p *PsqlStorage) CreateUser(ctx context.Context, user models.User) (models.UserID, error) {
	slog.DebugContext(ctx, "PSQL CreateUser")

	exists, err := p.userWithLoginExists(ctx, user)
	if err != nil {
		return models.UserID(""), fmt.Errorf("can't check user: %w", err)
	}
	if exists {
		slog.DebugContext(ctx, "User exist")
		return models.UserID(""), models.ErrUserLoginAlreadyExists
	}

	user.ID = models.UserID(utils.GenerateUUID())
	slog.DebugContext(ctx, "Gen ID", "id", user.ID)

	_, err = p.db.NamedExec(`
	INSERT INTO users
	(id, login, password, email,
	name, surname, phone, dob, created_at, updated_at)
	VALUES
	(:id, :login, :password, :email,
	:name, :surname, :phone, :dob, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`, &user)
	if err != nil {
		return models.UserID(""), fmt.Errorf("can't insert user: %w", err)
	}

	return user.ID, nil
}

func (p *PsqlStorage) UpdateUser(ctx context.Context, user models.User) error {
	slog.DebugContext(ctx, "PSQL UpdateUser")

	res, err := p.db.NamedExec(`
	UPDATE users
	SET email=:email,
	name=:name, surname=:surname, phone=:phone, dob=:dob, updated_at=CURRENT_TIMESTAMP
	WHERE id=:id`, &user)
	if err != nil {
		return fmt.Errorf("can't insert user: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("can't get rows affected")
	}

	if affected == 0 {
		return models.ErrUserNotFound
	}

	return nil
}

func (p *PsqlStorage) GetUserByID(ctx context.Context, id models.UserID) (models.User, error) {
	slog.DebugContext(ctx, "PSQL GetUserByID")

	var user models.User
	err := p.db.Get(&user, "SELECT * FROM users WHERE id=$1", id)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return models.User{}, models.ErrUserNotFound
		default:
			return models.User{}, fmt.Errorf("can't select: %w", err)
		}
	}
	return user, nil
}

func (p *PsqlStorage) GetUserByLogin(ctx context.Context, login string) (models.User, error) {
	slog.DebugContext(ctx, "PSQL GetUserByLogin")

	var user models.User
	err := p.db.Get(&user, "SELECT * FROM users WHERE login=$1", login)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return models.User{}, models.ErrUserNotFound
		default:
			return models.User{}, fmt.Errorf("can't select: %w", err)
		}
	}
	slog.Debug("PSQL GET User", "user", user)
	return user, nil
}
