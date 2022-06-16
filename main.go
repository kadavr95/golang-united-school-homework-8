package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type Arguments map[string]string

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

var (
	idFlag        = flag.String("id", "", "")
	operationFlag = flag.String("operation", "", "")
	itemFlag      = flag.String("item", "", "")
	filenameFlag  = flag.String("filename", "", "")
)

func Perform(args Arguments, writer io.Writer) error {
	switch args["operation"] {
	case "":
		return fmt.Errorf("-operation flag has to be specified")
	case
		"add",
		"list",
		"findById",
		"remove":
		if args["fileName"] == "" {
			return fmt.Errorf("-fileName flag has to be specified")
		}
		file, err := os.OpenFile(args["fileName"], os.O_RDWR|os.O_CREATE, 644)
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {

			}
		}(file)
		if err != nil {
			return fmt.Errorf("file can't be opened")
		}
		fileContents, err := ioutil.ReadAll(file)
		if err != nil {
			return fmt.Errorf("file can't be read")
		}
		switch args["operation"] {
		case "add":
			if args["item"] == "" {
				return fmt.Errorf("-item flag has to be specified")
			}
			var users []User
			if len(fileContents) > 0 {
				err := json.Unmarshal(fileContents, &users)
				if err != nil {
					return fmt.Errorf("json parse error")
				}
			}
			var newUser User
			err = json.Unmarshal([]byte(args["item"]), &newUser)
			if err != nil {
				return fmt.Errorf("new item json parse error")
			}
			notDuplicateFlag := true
			for _, val := range users {
				if val.Id == newUser.Id {
					_, err = writer.Write([]byte(fmt.Sprintf("Item with id %s already exists", newUser.Id)))
					notDuplicateFlag = false
					if err != nil {
						return fmt.Errorf("error to output Existing id message")
					}
				}
			}
			if notDuplicateFlag {
				users = append(users, newUser)
				fileData, err := json.Marshal(users)
				if err != nil {
					return fmt.Errorf("error converting structure to json")
				}
				err = os.WriteFile(args["fileName"], fileData, 644)
				if err != nil {
					return fmt.Errorf("error writing to file")
				}
			}
		case "list":
			_, err := writer.Write(fileContents)
			if err != nil {
				return fmt.Errorf("file contents can't be returned")
			}
		case "findById":
			if args["id"] == "" {
				return fmt.Errorf("-id flag has to be specified")
			}
			var users []User
			if len(fileContents) > 0 {
				err := json.Unmarshal(fileContents, &users)
				if err != nil {
					return fmt.Errorf("json parse error")
				}
			}
			for _, val := range users {
				if val.Id == args["id"] {
					result, err := json.Marshal(val)
					if err != nil {
						return fmt.Errorf("error converting result to json")
					}
					_, err = writer.Write(result)
					if err != nil {
						return fmt.Errorf("error to output Existing id message")
					}
				}
			}
		case "remove":
			if args["id"] == "" {
				return fmt.Errorf("-id flag has to be specified")
			}
			var users []User
			if len(fileContents) > 0 {
				err := json.Unmarshal(fileContents, &users)
				if err != nil {
					return fmt.Errorf("json parse error")
				}
			}
			notFoundFlag := true
			for i, val := range users {
				if val.Id == args["id"] {
					notFoundFlag = false
					users[i] = users[len(users)-1]
					users = users[:len(users)-1]
					fileData, err := json.Marshal(users)
					if err != nil {
						return fmt.Errorf("error converting structure to json")
					}
					err = os.WriteFile(args["fileName"], fileData, 644)
					if err != nil {
						return fmt.Errorf("error writing to file")
					}
				}
			}
			if notFoundFlag {
				_, err = writer.Write([]byte(fmt.Sprintf("Item with id %s not found", args["id"])))
				if err != nil {
					return fmt.Errorf("error to output Existing id message")
				}
			}
		}
		return nil
	default:
		return fmt.Errorf("Operation %s not allowed!", args["operation"])
	}
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}

func parseArgs() Arguments {
	flag.Parse()
	return Arguments{
		"id":        *idFlag,
		"operation": *operationFlag,
		"item":      *itemFlag,
		"fileName":  *filenameFlag,
	}
}
