# Snippet box 
# Snippetbox

Snippetbox is a web application written in Go that allows users to create, share, and manage code snippets. It's built with modern web development practices as taught in Alex Edwards' book "Let's Go". This application demonstrates the use of Go for building robust and efficient web applications.

## Features

- **User Authentication**: Sign up, log in, and manage user sessions.
- **CRUD Operations**: Create, read, update, and delete snippets.
- **Database Integration**: Store and retrieve snippets from a MySQL or PostgreSQL database.
- **Template Rendering**: Dynamically generate HTML pages using Go's template package.
- **Security Best Practices**: Protect against XSS, CSRF, and SQL injection attacks.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

Ensure you have the following installed:

- Go (version 1.13 or later)
- MySQL or PostgreSQL

### Installing

Clone the repository to your local machine:

```bash
git clone https://github.com/yourusername/snippetbox.git
cd snippetbox
```

Install the necessary Go packages:
```bash
go mod tidy
```

### Running the application
Start the server:
```bash
go run ./cmd/web
````
The application will be available at http://localhost:8080

### Build with
* Go - The Go Programming Language
* MySQL - Database
* Bootstrap - Front-end framework

