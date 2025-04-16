package pgxrepo_test

import (
	"context"
	"rest-service/internal/entities"
	"rest-service/internal/infrastructure/repository/pgxrepo"
	"testing"

	"github.com/google/uuid"
)

func TestCreateAndGetUser(t *testing.T) {
	authRepo := pgxrepo.NewPgxAuthRepository(TestPool)
	validUser := &entities.User{ID: uuid.New(), RefreshToken: "refreshTokenHash"}

	t.Run("Create and get valid user", func (t *testing.T) {
		if err := authRepo.CreateUser(context.Background(), validUser); err != nil {
			t.Errorf("got an error while creating user; err: %v", err)
		}

		fetchedUser, err := authRepo.GetUser(context.Background(), validUser.ID.String())
		if err != nil {
			t.Errorf("got an error while fetching user")
		}
		if fetchedUser.RefreshToken != validUser.RefreshToken {
			t.Errorf("invalid data being saved to db")
		}		
	})
	
	t.Run("Get invalid user", func (t *testing.T){
		invalidUser := &entities.User{ID: uuid.New()}

		_, err := authRepo.GetUser(context.Background(), invalidUser.ID.String())
		if err == nil {
			t.Errorf("expected an error while getting non-existent user, got nil")
		}
	})
}

func TestUpdateUser(t *testing.T) {
	authRepo := pgxrepo.NewPgxAuthRepository(TestPool)
	validUser := &entities.User{ID: uuid.New(), RefreshToken: "refreshToken"}


	t.Run("Update valid user", func (t *testing.T) {
		if err := authRepo.CreateUser(context.Background(), validUser); err != nil {
			t.Errorf("got an error while creating user; err: %v", err)
		}

		updatedUser := &entities.User{ID: validUser.ID, RefreshToken: "newToken"}
		if err := authRepo.UpdateUser(context.Background(), updatedUser); err != nil {
			t.Errorf("got an error while updating user; err: %v", err)
		}

		fetchedUser, err := authRepo.GetUser(context.Background(), validUser.ID.String())
		if err != nil {
			t.Errorf("got an error while fetching user; err: %v", err)
		}
		if fetchedUser.RefreshToken != updatedUser.RefreshToken {
			t.Errorf("invalid data being set to user")
		}
	})

	t.Run("Update invalid user", func(t *testing.T) {
		invalidUser:= &entities.User{ID: uuid.New(), RefreshToken: "token"}

		if err := authRepo.UpdateUser(context.Background(), invalidUser); err == nil {
			t.Errorf("expected an error while updating non-existent user, got nil")
		}
	})
}