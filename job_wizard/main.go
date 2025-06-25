package main
// Main module for JobWizard backend server
// Created by Sally Goldin, 18 June 2025

import (
    "fmt"
    "flag"
    //"log"
    "os"
    "path/filepath" 
    "github.com/segoldin/JobWizard/job_wizard/dbaccess"
    //"github.com/segoldin/JobWizard/job_wizard/api"
    "github.com/segoldin/JobWizard/job_wizard/api/middlewares"        
    "github.com/joho/godotenv"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"    
)


var (
    server         bool
    task           string
    user           user_info
    job            job_info
)

func main() {
    flag.BoolVar(&server, "server", false, "Specify as true to expose REST API")    
    flag.StringVar(&task, "task", "", "Task to perform")
    // see validate.go for a list of defined tasks
    // arguments for register
    flag.StringVar(&user.email,"email","","Email of user registering")
    flag.StringVar(&user.first,"first","","First name of user registering")
    flag.StringVar(&user.last,"last","","Last name of user registering")
    flag.StringVar(&user.phone,"phone","","10 digit phone number of user registering")   
    flag.IntVar(&user.education,"education",0,"Education of user registering - 0 to 4 (doctoral)") 
    // arguments for create (job) and modify job
    flag.StringVar(&job.creator,"creator","","Email of user creating the job")
    flag.StringVar(&job.title,"title","","Job title, in quotes - 64 chars max")
    flag.StringVar(&job.description,"description","","Job description, in quotes - 1024 chars max")
    flag.IntVar(&job.min_education,"min_education",0,"Minimum education level required - integer from 0 to 4")   
    flag.IntVar(&job.min_experience,"min_experience",0,"Minimum years of experience desired - integer")
    flag.IntVar(&job.salary,"salary",0,"Monthly salary offered - integer, max 1 million")

    flag.Parse()
    if server {
        setupAPI()
    } else {
        commandLineFunction()
    }
}

func setupAPI() {
    err := recordPid() // create a pid file so we can later kill the process
    if err != nil {
        fmt.Printf("Error writing PID file: %v\n", err)
        os.Exit(1)
    }
    godotenv.Load(".env_jobwizard")
    e := echo.New()
    //e.Use(middleware.Logger())
    e.Use(middleware.Recover())    

    e.Use(middleware.CORS())
    middlewares.InitCorsMiddleware(e)
/***
    _privateAPI := e.Group("/api")
    _ = _privateAPI
    api.ApplicationPrivateRoute(_privateAPI)
    e.Logger.Fatal(e.Start(":" + os.Getenv("JOBWIZARD_API_PORT")))
    **/
}

func commandLineFunction() {
    dbOk := dbaccess.CheckConnection()

    if !dbOk {
        fmt.Println("Connection to DB failed")
        os.Exit(1)
    }
    valid, msg := validateTaskArgs(task,user,job)
    if !valid {
        jsonErrorOutput(msg)
        os.Exit(1)
    }
    task_index := findTask(task)  // we have already validated the task above
    jsonResponse := dispatch(task_index)
    fmt.Println(jsonResponse)
}

// Figure out what db service/function to call to handle the task
// We assume that dispatch() knows which structure holds the appropriate arguments
// for the relevant task
func dispatch(task_index int) (jsonResponse string) {
    var err error
    switch(task_index) {
        case 0:
            err = dbaccess.RegisterUser(user.email,user.first,user.last,user.phone,user.education)
            if err != nil {
                jsonResponse = fmt.Sprintf("{ \"error\" : \"%v\" }\n",err)
            } else {
                jsonResponse = fmt.Sprintf("{ \"success\" : \"Registered user %s\"}\n",user.email)  
            }
            break
        case 1:
            fmt.Println("Not yet implemented")
    }
    return jsonResponse
}
// Output an error message to the terminal
// in JSON format
func jsonErrorOutput(msg string ) {
    fmt.Printf("{ \"error\" : \"%s\" }\n", msg)
}

// Get the current process ID and write to a file in
// the directory holding the executable,
// so it can be used later to kill the process
// if we have a new deployment
func recordPid() error {
    pid := os.Getpid()
    msg := fmt.Sprintf("Process Id is %d", pid)
    text := fmt.Sprintf("kill -9 %d\n", pid)
    fmt.Println(msg)
    ex, err := os.Executable()
    if err != nil {
        return err
    }
    execDir := filepath.Dir(ex)
    pidfile := fmt.Sprintf("%s/kill_jobwizard.sh", execDir)
    err = os.WriteFile(pidfile, []byte(text), os.ModePerm)
    if err != nil {
        return err
    }
    return nil
}
