# User Auth API

This repository contains the code for the User Authentication API, built using Golang and Gin. The API serves as a backend service for basic user authentication.

## Start service for testing
To start testing the service, make sure you have Docker and docker-compose installed. Then, start the MySQL container and API service by running the following command:

```bash
docker-compose up -d
```

Once the services are up and running, the API service will be available at http://0.0.0.0:8080. Additionally, the MySQL instance will be accessible at the host `0.0.0.0` with port `3306`, and the `root` user's password will be `example`

Please note that the current docker-compose.yaml configuration is not intended for production deployment.
## Key Opponents

### Auth Middleware
AuthMiddleware validates user credentials using basic authentication. If the credentials are valid, the middleware sets a userKey in the context for future requests.

The authentication implementation in this project uses basic authentication, which is a simple mechanism for sending user credentials over the network. When a client sends a request to a protected endpoint, the AuthMiddleware function extracts the Authorization header from the request and removes the "Basic " prefix. It then decodes the remaining header value from base64 to retrieve the username and password.

The middleware then checks whether the provided credentials match a user record in the database. If the user is found and the password matches the stored hash, the middleware sets a userKey in the context for future requests. If the credentials are invalid, the middleware returns an unauthorized error response to the client.

### Error Handling
This project includes error handling for common error scenarios, such as bad requests, unauthorized access, and server errors. Appropriate error responses are returned to the client along with an error message.


## API Endpoints
### POST /users

Create a new user with the provided details.

Request

```json
{
    "user_id": "john_doe",
    "password": "password"
}
```

Response
```json
{
    "message": "Account successfully created!",
    "user": {
        "user_id": "john_doe",
        "nickname": "John"
    }
}
```

### GET /users/:user_id
Retrieve user details for a given user_id.

Response


```json
{
    "message": "User details by user_id",
    "user": {
        "user_id": "john_doe",
        "nickname": "John",
        "comment": "A retail enthusiast"
    }
}
```

###  PATCH /users/:user_id
Update user details for a given user_id.

Request
```json
{
    "nickname": "Johnny",
    "comment": "A retail enthusiast and programmer"
}
```

Response
```json
{
    "message": "User successfully updated!",
    "user": {
        "user_id": "john_doe",
        "nickname": "Johnny",
        "comment": "A retail enthusiast and programmer"
    }
}
```

### DELETE /users/:user_id

Delete a user with a given user_id.

Response

```json
{
    "message": "Account and user successfully removed!"
}
```
