// Package util
package util

import (
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHttpJson(t *testing.T) {
	v := new(cyeamSearch)
	err := HttpJson("GET", cyeamSearchURL, "", nil, &v)
	assert.Nil(t, err)
	assert.True(t, len(v.Docs) > 0)
}

const cyeamSearchURL = `http://www.cyeam.com/s?t=golang`

type cyeamSearch struct {
	Docs []struct {
		Des    string `json:"des"`
		Figure string `json:"figure"`
		Link   string `json:"link"`
		Title  string `json:"title"`
	} `json:"docs"`
	Summary struct {
		D   int64  `json:"d"`
		Num int64  `json:"num"`
		Q   string `json:"q"`
	} `json:"summary"`
}

func TestHttpXml(t *testing.T) {
	v := new(cyeamBlog)
	err := HttpXml("GET", cyeamBlogURL, "", nil, &v)
	assert.Nil(t, err)
	assert.True(t, len(v.Channel.Item) > 0)
}

const cyeamBlogURL = `http://blog.cyeam.com/rss.xml`

type cyeamBlog struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Channel struct {
		Text  string `xml:",chardata"`
		Title struct {
			Text string `xml:",chardata"`
		} `xml:"title"`
		Description struct {
			Text string `xml:",chardata"`
		} `xml:"description"`
		Link []struct {
			Text string `xml:",chardata"`
		} `xml:"link"`
		LastBuildDate struct {
			Text string `xml:",chardata"`
		} `xml:"lastBuildDate"`
		PubDate struct {
			Text string `xml:",chardata"`
		} `xml:"pubDate"`
		Ttl struct {
			Text string `xml:",chardata"`
		} `xml:"ttl"`
		Item []struct {
			Text  string `xml:",chardata"`
			Title struct {
				Text string `xml:",chardata"`
			} `xml:"title"`
			Figure struct {
				Text string `xml:",chardata"`
			} `xml:"figure"`
			Info struct {
				Text string `xml:",chardata"`
			} `xml:"info"`
			Description struct {
				Text string `xml:",chardata"`
			} `xml:"description"`
			Link struct {
				Text string `xml:",chardata"`
			} `xml:"link"`
			Guid struct {
				Text string `xml:",chardata"`
			} `xml:"guid"`
			PubDate struct {
				Text string `xml:",chardata"`
			} `xml:"pubDate"`
		} `xml:"item"`
	} `xml:"channel"`
}
