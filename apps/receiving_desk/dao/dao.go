package dao

import "gorm.io/gorm"

func (receivingDeskDao *ReceivingDeskDao) BeginTransaction() *gorm.DB {
	// Begin a new transaction
	tx := receivingDeskDao.Db.Begin()
	if tx.Error != nil {
		return nil // Handle error if needed
	}
	return tx
}
