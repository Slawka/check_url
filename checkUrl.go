package main

import (
	"os/exec"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	"os"
)

func main() {
    LOG_FILE := "/var/log/httpd/checkHttpd.log"
    // open log file
    logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
    if err != nil {
        log.Panic(err)
    }
    defer logFile.Close()

    // Set log out put and enjoy :)
    log.SetOutput(logFile)

	end := ""
	n :=0
	for (true){
		time.Sleep(10)
		if(n>=3){
			log.Println("n>=3")
			cmd := exec.Command("/usr/bin/systemctl", "restart", "httpd")
			cmd.Run()
			// cmd := exec.Command("/usr/bin/echo", "`/usr/bin/date +\"%Y/%M/%d %H:%m:%S\"`\" Http reboot \"", ">>","/var/log/httpd/restart-nginx.log")
			// cmd.Run()
			break
		}
		if(!check()){
			end = "FALSE"
			n++
			// fmt.Println(check())
		}else{
			// log.Println(check())
			end = "TRUE"
			break
		}
	}
	log.Println("END Cycle Check Site:" + end)
}

func check() bool {
	r, err := http.Get("https://volveter.ru")
	// r, err := http.Get("https://volveter.ru/error.html")
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	// log.Println(string(b[:100]))
	return strings.Contains(string(b), "Группа компаний «Вольный Ветер»")
}
