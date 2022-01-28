package protocol

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

const (
	PackageLength = 4
	HeaderLength  = 2
)

//封包
func Pack(data []byte) ([]byte, error) {
	// 包头：data长度和header头
	var (
		packLen   = make([]byte, PackageLength)
		headerLen = make([]byte, HeaderLength)
	)
	binary.BigEndian.PutUint32(packLen, uint32(HeaderLength+len(data)))
	binary.BigEndian.PutUint16(headerLen, 0x12)
	packBuf := bytes.NewBuffer(packLen)
	packBuf.Write(headerLen)

	// 包体：data
	packBuf.Write(data)
	return packBuf.Bytes(), nil
}

//解包
func Unpack(reader *bufio.Reader) (string, error) {
	var (
		prefixLen = PackageLength + HeaderLength
	)
	//读取前缀字节（这里是定义的6个字节，包长度4个和header2个）
	prefixByte, err := reader.Peek(prefixLen)
	if err != nil {
		return "", err
	}

	//读取PackageLength长度
	bodyLen := binary.BigEndian.Uint32(prefixByte[:PackageLength])

	//创建一个用于读取数据的buffer
	prefixBuff := bytes.NewBuffer(prefixByte)
	err = binary.Read(prefixBuff, binary.BigEndian, &bodyLen)
	if err != nil {
		return "", err
	}

	// 当前reader可以读取的字节数小于前缀字节数，说明数据丢失，返回error
	if uint32(reader.Buffered()) < bodyLen+PackageLength {
		return "", err
	}

	packet := make([]byte, uint32(bodyLen+PackageLength))
	_, err = reader.Read(packet)
	if err != nil {
		return "", err
	}
	return string(packet[prefixLen:]), nil
}
