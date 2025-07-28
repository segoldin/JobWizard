package helper
// JobWizard demo application
// validation functions for command line arguments
// Created by Sally Goldin 2025-06-23
import (
	"fmt"
	"regexp"
    "strings"
    "time"
    "strconv"
    "github.com/segoldin/JobWizard/job_wizard/data"
    "github.com/segoldin/JobWizard/job_wizard/dbaccess"      
)

var tasklist = [...]string{"register","create","search","detail","offered","applied","modify","submit","candidates"} 

const (
	timeFormatString = "2006-01-02 15:04 +700"
	dateOnlyString = "2006-01-02"
)

// Find the specified task in the task list. Return its index (0...) or -1 if not found
func FindTask(task string) (index int) {
	index = -1
	for i,t := range tasklist {
		if strings.EqualFold(task,t) {
			index = i
			break
		} 
	}
	return index
}

// Pass all structs used for arguments 
// Note that some fields are used by multiple tasks
// We pass pointers so that any changes or copying gets preserved in the caller
func ValidateTaskArgs(task string, user *data.User_info, job *data.Job_info, filter *data.Search_criteria, submission *data.Submission) (bOk bool, msg string) {
	bOk = true
	taskIndex := FindTask(task)
	if taskIndex < 0 {
		return false, "Invalid task specified"
	}
	// set all email addresses to lower case
	user.Email = strings.ToLower(user.Email)
	job.Creator = strings.ToLower(job.Creator)
	switch taskIndex {
		case 0:
			bOk, msg = ValidateUserInfo(user)
			break
		case 1:
			bOk, msg = ValidateJobInfo(job, true) 
			break
		case 2:
			// can only define a command line arg once, so we copy from other structs
			filter.User_email = user.Email
			filter.Experience = job.Min_experience
			filter.Education = job.Min_education
			filter.Salary = job.Salary
			bOk, msg = ValidateSearchCriteria(filter)
			break
		case 3: 
			// detail
			job.Creator = user.Email
			bOk, msg = ValidateDetailRequest(job)
			break
		case 4: 
			// jobs offered search
			// will use job.Creator
			bOk, msg = ValidateOfferedAppliedRequest(job)
			break
		case 5: 
			// jobs applied search
			job.Creator = user.Email
			bOk, msg = ValidateOfferedAppliedRequest(job)
			break 
		case 6:
			bOk, msg = ValidateJobInfo(job, false) 
			break
		case 7: // submit a job application
			submission.Email = user.Email
			submission.Job_id = job.Job_id
			bOk, msg = ValidateJobSubmission(submission) 
			break
		case 8: // candidates
			// same arguments as detail request
			// will use job.Creator
			bOk, msg = ValidateDetailRequest(job)
			break											 				 			
	} 
	return bOk,msg 
}

// Check that all information needed to create a user is specified,
// and that the individual field values have valid format
func ValidateUserInfo(user *data.User_info) (bOk bool, msg string) {
	bOk, msg = validateEmail(user.Email)
	if bOk {
		bOk, msg = validateFirstLastName(user.First, "first")
	}	
	if bOk {
		bOk, msg = validateFirstLastName(user.Last, "last")
	}	
	if bOk {
		bOk, msg = validatePhone(user.Phone)
	}	
	if bOk {
		bOk, msg = validateEducation(user.Education)
	}
	return bOk, msg
}

// Check that all information needed to create a job is specified,
// and that the individual field values have valid format
// If "is_create" then we are creating a new job and all fields are required
// Otherwise, the user can specify only the values that are to be changed
func ValidateJobInfo(job *data.Job_info, is_create bool) (bOk bool, msg string) {
	bOk, msg = validateEmail(job.Creator)
	if bOk {
		bRegistered, _ := dbaccess.IsRegisteredUser(job.Creator)
		if !bRegistered {
			bOk = false
			msg = "Unknown user email"
		}
	}
	if bOk && !is_create {
		idstring := job.Job_id
		idval, err := strconv.Atoi(idstring)
		if (err != nil) || (idval <= 0) {
			bOk = false 
			msg = "Invalid job ID specified"
		} else {
			bOk = true
			msg = ""
		}
	}
	if bOk && is_create {
		bOk, msg = validateNonEmpty(job.Title,"title")
	}
	if bOk {
		bOk, msg = validateLength(job.Title,64,"title")
	}		
	if bOk && is_create {
		bOk, msg = validateNonEmpty(job.Description, "description")
	}
	if bOk {
		bOk, msg = validateLength(job.Description,1024,"description")
	}			
	if bOk && (is_create || job.Min_education != 0) {
		bOk, msg = validateEducation(job.Min_education)
	}
	if bOk && (is_create || job.Min_experience != 0) {
		bOk, msg = validateExperience(job.Min_experience)
	}
	if bOk && (is_create || job.Salary != 0) {
		bOk, msg = validateSalary(job.Salary)
	}		
	return bOk, msg
}

// Check the specified search criteria
// All are optional except for the user, but numeric values have limits
// If nothing is specified, the search will return all jobs
func ValidateSearchCriteria(filter *data.Search_criteria) (bOk bool, msg string) {
	bOk, msg = validateEmail(filter.User_email)
	if bOk {
		bRegistered, _ := dbaccess.IsRegisteredUser(filter.User_email)
		if !bRegistered {
			bOk = false
			msg = "Unknown user email"
		}
	}	
	if bOk && filter.Posted != ""{
		bOk, msg = validateDate(filter.Posted)
		if bOk {
			// add a time for DB search
			filter.Posted += " 00:01 +700"
		}
	}
	if bOk && (filter.Experience != 0) {
		bOk, msg = validateExperience(filter.Experience)
	}		
	if bOk && (filter.Education != 0) {
		bOk, msg = validateEducation(filter.Education)
	}
	if bOk && (filter.Salary != 0) {
		bOk, msg = validateSalary(filter.Salary)
	}		
	// no constraints on keyword criterion
	return bOk, msg
}

// check to see that the ID is set and is a positive integer
func ValidateDetailRequest(job *data.Job_info) (bOk bool, msg string) {	
	bOk, msg = validateEmail(job.Creator) // not really the creator... just use this field
	if bOk {
		bRegistered, _ := dbaccess.IsRegisteredUser(job.Creator)
		if !bRegistered {
			bOk = false
			msg = "Unknown user email"
		}
	}	
	if bOk {
		idstring := job.Job_id
		idval, err := strconv.Atoi(idstring)
		if (err != nil) || (idval <= 0) {
			bOk = false 
			msg = "Invalid job ID specified"
		} else {
			bOk = true
			msg = ""
		}
	}
	return bOk, msg
}

// check to see that the ID is set and is a positive integer
func ValidateJobSubmission(submission *data.Submission) (bOk bool, msg string) {	
	bOk, msg = validateEmail(submission.Email) 
	if bOk {
		bRegistered, _ := dbaccess.IsRegisteredUser(submission.Email)
		if !bRegistered {
			bOk = false
			msg = "Unknown user email"
		}
	}	
	if bOk {
		idstring := submission.Job_id
		idval, err := strconv.Atoi(idstring)
		if (err != nil) || (idval <= 0) {
			bOk = false 
			msg = "Invalid job ID specified"
		} else {
			bOk = true
			msg = ""
		}
	}
	return bOk, msg
}

// Specialized searches
// The only required argument is the email, which is interpreted differently
// depending on the task
func ValidateOfferedAppliedRequest(job *data.Job_info) (bOk bool, msg string) {	
	bOk, msg = validateEmail(job.Creator) // not really the creator... just use this field
	if bOk {
		bRegistered, _ := dbaccess.IsRegisteredUser(job.Creator)
		if !bRegistered {
			bOk = false
			msg = "Unknown user email"
		}
	}	
	return bOk, msg
}

// check to see if the passed date assumed to be in form YYYY-MM-DD is valid
func validateDate(datestring string) (bOk bool, msg string) {
	_, err := time.Parse(dateOnlyString,datestring)
	if err != nil {
		bOk = false 
		msg = "Invalid date"
	} else {
		bOk = true
		msg = ""
	}
	return bOk, msg
}

// Check simply to see if the string passed is not empty
// Use the label to construct an error message if it is
func validateNonEmpty(parameter string, label string) (bOk bool, msg string) {
	trimmed := strings.TrimSpace(parameter)
	if trimmed == "" || strings.HasPrefix(trimmed,"-") {
		return false, label + " must not be blank"
	}
	return true,""
}
// Be sure that a parameter is not too long
// Use label to construct an error message if it is
func validateLength(parameter string, maxlen int, label string) (bOk bool, msg string) {
	if len(parameter) > maxlen {
		msg = fmt.Sprintf("%s must be %d characters or less", label, maxlen)
		return false, msg
	}
	return true,""
}

// Do a simple validation of the email 
// This is not guaranteed to match every valid email but works in most cases
func validateEmail(email_addr string) (bOk bool, msg string) {
	trimmed := strings.TrimSpace(email_addr)
	if trimmed == "" {
		return false, "Missing user email"
	}
	bOk, msg = validateLength(email_addr, 32, "Email")
	if !bOk {
		return bOk, msg
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
	bOk, msg = validateLength(name, 32, which + " name")
	if !bOk {
		return bOk, msg
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

