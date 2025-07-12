-- Script for initializing the JobWizard database
-- 
-- Created by Sally Goldin, 18 June 2025
--

-- Look up table provides text for education level 
CREATE TABLE IF NOT EXISTS lu_education (
	value int,
	title VARCHAR(24)
);

-- Populate this lookup table
INSERT INTO lu_education (value,title) values (0,'Unknown');
INSERT INTO lu_education (value,title) values (1,'High school diploma');
INSERT INTO lu_education (value,title) values (2,'Bachelors degree');
INSERT INTO lu_education (value,title) values (3,'Masters degree');
INSERT INTO lu_education (value,title) values (4,'Doctoral degree');

-- Registered users
CREATE TABLE IF NOT EXISTS user (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_email varchar(32),  -- user email, which is the identifier for a user	
	first_name varchar(32),  
	last_name  varchar(32),  
	phone varchar(16),
	max_education integer,
	created varchar(32)         -- always a good idea to save a time stamp
	                            -- but Go seems to have trouble with sqlite datetime    
);

-- Job descriptions

CREATE TABLE IF NOT EXISTS job (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	created_by varchar(32),  -- email of creating user (employer)	
	title varchar(64),  
	description varchar(1024),  
	min_education integer,
	min_years_experience integer,
	salary integer,
	hired_person varchar(32) default '',    -- email of person hired
	is_open integer default 1,   -- 1 if still open, 0 if filled
	created varchar(32)             -- always a good idea to save a time stamp   
);

CREATE TABLE IF NOT EXISTS job_application (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	job_id int,              -- job ID with leading zeros
	user_email varchar(32),  -- user who has applied
	apply_time varchar(32),
	UNIQUE(job_id, user_email)  -- You can only apply once for a job
);

