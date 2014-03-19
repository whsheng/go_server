package main

import (
	"log"
	"conn"
	"command"
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"time"
	"handler"
	"table"
//	"runtime/pprof"
//	"os"
//	"encoding/json"
//	"encoding/binary"
//	"bytes"
//	"reflect"
)

// register message handler.
var _ = &handler.InitializeMessageHandler{}


//////////////////////////////////////////////////
func RunTimer( timeout int64, ch chan bool ) {
	go func () {
		for {
			time.Sleep( time.Duration( time.Millisecond * time.Duration(timeout) ) );
			ch <- true;
		}
	}();
}

func TestTimer( ms int64 ) {

	timer_ch := make( chan bool, 1 )
	RunTimer( ms, timer_ch )

	msg_ch := make( chan int, 10 )

	for {
		select {
		case msg := <-msg_ch:
				fmt.Println( "msg: ", msg );
			case <-timer_ch:
				fmt.Println( "timeout" );
		}
	}
}

func main() {

	testdata_mgr := table.GetTestDataManager()

	//fmt.Printf( "%v\n", testdata_mgr.FindEx( 3, 1 ) )
	fmt.Println( testdata_mgr.FindEx( 3, 3 ) )
	fmt.Println( testdata_mgr.FindEx( 8, 1 ) )

	item_mgr := table.GetItemBaseManager()

	fmt.Printf( "%+v\n", item_mgr.Find( 20208 ) )

	/////////////////////////////////////////////////////////////////////////////
	log.Println( "Hello Go." );

	////////////////////////////////////////////////////////////////////////////////
	// build message.
	msg := &Cmd_Test.TestMessage{
		Name: []byte("linbo"),
		Age: proto.Int32(25),
		Desc: []byte("descriptor !!!") };

	// msg := &Cmd_Test.RequestLogin{ Username: []byte("limpo"), Password: []byte("pwd") }

	////////////////////////////////////////////////////////////////////////////////
	// encode the message to byte buffer.
	msg_buff, err := conn.EncodeMessage( msg, false, false );

	if err != nil {
		fmt.Println( "encode err=> ", err.Error() );
		return
	}
	////////////////////////////////////////////////////////////////////////////////
	// decode message from byte buffer.
	consumeSize, packets, err := conn.ParsePacket( msg_buff[:] );

	if err != nil {
		fmt.Println( "parse err=> ", err.Error() );
		return
	}

	fmt.Println( "size: ", consumeSize, " packets: ", len(packets) );

	////////////////////////////////////////////////////////////////////////////////
	// call callback on per message.
	//conn.CallOnMessage( packets, nil );


	server := conn.NewServer( "0.0.0.0", 4321 );

	if server == nil {
		return
	}

	server.Run()
}

