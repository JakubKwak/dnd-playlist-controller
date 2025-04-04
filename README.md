# dnd-playlist-controller
Small app to allow playlist switching using hotkeys for running my DND campaign

## Requirements
- Windows OS
- Spotify Dev application ID and Secret
- If you plan to compile the binary yourself, you'll need Go 1.22.1

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


To build and run manually:

### WINDOWS:

Run run-win.bat from console or by double-clicking.
To rebuild, use build-win.bat.

Or rebuild and run manually with:

```
go build -o bin\win-dpc.exe .\cmd\dnd-playlist-controller-win
.\bin\win-dpc.exe
```

### MAC:

Run run-mac.sh from a terminal window.
To rebuild, use build-mac.sh.

Or rebuild and run manually with:

```
go build -a -o bin/mac-dpc ./cmd/dnd-playlist-controller-mac
./bin/mac-dpc
```


First time you run it, it will print a URL in the console window. Copy paste it into your browser and login with your Spotify account. Once done and you see an empty webpage, you can close the browser.

To avoid having to do this every time I save the auth token in a json file in the repo root folder (token.json). It's not very secure, but fine for my use case. If the auth ever expires or breaks for some reason just delete that file
