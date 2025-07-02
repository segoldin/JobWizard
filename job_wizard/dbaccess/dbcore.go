package dbaccess
// This module holds the core functions for accessing the database
// in the JobWizard demo application
// Created by Sally Goldin, 18 June 2025

import (
    "database/sql"
    "fmt"
    "os"
    "strconv"
    "time"
    "github.com/segoldin/JobWizard/job_wizard/data"      
    "github.com/joho/godotenv"     
     _ "github.com/mattn/go-sqlite3"
)

const timeFormatString = "2006-01-02 15:04 +700"
const formatWithoutZone = "2006-01-02 15:04"

// uppercase alphabet leaving out I and O
var (
    db *sql.DB
    dbname = ""
)
//**************** Private Functions *******************************//

// Connect to the database if not already done
// Use module global var for the db connection
// Return the db connection and error
func connectDb(dbname string) (dbconn *sql.DB, err error) {
    if db != nil {
        return db,nil
    }
    if dbname == "" {
        godotenv.Load(".env_jobwizard")
        dbname = os.Getenv("JOBWIZARD_DB_NAME")
    }    
    dbconn, err = sql.Open("sqlite3", dbname)
    if err != nil {
        msg := fmt.Sprintf("Error opening the database - %v\n",err)
        return nil,fmt.Errorf(msg)
    }
    //fmt.Printf("In connectDb - dbname is %s\n",dbname)
    return dbconn,nil
}


//******** Exported Functions *****************************//

func CheckConnection() bool {
    _,err := connectDb(dbname)
    if err != nil {
        return false
    } else {
        return true
    }    
}

func IsRegisteredUser(user_email string) (bRegistered bool, err error) {
    db, err = connectDb(dbname)
    if err != nil {
        return false, err
    }     
    sqlcmd := fmt.Sprintf("SELECT id FROM user WHERE user_email='%s'",user_email)
    row := db.QueryRow(sqlcmd)
    var id int
    err = row.Scan(&id)
    if err != nil {
        bRegistered = false
    } else {
        bRegistered = true
    }
    return bRegistered, nil

}

// Function to create a new user, implementing the Register use case
// If user email already exists, will return an error
func RegisterUser(user_email string, first_name string, last_name string, phone string, education int) (err error) {
    db,err = connectDb(dbname)
    if err != nil {
        return err
    } 
    sqlcmd := fmt.Sprintf("SELECT id from user where user_email='%s'",
            user_email)
    rows,err := db.Query(sqlcmd)
    if err != nil {
        return err
    }
    rowcount := 0
    var idval int
    for rows.Next() {
        rowcount++
        err = rows.Scan(&idval)
        if err != nil {
            rows.Close()
            return err
        }
        if rowcount == 1 {
            break;
        }
    }
    rows.Close() // need to explicitly close before insert/update or DB will be locked
    if rowcount > 0 {
        return fmt.Errorf("Email is not unique; user not created")
    }
    now := time.Now()
    nowstring := now.Format(timeFormatString) 
    sqlcmd = fmt.Sprintf("INSERT INTO user (user_email, first_name, last_name, phone, max_education, created) values ('%s','%s','%s','%s',%d,'%s')",
                                user_email, first_name, last_name, phone, education, nowstring)
    _,err = db.Exec(sqlcmd)
    if err != nil {
        return err
    }
    return nil
}

// Function to create a new job, implementing the Create Job use case
// Returns the ID (autoincrement) of the job, transformed into a string with leading zeros
func CreateJob(creator_email string, title string, desc string, education int, experience int, salary int) (job_id string, err error) {
    db,err = connectDb(dbname)
    if err != nil {
          return "", err
    }
    // do this in a transaction in case somebody else is also creating a job
    tx, err := db.Begin()
    if err != nil {
        return "", err
    }    
    now := time.Now()
    nowstring := now.Format(timeFormatString)     
    sqlcmd := 
      fmt.Sprintf("INSERT INTO job (created_by, title, description, min_education, min_years_experience, salary, created) values ('%s','%s','%s', %d, %d, %d,'%s')",
       creator_email, title, desc, education, experience, salary, nowstring) 
    _,err = tx.Exec(sqlcmd)
    if err != nil {
        tx.Rollback()
        return "",err
    }
    sqlcmd = "SELECT MAX(id) FROM job"
    var id int
    row := tx.QueryRow(sqlcmd)
    err = row.Scan(&id)  
    if err != nil {
        tx.Rollback()
        return "", err
    }
    tx.Commit()
    job_id = fmt.Sprintf("%05d",id)
    return job_id, nil
}

// Function to search for jobs based on criteria, implementing the Search Jobs use case
// Returns an array of job summary structures in posted date order (descending) or error
func SearchJobs(posted_criterion string, min_experience int, min_education int, salary int, keyword string) (summaries []data.Job_summary, err error) {
    var added_where = false   
    sqlcmd := "SELECT id,title,is_open,created FROM job "
    if posted_criterion != "" {
        if !added_where {
            sqlcmd += " where "
            added_where = true
        }
        clause := fmt.Sprintf(" created >= '%s' ",posted_criterion)
        sqlcmd += clause
    }
   if min_experience != 0 {
        if !added_where {
            sqlcmd += " where "
            added_where = true
        } else {
            sqlcmd += " and "
        }
        clause := fmt.Sprintf(" min_years_experience <= %d ",min_experience)
        sqlcmd += clause
    }    
    if min_education != 0 {
        if !added_where {
            sqlcmd += " where "
            added_where = true
        } else {
            sqlcmd += " and "
        }
        clause := fmt.Sprintf(" min_education <= %d ",min_education)
        sqlcmd += clause
    }
    if salary != 0 {
        if !added_where {
            sqlcmd += " where "
            added_where = true
        } else {
            sqlcmd += " and "
        }
        clause := fmt.Sprintf(" salary >= %d ",salary)
        sqlcmd += clause
    }
    if keyword != "" {
        if !added_where {
            sqlcmd += " where "
            added_where = true
        } else {
            sqlcmd += " and "
        }
        clause := fmt.Sprintf(" title like '%%%s%%' ",keyword)
        sqlcmd += clause
    }
    sqlcmd += " order by created desc"   
    summaries, err = doSearchOperation(sqlcmd)
    return summaries, err
}

// Function to get all the detail for a particular job, implementing the Show Job Detail use case
// Returns a filled in job structure if the job id is found
func GetJobDetail(job_id string) (foundjob data.Job_info, err error) {
    db,err = connectDb(dbname)
    if err != nil {
        return foundjob, err
    }     
    id, _ := strconv.Atoi(job_id)  // we already validated this 
    sqlcmd := "SELECT created_by, title, description, min_education, min_years_experience, salary, is_open, created"
    whereclause := fmt.Sprintf(" from job where id = %d", id);
    sqlcmd = sqlcmd + whereclause
    row := db.QueryRow(sqlcmd)
    err = row.Scan(&foundjob.Creator,&foundjob.Title,&foundjob.Description,
         &foundjob.Min_education,&foundjob.Min_experience,&foundjob.Salary,&foundjob.Is_open,&foundjob.Date_posted)
    if err != nil {
        if err == sql.ErrNoRows {
            return foundjob, fmt.Errorf("No matching job found")
        } else {
            return foundjob, err
        }
    }
    foundjob.Job_id = fmt.Sprintf("%05d",id)
    foundjob.Date_posted = foundjob.Date_posted[0:10]
    return foundjob, nil 
}

// Function to search for jobs offered by a particular user
func SearchOfferedJobs(user_email string) (summaries []data.Job_summary, err error) {
    sqlcmd := fmt.Sprintf("SELECT id,title,is_open,created FROM job where created_by='%s'",user_email)
    sqlcmd += " order by created desc"
    summaries, err = doSearchOperation(sqlcmd)
    return summaries, err
}

// Function to search for jobs applied to by a particular user
func SearchAppliedJobs(user_email string) (summaries []data.Job_summary, err error) {
    sqlcmd := "SELECT j.id,j.title,j.is_open,j.created FROM job j, job_applications ja "
    sqlcmd += " where j.id=ja.job_id and "
    sqlcmd += fmt.Sprintf("ja.user_email='%s'",user_email)
    sqlcmd += " order by created desc"
    summaries, err = doSearchOperation(sqlcmd)
    return summaries, err
}

// This function is a factorization that handles searching for jobs 
// and returning summaries
// It is called by three different tasks, which use different criteria/queries
// but otherwise handle the return information the same way
func doSearchOperation(sqlcmd string) (summaries []data.Job_summary, err error) {    
    db,err = connectDb(dbname)
    if err != nil {
        return summaries, err
    } 
    rows,err := db.Query(sqlcmd)
    if err != nil {
        return summaries, err
    }
    defer rows.Close()
    var idval int
    var title string
    var is_open bool
    var posted string
    for rows.Next() {
        err = rows.Scan(&idval,&title,&is_open,&posted)
        if err != nil {
            rows.Close()
            return summaries, err
        }
        var job data.Job_summary
        job.Job_id = fmt.Sprintf("%05d",idval)
        job.Title = title
        job.Is_open = is_open
        job.Date_posted = posted[0:10] 
        summaries = append(summaries,job)
    }
    return summaries, nil
}
