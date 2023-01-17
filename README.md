## CarPark service

This is a learning project from the skillsmart coding school.

The main idea of the project was the absence of the initial specification.
Assignments were issued sequentially with limited time to complete, and it together with my still little experience in software design, the project turned into a monster during development.

I'm not going to refactor this project, but rather leave it as a reference example of bad design and a reminder of the importance of thinking through the architecture and specification.

Well, to tell the truth, it's easier to completely rewrite it than to refactor it. :trollface:

Summary of concept is service for some outsource drivers company managers. Manager can create vehicle, update it and gather milage reports, rides history with thier
geotrack on generated map. Manager can delete vehicles too (CRUD, I guess)

There are two utils to generate fake vehicles and geodata in **mock-utils**.

Also one of the tasks was to implement telegram bot for managers, you can find it in telegram-bot submodule. And I truly believe that it is designed way better then main service :D

CarPark depends on here.com geoservice (rest api for ride tracks on map) and openstreetmap used in fake-vechicle utility.

## Current state:

The project is "ready", in terms of all tasks are completed.
But, as I said already in paragraph above — it is terribly unfinished. 

Anyway it has:

- Backend with excess rest api (returning json's) serves multiuser service with... basic auth, which I will never use again. DBMS is PostgreSQL with Postgis, for storing and representing geopoints.
- Fronend interface of CarPark is build with go-templates with bootstrap, serving on /view routes.

There is also gorilla/srcf package in use for protection UI forms from cross-site scripting. But of course security of this service not a topic of discuss at all, at least because of basic auth... I should have used normal JWT tokens :man_shrugging:

Here is several UI screenshot of service with mock data.

![](https://github.com/NaNameUz3r/car-park/blob/main/misc/ui.gif)

## Watch this mess yourself!

You can find needed env variables example in project root (.env-example)

If tou want to touch this creature by your hands check you Docker directory:

```
docker-compose up
```

And log in with "admin2" or "admin1" username, and "qwerty" password.
There is, actually admin1 with same password. Explore :D

The dump contains a bunch of rides for 2023-2024 years for vehicles with ID's  25184,25185 and 25186

Currently, to be able to view ride routes on the map you should have here.com api token and pass it in docker-compose as env "HERE_API_KEY"