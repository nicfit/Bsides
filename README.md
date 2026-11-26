
``bash
# Load a playlist and print it to the screen
$ Bsides playlist.m3u

# Create playlists
$ Bsides new m3u -x --with-ext extinf "/Music/Type O Negative/1993 - Bloody Kisses/"
$ Bsides new m3u --recursive  "/Music/Type O Negative/" > allTypeO.m3u

# Iterate playlists
$ Bsides iter allTypeO.m3u
$ Bsides iter allTypeO.m3u --reverse
$ Bsides iter allTypeO.m3u --shuffle
$ Bsides iter allTypeO.m3u --shuffle --repeat=3
```