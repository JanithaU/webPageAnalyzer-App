# Web Page Analyzer App

This is a Go-based web page analyzer that allows you to analyze web pages for various details such as HTML version, heading structure, internal and external links, and more.

## Prerequisites

Before you begin, ensure that you have the following installed on your local machine:

- **[Go Programming Language](https://golang.org/dl/)** (version 1.16 or later)
- **[Git](https://git-scm.com/)** (to clone the repository)

## Steps to Run the Application

Follow these steps to set up and run the application locally:

### 1. Install Go

If you don't have Go installed, download and install it by following the instructions for your platform:

- Visit the Go download page: [https://golang.org/dl/](https://golang.org/dl/)
- Follow the installation instructions provided for your operating system.

Once Go is installed, verify the installation by running the following command in your terminal:

```bash
go version
```


### 2. Clone the Repository
Clone this repository to your local machine using Git. Run the following command in your terminal:

```bash
git clone https://github.com/JanithaU/webPageAnalyzer-App.git
```

### 3. Navigate to the Project Directory
``` bash
cd webPageAnalyzer-App
```
### 4. Download Dependencies
``` bash
go mod tidy
```


### 5. Run the Application
```bash
go run cmd/web/main.go
```

## OR to run application as a docker
#### build docker 
```
docker build -t go-web-analyzer .
```

#### run docker image
 ``` 
 docker run -p 8080:8080 go-web-analyzer
 ```

## OR build and run via the Makefile
#### build and run 
```
make run
```

#### build only (with dependancies)
```
make build
```
