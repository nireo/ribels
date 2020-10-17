# ribels

A recode of the [bibels](https://github.com/eliaasi/bibels) in go.

## Running

The program starts from the `main.go` file. The `utils` folder holds different utilities like: osu api functions, name formatting and database related utils.

To run the bot you will need the [go programming language](https://golang.org/) and to have a PostgreSQL database setup. Before running the code you will need to setup a `.env` file. Here's an example setup:

```
DISCORD_TOKEN=your discord token
OSU_KEY=your osu api key
DATABASE_USER=username of database's owner
DATABASE_HOST=where the database is hosted
DATABASE_NAME=name of the database
DATABASE_PORT=the port on the host where the database is hosted
```

We just need to install a go package to interact with the discord api.

```
go get -u github.com/bwmarrin/discordgo
```

After that we can just run the bot:

```
go run main.go
```
