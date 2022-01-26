# spongebot

Spongebot is a small discord bot written in [Go] that creates "Spongebob" memes and text.  

## Installation using [Docker]

Clone repo.
```bash
git clone https://github.com/davidparks11/spongebot.git
```

Change dir into cloned repo and build the bot with docker.
```bash
cd spongebot
docker build -t spongebot .
```

Start the bot. A bot token must be provided as the first argument. 
```bash
docker run --rm -d --name spongebot spongebot [BOT_TOKEN]
```

### Stopping the bot
```bash
docker container stop spongebot
```

## Permissions
The bot must be able to reply to users, read content of messages, and create DMs with users. 

## Usage
For text replies via spongebot, @ the bot in a reply of a message

- Right click message
- Left click 'Reply'
- Type the following in the reply message
```
@yourBotsNameHere
```

For meme replies via spongebot, @ the bot in a reply of a message

- Right click message
- Left click 'Reply'
- Type the following in the reply message
```
@yourBotsNameHere --meme
```

For Usage, @ the bot in a normal channel message - you will be sent a direct message with the usage

## License
[MIT](https://choosealicense.com/licenses/mit/)

[Go]: https://go.dev/
[Docker]: (https://docs.docker.com/get-docker/)