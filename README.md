# dnd-playlist-controller
Small app to allow playlist switching using hotkeys for running my DND campaign

## Requirements
- Windows OS
- Go 1.22.1 to build the binary
- Spotify Dev application ID and Secret

## Usage
Create a .env file in the root of this project, and populate it with SPOTIFY_ID and SPOTIFY_SECRET

Create a hotkeys.json file (see hotkeys-example.json) to bind hotkeys to playlists. Format:
```
[
    {
        "id": 6, /* Unique integer ID */
        "uri": "4pxn3LHWRCd...", /* Spotify playlist URI */
        "modifiers": [ /* Array of modifiers */
            "Ctrl",
            "Alt",
            "Shift",
            "Win"
        ],
        "key": "6" /* Main key */
    },
    ...
]
```
Run build.bat to build the binary, then run.bat to launch the program.

To build and run manually:

```
go build -o bin\dnd-playlist-controller.exe .\cmd\dnd-playlist-controller\
.\bin\dnd-playlist-controller.exe
```