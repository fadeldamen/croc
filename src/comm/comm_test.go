package comm

import (
	"crypto/rand"
	"net"
	"testing"
	"time"

	log "github.com/cihub/seelog"
	"github.com/stretchr/testify/assert"
)

func TestComm(t *testing.T) {
	defer log.Flush()
	token := make([]byte, 40000000)
	rand.Read(token)

	port := "8001"
	go func() {
		log.Debugf("starting TCP server on " + port)
		server, err := net.Listen("tcp", "0.0.0.0:"+port)
		if err != nil {
			log.Error(err)
		}
		defer server.Close()
		// spawn a new goroutine whenever a client connects
		for {
			connection, err := server.Accept()
			if err != nil {
				log.Error(err)
			}
			log.Debugf("client %s connected", connection.RemoteAddr().String())
			go func(port string, connection net.Conn) {
				c := New(connection)
				err = c.Send([]byte("hello, world"))
				assert.Nil(t, err)
				data, err := c.Receive()
				assert.Nil(t, err)
				assert.Equal(t, []byte("hello, computer"), data)
				data, err = c.Receive()
				assert.Nil(t, err)
				assert.Equal(t, []byte{'\x00'}, data)
				data, err = c.Receive()
				assert.Nil(t, err)
				assert.Equal(t, token, data)
			}(port, connection)
		}
	}()

	time.Sleep(300 * time.Millisecond)
	a, err := NewConnection("localhost:" + port)
	assert.Nil(t, err)
	data, err := a.Receive()
	assert.Equal(t, []byte("hello, world"), data)
	assert.Nil(t, err)
	assert.Nil(t, a.Send([]byte("hello, computer")))
	assert.Nil(t, a.Send([]byte{'\x00'}))

	assert.Nil(t, a.Send(token))

}
