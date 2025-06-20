package main
// Main module for JobWizard backend server
// Created by Sally Goldin, 18 June 2025

import (
    "fmt"
    "flag"
    //"log"
    "os"
    "path/filepath"    
    "strings"
    "strconv"
    "github.com/segoldin/JobWizard/job_wizard/dbaccess"
    //"github.com/segoldin/JobWizard/job_wizard/api"
    //"github.com/segoldin/JobWizard/job_wizard/api/middlewares"        
    "github.com/joho/godotenv"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"    
)

var (
    task           string
 //   user_email     string
 //   first_name     string
 //   last_name      string
 //   phone          string
 //   education      int
)

func main() {
    flag.StringVar(&task, "task", "", "Task to perform")
    flag.Parse()
    if task != "" {
        commandLineFunction()
    } else {
        setupAPI()
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
    fmt.Println("Called commandLineFunction()")
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
