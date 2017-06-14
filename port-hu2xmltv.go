package main

import (
	"encoding/xml";
	"os";
	"io";
	"github.com/tv42/slug"
	"port-hu2xmltv/port"
)

func main() {
	type Channel struct {
		Id string `xml:"id,attr"`
		DisplayName string `xml:"display-name"`
	}

	type Programme struct {
		Start string `xml:"start,attr"`
		Stop string `xml:"stop,attr"`
		ChannelId string `xml:"channel,attr"`
		Title string `xml:"title"`
		Description string `xml:"desc"`
		Url string `xml:"url"`
	}

	type XMLTv struct {
		XMLName xml.Name `xml:"tv"`
		Url string `xml:"source-info-url,attr"`
		Name string `xml:"source-info-name,attr"`
		Generator string `xml:"generator-info-name,attr"`
		Channel []Channel `xml:"channel"`
		Programme []Programme `xml:"programme"`
	}

	if len(os.Args) < 3 {
		panic("Usage: port-hu2xmltv OUTPUT_FILE CHANNEL_NAME1 [CHANNEL_NAME2, ...]")
	}

	channelNames := os.Args[2:]
	channelIds := port.GetChannelIds(channelNames)
	shows := port.GetShowsByChannelIds(channelIds)

	xmltv := &XMLTv{}
	xmltv.Url = "http://xmltv.dzsubek.info"
	xmltv.Name = "Dzsubek Port Grab"
	xmltv.Generator = "xmltv/0.1-dzsubek-port-go"

	for channelName, programs := range shows {
		channel :=  Channel{
			Id: slug.Slug(channelName) + ".xmltv.dzsubek.info",
			DisplayName: channelName,
		}
		xmltv.Channel = append(xmltv.Channel, channel)

		for _, program := range programs {
			xmltv.Programme = append(xmltv.Programme, Programme{
				Start: program.Start.Format("20060102150400"),
				Stop: program.End.Format("20060102150400"),
				ChannelId: channel.Id,
				Title: program.Title,
				Description: program.Description,
				Url: program.Url,
				})
		}
	}

	file, _ := os.Create(os.Args[1])
	xmlWriter := io.Writer(file)
	xmlWriter.Write([]byte(xml.Header))
	xmlEncoder := xml.NewEncoder(xmlWriter)
	xmlEncoder.Indent("", "    ")
	if err := xmlEncoder.Encode(xmltv); err != nil {
		panic(err)
	}
}