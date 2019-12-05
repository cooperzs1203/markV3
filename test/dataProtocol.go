/**
* @Author: Cooper
* @Date: 2019/12/5 18:38
 */

package test

import (
	"errors"
	"log"
)

type MDataProtocol interface {
	EnPack(...[]byte) ([]byte, error)
	DePack([]byte)
	CompletedDataChan() chan Message
}

// ===============

func NewDataProtocol(tag string, options ...uint32) (MDataProtocol, error) {
	// 标志不可为空
	if tag == "" {
		return nil, errors.New("tag can not be empty")
	}

	// 自定义长度组不可包含长度0
	for _, option := range options {
		if option == uint32(0) {
			return nil, errors.New("option can not be zero")
		}
	}

	os := []uint32{uint32(len(tag))} // 追加tag长度
	os = append(os, options...)      // 追加自定义长度组
	os = append(os, uint32(0))       // 追加实际数据长度占位，默认为0

	dp := &dataProtocol{
		tag:     tag,
		options: os,
		buffer:  make([]byte, 0),
		cmChan:  make(chan Message, 0),
	}
	return dp, nil
}

type dataProtocol struct {
	tag     string
	options []uint32
	buffer  []byte
	cmChan  chan Message
}

// 封包
func (dp *dataProtocol) EnPack(values ...[]byte) ([]byte, error) {
	packData := make([]byte, 0)

	// values 在最前面加一个 tag
	if string(values[0]) != dp.tag {
		tags := [][]byte{[]byte(dp.tag)}
		values = append(tags, values...)
	}

	// values数量与options数量对不上
	if len(values) != len(dp.options) {
		return nil, errors.New("length of values not match length options")
	}

	var err error
	for index, value := range values {
		if dp.options[index] != uint32(0) && (uint32(len(value)) != dp.options[index]) { // 如果option 为0 ， 则为预留实际长度占位
			err = errors.New("value length not match option")
			break
		} else {
			packData = append(packData, value...)
		}
	}

	if err != nil {
		return nil, err
	}

	return packData, nil
}

// 解包
func (dp *dataProtocol) DePack(originData []byte) {
	dp.buffer = append(dp.buffer, originData...)
	length := uint32(len(dp.buffer))

	i := uint32(0)
	for ; i < length; i++ {
		// 不足以拿到完整头部
		if i+dp.totalPrefixLength() > length {
			break
		}
		// 非tag标志
		if string(dp.buffer[i:i+dp.tagLength()]) != dp.tag {
			continue
		}
		// 默认取options最后一个值为记录数据长度的长度
		dataLength := ByteToInt(dp.buffer[i+dp.totalPrefixLength()-dp.lordLength() : i+dp.totalPrefixLength()])
		// 无法取得完整数据(包括头部和真实数据)
		if i+dp.totalPrefixLength()+dataLength > length {
			break
		}
		completedData := dp.buffer[i : i+dp.totalPrefixLength()+dataLength]

		// 将完整数据转成结构
		msg := newMessage(dp.options)
		if err := msg.Unmarshal(completedData); err != nil {
			log.Println(err)
			break
		}

		// 发送到外界
		dp.cmChan <- msg
		// 游标移到取出数据的最后一位，循环会+1
		i += dp.totalPrefixLength() + dataLength - 1
	}

	if i != length { // 有剩余不完整数据，积攒到下次粘合
		leftData := dp.buffer[i:]
		dp.buffer = make([]byte, 0)
		dp.buffer = append(dp.buffer, leftData...)
	} else { // 无剩余不完整数据，格式化缓存
		dp.buffer = make([]byte, 0)
	}
}

// 完整数据通道
func (dp *dataProtocol) CompletedDataChan() chan Message {
	return dp.cmChan
}

func (dp *dataProtocol) totalPrefixLength() uint32 {
	length := uint32(0)
	for _, l := range dp.options {
		length += l
	}

	return length
}

func (dp *dataProtocol) tagLength() uint32 {
	return dp.options[0]
}

// length of record data length , 最后一个是0，追加的实际数据占位 , 默认倒数第二位是实际数据长度记录位
func (dp *dataProtocol) lordLength() uint32 {
	return dp.options[len(dp.options)-2]
}

func IntToByte(n uint32) []byte {
	bs := make([]byte, 4)

	bs[3] = uint8(n)
	bs[2] = uint8(n >> 8)
	bs[1] = uint8(n >> 16)
	bs[0] = uint8(n >> 24)

	return bs
}

func ByteToInt(bs []byte) uint32 {
	x := uint32(0)
	x += uint32(int(bs[3]) * 1)
	x += uint32(int(bs[2]) * 256)
	x += uint32(int(bs[1]) * 256 * 256)
	x += uint32(int(bs[0]) * 256 * 256 * 256)
	return x
}
