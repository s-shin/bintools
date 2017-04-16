package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/urfave/cli"
)

func CommandErrorf(cause error, format string, formatArgs ...interface{}) error {
	msg := fmt.Sprintf("ERROR: "+format, formatArgs...)
	if cause == nil {
		return fmt.Errorf("%s", msg)
	}
	return fmt.Errorf("%s\n\n[Error Details]\n%s", msg, cause)
}

func CommandErrorInvalidLengthOfArguments(expectedLength int) error {
	return CommandErrorf(nil, "Just %d arguments should be given.", expectedLength)
}

func CommandErrorCannotOpenFile(cause error, filePath string) error {
	return CommandErrorf(cause, "Cannot open the file '%s'.", filePath)
}

func CommandErrorUnexpected(cause error) error {
	return CommandErrorf(cause, "Unexpected error occurred.")
}

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "file",
			Usage: "Read the search key from file. The first argument is interpreted as the file path.",
		},
	}

	app.ArgsUsage = "<search-key> <target-file>"

	app.Action = func(c *cli.Context) error {
		args := c.Args()
		if len(args) != 2 {
			return CommandErrorInvalidLengthOfArguments(2)
		}
		searchKey := args[0]
		targetFile := args[1]

		var keyBytes []byte
		if c.Bool("file") {
			file, err := os.Open(searchKey)
			if err != nil {
				return CommandErrorCannotOpenFile(err, searchKey)
			}
			keyBytes, err = ioutil.ReadAll(file)
			if err != nil {
				return CommandErrorUnexpected(err)
			}
		} else {
			keyBytes = []byte(searchKey)
		}

		target, err := os.Open(targetFile)
		if err != nil {
			return CommandErrorCannotOpenFile(err, targetFile)
		}

		idx := FindIndexByKey(target, keyBytes)
		fmt.Printf("[%d, %d] (%d)\n", idx, idx+len(searchKey), len(searchKey))
		fmt.Printf("dd if=%s bs=1 skip=%d count=%d 2>/dev/null; echo\n", targetFile, idx, len(searchKey))
		return nil
	}

	app.Run(os.Args)
}
