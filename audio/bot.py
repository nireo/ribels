from discord.ext import commands
import youtube_dl
import discord
import asyncio
from config import config

YOUTUBE_DL_OPTIONS = {
    "default_search": "ytsearch",
    "format": "bestaudio/best",
    "quiet": True,
    "extract_flat": "in_playlist"
}

class Song:
    def __init__(self, search):
        with youtube_dl.YoutubeDL(YOUTUBE_DL_OPTIONS) as ytdl:
            song = self.info(search)
            song_format = song["formats"][0]
            self.title = song["title"]
            self.by = song["uploader"] if "uploader" in song else ""
            self.stream_url = song_format["url"]
            self.video_url = song["webpage_url"]

    def info(self, search):
        with youtube_dl.YoutubeDL(YOUTUBE_DL_OPTIONS) as ytdl:
            info = ytdl.extract_info(search, download=False)
            return info

class MusicPlayer(commands.Cog):
    def __init__(self, bot):
        self.bot = bot
        # list of guilds with their current music status
        self.status_list = {}

    # Get the current guild's music state
    def get_guild_status(self, guild):
        if guild.id in self.status_list:
            return self.status_list[guild.id]
        else:
            self.status_list[guild.id] = GuildStatus()
            return self.status_list[guild.id]


    # Handle playing songs and then continuing to new songs
    def play_helper(self, client, status, song):
        status.playing = song
        def after_done(err):
            if len(status.queue) > 0:
                self.play_helper(client, status, status.queue.pop(0))
            else:
                asyncio.run_coroutine_threadsafe(client.disconnect(), self.bot.loop)
        client.play(
           discord.PCMVolumeTransformer(discord.FFmpegPCMAudio(song.stream_url), volume=status.volume),
           after=after_done
        )
    
    @commands.command(aliases=["s"])
    @commands.guild_only()
    async def stop(self, ctx):
        client = ctx.guild.voice_client
        status = self.get_guild_status(ctx.guild)
        if client and client.channel:
            await client.disconnect()
            status.queue = []
            status.playing = None
        else:
            await ctx.send("Not in voice channel")

    @commands.command()
    @commands.guild_only()
    async def tutorial(self, ctx):
        embed = discord.Embed(title="Help", description="List of all the supported commands!")
        embed.add_field(name=";vol <volume:int>", "Change the volume of the music player", False)
        embed.add_field(";play <url>", "Play a given song from a youtube url", False)
        embed.add_field(";pause", "Stop the current song (resume with same command)", False)
        embed.add_field(";playing", "Get information about the current song", False)
        embed.add_field(";queue", "Get the song queue for the current guild", False)
        embed.add_field(";clear", "Clear the hole queue")
        await ctx.send("", embed=embed)

    @commands.command(aliases=["current"])
    @commands.guild_only()
    async def playing(self, ctx):
        client = ctx.guild.voice_client
        status = self.get_guild_status(ctx.guild)
        if client and client.channel:
            await ctx.send(status.playing.get_message())
        else:
            await ctx.send("Not in a voice channel")

    @commands.command(aliases=["continue", "resume"])
    @commands.guild_only()
    async def pause(self, ctx):
        client = ctx.guild.voice_client
        if client.is_paused():
            client.resume()
        else:
            client.pause()


    @commands.command(aliases=["playlist", "q"])
    @commands.guild_only()
    async def queue(self, ctx):
        client = ctx.guild.voice_client
        status = self.get_guild_status(ctx.guild)
        if client and client.channel:
            if len(status.queue) == 0:
                await ctx.send("Queue is currently empty")
                return 
            content, index = "", 1
            for song in status.queue:
                content += str(index) + song.title
                index += 1

            await ctx.send(content)
        else:
            await ctx.send("Not in voice channel")

    # Clear the queue
    @commands.command(aliases=["c"])
    @commands.guild_only()
    async def clear(self, ctx):
        status = self.get_guild_status(ctx.guild)
        status.queue = []

    @commands.command(aliases=["vol"])
    @commands.guild_only()
    async def volume(self, ctx, vol: int):
        status = self.get_guild_status(ctx.guild)
        if vol < 0:
            vol = 0
        if vol > 150:
            vol = 150

        client = ctx.guild.voice_client
        status.volume = float(vol) / 100.0
        client.source.volume = status.volume
        await ctx.send(f"Volume set to: `{float(vol)}%`")
    
    @commands.command(aliases=["p", "song"])
    @commands.guild_only()
    async def play(self, ctx, *, query):
        client = ctx.guild.voice_client
        status = self.get_guild_status(ctx.guild)

        # if already in a voice channel
        if client and client.channel:
            try:
                song = Song(query)
            except youtube_dl.DownloadError:
                await ctx.send("Error loading video")
                return
            status.queue.append(song)
            await ctx.send("Added to queue")
        else:
            if ctx.author.voice is not None and ctx.author.voice.channel is not None:
                channel = ctx.author.voice.channel
                try:
                    song = Song(query)
                except youtube_dl.DownloadError:
                    await ctx.send("Error loading video")
                    return
                client = await channel.connect()
                self.play_helper(client, status, song)
                await ctx.send(f"Now playing {song.title}")
            else:
                await ctx.send("You need to be in a voice channel!")

class GuildStatus:
    def __init__(self):
        self.volume = 1.0
        self.playing = None
        self.queue = []

bot = commands.Bot(command_prefix=";")

@bot.event
async def on_ready():
    print(f"Logged in as {bot.user.name}")

bot.add_cog(MusicPlayer(bot))
bot.run(config["discord_token"])
