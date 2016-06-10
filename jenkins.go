package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "log"
    "encoding/json"
    "os"
    "bytes"
)

type Job struct {
    Name string
    Url string
    Color string
}

type JenkinsHome struct {
    Jobs []Job
    Mode string
}

func Ls() {
    
    url := fmt.Sprintf("%s/api/json", os.Getenv("JENKINS_CMD_URL"))

    resp, err := http.Get(url)
    if err != nil {
		log.Fatal(err)
    }
    
    //read http response
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

    // Parse json
	var x JenkinsHome
	e := json.Unmarshal(body, &x)
	if e != nil {
	    log.Fatal(e)
	}

    for jobIdx := range x.Jobs {
        fmt.Println(x.Jobs[jobIdx].Name)
    }
}

func Build(jobName string) {
    
    url := fmt.Sprintf("%s/job/%s/build", os.Getenv("JENKINS_CMD_URL"), jobName)
    //fmt.Println(url)
    var body []byte
    resp, err := http.Post(url, "", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

    
}

type Configuration struct {
    Url string
}

func GetConfig() Configuration {
    file, openErr := os.Open( fmt.Sprintf("%s/.jenkinsconfig", os.Getenv("HOME")) )
    if openErr != nil {
        log.Fatal(openErr)
    }
    decoder := json.NewDecoder(file)
    configuration := Configuration{}
    err := decoder.Decode(&configuration)
    if err != nil {
      fmt.Println("error:", err)
    }
    
    return configuration
}

func SetJenkinsUrl(url string) {
    
    config := GetConfig()
    config.Url = url
    
    b, err := json.MarshalIndent(config, "", "  ")
    if err != nil {
       log.Fatal(err)

    }
    
    ioutil.WriteFile( fmt.Sprintf("%s/.jenkinsconfig", os.Getenv("HOME")), b,  0644)
    os.Setenv("JENKINS_CMD_URL", url)
}

func main() {
    
    // Load the configuration in
    config := GetConfig()
    os.Setenv("JENKINS_CMD_URL", config.Url)

    
    args := os.Args[1:]
    
    switch args[0] {
        case "ls":
            Ls()
        case "build":
            if len(args) < 2 {
                fmt.Println("please specify a job to build")
            } else {
                Build(args[1])
            }
        case "set-url":
            if len(args) < 2 {
                fmt.Println("please specify a url to set")
            } else {
                SetJenkinsUrl(args[1])
            }
        case "get-url":
            fmt.Println(os.Getenv("JENKINS_CMD_URL"))
        case "env":
            fmt.Println("export ")
        default:
            fmt.Println("unknown command: ", args[0])
    }
    
}
