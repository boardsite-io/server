# Boardsite API
![Build Status](https://github.com/heat1q/boardsite/workflows/Boardsite%20CI/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/heat1q/boardsite)](https://goreportcard.com/report/github.com/heat1q/boardsite)

HTTP API for Boardsite, a whiteboard application build with web technologies.

```
# Sessions
POST /board/create # create session
GET /board/{id} # join session, upgrade protocol
DELETE /board/{id} # close session


# Pages
GET /board/{id}/pages # get all pages for session
POST /board/{id}/pages # add a page
{
    "pageId": {pageId}
    ...
}
PUT /board/{id}/pages/{pageId} # clear page
{}
DELETE /board/{id}/pages/{pageId} # delete a page
```
