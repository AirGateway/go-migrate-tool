// go-migrate-tool - Go migration tool for mongodb
//
// Author: Dmitry Fedorov <klka1@live.ru>

package main

import (
	"fmt"
	"github.com/docopt/docopt.go"
	"github.com/kLkA/go-migrate-tool/modules"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"time"
	"bufio"
	"strings"
)

const (
	version           = "0.1.0-beta"
	defaultTimeFormat = "060102_150405"
)



func main() {

	arguments := getCmdArguments()

	if arguments["--new"] != nil {
		err := createMigration(arguments["--new"].(string))
		if err != nil {
			modules.Log.Error(err)
		}

		return
	}
}

func getCmdArguments() map[string]interface{} {
	documentation := `Go Migration Tool.

	Usage:
		go-migrate-tool [options]

	Options:
		-h --help         Show this screen
		--new=name        Create migration
	`
	/*
		--up=limit        Apply migrations
		--down=limit      Revert migrations
	*/
	arguments, err := docopt.ParseArgs(documentation, nil, version)
	if err != nil {
		modules.Log.Error("test")
	}

	return arguments
}

func createMigration(name string) error {
	filename := "m" + time.Now().Format(defaultTimeFormat) + "_" + name + ".json"
	folder := "migrations"

	dir, _ :=  os.Getwd()
	path := dir + "/migrations/" + filename


	if !askForConfirmation(fmt.Sprintf("Create new migration \"%s\"", path)) {
		return errors.New("confirmation failed")
	}

	os.Mkdir(folder, 0755)

	if _, err := os.Stat(folder + "/" + filename); err == nil {
		return errors.New(fmt.Sprintf("file %s already exists", filename))
	}

	ioutil.WriteFile(folder+"/"+filename, []byte("[]"), 0744)
	return nil
}

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			modules.Log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}