package binancecollector

import "fmt"

// Split2 is new try.
func Split2() {
	rows, err := db.Table("ETHBTC").Rows() // (*sql.Rows, error)
	if err != nil {
		panic("db.Table.Rows err: " + err.Error())
	}
	defer rows.Close()

	if rows.Next() {
		var t trade
		db.ScanRows(rows, &t)

		t.Symbol = "ETHBTC"

		fmt.Println(t)
	}
}
