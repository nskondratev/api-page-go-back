# Websockets API
## URL
```
{backend_url}/api/ws
```

## Events
* [Event created: `ap_event_created`](#ap_event_created)
* [Event updated: `ap_event_updated`](#ap_event_updated)
* [Event created: `ap_event_deleted`](#ap_event_deleted)
* [Page created: `ap_page_created`](#ap_page_created)
* [Page updated: `ap_page_updated`](#ap_page_updated)
* [Page deleted: `ap_page_deleted`](#ap_page_deleted)

## ap_event_created
Event is emitted when some event is created. Example:

```json
{
  "event": "ap_event_created",
  "data": {
    "event": {
      "id": 1,
      "label": "Event 1",
      "constant": "EVENT_1",
      "value": "event_1",
      "description": "New event",
      "type": "frontend",
      "fields": [
        {
          "id": 1,
          "type": "Integer",
          "key": "id",
          "required": true,
          "description": "Primary key"
        }  
      ]
    }
  }
}
```

## ap_event_updated
Event is emitted when some event is updated. Example:

```json
{
  "event": "ap_event_updated",
  "data": {
    "event": {
      "id": 1,
      "label": "Event 1 updated",
      "constant": "EVENT_1",
      "value": "event_1",
      "description": "New event",
      "type": "frontend",
      "fields": [
        {
          "id": 1,
          "type": "Integer",
          "key": "id",
          "required": true,
          "description": "Primary key"
        }  
      ]
    }
  }
}
```

## ap_event_deleted
Event is emitted when some event is deleted. Example:

```json
{
  "event": "ap_event_deleted",
  "data": {
    "id": 1
  }
}
```

## ap_page_created
Event is emitted when some page is created. Example:

```json
{
  "event": "ap_page_created",
  "data": {
    "page": {
      "id": 1,
      "title": "Page 1",
      "text": "Page 1 Text"
    }
  }
}
```

## ap_page_updated
Event is emitted when some page is updated. Example:

```json
{
  "event": "ap_page_updated",
  "data": {
    "page": {
      "id": 1,
      "title": "Page 1 updated",
      "text": "Page 1 Text"
    }
  }
}
```

## ap_page_deleted
Event is emitted when some page is deleted. Example:

```json
{
  "event": "ap_page_deleted",
  "data": {
    "id": 1
  }
}
```
