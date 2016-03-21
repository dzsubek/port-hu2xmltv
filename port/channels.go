package port

import "net/http"
import "encoding/json"
import "time"
import "strings"

const channelsUrl = "http://www.port.hu/pls/w/tv_api.init?i_page_id=1"

type PortInitResponse struct {
	Date time.Time `json:"date"`
	Services struct {
		IsFavoriteAvailable bool `json:"isFavoriteAvailable"`
		IsPersonalSettingsAvailable bool `json:"isPersonalSettingsAvailable"`
		IsNotificationAvailable bool `json:"isNotificationAvailable"`
	} `json:"services"`
	Channels []struct {
		ID int `json:"id"`
		Article string `json:"article"`
		Name string `json:"name"`
		Link string `json:"link"`
		Logo string `json:"logo"`
		GroupName string `json:"groupName"`
		GroupID int `json:"groupId"`
		SponzorationColor interface{} `json:"sponzoration_color"`
		Language string `json:"language"`
		Address string `json:"address"`
		Phone string `json:"phone"`
		Fax interface{} `json:"fax"`
		Web string `json:"web"`
		Email string `json:"email"`
		Description interface{} `json:"description"`
	} `json:"channels"`
	Favorite interface{} `json:"favorite"`
	Days []int `json:"days"`
	ShowType []struct {
		ID int `json:"id"`
		Name string `json:"name"`
	} `json:"showType"`
	AgeLimit []struct {
		ID int `json:"id"`
		Name string `json:"name"`
		Short string `json:"short"`
		Shorter string `json:"shorter"`
		PictogramURL string `json:"pictogram_url"`
		PictogramURLMobile string `json:"pictogram_url_mobile"`
		DvbRating int `json:"dvb_rating"`
		VideoFlashURL string `json:"video_flash_url"`
	} `json:"ageLimit"`
	SoundQuality []struct {
		ID int `json:"id"`
		Name string `json:"name"`
		Short string `json:"short"`
	} `json:"soundQuality"`
}

type Services struct {
	isFavoriteAvailable bool
	isNotificationAvailable bool
	isPersonalSettingsAvailable bool
}

func GetChannelIds(channels []string) []int {
	init := GetPortInit()

    var result []int
    for _,v := range init.Channels {
    	if (InArray(v.Name, channels)) {
    		result = append(result, v.ID)
    	}
    }

    return result
}

func GetPortInit() PortInitResponse {
	response, err := http.Get(channelsUrl)
	if err != nil {
        panic(err)
    }
    defer response.Body.Close()
    
    var init PortInitResponse
    err = json.NewDecoder(response.Body).Decode(&init)
    if err != nil {
        panic(err)
    }

    return init
}

func InArray(search string, list []string) bool {
	for _, v := range list {
        if strings.ToLower(v) == strings.ToLower(search) {
            return true
        }
    }
    return false
}