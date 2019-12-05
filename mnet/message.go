/**
* @Author: Cooper
* @Date: 2019/12/4 22:03
 */

package mnet

import (
	"markV3/mface"
	"markV3/mtool"
)

func NewMessage(connId string, data []byte) mface.MMessage {

}

func newMessage(connId string, data []byte, headLength uint32, headMsgLegth uint32, totalHeaderLegth uint32) mface.MMessage {
	m := &message{
		connId:     connId,
		head:       "",
		msgId:      "",
		dataLength: 0,
		data:       data,
	}

	m.parsing(headLength, headMsgLegth, totalHeaderLegth)

	return m
}

type message struct {
	connId     string
	head       string
	msgId      string
	dataLength uint32
	data       []byte
}

func (m *message) parsing(headLength uint32, headMsgLegth uint32, totalHeaderLegth uint32) {
	originData := m.data
	m.head = string(originData[:headLength])
	m.msgId = string(originData[headLength:headMsgLegth])
	m.dataLength = mtool.ByteToInt(originData[headMsgLegth:totalHeaderLegth])
	m.data = originData[totalHeaderLegth:]
}

func (m *message) MsgID() string {
	return m.msgId
}

func (m *message) RouteID() string {
	return m.msgId
}

func (m *message) ConnID() string {
	return m.connId
}

func (m *message) Marshal() []byte {
	bs := make([]byte , 0)

	bs = append(bs , []byte(m.head)...)
	bs = append(bs , []byte(m.msgId)...)
	bs = append(bs , mtool.IntToByte(m.dataLength)...)
	bs = append(bs , m.data...)

	return bs
}
