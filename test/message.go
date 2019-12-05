/**
* @Author: Cooper
* @Date: 2019/12/5 21:11
 */

package test

import (
	"errors"
)

type Message interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	Data() []byte
	DataLength() uint32
	Value(int) ([]byte, error)
	ValuesLength() uint32
	SetValue(int ,[]byte) error
}

// ====================

func newMessage(options []uint32) Message {
	m := &message{
		options:    options,
		values:     make([][]byte, 0),
		dataLength: uint32(0),
		data:       make([]byte, 0),
	}
	return m
}

type message struct {
	options []uint32
	values  [][]byte

	dataLength uint32
	data       []byte
}

// 将自身数据化作 []byte ， 默认本身已经具有数据
func (m *message) Marshal() ([]byte, error) {
	packData := make([]byte, 0)

	var err error
	for index, value := range m.values {
		if m.options[index] != uint32(0) && (uint32(len(value)) != m.options[index]) { // 如果option 为0 ， 则为预留实际长度占位
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

func (m *message) Unmarshal(data []byte) error {
	if len(data) == 0 {
		return errors.New("data can not be empty")
	}

	var err error
	lastOption := uint32(0)
	for _, option := range m.options {
		if lastOption+option > uint32(len(data)) { // 超出长度
			err = errors.New("index out of range")
			break
		}
		value := make([]byte, 0)
		if option == uint32(0) {
			value = data[lastOption:]
		} else {
			value = data[lastOption : lastOption+option]
		}
		m.values = append(m.values, value)
		lastOption += option
	}

	if err != nil {
		return err
	}

	if len(m.values) > 2 {
		m.dataLength = ByteToInt(m.values[len(m.values)-2])
		m.data = m.values[len(m.values)-1]
	}

	return nil
}

func (m *message) Data() []byte {
	return m.data
}

func (m *message) DataLength() uint32 {
	return m.dataLength
}

func (m *message) Value(index int) ([]byte, error) {
	if index > len(m.values)-1 {
		return nil, errors.New("index out of range")
	}
	return m.values[index], nil
}

func (m *message) ValuesLength() uint32 {
	return uint32(len(m.values))
}

func (m *message) SetValue(index uint32 , value []byte) error {
	if index > len(m.values)-1 {
		return errors.New("index out of range")
	}

	// option 为零 默认为占位的实际数据
	if m.options[index] != 0 && (uint32(len(value)) != m.options[index]) {
		return errors.New("value length not match index of option")
	}

	m.values[index] = value
	return nil
}

func (m *message) Reply(rspData []byte) Message {
	rsp := newMessage(m.options)

	for i := uint32(0); i < m.ValuesLength(); i++ {
		rsp.SetValue(i , m.values[i])
	}



}