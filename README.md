# GitHub API Data Fetching and Service
This Go-based service fetches data from the GitHub public API, stores the fetched data in a PostgreSQL database, and provides mechanisms to efficiently query the stored data.

## Overview
**Objectives**
- Fetch data from the GitHub public API for a given repository, including commits and repository metadata.
- Store the fetched data in a PostgreSQL database.
- Continuously monitor the repository for changes and fetch new data at regular intervals.
- Provide efficient querying mechanisms to retrieve the top N commit authors and commits for a specific repository.
- Ensure code modularity, maintainability, and adherence to best practices.
**Features**
- GitHub API Data Fetching:
1. Fetch commit details (message, author, date, URL) for a given repository.
2. Fetch repository metadata (name, description, URL, language, forks count, stars count, open issues count, watchers count, created/updated dates).
3. Implement a mechanism to continuously monitor the repository for changes and fetch new data at regular intervals (e.g., every hour).
4. Avoid pulling the same commit twice and ensure the database mirrors the commits on GitHub.
5. Allow configuring the date to start pulling commits from.
6. Provide a mechanism to reset the data collection to start from a specific point in time.
## Getting Started
1. Prerequisites
- Go 1.22.6 or later: Download here
- PostgreSQL: Make sure PostgreSQL is installed and running on your machine.
1a. Setting Up PostgreSQL on macOS
- Install PostgreSQL:
```shell
brew install postgresql
```
-  Start the PostgreSQL service:
```shell
brew services start postgresql
```
-  Access PostgreSQL:
```shell
psql postgres
```
-  Create the database:
```shell
CREATE DATABASE github;
```
2. Installation
- Clone the Repository:
```shell
git clone https://github.com/Ayoseun/Github-service.git
```
- Navigate to the Project Directory:
```shell
cd github-service
```
- Install Dependencies:
```shell
go mod tidy
```
- Set Up Environment Variables:
Create a `.env` file in the project root directory:
```shell
touch .env
```
Next add your database connection URLs:
```shell
DATABASE_DEV_URL=postgres://username:password@localhost:5432/github
DATABASE_PROD_URL=<Your-Cloud-POstgres-URL>
```
Ensure the connection string format matches your database configuration.

3. Testing
The project includes test cases to verify functionality.

Run Tests:
```shell
go test ./tests/...
```
4. Usage

- Start the app by running the below command in the app directory in your terminal or command line interface :
```go
go run cmd/main.go
```
- The service provides the following endpoints:

**Retrieves the top N commit authors by commit count from the database.**
```shell
GET /top_authors/:n
```
Example url
```shell
http://localhost:8080/top_authors/5?page=1&limit=3
```
Parameters:
- n (required): The number of top authors to retrieve.
- page : The page number for pagination.
- limit : The number of records per page.
Response:
```json
[
    {
        "author": "chromium-autoroll",
        "count": 9
    },
    {
        "author": "chromium-internal-autoroll",
        "count": 3
    },
    {
        "author": "Lingqi Chi",
        "count": 2
    }
]
```
**Retrieves the commits for the given repository from the database.**
```shell
GET /commits/:repo
```
Example
```shell
http://localhost:8080/commits/chromium?page=1&limit=20
```
Parameters:
- repo : The name of the repository (e.g., chromium).
- page : The page number for pagination.
- limit : The number of records per page.
Response:
```json
{
    "current_page": 1,
    "total_pages": 2,
    "commits": [
        {
            "ID": 1,
            "CreatedAt": "2024-08-21T03:10:56.405498+01:00",
            "UpdatedAt": "2024-08-21T03:10:56.405498+01:00",
            "DeletedAt": null,
            "message": "Add APU backend support\n\nWe want to support running models with different accelerator backends.\nThis adds the backend type parameter in ChromeMLModelDescriptor, which\nincludes the original GPU backend, and a new APU backend. This also add\nsome parameters that are used by llm_engine for other type of backends\nlike APU (e.g., model_path, sentencepiece_mode_path).\n\nBug: b:351276861\nTest: CQ\nChange-Id: I0080f3b7e37c065b1545d764312079f869e0776e\nReviewed-on: https://chromium-review.googlesource.com/c/chromium/src/+/5782475\nReviewed-by: Clark DuVall <cduvall@chromium.org>\nReviewed-by: Yi Chou <yich@google.com>\nCommit-Queue: Howard Yang <hcyang@google.com>\nCr-Commit-Position: refs/heads/main@{#1344548}",
            "author": "Howard Yang",
            "date": "2024-08-21T02:56:27+01:00",
            "url": "https://github.com/chromium/chromium/commit/cb57d73200f18b50f218b2a6117fc4266b3d5e10"
        },
    ]
}

```

Run the project postman docs [here](https://documenter.getpostman.com/view/17643992/2sA3sAhnpi)

5. Continuous Monitoring and Data Fetching
The service is designed to continuously monitor the repository for changes and fetch new data at regular intervals (e.g., every hour). This is achieved by implementing a background task or a cron job that periodically calls the fetchRepositoryCommits and fetchRepositoryData functions.

6. Data Storage and Querying
The solution uses a PostgreSQL database to store repository details and commit data. The database schema is designed for efficient querying.

- Models:

`SavedCommit: Represents commit data.`
`Repository: Represents repository metadata.`
- Functions:

`GetTopNCommitAuthors: Retrieves the top N commit authors.`
`GetCommitsByRepository: Retrieves commits for a specific repository.`
7. Troubleshooting
Common Issues:

Error: Database connection failed
Ensure your .env file is correctly configured with the right database URL.
Verify that your PostgreSQL service is running.