package data
// Structure types used for passing information around
// Created by Sally Goldin 23 June 2025

// Used for registering a new user
type User_info struct {
    Email           string   `json:"email"`
    First           string   `json:"first"`
    Last            string   `json:"last"`   
    Phone           string   `json:"phone"`
    Education       int      `json:"education"`
}

// Used both for input and output - to create jobs as well as to return job detail
type Job_info struct {
    Job_id          string    `json:"job_id"`
    Creator         string    `json:"creator"`
    Title           string    `json:"title"`
    Description     string    `json:"description"`
    Min_education   int       `json:"min_education"`
    Min_experience  int       `json:"min_experience"`
    Salary          int       `json:"salary"` 
    Is_open         bool      `json:"is_open"`
    Date_posted     string    `json:"date_posted"`
}

// Used to return information from a job search
type Job_summary struct {
    Job_id          string    `json:"job_id"`
    Title           string    `json:"title"`
    Is_open         bool      `json:"is_open"`
    Date_posted     string    `json:"date_posted"`    
}

// Used to search for jobs
type Search_criteria struct {
    User_email       string         
    Posted           string 
    Experience       int    
    Education        int
    Salary           int
    Keyword          string
}

// Used to apply for a job
type Submission struct {
    Email           string   `json:"email"`
    Job_id          string   `json:"job_id"`   
}

// Used for returning information about an applicant
type Candidate struct {
    Email           string   `json:"email"`
    Name            string   `json:"name"`  // concatenated first and last name
    Phone           string   `json:"phone"`
    Applied_date    string   `json:"applied_date"`
}