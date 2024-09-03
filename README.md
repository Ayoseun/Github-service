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
**Getting Started**

1. Prerequisites
- Go 1.22.6 or later: Download here
- Docker: Ensure Docker is installed and running on your machine.

2. Installation
Clone the Repository:

```sh
git clone https://github.com/Ayoseun/Github-service.git
```
Navigate to the Project Directory:

```sh
cd github-service
```
Install Dependencies:
```sh
go mod tidy
```
3. Run the test

The test covers
- Commit test:
   - validates case for saving a single comit
   - validates case for saving a multiple comit

- GetTopNCommitAuthors test:
  - It tests different pagination scenarios (first page, second page, and all authors on a single page).
  - It ensures that the pagination works correctly by checking the authors on different pages.
  - It checks that the correct number of authors is returned for each page.

3. Docker Setup
- Build and Run Docker Containers:

```sh
make up
```
- Stop the Running Containers:

```sh
make down
```
- Restart Docker Containers:

```sh
make restart
```
- Clean Up Containers and Volumes:

```sh
make clean
```
- Tail Logs from the App Service:

```sh
make logs
```
4. Usage
Start the App: The app will start automatically when running make up. However, you can also start it manually by running the below command in the app directory:

```sh
go run cmd/main.go
```
***note that this will require you to already have a running postgres server***

Service Endpoints:

Retrieves the top N commit authors by commit count from the database.

```sh
GET /repositories/:repo/top-authors
```
Example URL:

```c
http://localhost:8080/repositories/chromium/top-authors/5?page=1&limit=2
```
- Parameters:

n (required): The number of top authors to retrieve.
page : The page number for pagination.
limit : The number of records per page.

***Note : the idea of seing N along side page and limit is to cover edge cases where here is a need for a large number for N.***
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
- Retrieves the commits for the given repository from the database.

```sh
GET /repositories/:repo/commits
```
- Example URL:

```sh
http://localhost:8080/repositories/chromium/commits?page=1&limit=20

```

- Parameters:

repo : The name of the repository (e.g., chromium).
page : The page number for pagination.
limit : The number of records per page.
- Response:

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
            "message": "Add APU backend support...",
            "author": "Howard Yang",
            "date": "2024-08-21T02:56:27+01:00",
            "url": "https://github.com/chromium/chromium/commit/cb57d73200f18b50f218b2a6117fc4266b3d5e10"
        }
    ]
}
```
5. Continuous Monitoring and Data Fetching
The service is designed to continuously monitor the repository for changes and fetch new data at regular intervals (e.g., every hour). This is achieved by implementing a background task or a cron job that periodically calls the fetchRepositoryCommits and fetchRepositoryData functions.

6. Data Storage and Querying
The solution uses a PostgreSQL database to store repository details and commit data. The database schema is designed for efficient querying.

- Models:

* SavedCommit: Represents commit data.
* Repository: Represents repository metadata.
- Functions:

GetTopNCommitAuthors: Retrieves the top N commit authors.
GetCommitsByRepository: Retrieves commits for a specific repository.

7. Troubleshooting
# WARNING- The database credentials are for testing only do not use in production
- Common Issues:

Error: Database connection failed:
Ensure your .env file is correctly configured with the right database URL.
Verify that your PostgreSQL service is running.

