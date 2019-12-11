package binancecollector

import (
	"log"
)

// func save(trades []*trade) {
// 	valueStrings := []string{}
// 	valueArgs := []interface{}{}
// //
// 	for _, t := range trades {
// 		valueStrings = append(valueStrings, "(?, ?, ?,?,?,?)")
// 		valueArgs = append(valueArgs, t.ID)
// 		valueArgs = append(valueArgs, t.Price)
// 		valueArgs = append(valueArgs, t.Quantity)
// 		valueArgs = append(valueArgs, t.UTC)
// 		valueArgs = append(valueArgs, t.IsBuyerMaker)
// 		valueArgs = append(valueArgs, t.IsBestMatch)
// 	}
// //
// 	smt := fmt.Sprintf("INSERT INTO %s (id, price, quantity, utc, is_buyer_maker, is_best_match) VALUES ", trades[0].TableName())
// 	smt += "%s "
// //
// 	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))
// //
// 	tx := db.Begin()
// //
// 	if err := tx.Exec(smt, valueArgs...).Error; err != nil {
// 		tx.Rollback()
// 		log.Fatal("tx.Exec err:", err)
// 	}
// //
// 	if err := tx.Commit().Error; err != nil {
// 		log.Fatal("tx.Commit err:", err)
// 	}
// }

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
	log.Printf("Save %d data\n", len(trades))
}
