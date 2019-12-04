/**
* @Author: Cooper
* @Date: 2019/12/4 21:29
 */

package mnet

import (
	"markV3/mface"
	"markV3/mtool"
)

func newDataProtocol(head string, msgIdLength uint32, connId string, completedMsgCS uint64) mface.MDataProtocol {
	dp := &dataProtocol{
		connId:           connId,
		head:             head,
		headL:            uint32(len(head)),
		MsgIdL:           msgIdLength,
		LDataL:           4,
		buffer:           make([]byte, 0),
		completedMsgChan: make(chan mface.MMessage, completedMsgCS),
	}

	return dp
}

type dataProtocol struct {
	connId string

	head   string
	headL  uint32
	MsgIdL uint32
	LDataL uint32

	buffer           []byte
	completedMsgChan chan mface.MMessage
}

func (dp *dataProtocol) Unmarshal(data []byte) {
	dp.buffer = append(dp.buffer, data...)
	length := uint32(len(dp.buffer))

	i := uint32(0)
	for ; i < length; i++ {
		if i+dp.totalHeaderLegth() > length {
			break
		}
		if string(dp.buffer[i:i+dp.headL]) != dp.head {
			continue
		}
		dataLength := mtool.ByteToInt(dp.buffer[i+dp.headMsgLegth() : i+dp.totalHeaderLegth()])
		if i+dp.totalHeaderLegth()+dataLength > length {
			break
		}

		totalData := dp.buffer[i : i+dp.totalHeaderLegth()+length]
		msg := newMessage(dp.connId, totalData, dp.headLegth(), dp.headMsgLegth(), dp.totalHeaderLegth())
		dp.completedMsgChan <- msg

		i = i + dp.totalHeaderLegth() + length - 1

	}

	if i != length {
		leftData := dp.buffer[i:]
		dp.buffer = make([]byte, 0)
		dp.buffer = append(dp.buffer, leftData...)
	} else {
		dp.buffer = make([]byte, 0)
	}
}

func (dp *dataProtocol) CompletedMessageChan() chan mface.MMessage {
	return dp.completedMsgChan
}

func (dp *dataProtocol) headLegth() uint32 {
	return dp.headL
}

func (dp *dataProtocol) headMsgLegth() uint32 {
	return dp.headL + dp.MsgIdL
}

func (dp *dataProtocol) totalHeaderLegth() uint32 {
	return dp.headL + dp.MsgIdL + dp.LDataL
}
