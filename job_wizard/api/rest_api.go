package api

// This module provides the Echo framework implementation
// for the JobWizard REST API. The Echo functions will call
// appropriate functions in dbaccess package in order to retrieve
// or update information in the system. The REST API is responsible
// for parsing arguments and for formatting responses in correct JSON
// but knows almost nothing about the conceptual model or the business
// logic of the application
// Copyright 2025 by CMKL University
// Created by Sally Goldin, 28 July 2025

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
    "github.com/segoldin/JobWizard/job_wizard/data"    
    "github.com/segoldin/JobWizard/job_wizard/dbaccess"
    "github.com/segoldin/JobWizard/job_wizard/helper"   	
	"github.com/labstack/echo/v4"
)

// Endpoints provided
func ApplicationPrivateRoute(_echo *echo.Group) {
	_echo.POST("/register", postRegisterUser)
	_echo.POST("/job/create", postCreateJob)
	_echo.GET("/search", getSearchJobs)
	_echo.GET("/search/detail",getSearchJobDetail)
	_echo.GET("/search/offered",getSearchJobsOffered)
	_echo.GET("/search/applied",getSearchJobsApplied)
	_echo.GET("/search/candidates",getSearchJobCandidates)				
	_echo.PUT("/job/modify",putModifyJob)
	_echo.POST("/job/submit", postSubmitJob)
}

/**************  Endpoint Implementations *************************/

// Implementation for /search API endpoint
func getSearchJobs(c echo.Context) error {
	var criteria data.Search_criteria
	var err error
	// copy parameters to struct used for validation
	criteria.User_email = strings.ToLower(c.QueryParam("email"))
   	criteria.Posted = c.QueryParam("posted")
   	tmpstring := c.QueryParam("experience") 
    if len(tmpstring) > 0 {
    	criteria.Experience, err = strconv.Atoi(tmpstring)
    	if err != nil {
    		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Invalid experience format - must be integer",
			})	
    	}   
    }
   	tmpstring = c.QueryParam("education") 
    if len(tmpstring) > 0 {
    	criteria.Education, err = strconv.Atoi(tmpstring)
    	if err != nil {
    		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Invalid education format - must be integer from 0 to 4",
			})	
    	}    	   
    }
   	tmpstring = c.QueryParam("salary") 
    if len(tmpstring) > 0 {
    	criteria.Salary, err = strconv.Atoi(tmpstring)
    	if err != nil {
    		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Invalid salary format - must be integer less than one million",
			})	
    	}    	   
    }
    criteria.Keyword = c.QueryParam("keyword")
    bOk, msg := helper.ValidateSearchCriteria(&criteria)
	if !bOk {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": msg,
		})		
	}
	jobs, err := dbaccess.SearchJobs(criteria.Posted,criteria.Experience,criteria.Education,criteria.Salary,criteria.Keyword)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": fmt.Sprintf("%v",err),
		})
	}
	if len(jobs) == 0 {
		return c.JSON(http.StatusOK, echo.Map{
			"warning": "No matching jobs found",
		})		
	}
	return c.JSON(http.StatusOK, jobs)
}

// Implementation for /search/detail API endpoint
func getSearchJobDetail(c echo.Context) error {
	var job data.Job_info
	var err error
	// copy parameters to struct used for validation
	job.Creator = strings.ToLower(c.QueryParam("email"))
   	job.Job_id = c.QueryParam("job_id")
	bOk, msg := helper.ValidateDetailRequest(&job)
	if !bOk {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": msg,
		})
	}
	foundjob, err := dbaccess.GetJobDetail(job.Job_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": fmt.Sprintf("%v",err),
		})
	}
	return c.JSON(http.StatusOK, foundjob)
}

// Implementation for /search/offered API endpoint
func getSearchJobsOffered(c echo.Context) error {
	// copy parameters to struct used for validation
	var job data.Job_info	
	job.Creator = strings.ToLower(c.QueryParam("email"))
	bOk, msg := helper.ValidateOfferedAppliedRequest(&job)
	if !bOk {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": msg,
		})
	}
	jobs, err := dbaccess.SearchOfferedJobs(job.Creator)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": fmt.Sprintf("%v",err),
		})
	}
	if len(jobs) == 0 {
		return c.JSON(http.StatusOK, echo.Map{
			"warning": "No matching jobs found",
		})		
	}
	return c.JSON(http.StatusOK, jobs)
}

// Implementation for /search/applied API endpoint
// Note this is currently almost identical to getSearchJobsOffered
// However, since both functions are short, we're implementing them
// separately for now.
func getSearchJobsApplied(c echo.Context) error {
	// copy parameters to struct used for validation
	var job data.Job_info	
	job.Creator = strings.ToLower(c.QueryParam("email"))
	bOk, msg := helper.ValidateOfferedAppliedRequest(&job)
	if !bOk {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": msg,
		})
	}
	jobs, err := dbaccess.SearchAppliedJobs(job.Creator)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": fmt.Sprintf("%v",err),
		})
	}
	if len(jobs) == 0 {
		return c.JSON(http.StatusOK, echo.Map{
			"warning": "No matching jobs found",
		})		
	}
	return c.JSON(http.StatusOK, jobs)
}

// Implementation for /search/candidates API endpoint
// Returns a list of users who have applied for a job
func getSearchJobCandidates(c echo.Context) error {
	// copy parameters to struct used for validation
	var job data.Job_info
	var err error
	// copy parameters to struct used for validation
	job.Creator = strings.ToLower(c.QueryParam("email"))
   	job.Job_id = c.QueryParam("job_id")
	bOk, msg := helper.ValidateDetailRequest(&job)
	if !bOk {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": msg,
		})
	}
	candidates,err := dbaccess.SearchCandidates(job.Creator,job.Job_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": fmt.Sprintf("%v",err),
		})
	}
	if len(candidates) == 0 {
		return c.JSON(http.StatusOK, echo.Map{
			"warning": "No candidates found",
		})			
	}
	return c.JSON(http.StatusOK, candidates)	
}

// Implementation for /register API endpoint
func postRegisterUser(c echo.Context) (err error) {
	input := new(data.User_info)
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}
	bOk, msg := helper.ValidateUserInfo(input)
	if !bOk {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": msg,
		})		
	}
	err = dbaccess.RegisterUser(input.Email,input.First,input.Last,input.Phone,input.Education)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})				
	}
	return c.JSON(http.StatusOK, echo.Map{
			"registered" : input.Email,
		})
}

// Implementation for /job/create API endpoint
func postCreateJob(c echo.Context) (err error) {
	input := new(data.Job_info)
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}
	bOk, msg := helper.ValidateJobInfo(input, true)
	if !bOk {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": msg,
		})		
	}
	job_id, err := dbaccess.CreateJob(input.Creator,input.Title,input.Description,input.Min_education, input.Min_experience, input.Salary) 
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})				
	}
	return c.JSON(http.StatusOK, echo.Map{
			"created_job" : job_id,
		})
}

// Implementation for /job/modify API endpoint
func putModifyJob(c echo.Context) (err error) {
	input := new(data.Job_info)
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}
	bOk, msg := helper.ValidateJobInfo(input, false)
	if !bOk {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": msg,
		})		
	}
	job_id, err := dbaccess.ModifyJob(input.Creator,input.Job_id,input.Title,input.Description,input.Min_education, input.Min_experience, input.Salary, input.Is_open) 
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})				
	}
	return c.JSON(http.StatusOK, echo.Map{
			"modified_job" : job_id,
		})
}

// Implementation for /job/submit API endpoint
func postSubmitJob(c echo.Context) (err error) {
	input := new(data.Submission)
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}
	bOk, msg := helper.ValidateJobSubmission(input)
	if !bOk {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": msg,
		})		
	}
	job_id, err := dbaccess.SubmitJobApplication(input.Email,input.Job_id) 
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})				
	}
	return c.JSON(http.StatusOK, echo.Map{
			"applied_for_job" : job_id,
		})
}
