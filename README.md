This is a bot for my twitch channel

[whitegrimreaper_](https://twitch.tv/whitegrimreaper_)

Using Golang because I'm a back-end engineer and as such I hate JS as much as it deserves.

## Requirements:
(note that this is not intended to be used by anyone else yet as it is nowhere near feature complete)
 - A bot account with moderator privileges
 - SQLite
 - vague knowledge of how to use commandline

## Current functions:
 - !pigeon
    - Insults pigeonmob
 - !check
    - Checks your current points. There is no way to get them yet because I haven't fully connected the twitch event webhook

## Planned functionality:
 - integration with ?? API to track donations
    - Streamlabs wants a lot of info so might just use Twitch API and other stuff
 - integration with Twitch APIs to track bits, follows, subs, etc.
 - checking and controlling points gained from the above
 - spending those points to add activities to the subathon queue
 - Data storage using gorm and sqlite
    - using gorm means I can scale out of sqlite if i need to 
