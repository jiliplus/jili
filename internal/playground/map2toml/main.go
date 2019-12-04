package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/pelletier/go-toml"
)

func main() {
	data := make(map[string]int, 2)
	data["one"] = 1
	data["two"] = 2

	bytes, err := toml.Marshal(data)
	if err != nil {
		log.Fatal("toml.Marshal err:", err)
	}

	fmt.Println("bytes    ", bytes)

	fileName := "TradeID.toml"

	err = ioutil.WriteFile(fileName, bytes, 0666)
	if err != nil {
		log.Fatal("ioutil.WriteFile err:", err)
	}

	tree, err := toml.LoadFile(fileName)
	if err != nil {
		log.Fatal("toml.LoadFile err:", err)
	}

	dataBack := make(map[string]int, 2)

	err = tree.Unmarshal(&dataBack)
	if err != nil {
		log.Fatal("tree.Unmarshal err:", err)
	}

	fmt.Println("dataBack[\"one\"]=", dataBack["one"])
	fmt.Println("dataBack[\"two\"]=", dataBack["two"])

	os.Remove(fileName)
}
