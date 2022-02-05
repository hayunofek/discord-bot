# Discord Bot

![Discord Bot](/docs/discordgo-music.png)

----

Discord Bot, is a discord bot that was written using Go and the module Discord-Go, with a main purpose, to play music. As of right now, this bot is very minimal and just contains the command play. In the future I might add some new features.

# How to use
1. Download the relevant binary from [Releases](https://github.com/hayunofek/discord-bot/releases)
1. Extract the binary on your computer
1. Create a file named `token` which contains **only** the discord token. Information about how to get this token can be found [further in this document](#how-to-get-a-discord-bot-token).
1. Add your bot to your server. Information about how to add a discord bot to your server can be found [further in this document](#how-to-add-a-bot-to-my-server).
1. Run the binary you downloaded, and enjoy the bot!

## Commands
The default prefix command is `!`, which is defined in `cmd/discord_command.go` as the constant `PREFIX_SIGN`.
### The Play Command
Usage: `!play [youtube_url]`  
Example: `!play https://www.youtube.com/watch?v=klZNvJArVSE`

## How to get a discord bot token
1. Go into [Discord Developer Portal](https://discord.com/developers/applications) and sign in.
2. Create an application by click the **New Application** blue button.
3. Click the **Bot** blade on the left side of the page.
4. Click **Add Bot** on the right side of the page.
5. Click on **Click to Reveal Token** and take that token

## How to add a bot to my server
1. Go into [Discord Developer Portal](https://discord.com/developers/applications) and sign in.
2. Click on the **General Information** blade located on the left side of the page.
3. Copy the **Application Id** placed around the center of the page.
4. Replace the application id in the following url `https://discord.com/api/oauth2/authorize?client_id=[application_id]&permissions=0&scope=bot%20applications.commands`
5. Go into the url you just created, and add the bot
