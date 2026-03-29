# Gator - a Blog Aggregator in Go
This is my repo for the Boot.dev Guided Project (currently Course 17) called 'Build a Blog Aggregator in Go.'

## Welcome to the Blog Aggregator (from the first lesson in the first chapter for the course on Boot.dev)
We're going to build an [RSS](https://en.wikipedia.org/wiki/RSS) feed aggregator in Go! We'll call it "Gator", you know, because aggreGATOR 🐊. Anyhow, it's a CLI tool that allows users to:
- Add RSS feeds from across the internet to be collected
- Store the collected posts in a PostgreSQL database
- Follow and unfollow RSS feeds that other users have added
- View summaries of the aggregated posts in the terminal, with a link to the full post

RSS feeds are a way for websites to publish updates to their content. You can use this project to keep up with your favorite blogs, news sites, podcasts, and more!

## Prerequisites
### PostgreSQL (an SQL database)
You will need a working [PostgreSQL ](https://www.postgresql.org/) installation in order to use this command line application. Installation, configuration, and security of PostgreSQL is beyond the scope of this project, so please refer to the linked site as needed. This project *does* provide the schema files necessary to create the database used by this application.
### Go
The Go language and development support files will be needed if you want to update or modify the code for this project. If you only want to run the application, you will not need Go or its support files installed.

## Get the code from GitHub
Clone this repo by running the following command (on a Linux host) while in the directory where you want the project:
```
git clone https://github.com/KobayashiComputing/gator.git
```

## Build the gator binary
CD into the directory created by in the previous step and build the binary:
```
cd gator
go build
```

## Install the gator binary
In the root directory of the repo, use the following command to install the binary for your user:
```
go install
```

## Create the database and its tables
Coming soon...

## Create the gator config file in the user's home directory
Gator uses a JSON (JavaScript Object Notation) configuration file to store the PostgreSQL command to access the database as well as a variable (current_user_name) that contains the username of the currently "logged in" (to the application) user.

The format of the file is:
```
{
    "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
    "current_user_name": ""
}
```

1. Create a file in the home directory named .gatorconfig.json
2. Copy the lines above - including the braces ("{}") to the file
3. Adjust the postgres command to meet your isntallation

## Run gator
Gator is a command line application, which means that it is run from and interacts with the user from a command line or "terminal." A "help" command (should be) coming soon, but until then here is a list of supported commands:

1. User Management
    - register <username>
    - login <username>
    - users
2. RSS Feed Management
    - addfeed <feed name> <url of feed>
    - feeds
    - follow <url of feed>
    - following
    - unfollow <url of feed>
    - browse
3. Get or refresh the currently added feeds
    - agg
4. Reset the database - DANGER: This will delete all records in the database
    - reset

