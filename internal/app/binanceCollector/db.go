package binancecollector

import (
	"fmt"
	"log"
	"sync"

	"github.com/jinzhu/gorm"
)

// TODO: 删除此处内容
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

var sMu sync.Mutex

func save2(db *gorm.DB, trades []*trade) {
	sMu.Lock()
	defer sMu.Unlock()
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
			log.Fatal("tx.Create err:", err, t.TableName())
		}
	}
	if err := tx.Commit().Error; err != nil {
		log.Fatal("tx.Commit err:", err)
	}
}

func save2disk(trades []*trade) {
	if len(trades) == 0 {
		return
	}
	trade := trades[0]
	fileName := trade.monthDBName()
	//
	mDBmu.Lock()
	db, ok := MDB[fileName]
	mDBmu.Unlock()
	//
	if !ok {
		// initial db
		var err error
		db, err = gorm.Open("sqlite3", fileName)
		if err != nil {
			panic("failed to connect database")
		}
		//
		mDBmu.Lock()
		MDB[fileName] = db
		mDBmu.Unlock()
		//
		fmt.Printf("%s 数据库已经打开\n", fileName)
	}
	if !db.HasTable(trade.TableName()) {
		sMu.Lock()
		db = db.CreateTable(trade)
		sMu.Unlock()
	}
	save2(db, trades)
	msg := fmt.Sprintf("save %d %s into %s", len(trades), trades[0].Symbol, fileName)
	bc.Info(msg)
	log.Println(msg)
}
