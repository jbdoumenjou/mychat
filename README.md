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

## Register a User - POST /register

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

## Send a Message - POST /messages

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

## List Chats for a User - GET /chats?phoneNumber={phoneNumber}

List the chats for the provided phone number.

For the sake of simplicity, and security,
if the phone number is not provided, the API will return an empty list.
The number is passed as a query parameter but should be retrieved from the user token,
or in a way that does not expose it.

```bash
curl "http://localhost:8080/chats?phoneNumber=%2B3306666666"
```

| Status Code                 | 	Description                                             |
|-----------------------------|----------------------------------------------------------|
| 200 (ok)                    | return the list successfully.                            | 
| 400 (Bad Request)           | Invalid input (e.g., missing/invalid fields).            |
| 500 (Internal Server Error) | A server-side error occurs while processing the request. |

## List Messages for a Chat - GET /chats/{chat_id}/messages

Retrieve the messages for the specified chat.

For simplicity and clarity,
this endpoint uses a hierarchical structure to associate messages with a specific chat.
This approach avoids ambiguity and focuses on retrieving messages belonging to a single chat.

```bash
curl "http://localhost:8080/chats/3163f560-f246-4e68-8551-cb702f8a017a/messages"
```

| Status Code                 | 	Description                                             |
|-----------------------------|----------------------------------------------------------|
| 200 (ok)                    | return the list successfully.                            | 
| 400 (Bad Request)           | Invalid input (e.g., missing/invalid fields).            |
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

The first step was to create the project structure and CI. As I used one of my existing project as a template, this step took around 30 min.

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

## Step 4: Create the API - list the chats for a user

Now that we can send messages, we need to list the chats for a user.
For this step, I added 
* a chat repository that is in charge of managing the chats.
* a chat handler that is in charge of managing the chat API.

I keep let the phoneNumber in the query parameter, but it should be in the token, or in a way that does not expose it.
We should use a pagination system to avoid returning all the chats at once.

## Step 5: Create the API - list the messages for a chats

I considered both `/chats/{chat_id}/messages` and `/messages?chat={chat_id}`.
I chose `/chats/{chat_id}/messages` for its simplicity and clarity.
This approach avoids confusion about which query parameters are available.
While the query-based path /messages?chat={chat_id} is more suited for searches or complex filtering,
the hierarchical path `/chats/{chat_id}/messages` aligns better with the current lifecycle of a chat app,
where the primary goal is to retrieve messages belonging to a specific chat.

For this step, I keep the same line of thought as the previous steps, and the in memory database.
We should have more tests, especially unit test on repositories to migrate to a real database.