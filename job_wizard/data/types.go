package data
// Structure types used for passing information around
// Created by Sally Goldin 23 June 2025

type User_info struct {
    Email           string   `json:"email"`
    First           string   `json:"first"`
    Last            string   `json:"last"`   
    Phone           string   `json:"phone"`
    Education       int      `json:"education"`
}

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

type Job_summary struct {
    Job_id          string    `json:"job_id"`
    Title           string    `json:"title"`
    Is_open         bool      `json:"is_open"`
    Date_posted     string    `json:"date_posted"`    
}

type Search_criteria struct {
    User_email       string         
    Posted           string 
    Education        int
    Salary           int
    Keyword          string
}