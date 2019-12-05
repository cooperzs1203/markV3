/**
* @Author: Cooper
* @Date: 2019/12/4 20:32
 */

package mnet

import (
	"fmt"
	"log"
	"markV3/mface"
	"markV3/mtool"
	"net"
	"time"
)

func newConnection(conn *net.Conn, cm mface.MConnManager) (mface.MConnection, error) {
	c := &connection{
		id:           mtool.GetRandString(32),
		status:       Serve_Status_UnStarted,
		conn:         conn,
		cm:           cm,
		responseChan: make(chan mface.MMessage, 0),
		dp:           nil,
	}

	if err := c.Load(); err != nil {
		return nil, err
	}

	return c, nil
}

type connection struct {
	id     string
	status int
	conn   *net.Conn
	cm     mface.MConnManager
	dp     mface.MDataProtocol

	responseChan chan mface.MMessage
}

func (c *connection) Load() error {
	log.Printf("[%s] Load", c.id)
	c.status = Serve_Status_Load

	c.responseChan = make(chan mface.MMessage, c.cm.Server().Config().ConnResponseCS())
	c.dp = NewDataProtocol("HEAD", 10, 4, c.id, c.cm.Server().Config().DPCompletedCS())

	return nil
}

func (c *connection) Start() error {
	log.Printf("[%s] Start", c.id)
	c.status = Serve_Status_Starting

	go c.startReadData()
	go c.startAcceptRequest()
	go c.startReplyResponse()

	c.status = Serve_Status_Running

	return nil
}

func (c *connection) Stop() error {
	log.Printf("[%s] Stop", c.id)
	return nil
}

func (c *connection) StartEnding() error {
	log.Printf("[%s] Start Ending", c.id)
	// 1.close read data
	c.status = Serve_Status_Ending

	// 2.close data protocol
	close(c.dp.CompletedMessageChan())

	for {
		if len(c.dp.CompletedMessageChan()) == 0 {
			break
		}
	}

	return nil
}

func (c *connection) OfficialEnding() error {
	// 1.stop accept response
	close(c.responseChan)

	// 2.wait for response channel empty
	for {
		if len(c.responseChan) == 0 {
			break
		}
	}

	// 3.close *net.Conn
	if err := (*c.conn).Close(); err != nil {
		log.Printf(fmt.Sprintf("[%s] close net.Conn error : %+v", c.id, err))
	}

	c.status = Serve_Status_Stopped
	log.Printf("[%s] Official Ending", c.id)

	return nil
}

func (c *connection) Reload() error {
	log.Printf("[%s] Reload", c.id)
	c.status = Serve_Status_Reload

	newDP := NewDataProtocol("HEAD", 10, 4, c.id, c.cm.Server().Config().DPCompletedCS())

	close(c.dp.CompletedMessageChan())

	for {
		if len(c.dp.CompletedMessageChan()) == 0 {
			break
		}
	}

	c.dp = newDP

	newResponseChan := make(chan mface.MMessage, c.cm.Server().Config().ConnResponseCS())

	close(c.responseChan)

	for {
		if len(c.responseChan) == 0 {
			break
		}
	}

	c.responseChan = newResponseChan

	c.status = Serve_Status_Running

	return nil
}

func (c *connection) ID() string {
	return c.id
}

func (c *connection) ResponseChan() chan mface.MMessage {
	return c.responseChan
}

func (c *connection) Status() int {
	return c.status
}

// =========== private methods ===========

func (c *connection) startReadData() {
	for {
		buffer := make([]byte, 1024)
		_ = (*c.conn).SetReadDeadline(time.Now().Add(time.Second * time.Duration(c.cm.Server().Config().ConnReadTimeOut()))) // 设置此次读取超时时间
		cnt, err := (*c.conn).Read(buffer)
		_ = (*c.conn).SetReadDeadline(time.Time{}) // 取消此次读取超时时间
		if err != nil {
			if c.status >= Serve_Status_Ending {
				break
			} else {
				continue
			}
		}

		c.dp.Unmarshal(buffer[:cnt])
	}
}

func (c *connection) startAcceptRequest() {
	for {
		request, ok := <-c.dp.CompletedMessageChan()
		if !ok {
			if c.status >= Serve_Status_Ending {
				break
			} else {
				continue
			}
		}

		c.cm.RequestChan() <- request
	}
}

func (c *connection) startReplyResponse() {
	for {
		response, ok := <-c.responseChan
		if !ok {
			if c.status >= Serve_Status_Ending {
				break
			} else {
				continue
			}
		}

		cnt, err := (*c.conn).Write(response.Marshal())
		log.Printf(fmt.Sprintf("[%s] get reply response len : %d  error : %+v", c.id, cnt, err))
	}
}
