### Car park service

This is a learning project from the skillsmart coding school.




#### Current state:


Gin server at 8888 port, with one Vehicle model stored in PostgreSQL database.

##### View

View with vehicles table is accessible with GET method at `:8888/view/vehicles` URI.

There is no ui for adding new records service tight now, but basic api already implemented:


##### Api

POST at `:8888/api/vehicles` with json like:

```
{
"description": "Very old cheap car",
"price": 150,
"mileage": 99999,
"manufactured": 1998
}
```

Will add new record in database.

- PUT `"/vehicles/:id` with same json structure will update record with this ID.

- DELETE `"/vehicles/:id` will delete record with this ID.
