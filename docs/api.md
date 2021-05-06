# API
## General
All messages have the following structure:
```
interface Message {
    type: string
    sender?: string
    content?: any
    error?: any
}
```
The message type gives information on the content. The sender can give info on the origin of the message (with the help of `userId`). The content can be any JSON serializable value. The error field is solely populated by the server in case of failure.

## Routes
Accepted Content-Types: `application/json`, `plain/text`
 Routes | Methods | Description | Request Content | Response Content
 -------|---------|-------------|--------------|--------------
 `/b/create` | `POST` | Create a new session | - | `string`
 `/b/{id}` | `DELETE` | Close and clear the sesion | - | -
 `/b/{id}/users` | `GET` | Get all connected users | - | `{${id}: any}`
 `/b/{id}/users` | `POST` | Register a new user for the session | `{alias: string, color: string}` | `{id: string, alias: string, color: string}`
 `/b/{id}/users/{userId}/socket` | `GET` | Join a session with ID `{id}` as user `{userId}` and upgrade to websocket protocol if successful | - | -
 `/b/{id}/pages` | `GET` | Return all page IDs of the session in order | - | `string[]`
 `/b/{id}/pages` | `POST` | Add a page with ID and an index to denote the position | `{pageId: string, index: number}` | -
 `/b/{id}/pages/{pageId}` | `GET` | Get all data on the page `{pageId}` | - | `Stroke[]`
 `/b/{id}/pages/{pageId}` | `PUT` | Update page `${pageId}` | `{clear: bool, meta: any}` | -
 `/b/{id}/pages/{pageId}` | `DELETE` | Delete a page | - | -
 `/b/{id}/attachments` | `POST` | Upload file via MIME `multipart/form-data` with key `file`. Returns `{attachId}` on success | any blob | `string`
 `/b/{id}/attachments/{attachId}` | `GET` | Fetch file | - | any blob

## WS Message Content
### Stroke 
**Message Type**: `stroke`
```
{
    strokeType: number
    strokeId?: string
    userId?: string
    pageId?: string
    x?: number
    y?: number
    points?: number[]
    style?: {
        color: string
        width: number
        opacity: number
    }
}
```

### User Connected/Disconnected
**Message Type**: `{userconn, userdisc}`
```
{
    id: string
    alias: string
    color: string
}
```

### Page Sync/Clear
**Mesage Type**: `pagesync`
```
{
    pageRank: string[]
    meta: any
}
```
**Mesage Type**: `pageupdate`
```
{
    pageId: string
    clear: bool
    meta: any
}
```

### Mouse Move Event
**Mesage Type**: `mmove`
```
{
    x: number
    y: number
}
```
