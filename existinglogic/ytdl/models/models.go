package models

type VideoLinkParsed struct {
	Url           string
	VideoId       string
	Playlist      string
	PlaylistIndex string
	Channel       string
}

type PlaylistLinkParsed struct {
	Url        string
	PlaylistId string
}
