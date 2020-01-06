package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/mattn/go-sqlite3"
)

func main() {
	srcName := "./foo.sqlite3"
	os.Remove(srcName)
	makeFoo(srcName)

	db := openToMemory(srcName)
	defer db.Close()

	showFooTable(db)
}

func openToMemory(srcName string) *sql.DB {
	sqlite3conn := make([]*sqlite3.SQLiteConn, 0, 2)
	// fmt.Println(cap(sqlite3conn))
	sql.Register("sqlite3_with_hook_example",
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				sqlite3conn = append(sqlite3conn, conn)
				return nil
			},
		})

	srcDb, err := sql.Open("sqlite3_with_hook_example", srcName)
	if err != nil {
		log.Fatal(err)
	}
	defer srcDb.Close()
	srcDb.Ping()

	destDb, err := sql.Open("sqlite3_with_hook_example", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	// do NOT close destDB
	destDb.Ping()

	src, dest := sqlite3conn[0], sqlite3conn[1]

	copyDB(dest, src)

	return destDb
}

func copyDB(dst, src *sqlite3.SQLiteConn) {
	backup, err := dst.Backup("main", src, "main")
	if err != nil {
		return
	}
	defer backup.Finish()
	backup.Step(-1)
}

func showFooTable(db *sql.DB) {
	rows, err := db.Query("select id, name from foo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%02d %s\n", id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func makeFoo(srcName string) {
	db, err := sql.Open("sqlite3", srcName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table foo (id integer not null primary key, name text);
	delete from foo;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for i := 0; i < 10; i++ {
		_, err = stmt.Exec(i, fmt.Sprintf("hello world %03d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()
}
