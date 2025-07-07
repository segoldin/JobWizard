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
    flag.Usage = customUsage
    flag.Parse()
    if help {
        if task != "" {
            customUsageTask(task)
        } else {
            customUsage()
        }
    } 
    if server {
        setupAPI()
    } else {
        commandLineFunction()
    }
}

// display information about the arguments for each task
// or for a single task as specified by task_name
func customUsage() {
    fmt.Println("\nGeneral usage: ./job_wizard -task <taskname> [arguments...]")
    fmt.Println("\tWrites results to standard output in JSON format\n")
    fmt.Println("Available tasks: ")
    fmt.Println("\tregister\tCreate a new user in the database")
    fmt.Println("\tcreate\t\tCreate a new job posting")
    fmt.Println("\tsearch\t\tGeneral search for jobs")
    fmt.Println("\tdetail\t\tSee detailed information about a selected job")
    fmt.Println("\toffered\t\tSearch for jobs created by me")
    fmt.Println("\tapplied\t\tSearch for jobs that I have applied for")
    fmt.Println("\tmodify\t\tModify a job created by me")
    fmt.Println("\tsubmit\t\tSubmit an application for a job\n")
    fmt.Println("For task-specific arguments, type ./job_wizard -help=true -task <task_name>\n")
    os.Exit(0)                      
}


// display information about the arguments for a single task as specified by task_name
func customUsageTask(task_name string) {
    fmt.Println("\nGeneral usage: ./job_wizard -task <taskname> [arguments...]")
    fmt.Println("\tWrites results to standard output in JSON format\n")
    task_index := findTask(task_name)
    if (task_index < 0) {
        fmt.Printf("Unknown task '%s'\n",task_name)
        os.Exit(0)
    }
    switch(task_index) {
        case 0: // register
            fmt.Println("Register a new email as a JobWizard user")
            fmt.Println("Arguments for register task:")
            fmt.Println("\t-email <NEW email address>")
            fmt.Println("\t-first <first name>")
            fmt.Println("\t-last <last name>")
            fmt.Println("\t-phone <10 digit Thai phone>")
            fmt.Println("\t-education <integer 0 to 4>")
            fmt.Println("All arguments are required\n")
            fmt.Println("Example: ./job_wizard -task register -email sally@gmail.com -first Sally -last Goldin -phone 0987651122 -education 4\n")
            break;               
        case 1: // create
            fmt.Println("Create a new job posting in the JobWizard database")
            fmt.Println("Arguments for create task:")
            fmt.Println("\t-creator <email of registered user>")
            fmt.Println("\t-title <job title in quotes>")
            fmt.Println("\t-description <job description in quotes, up to 1024 chars>")
            fmt.Println("\t-min_education <integer 0 to 4>")
            fmt.Println("\t-min_experience <integer 0 to 75>")
            fmt.Println("\t-salary <monthly salary in baht, 0 means unspecified>")
            fmt.Println("Creator, title and description are required\n")
            fmt.Println("Example: ./job_wizard -task create -creator sally@gmail.com -title \"Front End Developer\" -description \"Build user interfaces for enterprise web applications\" -min_education 2 -salary 35000\n")         
            break
        case 2: // search
            fmt.Println("Search for jobs based on criteria, and print summaries")
            fmt.Println("Arguments for search task:")
            fmt.Println("\t-email <email of registered user>")          
            fmt.Println("\t-min_education <integer 1 to 4>")
            fmt.Println("\t-min_experience <integer >")
            fmt.Println("\t-salary <monthly salary in baht>")          
            fmt.Println("\t-posted <date: YYYY-MM-DD>")
            fmt.Println("\t-keyword <keyword to search for in title>")          
            fmt.Println("Only email is required\n")
            fmt.Println("Example: ./job_wizard -task search -email sally@gmail.com -salary 30000 -keyword Developer\n")
            break
        case 3: // detail
            fmt.Println("Return all detailed information for a specific job")
            fmt.Println("Arguments for detail task:")
            fmt.Println("\t-email <email of registered user>")
            fmt.Println("\t-job_id <show detail for what job>\n")
            fmt.Println("All arguments are required\n")    
            fmt.Println("Example: ./job_wizard -task detail -email sally@gmail.com -job_id 00003\n")
            break
        case 4: // offered
            fmt.Println("Return summaries for all jobs created/posted by a user")
            fmt.Println("Arguments for offered task:")
            fmt.Println("\t-creator <email of registered job creator>")
            fmt.Println("All arguments are required\n")    
            fmt.Println("Example: ./job_wizard -task offered -creator sally@gmail.com\n")
            break
        case 5: // applied
            fmt.Println("Return summaries for all jobs a user has applied for")
            fmt.Println("Arguments for applied task:")
            fmt.Println("\t-email <email of registered user>")
            fmt.Println("All arguments are required\n")    
            fmt.Println("Example: ./job_wizard -task applied -email sally@gmail.com\n")
        case 6: // modify job
            fmt.Println("Modify some attributes of a specific job")
            fmt.Println("Arguments for modify task:")
            fmt.Println("\t-creator <email of registered user>")
            fmt.Println("\t-job_id <modify what job>")            
            fmt.Println("\t-title <job title in quotes>")
            fmt.Println("\t-description <job description in quotes, up to 1024 chars>")
            fmt.Println("\t-min_education <integer 0 to 4>")
            fmt.Println("\t-min_experience <integer 0 to 75>")
            fmt.Println("\t-salary <monthly salary in baht, 0 means unspecified>")
            fmt.Println("\t-is_open=false")
            fmt.Println("Creator and job_id are required, changes any other attributes specified\n")
            fmt.Println("Example: ./job_wizard -task modify -creator sally@gmail.com -job_id 00002 -title \"User Experience Developer\" -salary 38000\n")         
            break
        case 7: // submit application for job
            fmt.Println("Apply for a particular job (submit application)")
            fmt.Println("Arguments for submit task:")
            fmt.Println("\t-email <email of registered user>")
            fmt.Println("\t-job_id <apply for what job>\n")                  
            fmt.Println("All arguments are required\n")    
            fmt.Println("Example: ./job_wizard -task submit -email sally@gmail.com -job_id 00014\n")
            break           
        default:
            fmt.Println("Invalid task specified\n")                     
    }
    os.Exit(0)
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
