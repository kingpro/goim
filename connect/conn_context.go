package connect

import (
	"io"
	"log"
	"net"
	"strings"
	"time"
)

const ReadDeadline = 10 * time.Minute

// ConnContext 连接上下文
type ConnContext struct {
	Codec *Codec      // 编解码器
	Info  interface{} // 附加信息
}

// Message 消息
type Message struct {
	Code    int    // 消息类型
	Content []byte // 消息体
}

func NewConnContext(conn *net.TCPConn) *ConnContext {
	codec := NewCodec(conn)
	return &ConnContext{Codec: codec}
}

// DoConn 处理TCP连接
func (c *ConnContext) DoConn() {
	defer RecoverPanic()

	c.HandleConnect()

	for {
		err := c.Codec.Conn.SetReadDeadline(time.Now().Add(ReadDeadline))
		if err != nil {
			log.Println(err)
			return
		}

		_, err = c.Codec.Read()
		if err != nil {
			c.HandleReadErr(err)
			return
		}

		for {
			message, ok := c.Codec.Decode()
			if ok {
				c.HandleMessage(message)
				continue
			}
			break
		}
	}
}

// HandleConnect 建立连接
func (c *ConnContext) HandleConnect() {
	return
}

// HandleMessage 处理消息
func (c *ConnContext) HandleMessage(message *Message) {

}

// HandleReadErr 读取conn错误
func (c *ConnContext) HandleReadErr(err error) {
	log.Println(err)
	// 客户端主动关闭连接或者异常程序退出
	if err == io.EOF {
		c.Codec.Conn.Close()
		return
	}
	str := err.Error()
	// SetReadDeadline 之后，超时返回的错误
	if strings.HasSuffix(str, "i/o timeout") {
		c.Codec.Conn.Close()
		return
	}
	// 服务器主动关闭连接
	if strings.HasSuffix(str, "use of closed network connection") {
		return
	}
}

// HandleActive
func (c *ConnContext) HandleActive(*ConnContext) {

}

// HandleInactive 监听到客户端停止活动
func (c *ConnContext) HandleInactive(*ConnContext) {

}

// Close 关闭TCP连接
func (c *ConnContext) Close(*ConnContext, error) {
	c.Codec.Conn.Close()
}
