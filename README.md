# Overview

This project is a sandbox for experimentation and learning,
designed to simulate the backend of a real-time messaging application.
It provides a simple framework to explore backend development concepts such as API design,
database interactions,
and system architecture.

The goal is to create a functional and flexible backend
that allows users to send and receive messages in private 1:1 chats while experimenting with new ideas and technologies.

# How to run

To check the available target from the Makefile, run the following command:

```bash
make help
```

You can run the project locally with debug logs by running the following command:

```bash
LOG_LEVEL=DEBUG go run ./cmd/main.go
```

# API

A [Bruno](https://www.usebruno.com/) collection is available in the `docs` folder.

## POST /register

Register a new user.

```bash
curl -X POST http://localhost:8080/register \
-H "Content-Type: application/json" \
-d '{"phoneNumber": "+1234567890"}'
```

| Status Code       | 	Description                                  |
|-------------------|-----------------------------------------------|
| 201 (Created)     | User registered successfully.                 | 
| 400 (Bad Request) | Invalid input (e.g., missing/invalid fields). |
| 409 (Conflict)    | 	Phone number already registered.             |

# CI/CD

This project use Github Actions to run the CI/CD pipeline. The pipeline is defined in the `.github/workflows` folder.
There are several workflows defined in this project:
* test.yml: Run the tests
* validate.yml: Run the linters
* release.yml: Create a new release (changelog, tag, and release notes)
    * :warning: This GH Action use [conventional commit](https://www.conventionalcommits.org/en/v1.0.0/) to generate the changelog.

# Project Development Journey

* The firs step was to create the project structure and CI. As I used one of my existing project as a template, this step took around 30 min.
* Then I started to implement the API. I started with the user API, as it is the most basic one.
  * I started to add a repository layer with a mongoDB implementation, but a realized that it was too much for this project at this step. So I decided to use a simple in-memory database.
  * I added some basic unit test to show how to test the API. I tried to not add too much tests, as the goal is to show how to test the API, not to test everything.
  * I added a bruno collection to show how to use it.
  * I added logs to have a better understanding of what is happening.