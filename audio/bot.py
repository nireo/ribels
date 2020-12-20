from discord.ext import commands
import youtube_dl
import discord
import asyncio

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

    def get_message(self):
        return f"Currently playing {self.title} by {self.by}"

class MusicPlayer(commands.Cog):
    def __init__(self, bot):
        self.bot = bot
        self.status_list = {}

    def get_guild_status(self, guild):
        if guild.id in self.status_list:
            return self.status_list[guild.id]
        else:
            self.status_list[guild.id] = GuildStatus()
            return self.status_list[guild.id]


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
    
    @commands.command()
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
bot.run("NjA0Mzk0MDc5NDIxNDY0NjAy.XTtT2Q.wlx3JvuWRTPFRb1WdZkA2hbwkqo")
