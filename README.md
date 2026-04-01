# Continuous Delivery in Agile Software Development -- Exercises

This repository contains four progressive exercises for the Master course **Continuous Delivery in Agile Software Development**.

## Overview

| Exercise | Topic | Branch |
|----------|-------|--------|
| 1 | Git Basics: PRs, Interactive Rebase, Unit Tests | `exercise/01-git-basics` |
| 2 | Microservice Architecture, Docker & GitHub Actions | `exercise/02-microservice-docker` |
| 3 | CI Pipeline: SonarCloud, Matrix Builds, Linting | `exercise/03-ci-pipeline` |
| 4 | Vulnerability Scanning & Kubernetes Deployment | `exercise/04-security-k8s` |

## Technology Stack

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
- **Language:** Go 1.22+
=======
- **Language:** Go 1.24+
>>>>>>> 7333620 (docs: update Go prerequisite to 1.24+)
=======
- **Language:** Go 1.26+
>>>>>>> 55ded5b (docs: update Go version references in README (1.22->1.26))
=======
- **Language:** Go 1.25+
>>>>>>> cb3b318 (docs: update Go prerequisite to 1.25+ in README)
=======
- **Language:** Go 1.25+
>>>>>>> cb3b318 (docs: update Go prerequisite to 1.25+ in README)
- **Web Framework:** Gorilla Mux
- **Database:** PostgreSQL
- **Containerization:** Docker & Docker Compose
- **CI/CD:** GitHub Actions
- **Code Quality:** SonarCloud, golangci-lint
- **Security:** Trivy, Snyk
- **Deployment:** Kubernetes (Minikube)

## Project: Product Catalog API

Throughout the exercises, you will build and evolve a RESTful Product Catalog API with CRUD operations, backed by PostgreSQL.

## Prerequisites

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
- Go 1.22+ installed
=======
- Go 1.24+ installed
>>>>>>> 7333620 (docs: update Go prerequisite to 1.24+)
=======
- Go 1.26+ installed
>>>>>>> 55ded5b (docs: update Go version references in README (1.22->1.26))
=======
- Go 1.25+ installed
>>>>>>> cb3b318 (docs: update Go prerequisite to 1.25+ in README)
=======
- Go 1.25+ installed
>>>>>>> cb3b318 (docs: update Go prerequisite to 1.25+ in README)
- Git 2.30+
- GitHub Account
- Docker Desktop (from Exercise 2)
- Minikube (Exercise 4)

## Getting Started

1. **Fork** this repository on GitHub (click the "Fork" button in the top right corner). All exercise branches will be included in your fork.
2. **Clone** your fork:

```bash
git clone https://github.com/<your-username>/CI-CD-MCM.git
cd CI-CD-MCM
```

3. Switch to the respective exercise branch:

```bash
git checkout exercise/01-git-basics
```

> **Important:** Do not clone the original repository directly — always work on your own fork so you can push changes and create Pull Requests.

Each exercise branch contains a detailed `README.md` with instructions.

## Author
- FH-Prof. Dr. Marc Kurz (marc.kurz@fh-hagenberg.at)
