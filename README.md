# mpd-mcp-server

`mpd-mcp-server` is a server that integrates MPD (Music Player Daemon) with MCP (Model Context Protocol). This project provides MPD operations as MCP tools, supporting features like music playback and playlist management.

## Features

- **MPD Integration**: Connects to an MPD server to perform operations like play, pause, skip, and manage playlists.
- **MCP Tool Support**: Exposes MPD operations as tools using the MCP protocol.
- **Flexible Configuration**: Allows configuration of the MPD server address and port using environment variables.

## Prerequisites

- Go 1.24.1 or later
- MPD server
- MCP protocol-compatible client

## Installation

```bash
go install github.com/gamoutatsumi/mpd-mcp-server@latest
```

## Usage

```bash
# Set environment variables (if needed)
export MPD_SERVER=localhost
export MPD_PORT=6600

# Start the server
mpd-mcp-server
```

## Available Tools

- `search`: Search for songs in the MPD database.
- `play`: Play a song from the playlist.
- `pause`: Pause the current song.
- `stop`: Stop the current song.
- `resume`: Resume the paused song.
- `next`: Skip to the next song.
- `previous`: Skip to the previous song.
- `get_status`: Get the current status of the MPD server.
- `get_current_song`: Get the currently playing song.
- `get_playlist`: Get the current playlist.
- `clear_playlist`: Clear the current playlist.
- `add_playlist`: Add a song to the current playlist.

## License

This project is licensed under the MIT License.

## Contributing

For bug reports or feature requests, please use GitHub Issues. Pull requests are also welcome.

---
