package port

import "fmt"
import "net/http"
import "encoding/json"
import "strings"
import "time"

const showsUrl = "http://port.hu/tvapi?channel_id[]=%s&date=%s"

type PortChannelDetails struct {
	Date string `json:"date"`
	DateFrom time.Time `json:"date_from"`
	DateTo time.Time `json:"date_to"`
	Channels []struct {
		ID string `json:"id"`
		Programs []struct {
			ID string `json:"id"`
			StartDatetime time.Time `json:"start_datetime"`
			StartTime string `json:"start_time"`
			StartTs int `json:"start_ts"`
			EndTime string `json:"end_time"`
			EndDatetime time.Time `json:"end_datetime"`
			IsChildEvent bool `json:"is_child_event"`
			Title string `json:"title"`
			SoundQuality interface{} `json:"sound_quality"`
			Italics interface{} `json:"italics"`
			EpisodeTitle interface{} `json:"episode_title"`
			Description interface{} `json:"description"`
			ShortDescription string `json:"short_description"`
			Highlight interface{} `json:"highlight"`
			IsRepeat bool `json:"is_repeat"`
			IsOverlapping bool `json:"is_overlapping"`
			FilmID string `json:"film_id"`
			FilmURL string `json:"film_url"`
			FavoriteURL string `json:"favorite_url"`
			DelCalendarURL string `json:"del_calendar_url"`
			HasReminder bool `json:"has_reminder"`
			ShowReminder bool `json:"show_reminder"`
			IsNotified bool `json:"is_notified"`
			ShowNotification bool `json:"show_notification"`
			MediaURL string `json:"media_url"`
			Media interface{} `json:"media"`
			HasVideo bool `json:"has_video"`
			AttributesText string `json:"attributes_text"`
			OuterLinks struct {
				FilmOuterLink interface{} `json:"film_outer_link"`
				WatchMovieLink interface{} `json:"watch_movie_link"`
				ExtraLink interface{} `json:"extra_link"`
			} `json:"outer_links"`
			Restriction struct {
				AgeLimit int `json:"age_limit"`
				Category string `json:"category"`
			} `json:"restriction"`
			Type string `json:"type"`
			IsLive bool `json:"is_live,omitempty"`
		} `json:"programs"`
		Article string `json:"article"`
		Name string `json:"name"`
		Domain string `json:"domain"`
		URL string `json:"url"`
		Logo string `json:"logo"`
		StreamURL interface{} `json:"stream_url"`
		StreamCtLinkurl interface{} `json:"stream_ct_linkurl"`
		Banners interface{} `json:"banners"`
		DateFrom time.Time `json:"date_from"`
		DateUntil time.Time `json:"date_until"`
		Cache string `json:"cache"`
	} `json:"channels"`
}

type Show struct {
	Title string
	Start time.Time 
	End time.Time
	Description string
	Url string
}

func GetShowsByChannelIds(channelIds []string) map[string][]Show {
	result := make(map[string][]Show)

	for _, channelId := range channelIds {
		data := GetPortShows(channelId)
		for _, channel := range data.Channels {
			for _, program := range channel.Programs {
				result[channel.Name] = append(result[channel.Name], Show{
						Title: program.Title,
						Description: strings.TrimSpace(program.ShortDescription),
						Url: program.FilmURL,
						Start: program.StartDatetime,
						End: program.EndDatetime,
					})
			}
		}
		time.Sleep(2 * time.Second)
	}

	return result
}

func GetPortShows(channelId string) PortChannelDetails {
	now := time.Now()
	showsUri := fmt.Sprintf(showsUrl, channelId, now.Format("2006-01-02"))

	response, err := http.Get(showsUri)
	if err != nil {
        panic(err)
    }
    defer response.Body.Close()

    var channel PortChannelDetails
	err = json.NewDecoder(response.Body).Decode(&channel)
    if err != nil {
		fmt.Println(showsUri)
        panic(err)
    }

    return channel
}
