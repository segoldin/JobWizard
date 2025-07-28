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
	//_echo.POST("/register", postRegisterUser)
	//_echo.POST("/job/create", postCreateJob)
	_echo.GET("/search", getSearchJobs)
	//_echo.GET("/search/detail",getSearchJobDetail)
	//_echo.GET("/search/offered",getSearchJobsOffered)
	//_echo.GET("/search/applied",getSearchJobsApplied)
	//_echo.GET("/search/candidates",getSearchJobCandidates)				
	//_echo.PUT("/job/modify",putModifyJob)
	//_echo.POST("/job/submit", postSubmitJob)
}

// structures for POST and PUT endpoints
// Just rename the structs used for command line processing

type UserRequest data.User_info

type JobRequest data.Job_info

type SubmissionRequest data.Submission

/**************  Endpoint Implementations *************************/

// Implementation for /search API endpoint
func getSearchJobs(c echo.Context) error {
	var criteria data.Search_criteria
	var err error
	// copy parameters to struct used for validation
	criteria.User_email = strings.ToLower(c.QueryParam("user_email"))
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

