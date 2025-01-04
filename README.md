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

| Status Code                 | 	Description                                             |
|-----------------------------|----------------------------------------------------------|
| 201 (Created)               | User registered successfully.                            | 
| 400 (Bad Request)           | Invalid input (e.g., missing/invalid fields).            |
| 409 (Conflict)              | Phone number already registered.                         |
| 500 (Internal Server Error) | A server-side error occurs while processing the request. |

## POST /messages

send a new message.

```bash
curl -X POST http://localhost:8080/message \
-H "Content-Type: application/json" \
-d '{"sender": "+1234567890", "receiver": "+0987654321", "content": "Hello, World!"}'
```

| Status Code                 | 	Description                                             |
|-----------------------------|----------------------------------------------------------|
| 201 (Created)               | User registered successfully.                            | 
| 400 (Bad Request)           | Invalid input (e.g., missing/invalid fields).            |
| 404 (Not Found)             | The sender or recipient phone number does not exist.     |
| 500 (Internal Server Error) | A server-side error occurs while processing the request. |


# CI/CD

This project use Github Actions to run the CI/CD pipeline. The pipeline is defined in the `.github/workflows` folder.
There are several workflows defined in this project:
* test.yml: Run the tests
* validate.yml: Run the linters
* release.yml: Create a new release (changelog, tag, and release notes)
    * :warning: This GH Action use [conventional commit](https://www.conventionalcommits.org/en/v1.0.0/) to generate the changelog.

# Project Development Journey

## Step 1: Create the project structure and CI

The firs step was to create the project structure and CI. As I used one of my existing project as a template, this step took around 30 min.

## Step 2: Create the API - register a user

* Then I started to implement the API. I started with the user API, as it is the most basic one.
  * I started to add a repository layer with a mongoDB implementation, but a realized that it was too much for this project at this step. So I decided to use a simple in-memory database.
  * I added some basic unit test to show how to test the API. I tried to not add too much tests, as the goal is to show how to test the API, not to test everything.
  * I added a bruno collection to show how to use it.
  * I added logs to have a better understanding of what is happening.

## Step 3: Create the API - send a message

I added the message API. I reused the same in-memory database to store the messages.

I keep the same line of thought as the user API, adding some basic unit tests to show how to test the API.

I face issue thinking about how to store the messages.
To keep it simple a decided to have a chat repository that is in charge of managing chat, which are discussions between two users.
Then, the chat ID is used to store the messages. In that way, we don't have to think in sender/receiver, just in chat.
This step is very naive too, but should be enough for the next step, which is listing the messages.