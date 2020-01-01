package binancedata

// source is new try.
func source(channel chan<- *trade, symbol string, id int64) {
	rows, err := db.Table(symbol).Where("id > ?", id).Rows() // (*sql.Rows, error)
	if err != nil {
		panic(symbol + " db.Table.Rows err: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var t trade
		db.ScanRows(rows, &t)
		t.Symbol = symbol
		channel <- &t
	}
	close(channel)
}
