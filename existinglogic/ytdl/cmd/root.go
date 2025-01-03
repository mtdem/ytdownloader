package cmd

import (
	"fmt"
	"os"

	"github.com/gosuri/uiprogress"
	"github.com/spf13/cobra"

	"ytdl/downloader"
	"ytdl/video"
)

/*
TODO: listed below
	0. get solution operational
	1. support single video, multi video and playlist
	2. support to output dir
	3. MAYBE? support for mp4 format
*/

var links []string
var destination string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:       "ytdl",
	Short:     "Command line tool for converting YouTube videos to mp3/mp4 files.",
	Long:      "Command line tool for converting YouTube videos to mp3/mp4 files. ",
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{"names", "links", "ytdl"},
	Run: func(cmd *cobra.Command, args []string) {
		// Init progress bar.
		uiprogress.Start()
		// Handle links.
		errs := handleLinks(cmd, links)
		// Stdout errors and exit.
		if len(errs) > 0 {
			cmd.Printf("\nThe following issues occurred during execution:\n")
			for _, err := range errs {
				cmd.Printf(" - %v\n", err)
			}
			cmd.Printf(
				"\nErrors encountered\n",
			)
			os.Exit(1)
		}
	},
}

func init() {
	// Define flags.
	rootCmd.Flags().StringSliceVarP(
		&links, "links", "l", []string{},
		"List of YouTube video links which will be converted to mp3 and saved on your local.",
	)
	rootCmd.MarkFlagRequired("links")
	workingDir, _ := os.Getwd()
	rootCmd.Flags().StringP(
		"dst", "d", workingDir,
		"Output directory for mp3 files.",
	)
}

// Execute This is called by main.main().
// It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// handleLinks makes next magic:
// - Validates incoming links.
// - Gets playback streams:
//   - If video is not well-protected get stream url using regex.
//   - If video is well-protected get stream url using python port of youtube-dl.
//
// - Fetches metadata for video.
// - Downloads videos and saves them in temp files.
// - ffmpeg magic.
// - Cleans up tmp files.
func handleLinks(cmd *cobra.Command, links []string) []error {
	// Validate links. If at least one link is not valid we stop an execution.
	errs := video.ValidateLinks(links)
	if len(errs) > 0 {
		return errs
	}

	var errors []error

	videos := fetchPlaybackURLS(links, &errors)

	// Fetch metadata.
	channelFetchMetadata := make(chan video.ChannelMessage, len(videos))
	for _, _video := range videos {
		go video.FetchMetadata(_video, channelFetchMetadata)
	}
	for i := 0; i < len(videos); i++ {
		msg := <-channelFetchMetadata
		if msg.Error != nil {
			errors = append(errors, msg.Error)
		}
	}
	close(channelFetchMetadata)

	// Download and save temp video files.
	channelFetchVideo := make(chan video.ChannelMessage, len(videos))
	for _, _video := range videos {
		go video.FetchVideo(_video, channelFetchVideo)
	}
	for i := 0; i < len(videos); i++ {
		msg := <-channelFetchVideo
		if msg.Error != nil {
			errors = append(errors, msg.Error)
		}
	}
	// Cleanup file when main function is over.
	defer func(videos []*video.Video) {
		for _, v := range videos {
			err := os.Remove((*v.File).Name())
			if err != nil {
				errors = append(errors, err)
			}
		}
	}(videos)

	// Run ffmpeg and convert videos to mp3 files.
	dstDir, _ := cmd.Flags().GetString("dst")
	channelConvertVideoToAudio := make(chan video.ChannelMessage, len(videos))
	for _, _video := range videos {
		go video.ConvertVideoToAudio(_video, dstDir, channelConvertVideoToAudio)
	}
	for i := 0; i < len(videos); i++ {
		msg := <-channelConvertVideoToAudio
		if msg.Error != nil {
			errors = append(errors, msg.Error)
		}
	}

	return errors
}

// allow video.Video to support video.Playlist
func fetchPlaybackURLS(links []string, errors *[]error) []*video.Video {

	var videos []*video.Video

	channelFetchPlaybackURL := make(chan video.ChannelMessage, len(links))

	for _, link := range links {
		// Start go runtime thread.
		go video.FetchPlaybackURL(link, channelFetchPlaybackURL)
	}

	for i := 0; i < len(links); i++ {
		// Wait until all threads are done.
		msg := <-channelFetchPlaybackURL
		if msg.Error != nil {
			*errors = append(*errors, msg.Error)
		} else if msg.Result.HasStreamURL() {
			videos = append(videos, msg.Result)
		}
	}
	close(channelFetchPlaybackURL)

	return videos
}

func RunTestDownload() {
	// vidlink := "https://www.youtube.com/watch?v=N4-Sk506pgA&list=PLWU88pbzc3rMYGfcVxDVvsQIqTsizJ0YL&index=154&ab_channel=victorpardo"
	playlistLink := "https://www.youtube.com/watch?v=_i1wNr2gEVo&list=PLWU88pbzc3rMYGfcVxDVvsQIqTsizJ0YL&ab_channel=rumworld"
	testOutputDir := "/home/mo/dev/goproj/ytdownloader/testing"
	videoOnly := false
	includeAuthor := false

	errs := downloader.DownloadLink(playlistLink, testOutputDir, videoOnly, includeAuthor)
	if len(errs) > 0 {
		fmt.Println("errors encountered")
		for i, err := range errs {
			fmt.Printf("\terror (%d): '%s'\n", i+1, err.Error())
		}
	}
}
