# Chirpy

Chirpy is a robust social networking application backend built in Go. It allows users
to create, retrieve, and delete short messages, known as "chirps".

## Features

- **User registration and authentication**: Users can register and authenticate
  themselves using their credentials.
- **Create, retrieve, and delete chirps**: Users can create their own chirps, view
  all chirps, and delete their chirps.
- **User profile updates**: Users can update their profile information.
- **Token-based authentication with JWT**: The application uses JSON Web Tokens
  (JWT) for secure, token-based user authentication.
- **Refresh and Access Tokens**: The application uses a combination of refresh and
  access tokens to maintain user sessions and re-authenticate users without
  asking for their credentials.
- **Middleware for metrics and logging**: The application includes middleware for
  capturing metrics and logging information about requests and responses.

## Getting Started

### Prerequisites

- Go 1.16 or later
- A text editor or IDE (like Neovim)

### Installation

Follow the steps below to get the application up and running on your local
machine:

1. Clone the repository:

   ``` git clone https://github.com/Chaitanya-Shahare/chirpy.git ```

2. Navigate to the project directory:

   ``` cd chirpy ```

3. Install the dependencies:

   ``` go mod download ```

4. Copy the `.env.example` file and rename it to `.env`, then replace the
placeholder values with your actual data:

   ``` cp .env.example .env ```

5. Run the application:

   ``` go run main.go ```

## Usage

Once the server is running, you can interact with the API using a tool like
`curl` or Postman. The available endpoints are grouped as follows:

### User Endpoints

- `POST /api/users`: Register a new user. The request body should include
  `username`, `email`, and `password`.
- `POST /api/login`: Authenticate a user and receive a JWT. The request body
  should include `email` and `password`.
- `PUT /api/users`: Update a user's profile. This endpoint requires
  authorization. The request body should include any of the fields that need to
  be updated.

### Chirp Endpoints

- `POST /api/chirps`: Create a new chirp. This endpoint requires authorization.
  The request body should include `content`.
- `GET /api/chirps`: Retrieve all chirps. This endpoint does not require
  authorization.
- `GET /api/chirps/{chirpID}`: Retrieve a specific chirp. This endpoint does
  not require authorization.
- `DELETE /api/chirps/{chirpID}`: Delete a specific chirp. This endpoint
  requires authorization.

### Token Endpoints

- `POST /api/refresh`: Refresh a user's JWT. The request body should include
  `refreshToken`.
- `POST /api/revoke`: Revoke a user's refresh token. This endpoint requires
  authorization.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License.

