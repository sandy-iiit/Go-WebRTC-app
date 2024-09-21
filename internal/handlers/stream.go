package handlers

import (
	w "GoVideoChat-Project/pkg/webrtc"
	"fmt"
	"github.com/gofiber/websocket/v2"

	"github.com/gofiber/fiber/v2"
	"time"
)

func Stream(c *fiber.Ctx) error {
	ssuid := c.Params("ssuid")
	if ssuid == "" {
		c.Status(400)
		return nil
	}
	ws := "ws"
	//if os.Getenv("ENVIRONMENT") == "PRODUCTION" {
	//	ws = "wss"
	//}

	w.RoomsLock.Lock()
	if _, ok := w.Streams[ssuid]; ok {
		w.RoomsLock.Unlock()
		return c.Render("stream", fiber.Map{
			"StreamWebSocketAddr": fmt.Sprintf("%s://%s/stream/%s/websocket", ws, c.Hostname(), ssuid),
			"ChatWebSocketAddr":   fmt.Sprintf("%s://%s/stream/%s/chat/websocket", ws, c.Hostname(), ssuid),
			"ViewerWebSocketAddr": fmt.Sprintf("%s://%s/stream/%s/viewer/websocket", ws, c.Hostname(), ssuid),
			"Type":                "stream",
		}, "/layouts/main")

	}
	w.RoomsLock.Unlock()
	return c.Render("stream", fiber.Map{
		"NoStream": "true",
		"Leave":    "true",
	}, "/layouts/main")
}

func StreamWebSocket(c *websocket.Conn) {
	ssuid := c.Params("ssuid")
	if ssuid == "" {
		return
	}
	w.RoomsLock.Lock()
	if stream, ok := w.Streams[ssuid]; ok {
		w.RoomsLock.Unlock()
		w.StreamConn(c, stream.Peers)
		return
	}
	w.RoomsLock.Unlock()

}

func StreamViewerWebSocket(c *websocket.Conn) {
	ssuid := c.Params("ssuid")
	if ssuid == "" {
		return
	}
	w.RoomsLock.Lock()
	if stream, ok := w.Streams[ssuid]; ok {
		w.RoomsLock.Unlock()
		ViewerConn(c, stream.Peers)
		return
	}
}

func ViewerConn(c *websocket.Conn, p *w.Peers) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	defer c.Close()

	for {
		select {
		case <-ticker.C:
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write([]byte(fmt.Sprintf("%d", len(p.Connections))))
		}
	}
}
