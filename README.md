# dvembed
(D)iscord (v).redd.it media (embed)der

`dvembed` uses `youtube-dl` to download v.redd.it / vreddit media, at the highest available audio and video quality.
If the media file is smaller than the 8MB Discord file upload limit, the bot will post the file, meaning 
v.redd.it media can now be embedded nicely within Discord, without having to leave the app and view the 
thread.

###### Future
- If the media file exceeds 8 MB, `dvembed` will attempt to use alternate versions of the file hosted by v.redd.it.
- Some larger will not have alternate versions that do not exceed the 8 MB limit. The bot will attempt to compress and reduce the file size with `ffmpeg`.

## Usage
Currently, this is a private bot, and so you will need to create a Discord Application, with a Bot User, and then generate a token. Please visit the [Discord developer site](https://discord.com/developers/).

Run the bot with an executable (or directly from source with go run)

`dvembed -t SECRET-TOKEN`

or, to have run as a background process (Linux/Unix only):

`nohup dvembed -t SECRET-TOKEN &`

By default, output will be sent to the file `nohup.out`, but you can specify with 

`nohup dvembed -t SECRET-TOKEN ./dvembed.log`

### Token
Flag: `-t SECRET-TOKEN`

Token can also be added to `token.go` if building from source, so you can omit the flag.

## Requirements
- `Python (2.6, 2.7, 3.2+)`

- `youtube-dl`

- `ffmpeg`

###### Future
- The functionality that `youtube-dl` provides will be shifted to Go code, so will not be required in the future.

## Example
![Example](https://github.com/kepler471/dvembed/blob/master/example.png?raw=true)
