/*
 * @Description:
 * @Author: Zjy
 * @Date: 2023-08-23 20:32:52
 * @LastEditTime: 2023-08-23 20:33:55
 * @version: 1.0
 */
package relation

import (
	"miniDouyin/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init() {
	var err error
	DB, err = gorm.Open(postgres.Open(utils.DSN),
		&gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
			Logger:                 logger.Default.LogMode(logger.Info),
		})
	if err != nil {
		panic(err)
	}
}
