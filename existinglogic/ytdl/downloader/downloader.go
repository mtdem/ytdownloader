package downloader

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"ytdl/parser"
	"ytdl/rootpath"

	"github.com/kkdai/youtube/v2"
)

// TODO: meant to download multiple videos CONCURRENTLY
func DownloadVideos(links []string, outputDir string) []error {
	return []error{}
}

// There are cases where the link opens is already in a playlist
// In that scenario we can either both extract a playlist, or just the video in the playlist
func DownloadLink(link string, outputDir string, extractVideoOnly bool, includeAuthor bool) []error {
	videoLinkParsed, err := parser.ParseVideoUrl(link)
	if err != nil {
		return []error{err}
	}

	// early return check
	if len(videoLinkParsed.Playlist) == 0 {
		fmt.Println("video link does not contain a playlist id. force downloading a video...")
		err = DownloadVideo(link, outputDir, includeAuthor)
		if err != nil {
			return []error{err}
		}
		return nil
	}

	// video only option
	if extractVideoOnly {
		fmt.Println("extracting only the video")
		err = DownloadVideo(link, outputDir, includeAuthor)
		if err != nil {
			return []error{err}
		}
		return []error{}
	}

	// playlist option
	fmt.Println("extracting entire playlist")
	playlistLink, err := parser.ConvertVideoLinkToPlaylistLink(link)
	if err != nil {
		return []error{err}
	}
	//fmt.Println(fmt.Sprintf("playlist link: '%s'", playlistLink))

	errs := DownloadPlaylist(playlistLink, outputDir, includeAuthor)
	if len(errs) > 0 {
		return errs
	}
	return []error{}
}

// single execution of a video download process
func DownloadVideo(link string, outputDir string, includeAuthor bool) error {

	videoLinkParsed, err := parser.ParseVideoUrl(link)
	if err != nil {
		return err
	}

	client := youtube.Client{}

	video, err := client.GetVideo(videoLinkParsed.VideoId)
	if err != nil {
		return err
	}

	formats := video.Formats.WithAudioChannels()

	// TODO: maybe a better way to utilize the stream?
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		return err
	}

	// defer means handle it after function executes
	defer stream.Close()

	videoName, err := createVideoName(video, true, includeAuthor)

	fmt.Printf("\tDownloading %s...", videoName)
	outputPath := filepath.Join(outputDir, videoName)
	file, err := os.Create(outputPath)

	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, stream)

	if err != nil {
		return err
	}

	fmt.Printf("\tDownloaded %s successfully", videoName)
	return nil
}

// single execution of a playlist download
func DownloadPlaylist(link string, outputDir string, includeAuthor bool) []error {

	client := youtube.Client{}
	playlistLinkParsed, err := parser.ParsePlaylistUrl(link)
	if err != nil {
		return []error{err}
	}

	playlist, err := client.GetPlaylist(playlistLinkParsed.PlaylistId)

	if err != nil {
		return []error{err}
	}

	// Enumerate Playlist videos
	playlistFolderName, err := createPlaylistName(playlist, includeAuthor)
	playlistPath := filepath.Join(outputDir, playlistFolderName)
	err = rootpath.CreateDirectoryIfNotExists(playlistPath)

	if err != nil {
		panic(err)
	}

	header := fmt.Sprintf("Playlist: %s", playlistFolderName)
	fmt.Println(header)
	fmt.Println(strings.Repeat("=", len(header)) + "\n")

	fmt.Printf("Downloading Playlist to: %s...\n\n", playlistPath)

	var errors []error
	for index, entry := range playlist.Videos {

		// TODO: make this multithreaded
		err = downloadVideoForPlaylist(&client, index, entry, playlistPath, includeAuthor)
		if err != nil {
			errors = append(errors, err)
		}
	}

	fmt.Printf("\nFinished Playlist download to %s\n\n", playlistPath)
	if len(errors) > 0 {
		return errors
	}

	return []error{}
}

// This is split into its own separate function for robustness
func downloadVideoForPlaylist(client *youtube.Client,
	index int,
	entry *youtube.PlaylistEntry,
	playlistPath string, includeAuthor bool) error {

	displayIndex := index + 1

	fmt.Printf("\t(%d) Accessing video for entry: %s - %s\n", displayIndex, entry.Title, entry.Author)
	video, err := client.VideoFromPlaylistEntry(entry)
	if err != nil {
		printError(err, displayIndex, true)
		return err
	}

	format := &video.Formats[0]
	stream, _, err := client.GetStream(video, format)

	if err != nil {
		printError(err, displayIndex, true)
		return err
	}
	fmt.Printf("\t(%d) Video accessed for entry: %s - %s\n", displayIndex, entry.Title, entry.Author)

	fileName, err := createVideoName(video, true, includeAuthor)
	if err != nil {
		printError(err, displayIndex, true)
		return err
	}

	videoFilePath := filepath.Join(playlistPath, fileName)
	file, err := os.Create(videoFilePath)
	if err != nil {
		return err
	}

	fmt.Printf("\t(%d) Downloading '%s'\n", displayIndex, fileName)
	defer file.Close()
	_, err = io.Copy(file, stream)
	if err != nil {
		os.Remove(videoFilePath)
		printError(err, displayIndex, true)
		return err
	}

	fmt.Printf("\t(%d) Downloaded '%s'\n", displayIndex, fileName)

	return nil
}

func createPlaylistName(pl *youtube.Playlist, includeAuthor bool) (string, error) {
	plName := pl.Title

	if includeAuthor && pl.Author != "" {
		plName = plName + " - " + pl.Author
	}

	plName = strings.Replace(plName, "/", "", -1)
	plName = strings.Replace(plName, "\\", "", -1)

	plName = rootpath.RemoveInvalidFileNameChars(plName)
	// sanitize folder name??
	return plName, nil
}

// Creates the file name for the video
func createVideoName(vid *youtube.Video, includeExt bool, includeAuthor bool) (string, error) {
	videoStr := vid.Title

	if includeAuthor && vid.Author != "" {
		videoStr = videoStr + " - " + vid.Author
	}

	if includeExt {
		videoStr = videoStr + ".mp4"
	}

	videoStr = strings.Replace(videoStr, "/", "", -1)
	videoStr = strings.Replace(videoStr, "\\", "", -1)

	videoStr = rootpath.RemoveInvalidFileNameChars(videoStr)
	// sanitize folder name??
	return videoStr, nil
}

func printError(err error, index int, indent bool) {
	var errMessage strings.Builder

	if indent {
		errMessage.WriteString("\t")
	}

	if index > 0 {
		errMessage.WriteString(fmt.Sprintf("(%d) ", index))
	}

	errMessage.WriteString("ERROR: ")
	errMessage.WriteString(err.Error())

	fmt.Println(errMessage.String())
}
