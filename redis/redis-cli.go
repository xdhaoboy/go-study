package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("config file error: %s\n", err)
		os.Exit(1)
	}

	redisURL := viper.Get("redis.url")
	redisPasswd := viper.Get("redis.password")
	redisDb := viper.Get("redis.db")

	fmt.Printf("url:%s,passwd:%s,db:%d\n", redisURL, redisPasswd, redisDb)

	products, err := os.Open("product.txt")
	if err != nil {
		panic(err)
	}

	productBuf := bufio.NewReader(products)

	for {
		product, _, err := productBuf.ReadLine()
		if err != nil {
			if err == io.EOF { //读取结束，会报EOF
				break
			}
			panic(err)
		}
		productT := string(product)

		//fmt.Printf("%s begin to be set\n", productT)
		lmids, err := os.Open("lmid.txt")
		lmidBuf := bufio.NewReader(lmids)
		for {
			line, _, err := lmidBuf.ReadLine()
			if err != nil {
				if err == io.EOF { //读取结束，会报EOF
					break
				}
				panic(err)
			}
			lineT := string(line)
			t := strings.Split(lineT, "\t")
			lmid := t[0]
			version := t[1]

			client := redis.NewClient(&redis.Options{
				Addr:     redisURL.(string),    // Redis地址
				Password: redisPasswd.(string), // Redis账号
				DB:       redisDb.(int),        // Redis库
			})
			//fmt.Printf("hset %s %s {\"lmVersion\":\"%s\"}\n", product, lmid, version)
			err = client.HSet(productT, lmid, `{"lmVersion":"`+version+`"}`).Err()
			if err != nil {
				panic(err)
			}

			value, err := client.HGet(productT, lmid).Result()
			if err == redis.Nil {
				fmt.Println("key does not exist")
			} else if err != nil {
				panic(err)
			} else {
				fmt.Println(productT, lmid, value)
			}

			defer client.Close()
		}
		lmids.Close()
	}
	products.Close()

}
