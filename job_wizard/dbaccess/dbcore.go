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
        godotenv.Load(".env_evalrun")
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

// generate a new login token and store it in the database along with
// the date and time
func createNewToken(user_id string) (token string, err error) {
    // get the time, add it to the summed ascii values of the
    // username, then print in hex
    bytevals := []byte(user_id)
    var sum = 0
    for b := range bytevals {
        sum +=b
    }
    now := time.Now()
    nowstring := now.Format(timeFormatString)
    unixseconds := now.Unix()
    token = fmt.Sprintf("%x",unixseconds+int64(sum))
    token = "EVR" + token
    sqlcmd := fmt.Sprintf("UPDATE user SET latest_token='%s',last_login='%s' WHERE user_id='%s'",
                 token, nowstring, user_id)
    _, err = db.Exec(sqlcmd)
    if err != nil {
        return "", err
    }
    return token, nil
}

// generate a random string to use as the subject code
// We do this by randomly selecting 10 characters from an array and adding a 4 digit random number
func createRandomSubjectCode() (newcode string) {
    newcode = ""
    rand.Seed(time.Now().UnixNano())
    for i := 0; i < 10; i++ {
        index := rand.Intn(23)
        newcode += code_chars[index]
    }
    // now add a 6 digit suffix
    val := rand.Intn(999999)
    suffix := fmt.Sprintf("%04d",val)
    newcode += suffix
    //fmt.Printf("In createRandomSubjectCode code is |%s|\n",newcode)
    return newcode
}

// generate the next condition
// NB: We are currently not calling this 
// Instead the experimenter will choose a condition
func generateCondition() (condition string, condition_val int, err error) {
    db,err = connectDb(dbname)
    if err != nil {
        return "",0, err
    } 
    sqlcmd := fmt.Sprintf("SELECT condition_number FROM last_condition")
    row := db.QueryRow(sqlcmd)
    var c int
    err = row.Scan(&c)
    if err != nil {
        //fmt.Println("Error looking up condition_number")
        return "", 0, err
    }
   
    condition_val = ((c+1) % 3)

    sqlcmd = fmt.Sprintf("UPDATE last_condition SET condition_number=%d",condition_val)
    _, err = db.Exec(sqlcmd)
    if err != nil {
        return "", 0, err
    }
    // look up the string version in the DB and return
    sqlcmd = fmt.Sprintf("SELECT title FROM lu_condition WHERE value=%d",condition_val)
    row = db.QueryRow(sqlcmd)
    err = row.Scan(&condition)
    if err != nil {
        //fmt.Println("Error looking up condition name")
        return "", 0, err
    }
    return condition, condition_val, nil    
}

// Factorization - create sql for update
func createUpdateCommand(subject_id string, condition int, counselor_name string, orient_start string, counseling_start string, question_start string) (sqlcmd string) {
    sqlcmd = "UPDATE experiment SET "
    var addComma = false
    if counselor_name != "" {
        sqlcmd += fmt.Sprintf(" counselor_id = (SELECT id FROM counselor WHERE NAME='%s') ",counselor_name)
        addComma = true
    }
    if condition != -1 {
        if addComma {
            sqlcmd += ", "
        }
        addComma = true
        sqlcmd += fmt.Sprintf(" condition=%d",condition)
    }    
    if counseling_start != "" {
        if addComma {
            sqlcmd += ", "
        }
        addComma = true        
        sqlcmd += fmt.Sprintf(" counseling_start='%s'",counseling_start)
    }
    if question_start != "" {
        if addComma {
            sqlcmd += ", "
        }
        addComma = true       
        sqlcmd += fmt.Sprintf(" question_start='%s'",question_start)
    }
    sqlcmd += fmt.Sprintf(" WHERE subject_id='%s'",subject_id)
    //fmt.Printf("Constructed update command: |%s|\n",sqlcmd)
    return sqlcmd
}

// Factorization - check to see if the experiment exists and can be updated
// Also checks if this experiment is the most recently created. If so, we 
// will reset the condition if a subject withdraws
func checkExperimentStatus(subject_id string) (latest_experiment bool, err error) {
    sqlcmd := fmt.Sprintf("SELECT withdrawn, start_time, subject_end_time,counselor_end_time FROM experiment WHERE subject_id='%s'",
            subject_id)
    row := db.QueryRow(sqlcmd)
    var withdrawn int
    var start_time string
    var subject_end_time string
    var counselor_end_time string
    err = row.Scan(&withdrawn,&start_time,&subject_end_time,&counselor_end_time)
    if err != nil {
        return false, fmt.Errorf("Invalid subject ID; experiment not found")
    }
    if (subject_end_time != "") && (counselor_end_time != "") {
        return false, fmt.Errorf("Experiment is already finished; no updates possible.")
    }
    if withdrawn == 1 {
        return false, fmt.Errorf("Subject has withdrawn from the experiment; no updates possible.")        
    }
    // okay, we know it's open... is it the most recently created?
    // first, we need to reformat the start_time from the DB
    st, _ := time.Parse(time.RFC3339,start_time)
    formatted_start := st.Format(time.DateTime)    
    sqlcmd = "SELECT max(start_time) FROM experiment"
    row = db.QueryRow(sqlcmd)
    var max_start string
    err = row.Scan(&max_start)
    if err != nil {
        return false, err
    }
    //fmt.Printf("max_start is %s and formatted_start is %s\n",max_start,formatted_start)
    if max_start == formatted_start {
        latest_experiment = true
    } else {
        latest_experiment = false
    }
    //fmt.Printf("latest_experiment is %t\n",latest_experiment)
    return latest_experiment, nil
}

// Factorization - check to see if the counselor name has been set for this experiment
// If not, return an error. Called before we store a counselor's questionnaire responses

func checkCounselorSet(subject_id string) (err error) {
    sqlcmd := fmt.Sprintf("SELECT counselor_id FROM experiment WHERE subject_id='%s'",
            subject_id)
    row := db.QueryRow(sqlcmd)
    var counselor_id int
    err = row.Scan(&counselor_id)
    if err != nil {
        return fmt.Errorf("Invalid subject ID; experiment not found")
    }
    if counselor_id == 0 {
        return fmt.Errorf("Counselor has not been assigned for this experiment")
    }
    return nil
}

// check to see if the counselor name entered matches one in the DB
func checkCounselorValid(counselor_name string) (err error) {
    sqlcmd := fmt.Sprintf("select id from counselor where name='%s'",counselor_name)
    row := db.QueryRow(sqlcmd)
    var counselor_id int
    err = row.Scan(&counselor_id)
    if err != nil {
        return fmt.Errorf("specified counselor does not exist")
    }
    return nil
}

//******** Exported Functions *****************************//

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
    sqlcmd = fmt.Sprintf("INSERT INTO user (user_email, first_name, last_name, phone, max_education, created) values ('%s','%s','%s',%d,'%s')",
                                user_email, first_name, last_name, phone, max_education, nowstring)
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
    tx, err = db.Begin()
    if err != nil {
        return "", err
    }    
    now := time.Now()
    nowstring := now.Format(timeFormatString)     
    sqlcmd := 
      fmt.Sprintf("INSERT INTO job (created_by, title, description, min_education, min_years_experience, salary, created) values ('%s','%s','%s',%d,%d,%d,'%s')",
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

