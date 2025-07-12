mkdir windows
mkdir windows/database
env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -o windows/job_wizard.exe -a .
cp -p .env_jobwizard windows
cp -pr database/jobwizard_db windows/database
zip -r JobWizard.zip windows/*.* windows/.env* windows/database/*


