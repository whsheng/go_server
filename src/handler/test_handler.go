package handler

import (
	"conn"
	"fmt"
	"command"
	"code.google.com/p/goprotobuf/proto"
)

// Cmd.Test.TestMessage from command/test_message.proto
func processTestMessage( pb proto.Message, connection conn.IConnection ) {

	msg := pb.(*Cmd_Test.TestMessage);

	//fmt.Printf( "process TestMessageHandle, Name: %s, Age: %d, Desc: %s, Count: %d\n",
	//msg.Name, *msg.Age, msg.Desc, *msg.Count );

	*msg.Count += 1;

	connection.SendCmd( pb )
}

// Cmd.Test.RequestLogin from command/test_message.proto
func processRequestLogin( pb proto.Message, connection conn.IConnection ) {

	msg := pb.(*Cmd_Test.RequestLogin);

	fmt.Printf( "process RequestLogin, Username: %s, Password: %s, Conn: %v\n",
	msg.Username, msg.Password, connection );

	benchmark := &Cmd_Test.TestMessage{
		Name: []byte("linbo"),
		Age: proto.Int32(25),
		Desc: []byte("descriptor !!!"),
		Count: proto.Int32(1) };

	connection.SendCmd( benchmark )
}

/////////////////////////////////////////////////////////////////////////////////////////
func init() {
	// bind the message handler.

	conn.REGISTER_OPCODE_HANDLER( &Cmd_Test.TestMessage{}, processTestMessage );
	conn.REGISTER_OPCODE_HANDLER( &Cmd_Test.RequestLogin{}, processRequestLogin );

}
