A CRUD REST API built with Golang
=================================

Purpose
-------
This project is intended to showcase the use of a REST api built in GO.
Users are allowed to create and edit stories.
Stories need to be connected to a user's account.
Users are able to share what they have created. A form of expression.

# Overview

This project is a CRUD REST API connected to a PostgreSQL and Redis database.
The system supports Creating, Reading, Updating and Deleting stories.

## Technologies, Libraries and Tools
- Golang version >1.20.0
- Go Fiber
- Google uuid
- Golang PostgreSQL (lib/pq)
- Dot env (joho/godotenv)
- Visual Studio

# Usage

Simply run `git clone https://github.com/tolubydesign/angular-story-backend.git` and `go run main.go`.

## API Definition

```yml
Location: {{host}}/story
Location: {{host}}/stories
```

### Get Story

#### Get All Stories
```js
GET /stories
```

#### Get Stories Response

```js
200 Ok
```

```json
{
  "type": "success",
  "data": [
    {
      "id": [id],
      "title": "Title of story that we want to look at",
      "description": "A small synopsys of the story that is contained in this story",
      "content": {
        "children": [
          {
            "children": [
              {
                "description": "Description of events taking place in the story. This will describe what is happening.",
                "id": [id],
                "name": "Name for this part of the story."
              },
              {
                "children": [
                  {
                    "description": "lorem ipsum",
                    "id": [id],
                    "name": "lorem ipsum"
                  }
                ],
                "description": "Description of events taking place in the story. This will describe what is happening.",
                "id": [id],
                "name": "lorem ipsum..."
              },
            ],
            "description": "Lorem ipsum...",
            "id": [id],
            "name": "Lorem ipsum..."
          },
          {
            "children": [
              {
                "description": "Lorem ipsum...",
                "id": [id],
                "name": "Lorem ipsum..."
              }
            ],
            "description": "Lorem ipsum...",
            "id": [id],
            "name": "Lorem ipsum..."
          }
        ],
        "description": "Lorem ipsum...",
        "id": [id],
        "name": "Lorem ipsum..."
      }
    },
    ...
  ],
  "message": "Fetch all stories."
}
```
---

### Get a single Story
#### Get Story with ID
```js
GET /story/
```

#### Get Story Header
```json
{
  "id": [id],
}
```

#### Get Story Response
```js
200 Ok
```

```json
{
  "type": "success",
  "data": {
    "id": [id],
    "title": "text",
    "description": "text",
    "content": {
      "children": [
        {
          "description": "not set",
          "id": [id],
          "name": "not set"
        }
      ],
      "id": [id],
      "name": ""
    }
  },
  "message": "Fetch single story."
}
```
---

### Post New Story
#### Add Story
```js
POST /story
```

#### Post Story Request Body 
```js
{
  "title": "ADDED postman post request title",
  "description": "ADDED postman post request description",
  "content": {
    "children": [
      {
        "children": null,
        "description": "Quisque blandit magna vel lacinia fringilla. Mauris sit amet gravida tellus.",
        "id": "ebd00c42-841c-44f2-8e8e-bde095d502c6",
        "name": "Porttitor quis ultrices tortor"
      },
      {
        "children": [
          {
            "children": null,
            "description": "Quisque blandit magna vel lacinia fringilla. Mauris sit amet gravida tellus.",
            "id": "859da15f-8cbf-4d31-b799-0e1309726534",
            "name": "Porttitor quis ultrices tortor"
          }
        ],
        "description": "2 Quisque blandit magna vel lacinia fringilla. Mauris sit amet gravida tellus.",
        "id": "c37eeaea-23f1-448e-89bd-1c010605c90e",
        "name": "2 Porttitor quis ultrices tortor"
      }
    ],
    "description": "In aliquet nisi a.",
    "id": "7e0d122f-b295-4082-9d7c-242d7b2bd517",
    "name": "Nam blandit magna vel lacinia"
  }
}
```

#### Post Story Response
```js
200 Ok
```

```json
{
  "type": "success",
  "data": null,
  "message": "Database has been updated."
}
```
---


### Delete Story
#### Delete Story using id Request
```js
DELETE /story/{{id}}
```

#### Delete Story Header 
```json
{
  "id": [id],
}
```

#### Delete Story Response
```js
200 Ok
```

```json
{
  "type": "success",
  "data": null,
  "message": "Deleted story with id: [id]"
}
```
---


### Update Story
#### Update Story Request
```js
PUT /story
```

#### Update Story Header 
```json
{
  "id": [id],
  "title": "content",
  "description": "description"
}
```

#### Update Story Body
```js
{
  "content": {
    "children": [
      {
        "children": null,
        "description": "Lorem Ipsum...",
        "id": [id],
        "name": "Lorem Ipsum..."
      },
      {
        "children": null,
        "description": "Lorem Ipsum...",
        "id": [id],
        "name": "Lorem Ipsum..."
      }
    ],
    "description": "Lorem Ipsum...",
    "id": [id],
    "name": "Lorem Ipsum..."
  }
}
``` 

#### Update Story Response
```js
404 Error
```
or
```js
200 Ok
```
or
...

```json
{
  "type":"success",
  "data":null,
  "message":"Updated story with id: [id]"
}
```
---

## In Development
- Improved Error Response
- JWT token
- User Profiles
- User login method

## Disclaimer
This is an educational project. The source code is licensed under the MIT license.

## License
This library is distributed under the [MIT](LICENSE) license.
