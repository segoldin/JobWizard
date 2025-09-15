/* Example program showing how to invoke job_wizard
 * from a simple console program in order to execute a task,
 * then capture the output and print (without parsing).
 *
 * Created by Sally Goldin for SEN-210 on 15 September 2025
 */
#include <stdio.h>
#include <stdlib.h>
#include <strings.h>

// Open the redirect file and display the contents
void displayOutput(char* redirectFilename)
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
		textbuffer = calloc(size + 2, sizeof(char)); // +2 for terminating 0
		if (textbuffer == NULL)
		{
			printf("Error allocating space to store results\n");
			return;
		}
		if (fread(textbuffer,sizeof(char),size,pF) != size)
		{
			printf("Error reading data from result file\n");
			return;
		}
		printf("RESULTS\n");
		printf(textbuffer);
		free(textbuffer);
	}
}

/* main program executes a search with no arguments */
int main(int argc, char* argv)
{
	char* userEmail = "sally@cmkl.ac.th";
	char* outputFile = "output.txt";
	char jobwizardCmd[2048];
	int returnCode;

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
		printf("Success!\n");
		displayOutput(outputFile);
	}
}