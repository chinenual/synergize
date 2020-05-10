package main

import (
	"encoding/json"
	"os"
	"fmt"
	"log"
	"path"
	"strings"
	
	"github.com/chinenual/synergize/data"
)

func jsonize (filename string) {
	ext := strings.ToLower(path.Ext(filename))
	switch ext {
	case ".vce":
		vce,err := data.ReadVceFile(filename)
		if err != nil {
			log.Panic(err);
		}
		
		b,_ := json.MarshalIndent(vce, "", " ")
		result := string(b)
		
		fmt.Println(result)
	case ".crt":
		crt,err := data.ReadCrtFile(filename)
		if err != nil {
			log.Panic(err);
		}
		
		b,_ := json.MarshalIndent(crt, "", " ")
		result := string(b)
		
		fmt.Println(result)
	default:
		log.Panicf("ERROR: don't know what to do with %s\n", filename)
	}
}

func main() {
	for _, a := range os.Args[1:] {
		jsonize(a)
	}
}
