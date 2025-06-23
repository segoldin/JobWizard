package main
// JobWizard demo application
// validation functions for command line arguments
// Created by Sally Goldin 2025-06-23
import (
	"fmt"
	"regexp"
    "strings"
    "strconv"
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
			bOk, msg = validateJobInfo(job) 
			break
	} 
	return bOk,msg 
}

// Check that all information needed to create a user is specified,
// and that the individual field values have valid format
func validateUserInfo(user user_info) (bOk bool, msg string) {
	bOk, msg = validateEmail(user.email)
	fmt.Printf("bOk is %t\n",bOk)	
	if bOk {
		bOk, msg = validateFirstLastName(user.first, "first")
	}
	fmt.Printf("bOk(2) is %t\n",bOk)		
	if bOk {
		bOk, msg = validateFirstLastName(user.last, "last")
	}
	fmt.Printf("bOk(3) is %t\n",bOk)		
	if bOk {
		bOk, msg = validatePhone(user.phone)
	}
	fmt.Printf("bOk(4) is %t\n",bOk)		
	if bOk {
		bOk, msg = validateEducation(user.education)
	}
	return bOk, msg
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
func validateEducation(education_string string) (bOk bool, msg string) {
	if len(education_string) == 0 {
		return true,""
	}
	level, err := strconv.Atoi(education_string)
	if err != nil {
		bOk = false 
		msg = "Invalid education level"
	} else if (level < 0) || (level > 4) {
		bOk = false 
		msg = "Invalid education level"		
	} 
	return bOk, msg
}

// Check that all information needed to create a job is specified,
// and that the individual field values have valid format
func validateJobInfo(user job_info) (bOk bool, msg string) {
	return true,""
}