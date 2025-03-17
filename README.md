# gator
An rss feed aggregator project for a course on [boot.dev](https://www.boot.dev/)

## Description (from the lesson)
Gator is a multi-user CLI application. There's no server (other than the database), so it's only intended for local use.

## Installation
First of all, you need Postgres and Go installed to use the program. Then, you can use 'go build' to make an executable of the program to be used from a specific directory. Alternatively, you can use 'go install' to install the program as a CLI application.

## Config file
When you start the program for the first time, it will scan for a '.gatorconfig.json' file in the home directory. If it does not exist, it will create one. The config file includes to pieces of information: the connection string to the database and the current user. **When a new config file is created the current user defaults to an empty string. The first command you should use is 'gator register \<user>' to register a new user.**

## Commands
You can get a list of commands by running 'gator' without Here is the list of commands for gator:
- register (registers a new user)
- login (logs in an existing user)
- reset (deletes all users, feeds and posts)
- users (lists all users)
- agg (aggregates posts continuously)
- addfeed (adds a feed to the list of feeds and follows it)
- follow (follow a feed for the current user)
- following(lists feeds followed by the current user)
- unfollow (unfollows a feed)
- browse (lists posts from the feeds)