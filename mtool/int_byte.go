/**
* @Author: Cooper
* @Date: 2019/12/4 21:57
 */

package mtool

func IntToByte(n uint32) []byte {
	bs := make([]byte , 4)

	bs[3] = uint8(n)
	bs[2] = uint8(n >> 8)
	bs[1] = uint8(n >> 16)
	bs[0] = uint8(n >> 24)

	return bs
}

func ByteToInt(bs []byte) uint32 {
	x := uint32(0)
	x += uint32(int(bs[3])*1)
	x += uint32(int(bs[2])*256)
	x += uint32(int(bs[1])*256*256)
	x += uint32(int(bs[0])*256*256*256)
	return x
}