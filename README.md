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
    - dummy help command, provides no help, just like the uncaring universe
 - !check
    - Checks your current points. There is no way to get them yet because I haven't fully connected the twitch event webhook
 - !checkCost
    - Checks the point cost of each kill for a specific boss 
 - !addKills
    - Adds a number of boss kills for a specific boss to the queue, if the user has enough points
 - !checkKills
    - Checks if the calling user's kills have been done yet 

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
