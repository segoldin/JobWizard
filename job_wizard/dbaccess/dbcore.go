package dbaccess
// This module holds the core functions for accessing the database
// in the JobWizard demo application
// Created by Sally Goldin, 18 June 2025

import (
    "database/sql"
    "fmt"
    "os"
    "time"
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

