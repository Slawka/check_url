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
	Error      int
	Wait       int32
	Command    string
	Log_File   string
}

func main() {

	// defConfigFile := "./checkUrl_config.conf"
	ConfigFile := ""

	if len(os.Args) > 1 {
		ConfigFile = os.Args[1]
	} else {
		fmt.Println(`
		Create a config file and add as argument

Url = "http://example.com"
SearchText = "Site Ok"
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
		if n >= conf.Error {
			log.Println("n>=3")
			// cmd := exec.Command("/usr/bin/systemctl", "restart", "httpd")
			// cmd := exec.Command(strings.Split(conf.Command, " ")[0], strings.Join(strings.Split(conf.Command, " ")[1:], " "))
			out, err := exec.Command(strings.Split(conf.Command, " ")[0], strings.Join(strings.Split(conf.Command, " ")[1:], " ")).Output()
			// cmd := exec.Command("/usr/bin/systemctl", "restart nginx")
			// out, err := exec.Command("/usr/bin/cat", "/var/log/messages").Output()
			// out, err := exec.Command("/usr/bin/systemctl", "restart nginx").Output()
			// err := cmd.Start()
			// var outb, errb bytes.Buffer
			// cmd.Stdout = &outb
			// cmd.Stderr = &errb

			if err != nil {
				log.Println("Error Cmd: ", err)
			}

			log.Println("Out Cmd: ", string(out[:]))

			log.Println(strings.Split(conf.Command, " ")[0], strings.Join(strings.Split(conf.Command, " ")[1:], " "))

			// log.Println("out:", outb.String(), "err:", errb.String())
			// log.Println("out:", cmd.Stdout, "err:", err)
			// cmd := exec.Command("/usr/bin/echo", "`/usr/bin/date +\"%Y/%M/%d %H:%m:%S\"`\" Http reboot \"", ">>","/var/log/httpd/restart-nginx.log")
			// cmd.Run()
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
	// r, err := http.Get("https://volveter.ru/error.html")
	if err != nil {
		// panic(err)
		log.Println(err)
		return false
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	// log.Println(string(b[:100]))
	return strings.Contains(string(b), search)
}
