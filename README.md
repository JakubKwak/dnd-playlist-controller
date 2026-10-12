# dnd-playlist-controller
Small app to allow playlist switching using system hotkeys for running my DND campaign.

## Requirements
- Windows or Mac (Linux currently not supported)
- A Spotify Premium account (required by Spotify's API to control playback)

## Usage
1. Download the latest release for your OS from the **Releases** page and extract the `.zip`.
2. Edit `hotkeys.json` (see `hotkeys-example.json`) in the same folder as the executable to bind hotkeys to your Spotify playlists:
```json
[
    {
        "uri": "4pxn3LHWRCd...", /* Spotify playlist ID */
        "modifiers": [ /* Supported: Ctrl, Shift, Alt / Option, Win / Cmd */
            "Ctrl",
            "Alt",
            "Shift"
        ],
        "key": "6" /* Main key: 0-9, A-Z, or F1-F20 */
    }
]
```
3. Run the executable (`dnd-playlist-controller.exe` on Windows or `./dnd-playlist-controller` on Mac).

The first time you run the executable it will print a URL in the console window. Copy + paste it into your browser and log in with your Spotify account. Once done and you see the confirmation webpage, you can close the browser.

To avoid having to do this every time, the auth token is saved in a JSON file alongside the executable (`token.json`). If the auth ever breaks or you wish to switch to a different account, just delete that file.

---

## Building the Client from Source (Optional)
If you want to compile the client binary yourself, you'll need Go 1.22.1+:

### Windows:
```powershell
go build -o bin\dnd-playlist-controller.exe .\cmd\client
.\bin\dnd-playlist-controller.exe
```

### Mac:
```sh
go build -o bin/dnd-playlist-controller ./cmd/client
./bin/dnd-playlist-controller
```

---

## Self-Hosting the Auth Server (Optional)
By default, the client uses the official hosted authentication server, so you **do not** need to run the server or create a Spotify Developer application to use the app.

If you want to host your own OAuth relay server:

1. Create a Spotify application in the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard) and add your server's callback URL (e.g. `https://your-domain.com/callback`) to its Redirect URIs.
2. On your server, create a `.env` file (see `.env.example`) or pass the following environment variables:
   ```env
   SPOTIFY_ID=your_spotify_client_id
   SPOTIFY_SECRET=your_spotify_client_secret
   REDIRECT_URI=https://your-domain.com/callback
   LISTEN_ADDR=:8080
   ```
3. Run the server using Docker Compose (or build the Go binary directly):
   ```sh
   docker build -t dnd-playlist-controller-server:latest .
   docker compose up -d
   ```
   *(Or without Docker: `go build -o bin/auth-server ./cmd/server && ./bin/auth-server`)*
4. On your local machine, create a `.env` file alongside the client executable to point it at your server:
   ```env
   AUTH_SERVER_URL=https://your-domain.com
   ```
