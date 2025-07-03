package main
// Main module for JobWizard backend server
// Created by Sally Goldin, 18 June 2025

import (
    "encoding/json"
    "flag"
    "fmt"
    //"log"
    "os"
    "path/filepath" 
    "github.com/segoldin/JobWizard/job_wizard/data"    
    "github.com/segoldin/JobWizard/job_wizard/dbaccess"
    //"github.com/segoldin/JobWizard/job_wizard/api"
    "github.com/segoldin/JobWizard/job_wizard/api/middlewares"        
    "github.com/joho/godotenv"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"    
)

// Job search criteria

var (
    server         bool
    task           string
    help           bool
    taskhelp       bool
    user           data.User_info
    job            data.Job_info
    filter         data.Search_criteria
    submission     data.Submission
)


func main() {
    flag.BoolVar(&server, "server", false, "Specify as true to expose REST API")
    flag.BoolVar(&help, "help", false, "Specify as true to see general help")
    flag.BoolVar(&help, "taskhelp", false, "Specify as true to help for one task")               
    flag.StringVar(&task, "task", "", "Task to perform")
    // see validate.go for a list of defined tasks
    // arguments for register
    flag.StringVar(&user.Email,"email","","Email of user")
    flag.StringVar(&user.First,"first","","First name of user registering")
    flag.StringVar(&user.Last,"last","","Last name of user registering")
    flag.StringVar(&user.Phone,"phone","","10 digit phone number of user registering")   
    flag.IntVar(&user.Education,"education",0,"Education of user registering - 0 to 4 (doctoral)") 
    // arguments for create (job) and modify job
    flag.StringVar(&job.Creator,"creator","","Email of user creating the job")
    flag.StringVar(&job.Title,"title","","Job title, in quotes - 64 chars max")
    flag.StringVar(&job.Description,"description","","Job description, in quotes - 1024 chars max")
    flag.IntVar(&job.Min_education,"min_education",0,"Minimum education level required - integer from 0 to 4")   
    flag.IntVar(&job.Min_experience,"min_experience",0,"Minimum years of experience desired - integer")
    flag.IntVar(&job.Salary,"salary",0,"Monthly salary offered - integer, max 1 million")
    flag.BoolVar(&job.Is_open,"is_open",true,"Is the job still open?")    
    // arguments for search jobs
    //   uses "email" ==> user.Email
    flag.StringVar(&filter.Posted,"posted","","Posted date in format YYYY-MM-DD") 
    //   uses "min_education" ==> job.Min_education
    //   uses "salary"==> job.Salary
    flag.StringVar(&filter.Keyword,"keyword","","Keyword for title search")  
    // arguments for detail task
    flag.StringVar(&job.Job_id,"job_id","","Id of job to be displayed")
    flag.Parse()
    if help {
        usage("")
    } 
    if server {
        setupAPI()
    } else {
        commandLineFunction()
    }
}

// display information about the arguments for each task
// or for a single task as specified by task_name
func usage(task_name string) {
    fmt.Println("General usage: ./job_wizard -task <taskname> [arguments...]")
    fmt.Println("\tWrites results to standard output in JSON format\n")
    if task_name == "" {
        fmt.Println("Available tasks: ")
        fmt.Println("\tregister - create a new user in the database")
        fmt.Println("\tcreate - create a new job posting")
        fmt.Println("\tsearch - general search for jobs")
        fmt.Println("\tdetail - see detailed information about a selected job")
        fmt.Println("\toffered - search for jobs created by me")
        fmt.Println("\tapplied - search for jobs that I have applied for")
        fmt.Println("\tmodify - modify a job created by me")
        fmt.Println("\tsubmit - submit an application for a job\n")
        fmt.Println("For task arguments, type ./job_wizard -taskhelp <task_name>\n")
        return                       
    } 
    task_index := findTask(task_name)
    if (task_index < 0) {
        fmt.Printf("Unknown task '%s'\n",task_name)
        return
    }
    switch(task_index) {
        case 0: // register
          
        case 1: // create
         
        case 2: // search
  
        case 3: // detail
 
        case 4: // offered
 
        case 5: // applied

        case 6: // modify job
 
        case 7: // submit application for job        
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
    valid, msg := validateTaskArgs(task,&user,&job,&filter,&submission)
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
            err = dbaccess.RegisterUser(user.Email,user.First,user.Last,user.Phone,user.Education)
            if err != nil {
                jsonResponse = fmt.Sprintf("{ \"error\" : \"%v\" }\n",err)
            } else {
                jsonResponse = fmt.Sprintf("{ \"success\" : \"Registered user %s\"}\n",user.Email)  
            }
            break
        case 1:
            job_id, err := dbaccess.CreateJob(job.Creator,job.Title,job.Description,job.Min_education,
                             job.Min_experience,job.Salary)
            if err != nil {
                jsonResponse = fmt.Sprintf("{ \"error\" : \"%v\" }\n",err)
            } else {
                jsonResponse = fmt.Sprintf("{ \"job_id\" : \"%s\" }\n",job_id)  
            }
        case 2:
            summaries, err := dbaccess.SearchJobs(filter.Posted, filter.Experience, filter.Education, filter.Salary, filter.Keyword) 
            if err != nil {
                jsonResponse = fmt.Sprintf("{ \"error\" : \"%v\" }\n",err)
            } else if len(summaries) == 0 {
                jsonResponse = "{ \"warning\" : \"No matching jobs found\"}"
            } else {
                resp, err := json.Marshal(summaries)
                if err != nil {
                    jsonResponse = fmt.Sprintf("{ \"error\" : \"%v\" }\n",err)
                } else {
                    jsonResponse = string(resp)
                } 
            }
        case 3:
            return_job, err := dbaccess.GetJobDetail(job.Job_id)
            if err != nil {
                jsonResponse = fmt.Sprintf("{ \"error\" : \"%v\" }\n",err)
            } else {
                resp, err := json.Marshal(return_job)
                if err != nil {
                    jsonResponse = fmt.Sprintf("{ \"error\" : \"%v\" }\n",err)
                } else {
                    jsonResponse = string(resp)
                } 
            }
        case 4:
            summaries, err := dbaccess.SearchOfferedJobs(job.Creator) 
            if err != nil {
                jsonResponse = fmt.Sprintf("{ \"error\" : \"%v\" }\n",err)
            } else if len(summaries) == 0 {
                jsonResponse = "{ \"warning\" : \"No matching jobs found\"}"
            } else {
                resp, err := json.Marshal(summaries)
                if err != nil {
                    jsonResponse = fmt.Sprintf("{ \"error\" : \"%v\" }\n",err)
                } else {
                    jsonResponse = string(resp)
                } 
            }
       case 5:
            summaries, err := dbaccess.SearchAppliedJobs(job.Creator) 
            if err != nil {
                jsonResponse = fmt.Sprintf("{ \"error\" : \"%v\" }\n",err)
            } else if len(summaries) == 0 {
                jsonResponse = "{ \"warning\" : \"No matching jobs found\"}"
            } else {
                resp, err := json.Marshal(summaries)
                if err != nil {
                    jsonResponse = fmt.Sprintf("{ \"error\" : \"%v\" }\n",err)
                } else {
                    jsonResponse = string(resp)
                } 
            } 
        case 6: // modify job
            job_id, err := dbaccess.ModifyJob(job.Creator,job.Job_id,job.Title,job.Description,job.Min_education,
                             job.Min_experience,job.Salary,job.Is_open)
            if err != nil {
                jsonResponse = fmt.Sprintf("{ \"error\" : \"%v\" }\n",err)
            } else {
                jsonResponse = fmt.Sprintf("{ \"modified_job_id\" : \"%s\" }\n",job_id)
            }
        case 7: // submit application for job
            job_id, err := dbaccess.SubmitJobApplication(submission.Email,submission.Job_id)                       
            if err != nil && job_id == "" {
                jsonResponse = fmt.Sprintf("{ \"error\" : \"%v\" }\n",err)
            } else if err != nil {
                jsonResponse = fmt.Sprintf("{ \"warning\" : \"%v\" }\n",err)
            } else {
                jsonResponse = fmt.Sprintf("{ \"applied_job_id\" : \"%s\" }\n",job_id)
            }           
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
