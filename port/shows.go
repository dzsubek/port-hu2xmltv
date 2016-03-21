package port

import "fmt"
import "net/http"
import "encoding/json"
import "time"
import "strings"
import "github.com/bjarneh/latinx"

const showsUrl = "http://port.hu/pls/w/tv_api.event_list?i_channel_id=%d&i_datetime_from=%s&i_datetime_to=%s"

type PortChannelDetails struct {
	Date time.Time `json:"date"`
	DatetimeFrom string `json:"datetime_from"`
	DatetimeTo string `json:"datetime_to"`
	Channels []struct {
		Cache string `json:"cache"`
		Article string `json:"article"`
		Name string `json:"name"`
		Domain string `json:"domain"`
		Logo string `json:"logo"`
		StreamURL interface{} `json:"stream_url"`
		StreamCtLinkurl interface{} `json:"stream_ct_linkurl"`
		Banners interface{} `json:"banners"`
		Capture string `json:"capture"`
		Programs []struct {
			ID int `json:"id"`
			StartDatetime time.Time `json:"start_datetime"`
			StartTime string `json:"start_time"`
			EndTime string `json:"end_time"`
			EndDatetime time.Time `json:"end_datetime"`
			IsChildEvent bool `json:"is_child_event"`
			ShortDescription string `json:"short_description"`
			IsRepeat bool `json:"is_repeat"`
			Highlight string `json:"highlight"`
			Italics string `json:"italics"`
			IsFavorite bool `json:"is_favorite"`
			HasNotification bool `json:"has_notification"`
			FavoriteURL string `json:"favorite_url"`
			FilmID int `json:"film_id"`
			DelCalendarURL string `json:"del_calendar_url"`
			ShowNotification bool `json:"show_notification"`
			ShowFavorite bool `json:"show_favorite"`
			MediaURL string `json:"media_url"`
			SetNotificationURL string `json:"set_notification_url"`
			IsLive bool `json:"is_live"`
			Title string `json:"title"`
			FilmURL string `json:"film_url"`
			EpisodeTitle string `json:"episode_title"`
			SoundQuality interface{} `json:"sound_quality"`
			Description interface{} `json:"description"`
			HasVideo bool `json:"has_video"`
			Attributes []struct {
				ID string `json:"id"`
				Description interface{} `json:"description"`
				Name string `json:"name"`
				AttributeDescription string `json:"attribute_description"`
				AttributePictogram string `json:"attribute_pictogram"`
			} `json:"attributes"`
			AttributesText string `json:"attributes_text"`
			OuterLinks struct {
				FilmOuterLink interface{} `json:"film_outer_link"`
				WatchMovieLink interface{} `json:"watch_movie_link"`
				ExtraLink interface{} `json:"extra_link"`
			} `json:"outer_links"`
			Media interface{} `json:"media"`
			Restriction struct {
				AgeLimit int `json:"age_limit"`
				Category int `json:"category"`
			} `json:"restriction"`
			Type string `json:"type"`
		} `json:"programs"`
	} `json:"channels"`
}

type Show struct {
	Title string
	Start time.Time 
	End time.Time
	Description string
	Url string
}

func GetShowsByChannelIds(channelIds []int) map[string][]Show {
	result := make(map[string][]Show)

	for _, channelId := range channelIds {
		data := GetPortShows(channelId)
		for _, channel := range data {
			for _, program := range channel.Channels[0].Programs {
				result[channel.Channels[0].Name] = append(result[channel.Channels[0].Name], Show{
						Title: program.Title,
						Description: strings.TrimSpace(program.ShortDescription + " " +program.Highlight),
						Url: program.FilmURL,
						Start: program.StartDatetime,
						End: program.EndDatetime,
					})
				
			}
		}
	}

	return result
}

func GetPortShows(channelId int) map[string]PortChannelDetails {
	now := time.Now()
	showsUri := fmt.Sprintf(showsUrl, channelId, now.Format("2006-01-02"), now.Add(7 * 24 * time.Hour).Format("2006-01-02"))

	response, err := http.Get(showsUri)
	if err != nil {
        panic(err)
    }
    defer response.Body.Close()

    BodyReader := latinx.NewReader(latinx.ISO_8859_2, response.Body)
    
    var channel map[string]PortChannelDetails
    err = json.NewDecoder(BodyReader).Decode(&channel)
    if err != nil {
        panic(err)
    }

    return channel
}
