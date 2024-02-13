package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type Configtmp struct {
	Url        string
	Searchtext string
	TechBreak  string
	Error      int
	Wait       int32
	Command    string
	Log_File   string
}

func main() {

	ConfigFile := ""

	if len(os.Args) > 1 {
		ConfigFile = os.Args[1]
	} else {
		fmt.Println(`
Create a config file and add as argument

Url = "http://example.com"
SearchText = "Site Ok"
TechBreak = "Проводятся технические работы"
Error = 3
WaitReplay = 5
Command = "systemctl restart httpd"
Log_File = "/var/log/checkUrl.log"

$ checkUrl FileConfig.conf
`)
		os.Exit(2)
	}

	var conf Configtmp

	_, err := toml.DecodeFile(ConfigFile, &conf)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	LOG_FILE := conf.Log_File
	// open log file
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()

	// Set log out put and enjoy :)
	log.SetOutput(logFile)

	end := ""
	n := 0

	for true {
		
		time.Sleep(time.Duration(conf.Wait) * time.Second)
		
		if check(conf.Url, conf.TechBreak) {
			log.Println("TechBreak")
			break
		}

		if n >= conf.Error {
			log.Println("Error >=3")

			app := strings.Split(conf.Command, " ")[0]
			appArgs := strings.Split(conf.Command, " ")[1:]

			cmd := exec.Command(app, appArgs...)
			err := cmd.Run()

			if err != nil {
				log.Println("Error Cmd: ", err.Error())
			}

			log.Println(app, appArgs)

			break
		}

		if !check(conf.Url, conf.Searchtext) {
			end = "FALSE"
			n++
			// fmt.Println(check())
		} else {
			// log.Println(check())
			end = "TRUE"
			break
		}
	}
	log.Println("END Cycle Check Site:" + end)
}

func check(url string, search string) bool {
	r, err := http.Get(url)

	if err != nil {
		// panic(err)
		log.Println(err)
		return false
	}

	b, err := io.ReadAll(r.Body)

	if err != nil {
		panic(err)
	}

	return strings.Contains(string(b), search)
}
