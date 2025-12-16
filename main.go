package main

import (
	"bufio"
	"errors"
	"fmt"
	godb "go_db/db"
	"log"
	"os"
	"strings"
)

const MAX_TOKENS = 3

func main() {
	db, err := godb.NewHashDB("db.txt")
	if err != nil {
		log.Fatal(err)
	}

	in := bufio.NewScanner(os.Stdin)

scan:
	for in.Scan() {
		line := in.Text()
		tokens := strings.SplitN(line, " ", MAX_TOKENS)

		switch tokens[0] {
		case "get":
			if len(tokens) < 2 {
				break
			}

			key := tokens[1]
			value, err := db.Read(key)
			if err != nil {
				if errors.As(err, &godb.ERR_KNF) {
					break
				}
				log.Fatal(err)
			}

			fmt.Println(value)
		case "set":
			if len(tokens) < 3 {
				break
			}

			key := tokens[1]
			value := tokens[2]

			err := db.Write(key, value)
			if err != nil {
				log.Fatal(err)
			}
		case "exit":
			break scan
		}
	}

	if err := in.Err(); err != nil {
		log.Fatal(err)
	}
}
