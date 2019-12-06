package binancecollector

import "log"

func save(trades []*trade) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		log.Fatal("save tx err:", err)
	}

	for _, t := range trades {
		if err := tx.Create(t).Error; err != nil {
			log.Fatal("tx.Create err:", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Fatal("tx.Commit err:", err)
	}
}
