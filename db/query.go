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

func CreateImage(image *types.Image) error {
	_, err := Bun.NewInsert().
		Model(image).
		Exec(context.Background())
	return err
}

func UpdateImage(image *types.Image) error {
	_, err := Bun.NewUpdate().
		Model(image).
		WherePK().
		Exec(context.Background())

	return err
}

func GetImageByID(id int) (types.Image, error) {
	var image types.Image
	err := Bun.NewSelect().
		Model(&image).
		Where("id = ?", id).
		Scan(context.Background())

	return image, err
}

func GetImagesByUserID(userID uuid.UUID) ([]types.Image, error) {
	var images []types.Image
	err := Bun.NewSelect().
		Model(&images).
		Where("deleted = ?", false).
		Where("user_id = ?", userID).
		Order("created_at desc").
		Scan(context.Background())

	return images, err
}

func GetImagesByBatchID(batchID uuid.UUID) ([]types.Image, error) {
	var images []types.Image
	err := Bun.NewSelect().
		Model(&images).
		Where("batch_id = ?", batchID).
		Scan(context.Background())

	return images, err
}
