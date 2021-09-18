from discord.ext import commands
import youtube_dl
import discord
import asyncio
import os
from config import config

YOUTUBE_DL_OPTIONS = {
    "default_search": "auto",
    "format": "bestaudio/best",
    "quiet": True,
    "noplaylist": True,
    "nocheckcertificate": True,
    "ignoreerrors": False,
    "no_warnings": True,
    "source_address": "0.0.0.0"
}

ffmpeg_options = {
    'options': '-vn'
}

ytdl = youtube_dl.YoutubeDL(YOUTUBE_DL_OPTIONS)
class YTDLSource(discord.PCMVolumeTransformer):
    def __init__(self, source, *, data, volume=0.5):
        super().__init__(source, volume)

        self.data = data

        self.title = data.get('title')
        self.url = data.get('url')
        self.thumbnail = data.get('thumbnail')
        self.duration = data.get('duration')
        self.views = data.get('view_count')

    @classmethod
    async def from_url(cls, url, *, loop=None, stream=False):
        loop = loop or asyncio.get_event_loop()
        data = await loop.run_in_executor(None, lambda: ytdl.extract_info(url, download=not stream))

        if 'entries' in data:
            data = data['entries'][0]

        filename = data['url'] if stream else ytdl.prepare_filename(data)
        return cls(discord.FFmpegPCMAudio(filename, **ffmpeg_options), data=data)


async def is_currently_playing(ctx):
    client = ctx.guild.voice_client
    if client and client.channel and client.source:
        return True
    else:
        raise commands.CommandError("Currently not playing audio")


async def in_voice(ctx):
    voice = ctx.author.voice
    bot_voice = ctx.guild.voice_client
    if voice and bot_voice and voice.channel and bot_voice.channel and voice.channel == bot_voice.channel:
        return True
    else:
        raise commands.CommandError("You're not in a voice channel.")


class MusicPlayer(commands.Cog):
    def __init__(self, bot):
        self.bot = bot
        self.queue = []
        self.curr_playing = None
        self.guild_states = {}

    @commands.command(aliases=["disconnect", "leave"])
    @commands.guild_only()
    async def stop(self, ctx):
        """ Completely stop playing music and leave a channel. Also clears the queue. """
        client = ctx.guild.voice_client
        if client:
            await client.disconnect()
            after_playing()
            self.queue = {}
        else:
            await ctx.send("Not in a voice channel")

    def get_guild_state(self, guild):
        if guild.id in self.guild_states:
            return self.guild_states[guild.id]
        else:
            self.guild_states[guild.id] = GuildState()
            return self.guild_states[guild.id]

    @commands.command(aliases=["cp", "np"])
    async def currplaying(self, ctx):
        guild_state = self.get_guild_state(ctx.guild)
        await ctx.send("", embed=guild_state.now_playing.get_embed())

    def play_helper(self, client, state, song):
        state.now_playing = song
        source = discord.PCMVolumeTransformer(
            discord.FFmpegPCMAudio(song.stream_url), volume=state.volume)

        async def after_playing(err):
            if len(state.playlist) > 0:
                next_song = state.playlist.pop(0)
                self.play_helper(client, state, next_song)
            else:
                await client.disconnect()

        client.play(source, after=after_playing)
    
    @commands.command()
    async def testplay(self, ctx):
        pass
        

    @commands.command()
    @commands.guild_only()
    async def join(self, ctx, *, channel: discord.VoiceChannel):
        """ Make the bot join a voice channel without playing anything. """
        client = ctx.guild.voice_client
        if client is not None:
            return await client.move_to(channel)
        await channel.connect()


    @commands.command(aliases=["vol", "v"])
    @commands.guild_only()
    async def volume(self, ctx, volume:int):
        """ Set the volume of the music player """
        client = ctx.guild.voice_client
        if client is None:
            return await ctx.send("Not connected to a voice channel.")

        client.source.volume = volume / 100
        await ctx.send("Changed volume to {}%".format(volume))


    @commands.command(aliases=["continue", "resume"])
    @commands.guild_only()
    async def pause(self, ctx):
        """ Stop playing until the playing is resumed. """
        client = ctx.guild.voice_client
        if client.is_paused():
            client.resume()
        else:
            client.pause()


    @commands.command(aliases=["p", "song"])
    @commands.guild_only()
    async def play(self, ctx, *, url):
        """ Play a song from a given URL """

        # Check if if we're currently playing a song.
        if self.curr_playing == None:
            self.curr_playing = url
        else:
            self.queue.append(url)

        async with ctx.typing():
            player = await YTDLSource.from_url(url, loop=self.bot.loop)
            ctx.guild.voice_client.play(player, after=lambda e: print("Player error: %s" % e) if e else None)
        await ctx.send("Now plying: {}".format(player.title))

    @play.before_invoke
    @pause.before_invoke
    @volume.before_invoke
    @join.before_invoke
    @stop.before_invoke
    async def ensure_voice(self, ctx):
        """ 
        Check that the messager is in a voice channel so they cannot interact with 
        the playing.
        """
        if ctx.guild.voice_client is None:
            if ctx.author.voice:
                await ctx.author.voice.channel.connect()
            else:
                await ctx.send("You are not connected to a voice channel.")
                raise commands.CommandError("Author not connected to a voice channel.")
        elif ctx.guild.voice_client.is_playing():
            ctx.guild.voice_client.stop()

class GuildState:
    def __init__(self):
        self.volume = 1.0
        self.playlist = []
        self.now_playing = None
    
    def same_requester(self, user):
        return self.now_playing.user == user 

bot = commands.Bot(command_prefix=";")
def after_playing():
    """ Removes all the downloaded .webm files after playing. """
    CURR_DIR = os.path.dirname(os.path.realpath(__file__))
    for path in os.listdir(CURR_DIR):
        if path.endswith(".webm"):
            os.remove(path)

@bot.event
async def on_ready():
    print(f"Logged in as {bot.user.name}")

bot.add_cog(MusicPlayer(bot))
bot.run(config["discord_token"])
