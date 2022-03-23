package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/wmillers/blivedm-go/packet"
)

type Client struct {
	conn                *websocket.Conn
	roomID              string
	token               string
	host                string
	eventHandlers       *eventHandlers
	customEventHandlers *customEventHandlers
	Done                chan struct{}
}

func NewClient(roomID string) *Client {
	return &Client{
		roomID:              roomID,
		eventHandlers:       &eventHandlers{},
		customEventHandlers: &customEventHandlers{},
		Done:                make(chan struct{}),
	}
}

func (c *Client) Connect() error {
	if c.host == "" {
		info, err := getDanmuInfo(c.roomID)
		if err != nil {
			return err
		}
		c.host = fmt.Sprintf("wss://%s/sub", info.Data.HostList[0].Host)
		c.token = info.Data.Token
	}
	conn, _, err := websocket.DefaultDialer.Dial(c.host, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) Start() {
	c.sendEnterPacket()
	go func() {
		for {
			select {
			case <-c.Done:
				return
			default:
				msgType, data, err := c.conn.ReadMessage()
				if err != nil {
					_ = c.Connect()
					continue
				}
				if msgType != websocket.BinaryMessage {
					log.Error("packet not binary", data)
					continue
				}
				for _, pkt := range packet.DecodePacket(data).Parse() {
					go c.Handle(pkt)
				}
			}
		}
	}()
	go c.startHeartBeat()
}

func (c *Client) ConnectAndStart() error {
	err := c.Connect()
	if err != nil {
		return err
	}
	c.Start()
	return nil
}

func (c *Client) SetHost(host string) {
	c.host = host
}

func (c *Client) UseDefaultHost() {
	c.SetHost("wss://broadcastlv.chat.bilibili.com/sub")
}

func (c *Client) startHeartBeat() {
	pkt := packet.NewHeartBeatPacket()
	for {
		select {
		case <-c.Done:
			return
		case <-time.After(30 * time.Second):
			if err := c.conn.WriteMessage(websocket.BinaryMessage, pkt); err != nil {
				c.LogFatal(err)
			}
			log.Debug("send: HeartBeat")
		}
	}
}

func (c *Client) sendEnterPacket() {
	rid, err := strconv.Atoi(c.roomID)
	if err != nil {
		c.LogFatal("error roomID")
	}
	pkt := packet.NewEnterPacket(0, rid)
	if err := c.conn.WriteMessage(websocket.BinaryMessage, pkt); err != nil {
		c.LogFatal(err)
	}
	log.Debugf("send: EnterPacket: %v", pkt)
}

func (c *Client) LogFatal(v ...interface{}) {
	c.Stop()
	log.Panic(v...)
}

func (c *Client) LogFatalln(v ...interface{}) {
	c.Stop()
	log.Panicln(v...)
}

func LogFatal(v ...interface{}) {
	log.Panic(v...)
}

func (c *Client) Stop() {
	c.Done <- struct{}{}
}

func getDanmuInfo(roomID string) (*DanmuInfo, error) {
	url := fmt.Sprintf("https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo?id=%s&type=0", roomID)
	resp, err := http.Get(url)
	if err != nil {
		LogFatal(err)
		return nil, err
	}
	defer resp.Body.Close()
	result := &DanmuInfo{}
	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		LogFatal(err)
		return nil, err
	}
	return result, nil
}
