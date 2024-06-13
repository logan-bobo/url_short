# URL Short
A simple URL shortening service, designed to convert longer URLs into concise, unique keys. 
This service provides the capability to use the unique key to access the original URL seamlessly.

## Key Features
- URL Compression: The service accepts a URL and generates a unique key representing the corresponding URL.
- URL Redirection: Clients are can accessing the original link via the unique key provided to them.
- CRUD Operations: Clients can create, read, update and delete their own links.

## High Level Design
```mermaid
sequenceDiagram
    participant Client
    participant Server
    participant Destination Server
    Client ->> Server: visit short URL
    Server ->> Client: HTTP/301 permanent redirect to long URL
    Client ->> Destination Server: redirected to long URL
```

## Authentication Overview

Authentication is handled through the use of JSON Web Tokens (JWT).
Upon a valid login request to the `/api/v1/login` endpoint, the client receives:
- **Access Token**: Valid for 1 hour, used to access protected endpoints.
- **Refresh Token**: Valid for 60 days, used to obtain a new access token without requiring the user to log in again.

Clients use the access token to access endpoints that require authentication, such as 
`/api/v1/data/shorten`. When the access token expires, the client can obtain a new one from the 
`/api/v1/refresh` endpoint using the refresh token.

Finally, when the refresh token expires, the client must request a new set of tokens
(both access token and refresh token) by logging in again at the `/api/v1/login` endpoint.

The security considerations around the use of JWTs are:
- HTTPs should always be used as a transmission protocol between client and server
- Clients should look to securely store tokens for example using `HttpOnly` cookie (This would be communicated with the front end team).
- Access tokens have a short lifetime and refresh tokens can be revoked from the database. 
- The JWT signing secret my remain secure, I would look to store this in some secret storage platform such as 
Hashicorp Vault or AWS Secrets Manager.

```mermaid
sequenceDiagram
    participant Client
    participant Server

    Client->>Server: POST /api/v1/login (credentials)
    Server-->>Client: access token (1 hour) & refresh token (60 days)

    Client->>Server: POST /api/v1/data/shorten (access token)
    Server-->>Client: Data

    Note over Client: access token expires

    Client->>Server: POST /api/v1/refresh (refresh token)
    Server-->>Client: new access token (1 hour)

    Note over Client: refresh token expires
```
