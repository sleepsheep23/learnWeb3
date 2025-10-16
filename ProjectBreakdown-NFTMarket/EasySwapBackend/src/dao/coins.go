package dao

import (
	"github.com/ProjectsTask/EasySwapBackend/src/types/v1"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (d *Dao) AddCoinsMetadata(metadata []types.CoinsSearchApiResponse) error {
	for _, coin := range metadata {
		coinMeta := types.CoinsMetadata{
			Id:        coin.Id,
			Symbol:    coin.Symbol,
			Name:      coin.Name,
			ApiSymbol: coin.ApiSymbol,
		}
		if err := d.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&coinMeta); err != nil {
			// log the error and continue with the next coin
			continue
		}
	}
	return nil
}

func (d *Dao) ListUserFavoriteCoinMetadata(address string) ([]types.CoinsMetadata, error) {
	var coins []types.CoinsMetadata
	err := d.DB.Model(types.CoinsMetadata{}).
		Joins("JOIN cg_user_coins_mapping m ON cg_coins_metadata.id = m.coin_id").
		Joins("JOIN cg_users u ON m.user_id = u.id").
		Where("u.address = ?", address).
		Find(&coins).Error
	if err != nil {
		return nil, err
	}
	return coins, nil
}

func (d *Dao) AddFavoriteCoin(address string, coinId string) error {
	var user types.Users
	if err := d.DB.Where("address = ?", address).First(&user).Error; err != nil {
		return err // 如果用户不存在，这里会报 gorm.ErrRecordNotFound
	}

	mapping := types.UserCoinsMapping{
		UserId: user.Id,
		CoinId: coinId,
	}
	if err := d.DB.Create(&mapping).Error; err != nil {
		return err
	}

	return nil
}

func (d *Dao) RemoveFavoriteCoin(address string, coinId string) error {
	return d.DB.Transaction(func(tx *gorm.DB) error {
		var user types.Users
		if err := tx.Where("address = ?", address).First(&user).Error; err != nil {
			return err // 如果用户不存在，这里会报 gorm.ErrRecordNotFound
		}

		if err := tx.Where("user_id = ? AND coin_id = ?", user.Id, coinId).Delete(&types.UserCoinsMapping{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (d *Dao) ListAllCoinsAlerts() ([]types.UserCoinAlert, error) {
	var alerts []types.UserCoinAlert
	err := d.DB.Find(&alerts).Error
	if err != nil {
		return nil, err
	}
	return alerts, nil
}

func (d *Dao) ListUserCoinAlerts(address string) ([]types.UserCoinAlert, error) {
	var alerts []types.UserCoinAlert

	err := d.DB.Model(types.UserCoinAlert{}).
		Joins("JOIN cg_users u ON cg_user_coin_alerts.user_id = u.id").
		Where("u.address = ?", address).
		Find(&alerts).Error
	if err != nil {
		return nil, err
	}
	return alerts, nil
}

func (d *Dao) AddCoinPriceAlert(address string, coinId string, alertType string, targetPrice float64) error {
	// 查询 user_id
	var user types.Users
	if err := d.DB.Where("address = ?", address).First(&user).Error; err != nil {
		return err
	}

	alert := types.UserCoinAlert{
		UserID:      uint64(user.Id),
		CoinID:      coinId,
		AlertType:   alertType,
		TargetPrice: targetPrice,
		Enabled:     true,
	}

	// UPSERT 操作
	err := d.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "coin_id"}, {Name: "alert_type"}},
		UpdateAll: true,
	}).Create(&alert).Error

	if err != nil {
		return err
	}

	return nil
}

func (d *Dao) RemoveCoinPriceAlert(address string, coinId string, alertType string) error {
	return d.DB.Transaction(func(tx *gorm.DB) error {
		var user types.Users
		if err := tx.Where("address = ?", address).First(&user).Error; err != nil {
			return err // 如果用户不存在，这里会报 gorm.ErrRecordNotFound
		}

		if err := tx.Where("user_id = ? AND coin_id = ? AND alert_type = ?", user.Id, coinId, alertType).Delete(&types.UserCoinAlert{}).Error; err != nil {
			return err
		}

		return nil
	})
}
