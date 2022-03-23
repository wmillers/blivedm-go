package message

import (
	"bytes"

	"github.com/tidwall/gjson"
)

const (
	TextDanmaku = iota
	EmoticonDanmaku
)

type (
	Danmaku struct {
		Mode      int64
		FontSize  int64
		Color     int64
		Timestamp int64
		Rnd       int64
		UID_CRC32 string
		MsgType   int64
		Bubble    int64

		Msg string

		Uid          int64
		Uname        string
		Admin        bool
		Vip          bool
		Svip         bool
		Urank        int64
		MobileVerify bool
		UnameColor   string

		MedalLevel   int64
		MedalName    string
		MedalUpName  string
		MedalRoomId  int64
		MedalColor   int64
		SpecialMedal string

		UserLevel      int64
		UserLevelColor int64
		UserLevelRank  string

		OldTitle string
		Title    string

		PrivilegeType int64
	}

	CommonNoticeDanmaku struct {
		ContentSegments []struct {
			FontColor string `json:"font_color"`
			Text      string `json:"text"`
			Type      int    `json:"type"`
		} `json:"content_segments"`
		Dmscore   int   `json:"dmscore"`
		Terminals []int `json:"terminals"`
	}
)

func (d *Danmaku) Parse(data []byte) {
	sb := bytes.NewBuffer(data).String()
	Get := func(place string) gjson.Result {
		return gjson.Get(sb, place)
	}
	d = &Danmaku{
		Mode:           Get("info.0.1").Int(),
		FontSize:       Get("info.0.2").Int(),
		Color:          Get("info.0.3").Int(),
		Timestamp:      Get("info.0.4").Int(),
		Rnd:            Get("info.0.5").Int(),
		UID_CRC32:      Get("info.0.7").String(),
		MsgType:        Get("info.0.9").Int(),
		Bubble:         Get("info.0.10").Int(),
		Msg:            Get("info.1").String(),
		Uid:            Get("info.2.0").Int(),
		Uname:          Get("info.2.1").String(),
		Admin:          Get("info.2.2").Bool(),
		Vip:            Get("info.2.3").Bool(),
		Svip:           Get("info.2.4").Bool(),
		Urank:          Get("info.2.5").Int(),
		MobileVerify:   Get("info.2.6").Bool(),
		UnameColor:     Get("info.2.7").String(),
		MedalLevel:     Get("info.3.0").Int(),
		MedalName:      Get("info.3.1").String(),
		MedalUpName:    Get("info.3.2").String(),
		MedalRoomId:    Get("info.3.3").Int(),
		MedalColor:     Get("info.3.4").Int(),
		SpecialMedal:   Get("info.3.5").String(),
		UserLevel:      Get("info.4.0").Int(),
		UserLevelColor: Get("info.4.2").Int(),
		UserLevelRank:  Get("info.4.3").String(),
		OldTitle:       Get("info.5.0").String(),
		Title:          Get("info.5.1").String(),
		PrivilegeType:  Get("info.5").Int(),
	}
}
