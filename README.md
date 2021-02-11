# Boardsite
![Build Status](https://github.com/heat1q/boardsite/workflows/Boardsite%20CI/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/heat1q/boardsite)](https://goreportcard.com/report/github.com/heat1q/boardsite)

API for Boardsite, a whiteboard application built with web technologies. The API utilizes Websockets for near realtime communication between clients. 

# API
## Routes
Accepted Content-Types: `application/json`, `plain/text`
 Routes | Methods | Description | Request Body | Response Body
 -------|---------|-------------|--------------|--------------
 `/b/create` | `POST` | Create a new session | - | `{sessionId: string}`
 `/b/${id}` | `GET` | Join a session with ID `${id}` and upgrade to websocket protocol if successful | - | -
 `/b/${id}` | `DELETE` | Close and clear the sesion | - | -
 `/b/${id}/pages` | `GET` | Return all page IDs of the session in order | - | `{pageRank: string[]}`
 `/b/${id}/pages` | `POST` | Add a page with ID and an index to denote the position | `{pageId: string, index: number}` | -
 `/b/${id}/pages/${pageId}` | `GET` | Get all data on the page `${pageId}` | - | `Stroke[]`
 `/b/${id}/pages/${pageId}` | `PUT` | Clear all data on the page `${pageId}` | - | -
 `/b/${id}/pages/${pageId}` | `DELETE` | Delete a page | - | -

## Websocket
All data transmitted over the websocket is serializable to a single interface. We refer to this interface as `Stroke`.

```
interface Stroke {
    type: number
    id?: string
    pageId?: string
    x?: number
    y?: number
    points?: number[]
    style?: {
        color: string
        width: number
    }
    PageRank?: string[]
    pageClear?: string[]
}
```
The application defines the `type`, the server only knows the following types:
1. `type > 0`: relay and cache stroke in Redis
2. `type < 0`: only relay stroke
3. `type == 0`: relay and remove stroke from Redis