package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

const fileDefaultPermission = 0644

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

	validOperations = map[string]bool{
		"add":      true,
		"list":     true,
		"findById": true,
		"remove":   true,
	}
)

func Perform(args Arguments, writer io.Writer) error {
	if args["operation"] == "" {
		return fmt.Errorf("-operation flag has to be specified")
	}
	if !validOperations[args["operation"]] {
		return fmt.Errorf("Operation %s not allowed!", args["operation"])
	}

	if args["fileName"] == "" {
		return fmt.Errorf("-fileName flag has to be specified")
	}

	file, err := os.OpenFile(args["fileName"], os.O_RDWR|os.O_CREATE, fileDefaultPermission)
	defer file.Close()

	if err != nil {
		return fmt.Errorf("file can't be opened")
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("file can't be read")
	}

	switch args["operation"] {
	case "add":
		err = add(fileContents, writer, args)
		if err != nil {
			return err
		}
	case "list":
		err = list(fileContents, writer)
		if err != nil {
			return err
		}
	case "findById":
		err = findById(fileContents, writer, args)
		if err != nil {
			return err
		}
	case "remove":
		err = remove(fileContents, writer, args)
		if err != nil {
			return err
		}
	}

	return nil
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

func add(fileContents []byte, writer io.Writer, args Arguments) error {
	if args["item"] == "" {
		return fmt.Errorf("-item flag has to be specified")
	}

	users, err := parseUsers(fileContents)
	if err != nil {
		return err
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
			if err != nil {
				return fmt.Errorf("error to output Existing id message")
			}

			notDuplicateFlag = false
		}
	}

	if notDuplicateFlag {
		users = append(users, newUser)

		fileData, err := json.Marshal(users)
		if err != nil {
			return fmt.Errorf("error converting structure to json")
		}

		err = os.WriteFile(args["fileName"], fileData, fileDefaultPermission)
		if err != nil {
			return fmt.Errorf("error writing to file")
		}
	}

	return nil
}

func findById(fileContents []byte, writer io.Writer, args Arguments) error {
	if args["id"] == "" {
		return fmt.Errorf("-id flag has to be specified")
	}

	users, err := parseUsers(fileContents)
	if err != nil {
		return err
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

	return nil
}

func list(fileContents []byte, writer io.Writer) error {
	_, err := writer.Write(fileContents)
	if err != nil {
		return fmt.Errorf("file contents can't be returned")
	}

	return nil
}

func remove(fileContents []byte, writer io.Writer, args Arguments) error {
	if args["id"] == "" {
		return fmt.Errorf("-id flag has to be specified")
	}

	users, err := parseUsers(fileContents)
	if err != nil {
		return err
	}

	notFoundFlag := true

	for i, val := range users {
		if val.Id == args["id"] {
			users[i] = users[len(users)-1]
			users = users[:len(users)-1]

			fileData, err := json.Marshal(users)
			if err != nil {
				return fmt.Errorf("error converting structure to json")
			}

			err = os.WriteFile(args["fileName"], fileData, fileDefaultPermission)
			if err != nil {
				return fmt.Errorf("error writing to file")
			}

			notFoundFlag = false
		}
	}

	if notFoundFlag {
		_, err := writer.Write([]byte(fmt.Sprintf("Item with id %s not found", args["id"])))
		if err != nil {
			return fmt.Errorf("error to output Existing id message")
		}
	}
	return nil
}

func parseUsers(fileContents []byte) ([]User, error) {
	var users []User

	if len(fileContents) > 0 {
		err := json.Unmarshal(fileContents, &users)
		if err != nil {
			return nil, fmt.Errorf("json parse error")
		}
	}

	return users, nil
}
