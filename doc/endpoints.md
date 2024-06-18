## API Endpoints

### `GET /api/v1/healthz` 
Description: The health endpoint for the API used for health checks.

Response:
`200 OK`: The server is healthy and ready to respond to requests.

### `POST /api/v1/data/shorten` 
Description: Used to turn a long URL into a short URL.

Request:
```
{
    "long_url":"https://www.google.com/my/long/path"
}
```

Response:
```
{
    "short_url":"<short url hash>"
}
```

Parameters:
- Headers
    - `Authorization: Bearer <token>`

### `GET /api/v1/{shortUrl}`
Description: Redirects an unauthenticated client from the short URL to the long URL.

Parameters: 
- Path 
    - `shortUrl` a reference to a short URL in that is stored in the database.     

- `DELTE /api/v1/{shortUrl}` 

Description: An authenticated endpoint that will delete a short URL a user owns.

Parameters:
- Path
    - `shortUrl` a reference to a short URL that is stored in the database.
- Headers
    - `Authorization: Bearer <token>`

### `PUT /api/v1/{shortUrl}`
Description: Allows for the updating of a long URL based on a short URL

Parameters:
- Path 
    - `shortUrl` a reference to a short URL in the database
- Headers
    - `Authorization: Bearer <token>`

### `POST /api/v1/users`
Description: Creates a user to be used by a client

Request:
```
{
    "email":"<client email>",
    "password:"<client password>"
}
```

Response:
```
{
    "id":"<client id>",
    "email":"<client email>"
}
```

### `PUT /api/v1/users`
Description: Allows a user to update their email or password

Request:
```
{
    "email":"<client email>",
    "password:"<client password>"
}
```

Parameters:
- Headers
    - `Authorization: Bearer <token>`

Response:
```
{
    "id":"<client id>",
    "email":"<client email>"
}
```

### `POST /api/v1/login`
Description: Allows a client to login by returning a access and refresh token

Request:
```
{
    "email":"<client email>",
    "password:"<client password>"
}
```

Response:
```
{
    "id": "<client id>"
    "email":"<client email>"
    "token":"<client access token>"
    "refresh_token":"<client refresh token>"
}
```
### `POST /api/v1/refresh`
Description: Uses a refresh token to refresh an access token 

Parameters:
- Headers
    - `Authorization: Bearer <refresh token>` 

Response:
```
{
    "token":"<client access token>"
}
```
