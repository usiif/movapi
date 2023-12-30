Movapi is a comprehensive and structured JSON API designed to provide access to an extensive database of movie information.  
# Overview 
Essentially, Movapi is an accessible movie database that you can interact with via a user-friendly API. Built using **Go**, it ensures that safety and performance come first.  
# Routes

| Endpoint                           | Description                               | Method                          |
|------------------------------------|-------------------------------------------|---------------------------------|
| `/v1/healthcheck`               | Show application health and version info     | _GET_                                 |
| `/v1/movies`                    | Show details of all movies                   | _GET_                               |
| `/v1/movies`                   | Create a new movie                            | _POST_                                |
| `/v1/movies/:id`                | Show details of a specific movie             | _GET_                            |
| `/v1/movies/:id`              | Update details of a specific movie             | _PATCH_                                |
| `/v1/movies/:id`             | Delete a specific movie                         | _DELETE_                         |
| `/v1/users`                    | Register a new user                           | _POST_                                 |
| `/v1/users/activated`           | Activate a specific user                     | _PUT_                                |
| `/v1/users/password`            | Update password for a specific user          | _PUT_                                |
| `/v1/tokens/authentication`    | Generate a new authentication token           | _POST_                                |
| `/v1/tokens/password-reset`    | Generate a new password-reset token           | _POST_                                |

# Authentication Details 
To interact with Movapi's endpoints, you would have to be an authenticated activated user. To regiester a new user, hit the ``v1/users`` endpoint with a JSON request in the format  
```
{
  "name": "",
  "email": "",
  "password": ""
}
```
An authentication token will be sent to the registered email address. Please submit it using the following format to ``v1/users/activated``:
```
{"token": ""}
```
After that, you are free to use the API as you choose. To begin, you will be given _read_ privileges.  

## Authorization 
Whenever you log into your account, you will be sent a stateful token in the format: 
```
{
"authentication_token": {
    "token": "IEYZQUBEMPPAKPOAWTPV6YJ6RM",
    "expiry": "2021-04-17T11:03:36.767078518+02:00"
    }
}
```
For subsequent requests, You must include the provided token in an ``Authorization`` header in the format:
```
Authorization: Bearer IEYZQUBEMPPAKPOAWTPV6YJ6RM"
```
# Getting started 
To post a new movie, Follow the format: 
```
{
  "title": "The Holdover",
  "year": 2023,
  "runtime": 133,
  "genres": ["Comedy", "Drama"]
}
```
You can hit the ``v1/movies`` enpoint with the following filters: 
- title
- genres
- page
- page_size
- sort

Example request might look something like that: 
``/v1/movies?title=godfather&genres=crime,drama&page=1&page_size=5&sort=-year``  

In this example, the - before the year field in the sort query parameter indicates that the movies should be sorted in descending order based on the year attribute. Adjust the field as needed for your specific use case.
