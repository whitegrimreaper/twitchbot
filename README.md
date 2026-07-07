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
 - !help
    - now actually helps! run !help or !help <commandName>
 - !check
    - checks your current points
 - !checkCost <boss name>
    - checks the point cost per kill for a specific boss
 - !addKills <boss name> <# of kills>
    - spends your points to add kills for a specific boss to the queue
 - !checkKills
    - shows how many of your kills are still left in the queue
 - !howToEarn
    - explains how to earn points

## Planned functionality:
 - integration with ?? API to track donations
    - Streamlabs wants a lot of info so might just use Twitch API and other stuff
 - integration with Twitch APIs to track bits, follows, subs, etc.
    - Done! 
 - checking and controlling points gained from the above
    - Done!
 - spending those points to add activities to the subathon queue
    - Done! 
 - Data storage using gorm and sqlite
    - using gorm means I can scale out of sqlite if i need to
    - Done!

## Planned Functionality (including other packages)
 - webhosted interface for checking points and queue
    - Not done! Webhosting for a stream overlay exists, but no fancy actual interface. 
 - CLI for manually modifying DB (in case of bugs or moderator action)
    - Not done! Probably not bothering with this for a while, as it would introduce some mutex/locking issues and also i can just manually change the DBs with VSCode + extensions 
 - Streamdeck integration for ticking off kills?
   - Done! Not sure why I don't have a repo for this, but will probably make one soon™
