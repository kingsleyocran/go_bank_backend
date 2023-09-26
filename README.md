# Bank Backend with Go Lang (Gin-Gonic)

Welcome to the Bank Backend built with Go Lang using the Gin-Gonic framework and a PostgreSQL database. This backend serves as the foundation for a banking application, providing essential functionality for managing accounts, transactions, and user data. This README provides a guide on setting up and using the backend.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [API Documentation](#api-documentation)
- [Project Structure](#project-structure)
- [Contributing](#contributing)
- [License](#license)

## Features

- **User Management**: Create, update, and delete user accounts.
- **Account Management**: Add, view, and manage bank accounts for users.
- **Transaction Handling**: Record and manage transactions between accounts.
- **Security**: Protect sensitive data with encryption and authentication.

## Prerequisites

Before setting up the Bank Backend, ensure you have the following prerequisites installed:

- [Go Lang](https://golang.org/dl/): Go programming language.
- [PostgreSQL](https://www.postgresql.org/download/): PostgreSQL database server.
- [Git](https://git-scm.com/): Version control tool.

## Installation

1. Clone this repository:

   ```bash
   git clone https://github.com/kingsleyocran/go-bank-backend.git
   cd bank-backend
   ```

2. Create a PostgreSQL database for the application. You can use the PostgreSQL command line or a GUI tool like pgAdmin to create a database and user with appropriate privileges.

3. Set up environment variables. Create a `.env` file in the root directory with the following content:

   ```env
DB_DRIVER=postgres
DB_SOURCE=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable
SERVER_ADDRESS=0.0.0.0:8080
   ```

   Replace the placeholders with your database and secret key information.

4. Install required Go packages:

   ```bash
   go mod tidy
   ```

## Usage

1. Start the application:

   ```bash
   go run main.go
   ```

   The server should now be running at `http://localhost:8080`.

2. Access the API endpoints using an API client like [Postman](https://www.postman.com/) or [curl](https://curl.se/). Refer to the API documentation for available endpoints and request formats.

## API Documentation

The Bank Backend provides a comprehensive API for managing user accounts, bank accounts, and transactions. You can access the API documentation by navigating to `http://localhost:8080/swagger/index.html` when the application is running.

## Project Structure

- `main.go`: Application entry point.
- `db`: Database interaction logic, models and schemas.
- `api`: API routes and middleware.
- `utils`: Utility functions and helper methods.

## Contributing

Contributions to this project are welcome. Feel free to open issues and pull requests to help improve the Bank Backend.

## License

This Bank Backend is open-source software released under the [MIT License](LICENSE).
