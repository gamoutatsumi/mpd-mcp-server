package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/fhs/gompd/v2/mpd"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var version string

func newServer() *server.MCPServer {
	s := server.NewMCPServer("mpd-mcp-server", version, server.WithResourceCapabilities(true, true), server.WithLogging())

	searchTool := mcp.NewTool("search",
		mcp.WithDescription("Search for songs in the MPD database."),
		mcp.WithString("artist", mcp.Description("The artist to search for.")),
		mcp.WithString("album", mcp.Description("The album to search for.")),
		mcp.WithString("title", mcp.Description("The title to search for.")),
		mcp.WithString("genre", mcp.Description("The genre to search for.")),
		mcp.WithString("album_artist", mcp.Description("The album artist to search for.")),
	)
	playTool := mcp.NewTool("play",
		mcp.WithDescription("Play a song from the playlist."),
		mcp.WithNumber("pos", mcp.Description("The position of the song in the playlist.")),
	)
	pauseTool := mcp.NewTool("pause",
		mcp.WithDescription("Pause the current song."),
	)
	stopTool := mcp.NewTool("stop",
		mcp.WithDescription("Stop the current song."),
	)
	resumeTool := mcp.NewTool("resume",
		mcp.WithDescription("Resume the current song."),
	)
	nextTool := mcp.NewTool("next",
		mcp.WithDescription("Skip to the next song."),
	)
	previousTool := mcp.NewTool("previous",
		mcp.WithDescription("Skip to the previous song."),
	)
	getStatusTool := mcp.NewTool("get_status",
		mcp.WithDescription("Get the current status of the MPD server."),
	)
	getCurrentSongTool := mcp.NewTool("get_current_song",
		mcp.WithDescription("Get the current song being played."),
	)
	getPlaylistTool := mcp.NewTool("get_playlist",
		mcp.WithDescription("Get the current playlist."),
	)
	clearPlaylistTool := mcp.NewTool("clear_playlist",
		mcp.WithDescription("Clear the current playlist."),
	)
	addPlaylistTool := mcp.NewTool("add_playlist",
		mcp.WithDescription("Add a song to the current playlist."),
		mcp.WithString("uri", mcp.Description("The URI of the song to add.")),
	)

	s.AddTool(playTool, playHandler)
	s.AddTool(pauseTool, pauseHandler)
	s.AddTool(stopTool, stopHandler)
	s.AddTool(nextTool, nextHandler)
	s.AddTool(previousTool, previousHandler)
	s.AddTool(resumeTool, resumeHandler)
	s.AddTool(searchTool, searchToolHandler)
	s.AddTool(getStatusTool, getStatusHandler)
	s.AddTool(getCurrentSongTool, getCurrentSongHandler)
	s.AddTool(getPlaylistTool, getPlaylistHandler)
	s.AddTool(clearPlaylistTool, clearPlaylistHandler)
	s.AddTool(addPlaylistTool, addPlaylistHandler)
	return s

}

func run() error {
	mpdClient, err := connectMPD()
	if err != nil {
		return err
	}
	mpdClient.Close()

	s := newServer()
	if err := server.ServeStdio(s); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func searchToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	mpdClient, err := connectMPD()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to connect to MPD server: %v", err)), nil
	}
	defer mpdClient.Close()
	artist, _ := request.Params.Arguments["artist"].(string)
	album, _ := request.Params.Arguments["album"].(string)
	title, _ := request.Params.Arguments["title"].(string)
	genre, _ := request.Params.Arguments["genre"].(string)
	albumArtist, _ := request.Params.Arguments["album_artist"].(string)
	queries := []string{}
	if artist != "" {
		queries = append(queries, "artist", artist)
	}
	if album != "" {
		queries = append(queries, "album", album)
	}
	if title != "" {
		queries = append(queries, "title", title)
	}
	if genre != "" {
		queries = append(queries, "genre", genre)
	}
	if albumArtist != "" {
		queries = append(queries, "albumartist", albumArtist)
	}
	result, err := mpdClient.Search(queries...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to search: %v", err)), nil
	}
	m, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(fmt.Sprintln(string(m))), nil
}

func playHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	mpdClient, err := connectMPD()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to connect to MPD server: %v", err)), nil
	}
	defer mpdClient.Close()
	pos, ok := request.Params.Arguments["pos"].(float64)
	if !ok {
		return mcp.NewToolResultError("pos must be float64"), nil
	}
	if err := mpdClient.Play(int(pos)); err != nil {

		return mcp.NewToolResultError(fmt.Sprintf("failed to play song: %v", err)), nil
	}
	return mcp.NewToolResultText("Playing song with"), nil
}

func pauseHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	mpdClient, err := connectMPD()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to connect to MPD server: %v", err)), nil
	}
	defer mpdClient.Close()
	if err := mpdClient.Pause(true); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to pause: %v", err)), nil
	}
	return mcp.NewToolResultText("Paused playback"), nil
}

func stopHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	mpdClient, err := connectMPD()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to connect to MPD server: %v", err)), nil
	}
	defer mpdClient.Close()
	if err := mpdClient.Stop(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to stop: %v", err)), nil
	}
	return mcp.NewToolResultText("Stopped playback"), nil
}

func nextHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	mpdClient, err := connectMPD()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to connect to MPD server: %v", err)), nil
	}
	defer mpdClient.Close()
	if err := mpdClient.Next(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to skip to next track: %v", err)), nil
	}
	return mcp.NewToolResultText("Skipped to next track"), nil
}

func previousHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	mpdClient, err := connectMPD()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to connect to MPD server: %v", err)), nil
	}
	defer mpdClient.Close()
	if err := mpdClient.Previous(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to skip to previous track: %v", err)), nil
	}
	return mcp.NewToolResultText("Skipped to previous track"), nil
}

func resumeHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	mpdClient, err := connectMPD()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to connect to MPD server: %v", err)), nil
	}
	defer mpdClient.Close()
	if err := mpdClient.Pause(false); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to resume: %v", err)), nil
	}
	return mcp.NewToolResultText("Resumed playback"), nil
}

func getStatusHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	mpdClient, err := connectMPD()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to connect to MPD server: %v", err)), nil
	}
	defer mpdClient.Close()
	result, err := mpdClient.Status()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get status: %v", err)), nil
	}
	m, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(fmt.Sprintln(string(m))), nil
}

func getCurrentSongHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	mpdClient, err := connectMPD()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to connect to MPD server: %v", err)), nil
	}
	defer mpdClient.Close()
	result, err := mpdClient.CurrentSong()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get current song: %v", err)), nil
	}
	m, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(fmt.Sprintln(string(m))), nil
}

func getPlaylistHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	mpdClient, err := connectMPD()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to connect to MPD server: %v", err)), nil
	}
	defer mpdClient.Close()
	result, err := mpdClient.PlaylistInfo(-1, -1)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get playlist: %v", err)), nil
	}
	m, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(fmt.Sprintln(string(m))), nil
}

func clearPlaylistHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	mpdClient, err := connectMPD()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to connect to MPD server: %v", err)), nil
	}
	defer mpdClient.Close()
	if err := mpdClient.Clear(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to clear playlist: %v", err)), nil
	}
	return mcp.NewToolResultText("Cleared playlist"), nil
}

func addPlaylistHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	mpdClient, err := connectMPD()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to connect to MPD server: %v", err)), nil
	}
	defer mpdClient.Close()
	song, ok := request.Params.Arguments["uri"].(string)
	if !ok {
		return mcp.NewToolResultError("song must be string"), nil
	}
	if err := mpdClient.Add(song); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to add song to playlist: %v", err)), nil
	}
	return mcp.NewToolResultText("Added song to playlist"), nil
}

func connectMPD() (*mpd.Client, error) {
	server, ok := os.LookupEnv("MPD_SERVER")
	if !ok {
		server = "localhost"
	}
	port, ok := os.LookupEnv("MPD_PORT")
	if !ok {
		port = "6600"
	}
	mpdClient, err := mpd.Dial("tcp", fmt.Sprintf("%s:%s", server, port))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MPD server: %v", err)
	}
	return mpdClient, nil
}
