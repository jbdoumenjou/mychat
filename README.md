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

# CI/CD

This project use Github Actions to run the CI/CD pipeline. The pipeline is defined in the `.github/workflows` folder.
There are several workflows defined in this project:
* test.yml: Run the tests
* validate.yml: Run the linters
* release.yml: Create a new release (changelog, tag, and release notes)
    * :warning: This GH Action use [conventional commit](https://www.conventionalcommits.org/en/v1.0.0/) to generate the changelog.

# Project Development Journey

* The firs step was to create the project structure and CI. As I used one of my existing project as a template, this step took around 30 min.