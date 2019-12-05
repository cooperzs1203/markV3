/**
* @Author: Cooper
* @Date: 2019/12/4 22:07
 */

package mface

type MDataProtocol interface {
	Unmarshal([]byte)
	Marshal(string, []byte) []byte
	CompletedMessageChan() chan MMessage
}
