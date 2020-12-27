import discord
import youtube_dl
from bot import YOUTUBE_DL_OPTIONS

class Song:
    def __init__(self, url, user):
        self.url = url
        self.user = user

        with youtube_dl.YoutubeDL() as ydl:
            video = self.get_info_url(url)
            video_format = video["formats"][0]
            self.stream_url = video_format["url"]
            self.video_url = video["webpage_url"]
            self.title = video["title"]
            self.uploader = video["uploader"] if "uploader" in video else ""
            self.thumbnail = video["thumbnail"] if "thumbnail" in video else None
            self.user = user

    def get_info_url(self, url):
        with youtube_dl.YoutubeDL(YOUTUBE_DL_OPTIONS) as ydl:
            info = ydl.extract_info(url, download=True)
            song = None
            if "_type" in info and info["_type"] == "playlist":
                return self.get_info_url(info["entries"][0]["url"])
            else:
                song = info
            return song
