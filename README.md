# dvembed
(D)iscord (v).redd.it media (embed)der

`dvembed` uses `youtube-dl` to download v.redd.it media, at the highest available audio and video quality.
If the media file is smaller than the 8MB Discord file upload limit, the bot will post the file, meaning 
v.redd.it media can now be embedded nicely within Discord, without having to leave the app and view the 
thread.

###### Future
If the media file exceeds 8MB, `dvembed` will attempt to compress or reduce the file size with `ffmpeg`.

## Usage
`dvembed`, or `go run dvembed` to run without executable

## Requirements
- `Python (2.6, 2.7, 3.2+)`

- `youtube-dl`

- `ffmpeg`

###### Future
A script will setup a working environment if these programs are not available.

## Example
![Example](https://github.com/kepler471/dvembed/blob/master/example.png?raw=true)
