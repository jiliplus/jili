package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type pair struct {
	Name string
	ID   int
}

func main() {
	data := make([]pair, 2, 10)
	data[0] = pair{
		Name: "zero",
		ID:   0,
	}
	data[1] = pair{
		Name: "one",
		ID:   1,
	}

	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)

	err := enc.Encode(data)
	if err != nil {
		log.Fatal("encode err:", err)
	}

	fmt.Println("buf    ", buf)

	fileName := "test.tradeID"

	err = ioutil.WriteFile(fileName, buf.Bytes(), 0666)

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("os.Open err:", err)
	}
	defer file.Close()

	dataBack, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("ioutil read all err:", err)
	}

	fmt.Println("dataBack", dataBack)

	bufBack := bytes.NewBuffer(dataBack)
	dec := gob.NewDecoder(bufBack)

	// sliceBack := make([]pair, 0, 0)
	var sliceBack []pair
	fmt.Printf("%d,%d\n", len(sliceBack), cap(sliceBack))

	err = dec.Decode(&sliceBack)
	if err != nil {
		log.Fatal("dec.Decode err:", err)
	}

	fmt.Println("sliceBack[0]=", sliceBack[0])
	fmt.Println("sliceBack[1]=", sliceBack[1])
	fmt.Printf("%d,%d\n", len(sliceBack), cap(sliceBack))

	os.Remove(fileName)
}
