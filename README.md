
# gator
boot.dev guided project RSS aggregator

## Description
CLI made in go that reads RSS feeds and shows them in the terminal.
Learning project to test XML, database connections, handlers, etc 
Bugs

## Requirements
Go (Go run . <args>)
Docker (docker-compose up, for postgres database)
Port 4321 available in host, to connect with database.

## Bugs - Knows Improvements
- gator database is not created on container creation (add volume to docker-compose with the dbinit to execute on init and create gator database)
- probably would be better to containerize the full app and enter a go container to run the application, virtual network to connect them and not expose the database port on the host
- probably some error verifications are missing
- some functions need pretty print
- a guide (help) of commands
 
