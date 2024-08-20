# Golang-Assessment

## GitHub API Data Fetching and Service
This is a Go-based service that fetches data from the GitHub public API, stores the fetched data in a persistent database, and provides mechanisms to efficiently query the stored data.
1. Objective
The main objective of this project is to:

- Fetch data from the GitHub public API for a given repository, including commits and repository metadata.
- Store the fetched data in a persistent database (PostgreSQL in this case).
Continuously monitor the repository for changes and fetch new data at regular intervals.
- Provide efficient querying mechanisms to retrieve the top N commit authors and commits for a specific repository.
- Ensure code modularity, maintainability, and adherence to best practices.

2. Features

* GitHub API Data Fetching:

- Fetch commit details (message, author, date, URL) for a given repository.
- Fetch repository metadata (name, description, URL, language, forks count, stars count, open issues count, watchers count, created/updated dates).
- Implement a mechanism to continuously monitor the repository for changes and fetch new data at regular intervals (e.g., every hour).
- Avoid pulling the same commit twice and ensure the database mirrors the commits on GitHub.
- Allow configuring the date to start pulling commits from.
- Provide a mechanism to reset the data collection to start from a specific point in time.


3. Data Storage:

Design and create necessary tables to store repository details and commit data.
Ensure efficient querying of data using appropriate database indexes and schema design.


4. Querying Capabilities:

Provide a functionality to get the top N commit authors by commit count from the database.
Implement a way to retrieve commits of a repository by the repository name from the database.


5. Solution Overview
The solution is structured using a modular approach, with the following key components:

- cmd/main.go: This is the entry point of the application, responsible for setting up the Gin web framework, registering the routes, and starting the server.
- internal/: This directory contains the core functionality of the application.

- config/config.go: Handles the application configuration, such as the database connection details.
- database/database.go: Manages the database connection and migrations.
- handlers/: Contains the HTTP request handlers for the various endpoints.
- models/: Defines the data models for the application, such as SavedCommit and Repository.
- repository/: Implements the repository-layer logic for interacting with the database.


- pkg/github/github.go: This package contains the logic for interacting with the GitHub API to fetch the repository commits and metadata.
- README.md: The comprehensive documentation for the solution, covering all the requirements and guidelines.
- Unit Tests: The solution includes at least one unit test for a core function of the service.

### Getting Started

- Prerequisites:

1. Go 1.22.6 or later
2. PostgreSQL database


- Installation:

Clone the repository:
Copygit clone https://github.com/your-github-username/github-api.git

Navigate to the project directory:
Copycd github-api

Update the database connection details in the internal/config/config.go file.
Build and run the application:
Copygo build -o github-api ./cmd
./github-api




Usage
The service provides the following endpoints:

Fetch Repository Commits: GET /:repo/fetch_commits

Fetches the commits for the given repository and saves them to the database.


Fetch Repository Data: GET /fetch_repository/:repo

Fetches the repository metadata for the given repository and saves it to the database.


Get Top N Commit Authors: GET /top_authors/:n

Retrieves the top N commit authors by commit count from the database.
Supports pagination using the page and limit query parameters.


Retrieve Commits by Repository: GET /commits/:repo

Retrieves the commits for the given repository from the database.
Supports pagination using the page and limit query parameters.



Continuous Monitoring and Data Fetching
The service is designed to continuously monitor the repository for changes and fetch new data at regular intervals (e.g., every hour). This is achieved by implementing a background task or a cron job that periodically calls the fetchRepositoryCommits and fetchRepositoryData functions.
To reset the data collection to start from a specific point in time, you can add a new endpoint or a command-line utility that allows the user to update the starting date for the commit fetching process.
Data Storage and Querying
The solution uses a PostgreSQL database to store the repository details and commit data. The database schema is designed to ensure efficient querying of the data.
The SavedCommit and Repository models are used to represent the commit and repository data, respectively. The GetTopNCommitAuthors and GetCommitsByRepository functions in the repository package provide the necessary functionality to retrieve the top N commit authors and commits for a specific repository.
Unit Tests
The solution includes at least one unit test for a core function of the service, such as the GetTopNCommitAuthors function in the repository package.
Conclusion
This Go-based service provides a comprehensive solution for fetching data from the GitHub public API, storing the data in a persistent database, and offering efficient querying mechanisms. The modular design, adherence to best practices, and inclusion of unit tests ensure the maintainability and scalability of the application.