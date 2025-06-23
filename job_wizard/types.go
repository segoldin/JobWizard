package main
// Structure types used for command line functionality
// Created by Sally Goldin 23 June 2025

type user_info struct {
    email           string
    first           string 
    last            string 
    phone           string 
    education       string    
}

type job_info struct {
    creator         string 
    title           string 
    description     string 
    min_education   int
    min_experience  int
    salary          int
}
