# Feed Server

A small server to serve events and markets

## Endpoints

The following endpoints are available.

Please note that the data may not be consistenly formatted.

### /football/events

Returns an array of event ids (integers)

E.g.

```[1, 2, 3]```

### /football/events/{id}

`id` must be a valid integer

Returns a single event

```json
{
    "id": 1,
    "name": "Southampton v Bournemouth",
    "time": "2017-08-20:15:00:00Z",
    "markets": [
       101,
       102
    ]
}
```

### /football/markets/{id}

`id` must be a valid integer

Returns a single market

```json
{
    "id": "101",
    "type": "win-draw-win",
    "options": [
        {
            "id": "10101",
            "name": "Southampton",
            "odds": "3/5"
        },
        {
            "id": "10102",
            "name": "Draw",
            "odds": "4/5"
        },
        {
            "id": "10103",
            "name": "Bournemouth",
            "odds": "5/1"
        }
    ]
}
```
