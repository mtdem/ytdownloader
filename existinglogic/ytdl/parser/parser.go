package parser

import (
	"errors"
	"net/url"
	"ytdl/models"
)

const audioQualityMedium string = "AUDIO_QUALITY_MEDIUM"
const videoQualityHigh string = "hd720"
const videoQualityMedium string = "medium"
const videoQualityTiny string = "tiny"
const prefixLong string = "https://www.youtube.com/"
const prefixShort string = "https://youtu.be/"
const progressBarWidth int = 40

func ParseVideoUrl(link string) (*models.VideoLinkParsed, error) {
	youtubeUrl, err := url.ParseRequestURI(link)

	if err != nil {
		return nil, err
	}

	queryParams, err := url.ParseQuery(youtubeUrl.RawQuery)
	if err != nil {
		return nil, err
	}

	var parsedVideoObject models.VideoLinkParsed

	parsedVideoObject.Url = link
	if queryParams.Has("v") {
		parsedVideoObject.VideoId = queryParams.Get("v")
	}

	if queryParams.Has("ab_channel") {
		parsedVideoObject.Channel = queryParams.Get("ab_channel")
	}

	if queryParams.Has("list") {
		parsedVideoObject.Playlist = queryParams.Get("list")
	}

	if queryParams.Has("index") {
		parsedVideoObject.PlaylistIndex = queryParams.Get("index")
	}

	return &parsedVideoObject, nil
}

func ParsePlaylistUrl(link string) (*models.PlaylistLinkParsed, error) {

	youtubeUrl, err := url.ParseRequestURI(link)

	if err != nil {
		return nil, err
	}

	queryParams, err := url.ParseQuery(youtubeUrl.RawQuery)
	if err != nil {
		return nil, err
	}

	var parsedPlaylistObject models.PlaylistLinkParsed

	parsedPlaylistObject.Url = link
	if queryParams.Has("list") {
		parsedPlaylistObject.PlaylistId = queryParams.Get("list")
	}

	return &parsedPlaylistObject, nil
}

func ConvertVideoLinkToPlaylistLink(videoInPlaylistLink string) (string, error) {
	parsedLink, err := url.Parse(videoInPlaylistLink)
	if err != nil {
		return "", err
	}

	baseDomain := parsedLink.Scheme + "://" + parsedLink.Host
	// endPoint := parsedLink.Path
	queryParams := parsedLink.Query()

	if !queryParams.Has("list") {
		return "", errors.New("videoUrl passed in does not contain the 'link' parameter")
	}

	playListId := queryParams.Get("list")

	playlistLink := baseDomain + "/playlist" + "?" + "list=" + playListId

	return playlistLink, nil
}
