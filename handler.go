package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "log"

	"github.com/Jeffail/gabs/v2"
	_ "github.com/mattn/go-sqlite3"
	"launchpad.net/goamz/aws"
	_ "launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
	_ "launchpad.net/goamz/s3"
)

var dbname string = "classicmodels"
var tablename string = "productlines"
var endpoint string = ""
var pkey string = ""

func main() {
	db, err := sql.Open("sqlite3", "./lb_29.sqlite")

	// query
	rows, err := db.Query("SELECT clustername FROM lb_clusters")
	checkErr(err)
	var clustername string

	for rows.Next() {
		err = rows.Scan(&clustername)
		checkErr(err)
		fmt.Println(clustername)
	}

	rows.Close()
	db.Close()
	auth := aws.Auth{
		AccessKey: "",
		SecretKey: "",
	}
	useast := aws.USEast
	connection := s3.New(auth, useast)
	mybucket := connection.Bucket("offloadmysqldata")
	// res, err := mybucket.List("", "", "", 1000)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	res2, err := mybucket.Get("vivek.bhardwaj@codenation.co.in_replication_server_test1.json")
	jsonParsed, err := gabs.ParseJSON(res2)
	temp_endpoint, ok := jsonParsed.Path("ELASTICSEARCH.HOST").Data().(string)
	if ok {
		endpoint = temp_endpoint
	}
	// search_string := fmt.Sprintf(`"MAPPING.%v.%v.primary_key.1"`, dbname, tablename)
	// for _, child := range jsonParsed.S(search_string).Children() {
	// 	fmt.Println(child.Data().(string))
	// }
	// fmt.Println(search_string)
	// fmt.Println(jsonParsed.Path(search_string).String())
	temp_pkey := jsonParsed.Search("MAPPING", dbname, tablename, "primary_key", "0").String()
	// fmt.Println(temp_pkey)
	// if ok {
	pkey = temp_pkey
	// }
	fmt.Println(endpoint)
	fmt.Println(pkey)
	// fmt.Println(jsonParsed)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	res2, _ = prettyprint(res2)
	// fmt.Printf("%s", res2)
	// s := string(res2[:])
	// fmt.Println(s)
	// for _, v := range res2.Contents {
	// 	fmt.Println(v.Key)
	// }

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func prettyprint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}
