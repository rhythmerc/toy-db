package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	db, err := NewHashDB("db.txt")
	if err != nil {
		log.Fatal(err)
	}

	var knf *KeyNotFoundError

	in := bufio.NewScanner(os.Stdin)

scan:
	for in.Scan() {
		line := in.Text()
		tokens := strings.Split(line, " ")

		switch tokens[0] {
		case "get":
			{
				if len(tokens) < 2 {
					break
				}

				key := tokens[1]
				value, err := db.Read(key)
				if err != nil {
					if errors.As(err, &knf) {
						break
					}
					log.Fatal(err)
				}

				fmt.Println(value)
			}
		case "set":
			{
				if len(tokens) < 3 {
					break
				}

				key := tokens[1]
				value := tokens[2]
				for i := 3; i < len(tokens); i++ {
					value += " " + tokens[i]
				}

				err := db.Write(key, value)
				if err != nil {
					log.Fatal(err)
				}
			}
		case "exit":
			break scan
		}
	}

	if err := in.Err(); err != nil {
		log.Fatal(err)
	}
}
