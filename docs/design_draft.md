# Design

Implement simple web-server application for storing stores items catalog.

* Implement storage in memory as array. Implement storage as package.\
Storage should provide three methods:\
Add (one or more items),\
Search (search for produce by partial match of whole words),\
Fetch (get information about specific item).\
Delete (delete item).\
Storage package should be thread-safe.
* Implemet REST API package. Define REST API Request/Payload formats. Error formats. Use JSON.\
REST API should implement Add, Search, Fetch methods wrapping storage methods.
* Code should be testable. Cover code with unit tests.
* Add fixtures for predefined items in storage.
* Add build/deployment instructions. Implement Makefile for build/deploy commands. Implement Dockerfile.

## Storage Format
Store Table Column:
```
Produce Code - case-insensitive alphanumeric, sixteen characters long,
with dashes separating each four character group (example: A12T-4GH7-QPL9-3N4M).

Name - name of produce, alphanumeric and case insensitive (example: Lettuce, Gala Apple).

Unit price - price of produce, a number with up to 2 decimal places (example: $3.46).
```

## API Format
### Add Request
```
POST
/store/add
{
    "items": [
        {
            "code": "A12T-4GH7-QPL9-3N4M",
            "name": "Lettuce",
            "price": 3.46
        },
        {
            "code": "1111-2222-3333-4444",
            "name": "Pomegranate",
            "price": 5.01
        }
    ]
}
```
Add Response
```
{
    "item_codes": ["A12T-4GH7-QPL9-3N4M", "1111-2222-3333-4444"],
    "item_count": 2
}
```

### Search Request
```
POST
/store/search
{
    "search": "Red Grapefruit"
}
```
Search Response
```
{
    items: [
        {
            "code": "A12T-4GH7-QPL9-3N4M",
            "name": "Lettuce",
            "price": 3.46
        },
        {
            "code": "1111-2222-3333-4444",
            "name": "Pomegranate",
            "price": 5.01
        }
    ]
}
```

### Fetch Request
```
GET
/store/{code}
```
Fetch Response
```
{
    item: {
        "code": "A12T-4GH7-QPL9-3N4M",
        "name": "Lettuce",
        "price": 3.46
    }
}
```

### Delete Request
```
DELETE
/store/delete
{
    "item_codes": ["A12T-4GH7-QPL9-3N4M", "1111-2222-3333-4444"]
}
```
Delete Response
```
{
    "item_count": 2
}
```

### Error Response
```
{
    "code": 500,
    "error": "verbal name of error",
    "message": "error details"
}
```
