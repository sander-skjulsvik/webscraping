package utils

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

type printProgress struct {
	N, Common int
	Threshold float64
	Msg       string
}

func (p printProgress) printP(i int) {
	if float64(i) > p.Threshold {
		progress := float64(p.N) / float64(p.Common)
		p.Threshold += progress
		fmt.Printf("%s %f %% done \n", p.Msg, progress)
	}
}

func PrintStingArr(arr []*string) {
	for i, v := range arr {
		fmt.Printf("%d: %s\n", i, *v)
	}
}

func Ascii2Int(s string) (int, error) {
	intString := ""
	isDigit, _ := regexp.Compile(`\d`)
	c := ""
	for _, char := range s {
		c = string(char)
		if isDigit.MatchString(c) {
			intString += c
		}
	}
	r, err := strconv.Atoi(intString)
	if err != nil {
		log.Println("Str:", s, "Err:", err.Error())
	}
	return r, err
}

func readJsonArray(filePath string) []interface{} {
	fp, openErr := os.Open(filePath)
	if openErr != nil {
		fmt.Println(openErr)
		os.Exit(0)
	}
	byteVal, readErr := ioutil.ReadAll(fp)
	if readErr != nil {
		fmt.Println(readErr)
		os.Exit(0)
	}
	var r []interface{}
	json.Unmarshal([]byte(byteVal), &r)
	fp.Close()
	return r
}

func LogIfErr(e error, msg string) bool {
	if e != nil {
		log.Fatal(msg+"\n", e)
	}
	return e != nil
}

func LogIfFatal(e error, msg string) bool {
	if e != nil {
		log.Fatalf("%s %s", e)
	}
	return e != nil
}

func CleanKeysForMongoDb(key string) string {
	s := strings.Split(key, "")
	// Remove "." in the end of keys.

	if s[len(s)-1] == "." {
		s = s[:len(s)-1]
	}

	return strings.Join(s, "")
}

func write2csv(fileName string, data [][]string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Failed to create file: %s", err)
	}

	w := csv.NewWriter(file)
	w.WriteAll(data)

	if err := w.Error(); err != nil {
		log.Fatal("Error writing ")
	}

}

func printInterfaceArray(arr []interface{}) {
	for i, val := range arr {
		fmt.Printf("%d: %v", i, val)
	}
}

func printCursorValues(cur *mongo.Cursor) {
	for cur.Next(context.TODO()) {
		fmt.Println(cur.Current)
	}
}
