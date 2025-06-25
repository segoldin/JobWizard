package main
// JobWizard demo application
// validation functions for command line arguments
// Created by Sally Goldin 2025-06-23
import (
	"fmt"
	"regexp"
    "strings"
    //"strconv"
)

var tasklist = [...]string{"register","create"} 

// Find the specified task in the task list. Return its index (0...) or -1 if not found
func findTask(task string) (index int) {
	index = -1
	for i,t := range tasklist {
		if strings.EqualFold(task,t) {
			index = i
			break
		} 
	}
	return index
}

func validateTaskArgs(task string, user user_info, job job_info) (bOk bool, msg string) {
	bOk = true
	taskIndex := findTask(task)
	if taskIndex < 0 {
		return false, "Invalid task specified"
	}
	switch taskIndex {
		case 0:
			bOk, msg = validateUserInfo(user)
			break
		case 1:
			bOk, msg = validateJobInfo(job, true) 
			break
	} 
	return bOk,msg 
}

// Check that all information needed to create a user is specified,
// and that the individual field values have valid format
func validateUserInfo(user user_info) (bOk bool, msg string) {
	bOk, msg = validateEmail(user.email)
	if bOk {
		bOk, msg = validateFirstLastName(user.first, "first")
	}	
	if bOk {
		bOk, msg = validateFirstLastName(user.last, "last")
	}	
	if bOk {
		bOk, msg = validatePhone(user.phone)
	}	
	if bOk {
		bOk, msg = validateEducation(user.education)
	}
	return bOk, msg
}

// Check that all information needed to create a job is specified,
// and that the individual field values have valid format
// If "is_create" then we are creating a new job and all fields are required
// Otherwise, the user can specify only the values that are to be changed
func validateJobInfo(job job_info, is_create bool) (bOk bool, msg string) {
	bOk, msg = validateEmail(job.creator)
	if bOk && is_create {
		bOk, msg = validateNonEmpty(job.title,"title")
	}	
	if bOk && is_create {
		bOk, msg = validateNonEmpty(job.description, "description")
	}	
	if bOk && (is_create || job.min_education > 0) {
		bOk, msg = validateEducation(job.min_education)
	}
	if bOk && (is_create || job.min_experience > 0) {
		bOk, msg = validateExperience(job.min_experience)
	}
	if bOk && (is_create || job.salary > 0) {
		bOk, msg = validateSalary(job.salary)
	}		
	return bOk, msg
}


// Check simply to see if the string passed is not empty
// Use the label to construct an error message if it is
func validateNonEmpty(parameter string, label string) (bOk bool, msg string) {
	if parameter == "" {
		return false, label + " must not be blank"
	}
	return true,""
}

// Do a simple validation of the email 
// This is not guaranteed to match every valid email but works in most cases
func validateEmail(email_addr string) (bOk bool, msg string) {
	if email_addr == "" {
		return false, "Missing user email"
	}
	var regex = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
	bOk, err := regexp.MatchString(regex, email_addr)
	if !bOk {
		if err == nil {
			msg = "Invalid email address"
		} else {
			fmt.Sprintf(msg, "Error: %v", err)
		}
	}
	return bOk, msg
}

// Validate first or last name
// Both must non-empty and alphabetic
// The 'which' argument indicates which one we are checking 
func validateFirstLastName(name string, which string) (bOk bool, msg string) {
	if name == "" {
		return false, "Missing user " + which + " name"
	}	
	var regex = "^[a-zA-Z]+$"
	bOk, err := regexp.MatchString(regex, name)
	if !bOk {
		if err == nil {
			msg = "Invalid " + which + " name"
		} else {
			fmt.Sprintf(msg,"Error: %v",err)
		}
	}
	return bOk, msg	
}

// Validate phone number. Must be 10 digits starting with 0
func validatePhone(phone_num_string string) (bOk bool, msg string) {
	if phone_num_string == "" {
		return false, "Missing user phone number"
	}	
	var regex = "^0[0-9]{9}$"
	bOk, err := regexp.MatchString(regex, phone_num_string)
	if !bOk {
		if err == nil {
			msg = "Invalid phone number"
		} else {
			fmt.Sprintf(msg,"Error: %v",err)
		}
	}
	return bOk, msg
}

// Validate education level. If missing we will assume 0
// Allowed values are 0 through 4
func validateEducation(ed_level int) (bOk bool, msg string) {
	bOk = true
	if (ed_level < 0) || (ed_level > 4) {
		bOk = false 
		msg = "Invalid education level"		
	} 
	return bOk, msg
}

// Validate experience. If missing we will assume 0
// Screen for ridiculous values
func validateExperience(experience int) (bOk bool, msg string) {
	bOk = true
	if (experience < 0) || (experience > 75) {
		bOk = false 
		msg = "Invalid years of experience"		
	} 
	return bOk, msg
}

// Validate monthly salary. If missing we will assume 0 which means unspecified
// Screen for ridiculous values
func validateSalary(salary int) (bOk bool, msg string) {
	bOk = true
	if (salary < 0) || (salary > 1000000) {  // upper limit is 1 million baht/month
		bOk = false 
		msg = "Invalid salary"		
	} 
	return bOk, msg
}

