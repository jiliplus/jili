package main

func main() {
	srcName := "./data/btcusdt.sqlite3"
	//
	db := openToMemory(srcName)
	defer db.Close()
	//
	tickSrc(db, nil)
}
