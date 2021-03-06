package conn

import (
	"net"
	"fmt"
	"io"
	"time"
	"code.google.com/p/goprotobuf/proto"
)


const (
	MAX_RECV_QUEUE_SIZE = 1024
	MAX_SEND_QUEUE_SIZE = 1024
	MAX_RECV_BUFF_SIZE	= 65536
	MAX_SEND_TIMEOUT	= 5			// second
)

type SendQueue chan proto.Message
type RecvQueue chan []*Packet

// implement IConnection interface.
type Connection struct {
	server		IServer
	terminate	TerminateType	// enum TerminateType
	conn		net.Conn		// connection
	recv		RecvQueue		// received queue.
	send		SendQueue		// sending queue.
}

// create connection.
func NewConnection( server IServer, c net.Conn ) *Connection {
	task := new(Connection)
	task.server = server
	task.terminate = TerminateType_None
	task.conn = c

	task.recv = make( RecvQueue, MAX_RECV_QUEUE_SIZE )
	task.send = make( SendQueue, MAX_SEND_QUEUE_SIZE )

	go task.OnRecv()
	go task.OnSend()

	return task
}

// send message.
func (this *Connection) SendCmd( msg proto.Message ) {
	select {
		case <- time.After(time.Second * MAX_SEND_TIMEOUT):
			fmt.Println( "send message to timeout. ", &msg )
		case this.send <- msg:
			break;
	}
}

// close
func (this *Connection) Close() {
	if TerminateType_None != this.terminate {
		this.terminate = TerminateType_Active
	}

	this.conn.Close()
	close( this.recv )
	close( this.send )
}

// message callback
func (this *Connection) Process() {
	select {
		case packets := <-this.recv:
			CallOnMessage( packets, this )
		default:
			//fmt.Println( "no messaged." );
			break;
	}
}

func (this *Connection) CheckTerminate() bool {
	return TerminateType_None != this.terminate
}

// check error
func CheckError( conn *Connection, err error ) bool {
	if err == nil {
		return true
	}

	if err == io.EOF {
		conn.terminate = TerminateType_Passive
		fmt.Println( "remote colsed: ", err.Error() )
		return false;
	}

	if TerminateType_None == conn.terminate {
		conn.terminate = TerminateType_Error;
		fmt.Println( "Connection read error: ", err.Error() )
		return false;
	}

	fmt.Println( "connection closed." )
	return false
}

// receive handle.
func (this *Connection) OnRecv() {

	buffer := make( []byte, MAX_RECV_BUFF_SIZE )

	recvd_len := uint16(0);

	for {
		// block read from connection.
		rlen, err := this.conn.Read( buffer[recvd_len:] )

		if !CheckError( this, err ) {
			break;
		}

		if rlen == 0 {
			continue
		}

		recvd_len += uint16( rlen )

		read_size, packets, perr := ParsePacket( buffer[:recvd_len] )

		// parse message to faield.
		if perr != nil {
			fmt.Println( "Connection parse message to failed. ", perr.Error() )
			this.terminate = TerminateType_Error;
			break;
		}

		// cleanup recv buffer.
		if read_size != 0 {
			copy( buffer, buffer[read_size:recvd_len] )
			recvd_len -= uint16(read_size)
		}

		// push message(s) to chan.
		if len( packets ) != 0 {
			//fmt.Println( "Connection: push message, ", len( packets ) );
			this.recv<-packets
		}
	}
}

// sending handle.
func (this *Connection) OnSend() {

	for {
		msg, ok := <-this.send

		if false == ok {
			// channel closed.
			break;
		}

		msg_buf, err := EncodeMessage( msg, false, false )

		if err != nil {
			fmt.Println( "EncodeMessage to falied. ", err.Error() );
			break;
		}

		msg_len, offset := len( msg_buf ), 0

		for ; offset < msg_len; {
			sent, err := this.conn.Write( msg_buf[offset:] )

			if !CheckError( this, err ) {
				break;
			}

			offset += sent;
		}
	}
}

