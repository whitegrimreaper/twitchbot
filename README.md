This is a bot for my twitch channel

[whitegrimreaper_](https://twitch.tv/whitegrimreaper_)

Using Golang because I'm a back-end engineer and as such I hate JS as much as it deserves.

## Current functions:
 - !pigeon
    - Insults pigeonmob

## Planned functionality:
 - integration with ?? API to track donations
    - Streamlabs wants a lot of info so might just use Twitch API and other stuff
 - integration with Twitch APIs to track bits, follows, subs, etc.
 - checking and controlling points gained from the above
 - spending those points to add activities to the subathon queue
 - Data storage using some sort of local db
    - either Mongo if i want nosql or sqllite and gorm if SQL