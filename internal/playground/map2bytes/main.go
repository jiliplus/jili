package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	data := make(map[string]int, 2)
	data["one"] = 1
	data["two"] = 2

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

	mapBack := make(map[string]int, 2)

	err = dec.Decode(&mapBack)
	if err != nil {
		log.Fatal("dec.Decode err:", err)
	}

	fmt.Println("mapBack[\"one\"]=", mapBack["one"])
	fmt.Println("mapBack[\"two\"]=", mapBack["two"])

	os.Remove(fileName)
}
