package main

import (
  "net"
  "time"
  "fmt"
  "io/ioutil"
  "github.com/fatih/color"
  "gopkg.in/yaml.v2"
	"strconv"
)

var useSQL = false
var lastUp = int64(0)

type Config struct {
	Net struct {
    	Host string `yaml:"host"`
    	Port string `yaml:"port"`
    	PollRate string `yaml:"pollRate"`
    	Timepout string `yaml:"timeout"`
	}
	SQL struct {
    	Host string `yaml:"host"`
    	Port string `yaml:"port"`
	}
	LogAll bool `yaml:"logAll"`
}

func main() {

	dat, err := ioutil.ReadFile("./config.yml")
    if err != nil {
        color.New(color.FgRed).Add(color.Bold).Println("[CRITICAL ERROR] Your configeration file could not be found! Check for help")
        return
    }
    fmt.Print(string(dat) + "\n")
	color.New(color.FgBlue).Println("[BOOT] Loaded configeration file")

	var config Config
    err = yaml.Unmarshal(dat, &config)
    if err != nil {
		color.New(color.FgRed).Add(color.Bold).Println("[CRITICAL ERROR] Your configeration is not valid yml! Check for help")
		return
    }
	if config.Net.Host == "" || config.Net.Port == "" || config.Net.PollRate == ""  {
		color.New(color.FgRed).Add(color.Bold).Println("[CRITICAL ERROR] Your configeration is missing required fields! Check for help")
		return
	}

	if config.SQL.Host == "" {
		color.New(color.FgBlue).Println("[BOOT] SQL loging disabled (No SQL Credentails in confogeration file)")
	}
	color.New(color.FgGreen).Add(color.Bold).Println("[BOOT] Started checking connection using " + config.Net.Host + ":" + config.Net.Port + " every " + config.Net.PollRate + " second(s)")

	startChecker(config)	
}

func startChecker(config Config) {
	timeBetweenPolls, err := strconv.ParseInt(config.Net.PollRate, 10, 64)
	if err != nil {
		panic(err)
	}
	for range time.Tick(time.Second * time.Duration(timeBetweenPolls)) {
        isOnline := checkUptime(config.Net.Host, config.Net.Port, config.Net.Timepout)
		if isOnline {
			if config.LogAll {
				color.New(color.FgGreen).Println("[" + time.Now().Format("3:4 (5) pm") + "] Target is online. ")
			}
			if lastUp != 0 {
				offlineFor := time.Now().Unix() - lastUp
				str := strconv.FormatInt(offlineFor, 10)
				color.New(color.FgGreen).Add(color.Bold).Println("[CONNECTION UPDATE] [" + time.Now().Format("3:4 (5) pm") + "] Target is online. Offline for " + str + " second(s)")
				lastUp = 0
			}
		} else {
			if config.LogAll {
				color.New(color.FgRed).Println("[" + time.Now().Format("3:4 (5) pm") + "] Target is Offline. ")
			}
			if lastUp == 0 {
				color.New(color.FgRed).Add(color.Bold).Println("[CONNECTION UPDATE] [" + time.Now().Format("3:4 (5) pm") + "] Target is Offline.")
				lastUp = time.Now().Unix()
			}
		}
    }
}

func checkUptime(host string, port string, timeOut string) bool {
	timeToWaitBeforeTimeout, err := strconv.Atoi(timeOut)
    timeout := time.Duration(time.Duration(timeToWaitBeforeTimeout) * time.Second)
    _, err = net.DialTimeout("tcp", host+":"+port, timeout)
    if err != nil {
		return false
    } else {
		return true
    }
}