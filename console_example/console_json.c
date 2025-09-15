/* console_json.c 
 * Example program showing how to invoke job_wizard
 * from a simple console program in order to execute a task,
 * then capture the output and parse as JSON.
 * 
 * Uses https://github.com/whyisitworking/C-Simple-JSON-Parser
 *
 * Created by Sally Goldin for SEN-210 on 15 September 2025
 */
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "json.h"

// structure to hold one search result, a job summary
typedef struct 
{
	char job_id[8]; 
	char title[64];
	int is_open;     // boolean
	char date_posted[32];
} JOB_SUMMARY_T;

// Open the redirect file, reads and returns the content
// as a null terminated text string. Note that the returned
// string must be freed by the caller when no longer needed
char* readOutput(char* redirectFilename)
{
	FILE* pF = NULL;
	char* textbuffer = NULL;

	pF = fopen(redirectFilename,"r");
	if (pF != NULL)
	{
		// get the size of the file
		fseek(pF, 0L, SEEK_END);
		int size = ftell(pF);
		rewind(pF);
		textbuffer = calloc(size+2, sizeof(char)); // +2 for 0 term
		if (textbuffer == NULL)
		{
			printf("Error allocating space to store results\n");
			return NULL;
		}
		if (fread(textbuffer,sizeof(char),size,pF) != size)
		{
			printf("Error reading data from result file\n");
			return NULL;
		}
	}
	return (textbuffer);
}

/* parse the passed string and return the generic structure (json_element)
 * build by C-Simple-JSON-Parser.
 * Also returns a status via the integer pointer - 0 if okay, else -1
 */
typed(json_element) parseJsonString(char* rawJsonText, int* pStatus)
{
	// based on example in the C-Simple-JSON-Parser repo
 	typed(json_element) element;
 	*pStatus = 0;  // assume it will work

	result(json_element) element_result = json_parse(rawJsonText);
	if (result_is_err(json_element)(&element_result)) 
	{
    	typed(json_error) error = result_unwrap_err(json_element)(&element_result);
    	fprintf(stderr, "Error parsing JSON: %s\n", json_error_to_string(error));
    	*pStatus = -1;
  	}
  	else
  	{
  		element = result_unwrap(json_element)(&element_result);
  	}
  	return element;
}

/* Attempts to parse the string passed as rawJsonText into 
 * an array of structures that represent job summaries. 
 * This is a two step process. The first, which parses the JSON,
 * will be the same for all job_wizard results. The second extracts
 * job summary information from the generic structures produced
 * by the first.
 * Returns 0 for success, -1 for error
 * If successful, also allocates and returns an array of results.
 * This array must be freed by the caller.
 * Also sets the value of pJobCount
 */
int parseSearchResults(char* rawJsonText, JOB_SUMMARY_T** allResults,int *pJobCount)
{
	int status = 0;
	int i,j;
	JOB_SUMMARY_T * resultArray = NULL;
	
	typed(json_element) element = parseJsonString(rawJsonText,&status);
	if (status < 0)
	{
		return status; // do we need to free the element?
	}
	// we expect an array of jobs
	typed(json_array) *arr = element.value.as_array;
	*pJobCount = arr->count;
	printf("Found %d jobs\n",*pJobCount);
	resultArray = (JOB_SUMMARY_T*) calloc(*pJobCount,sizeof(JOB_SUMMARY_T));
	if (resultArray == NULL)
	{
		printf("Error allocating job summary structures\n");
		status = -1;
		return status;
	}
	*allResults = resultArray;
	for (j=0; j < *pJobCount; j++)
	{
		typed(json_element) element = arr->elements[j];
		typed(json_object) *obj = element.value.as_object;
		for (i = 0; i < obj->count; i++) 
		{
  			typed(json_entry) entry = *(obj->entries[i]);
  			typed(json_string) key = entry.key;
			typed(json_element_value) value = entry.element.value;
			if (strcmp(key,"job_id") == 0)
			{
				strcpy(resultArray[j].job_id,value.as_string);
			}
			else if (strcmp(key,"title") == 0)
			{
				strcpy(resultArray[j].title,value.as_string);
			}
			else if (strcmp(key,"is_open") == 0)
			{
				if (value.as_boolean)
					resultArray[j].is_open = 1;
				else
					resultArray[j].is_open = 0;
			}
			else if (strcmp(key,"date_posted") == 0)
			{
				strcpy(resultArray[j].date_posted,value.as_string);
			}
			else
			{
				printf("Unrecognized object key %s\n",key);
			}
		}
	}
	json_free(&element);

	return status;
}

/* Display the contents of a job structure as a line of text
 */
void printJobSummary(JOB_SUMMARY_T job)
{
	printf("%6s %32s %4s %12s\n",
		job.job_id,job.title,job.is_open? "t" : "f", job.date_posted);
}

/* main program executes a search with no arguments */
int main(int argc, char* argv)
{
	char* userEmail = "sally@cmkl.ac.th";
	char* outputFile = "output.txt";
	char jobwizardCmd[2048];
	char* resultText = NULL;  // holds results, must be freed after use
	int returnCode;
	int jobCount;
	int j;
	JOB_SUMMARY_T * searchResults = NULL; // array of structs allocated for parsed JSON

	// Create the command to execute
	// In your real console-based UI, you probably want different functions
	// that know about the arguments for different job_wizard commands
	sprintf(jobwizardCmd,"./job_wizard -task search -email %s > %s 2>&1",userEmail,outputFile);
	// note that 2>&1 redirects both standard output and standard error to the file output.txt
	returnCode = system(jobwizardCmd);
	if (returnCode != 0)
	{
		printf("Error %d executing job_wizard command\n");
		printf("Command: |%s|\n",jobwizardCmd);
	}
	else
	{
		printf("Successfully ran job_wizard!\n");
		resultText = readOutput(outputFile);
		if (resultText != NULL)
		{
			returnCode = parseSearchResults(resultText,&searchResults,&jobCount);
			if (returnCode != 0)
				printf("Error parsing JSON\n");
			if (searchResults != NULL)
			{
				for (j = 0; j < jobCount; j++)
					printJobSummary(searchResults[j]);
				free(searchResults);
			}
			free(resultText);
		}
	}
}
