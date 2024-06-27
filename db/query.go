package db

import (
	"context"
	"picturethisai/types"

	"github.com/google/uuid"
)

func CreateAccount(account *types.Account) error {
	_, err := Bun.NewInsert().NewInsert().
		Model(account).
		Exec(context.Background())

	return err
}

func GetAccountByUserID(userId uuid.UUID) (types.Account, error) {
	var account types.Account
	err := Bun.NewSelect().
		Model(&account).
		Where("user_id = ?", userId).
		Scan(context.Background())
	return account, err
}

func UpdateAccount(account *types.Account) error {

	_, err := Bun.NewUpdate().
		Model(account).
		WherePK().
		Exec(context.Background())

	return err
}

func GetImagesByUserID(userID uuid.UUID) ([]types.Image, error) {
	var images []types.Image
	err := Bun.NewSelect().
		Model(&images).
		Where("deleted = ?", false).
		Where("user_id = ?", userID).
		Scan(context.Background())

	return images, err
}
