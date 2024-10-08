# GitHub API Data Fetching and Service
This Go-based service fetches data from the GitHub public API, stores the fetched data in a PostgreSQL database, and provides mechanisms to efficiently query the stored data.

[Postman docs](https://documenter.getpostman.com/view/17643992/2sA3sAhnpi)

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
- Go 1.23 or later: Download [here](https://go.dev/doc/install)
- Docker: Ensure Docker is installed and running on your machine, you get docker [here](https://www.docker.com/products/docker-desktop/).

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

```sh
go test ./tests/...
```

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


### Service Endpoints:

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

- Retrieve a given repository metadata
```sh
GET /repositories/:repo/commits
```

- Example

```sh
http://localhost:8080/repositories/chromium/fetch
```

- Parameters:
repo : The name of the repository (e.g., chromium).

- Response:
```json

{
    "ID": 120360765,
    "CreatedAt": "0001-01-01T00:00:00Z",
    "UpdatedAt": "0001-01-01T00:00:00Z",
    "DeletedAt": null,
    "name": "chromium",
    "description": "The official GitHub mirror of the Chromium source",
    "html_url": "https://github.com/chromium/chromium",
    "language": "C++",
    "forks_count": 6884,
    "stargazers_count": 18642,
    "open_issues_count": 93,
    "watchers_count": 18642,
    "subscribers_count": 562,
    "created_at": "2018-02-05T20:55:32Z",
    "updated_at": "2024-09-03T17:41:53Z"
}

```
Add new repository to monitor.

```sh
GET /repositories/monitor/:owner
```
Example URL:

```c
http://localhost:8080/repositories/monitor/golang?repo=go&start_date=2023-01-01T00:00:00Z
```
- Query Parameters:

owner (required): The owner of the repo
repo : The repository to add.
start_date : The defined N history to begin pulling from.

- Response:
```json
{
    "statusCode":200,
    "message": "Repository added successfully"
}

```

Remove a repository from monitor service and delete it from DB.

```sh
DELETE /repositories/monitor/:owner
```
Example URL:

```c
http://localhost:8080/repositories/monitor?owner=chromium&repo=chromium
```
- Query Parameters:

owner (required): The owner of the repo
repo : The repository to add.

- Response:
```json
{
    "statusCode":200,
  "message": "Repository removed successfully"
}
```

Reset a repository commit collection.

```sh
GET /repositories/reset/:owner
```
Example URL:

```c
http://localhost:8080/repositories/reset/golang?repo=go
```
- Query Parameters:

owner (required): The owner of the repo
repo : The repository to add.


- Response:
```json
{
    "statusCode":200,
    "message": "Repository commits removed successfully"
}

```

5. Continuous Monitoring and Data Fetching
The service is designed to continuously monitor the repository for changes and fetch new data at regular intervals (e.g., every hour). This is achieved by implementing a background task or a cron job that periodically calls the fetchRepositoryCommits and fetchRepositoryData functions.

6. Data Storage and Querying
The solution uses a PostgreSQL database to store repository details and commit data. The database schema is designed for efficient querying.

- Models:

* Commit: Represent raw commits from reposiory on github
* SavedCommit: Represents commit data.
* PaginatedResponse: Represents he commits response in a paginated relay
* Repository: Represents repository metadata.
* TopAuthorsCount: Represents top N authors

- Functions Interface:
Commit functions
type CommitRepository interface {
	SaveCommit(commit *models.SavedCommit) error
	GetCommits(repositoryURL string, page, limit int) ([]models.SavedCommit, error)
	GetTotalCommits(repositoryURL string) (int64, error)
}

Repository functions
type RepositoryRepository interface {
	SaveRepository(repo *models.Repository) error
	GetTopNCommitAuthors(page, limit int) (models.TopAuthorsCount, error)
	GetRepositoryByName(repository string) (models.Repository, error)
}

Repository Monitoring interface
type CommitServiceInterface interface {
	FetchAndSaveCommits(owner, repo string, since time.Time) ([]models.Commit, error)
	FetchCommitsInRange(owner, repo string, from, to time.Time) ([]models.Commit, error)
}

7. Troubleshooting
# WARNING- The database credentials are for testing only do not use in production.

### Common Issues:
Error: Database connection failed:
Ensure your .env file is correctly configured with the right database URL.
Verify that your PostgreSQL service is running.

