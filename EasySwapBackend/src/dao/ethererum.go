package dao

import (
	appTypes "github.com/ProjectsTask/EasySwapBackend/src/types/v1"

	"github.com/ethereum/go-ethereum/core/types"
)

func (d *Dao) AddEthereumBlock(header *types.Header) error {
	record := appTypes.EthereumBlockHeader{
		Number:     header.Number.Uint64(),
		Hash:       header.Hash().Hex(),
		ParentHash: header.ParentHash.Hex(),
		Timestamp:  header.Time,
	}

	return d.DB.Model(appTypes.EthereumBlockHeader{}).Create(&record).Error
}

func (d *Dao) GetWatchAddressList() (map[string]interface{}, error) {
	var records []appTypes.WatchAddress
	err := d.DB.Model(appTypes.WatchAddress{}).Find(&records).Error
	if err != nil {
		return nil, err
	}

	addresses := make(map[string]interface{})
	for _, record := range records {
		addresses[record.Address] = struct{}{}
	}

	return addresses, nil
}
