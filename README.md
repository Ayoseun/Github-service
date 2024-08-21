# GitHub API Data Fetching and Service

 https://documenter.getpostman.com/view/17643992/2sA3sAhnpi

## Overview
This is a Go-based service that fetches data from the GitHub public API, stores the fetched data in a persistent database, and provides mechanisms to efficiently query the stored data.

### Objectives
The main objectives of this project are:

- Fetch data from the GitHub public API for a given repository, including commits and repository metadata.
- Store the fetched data in a persistent database (PostgreSQL).
- Continuously monitor the repository for changes and fetch new data at regular intervals.
- Provide efficient querying mechanisms to retrieve the top N commit authors and commits for a specific repository.
- Ensure code modularity, maintainability, and adherence to best practices.

### Features

GitHub API Data Fetching:

1. Fetch commit details (message, author, date, URL) for a given repository.
2. Fetch repository metadata (name, description, URL, language, forks count, stars count, open issues count, watchers count, created/updated dates).
3. Implement a mechanism to continuously monitor the repository for changes and fetch new data at regular intervals (e.g., every hour).
4. Avoid pulling the same commit twice and ensure the database mirrors the commits on GitHub.
5. Allow configuring the date to start pulling commits from.
6. Provide a mechanism to reset the data collection to start from a specific point in time.



## Getting Started

1. Prerequisites:

Go 1.22.6 or later
PostgreSQL database


2. Installation:

Clone the repository: 
```shell
git clone https://github.com/Ayoseun/Github-service.git
```
Navigate to the project directory: 
``` shell
cd github-service
```


3. Usage
The service provides the following endpoints:

- Fetch Repository Commits:
```shell
GET /:repo/fetch_commits
```
Fetches the commits for the given repository and saves them to the database.

- Fetch Repository Data:
```shell
GET /fetch_repository/:repo
```
Fetches the repository metadata for the given repository and saves it to the database.

- Get Top N Commit Authors:
```shell
GET /top_authors/:n
```
Retrieves the top N commit authors by commit count from the database.
Supports pagination using the page and limit query parameters.

- Retrieve Commits by Repository:
```shell
GET /commits/:repo
```
Retrieves the commits for the given repository from the database.
Supports pagination using the page and limit query parameters.

4. Continuous Monitoring and Data Fetching
The service is designed to continuously monitor the repository for changes and fetch new data at regular intervals (e.g., every hour). This is achieved by implementing a background task or a cron job that periodically calls the fetchRepositoryCommits and fetchRepositoryData functions.
To reset the data collection to start from a specific point in time, you can add a new endpoint or a command-line utility that allows the user to update the starting date for the commit fetching process.

5. Data Storage and Querying
The solution uses a PostgreSQL database to store the repository details and commit data. The database schema is designed to ensure efficient querying of the data.
The SavedCommit and Repository models are used to represent the commit and repository data, respectively. The GetTopNCommitAuthors and GetCommitsByRepository functions in the repository package provide the necessary functionality to retrieve the top N commit authors and commits for a specific repository.

6. Testing
The solution includes a test case
run
```shell
go test ./tests/...
```

