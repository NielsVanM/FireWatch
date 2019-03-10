# API Doc
This is the basic API documentation, a better way of documenting the API will be released at a later date. A general response to any request is in the format of:
```json
    {
        "succes": <true/false>,
        "status_code": "<status_string>",
        "data": "object with requested data"
    }
```

## Table of Contents
1. [Authentication](#authentication)
    
<a name="authentication"></a>
## Authentication
### Login

Endpoint: `/api/v1/login/`

Request: 
```json
    {
        "username": "<username>",
        "password": "<password>"
    }
```
Response:
```json
    {
        "success":true,
        "status_code":"okay",
        "data": {
            "token":"<token>"
        }
    }
```
Errors:
```json
    [
        "internal_error",
        "invalid_request",
        "invalid_credentials"
    ]
```


### Logout

Endpoint: `/api/v1/logout/`

Request:
```json
    {
        "token": "<token>"
    }
```
Response:
```json
    {
        "success": true,
        "status_code": "okay",
        "data":{}
    }
```
Errors:
```json
    [
        "internal_error",
        "invalid_request",
        "invalid_token"
    ]
```