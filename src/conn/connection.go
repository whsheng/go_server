package conn

import (
	"reflect"
	"errors"
	"fmt"
	"bytes"
	"encoding/binary"
	"code.google.com/p/goprotobuf/proto"
)

// defined the connection interface.
type IConnection interface {

	// send pb message to the remote.
	SendCmd( proto.Message )

	// close the connection.
	Close()

	// process message queue.
	Process()

	// check terminate
	CheckTerminate() bool
}

// packet define.
/*
struct Packet 
{
	// COMPRESS_MARKER | CRYPTO_MARKER | PACKET_LEGTH
	uint32	head;

	// checksum(opcode + message).
	uint32	checksum;

	// opcode first | second.
	uint8	first;
	uint8	second;
	
	// binary data.
	char	data[0];
}

*/

/////////////////////////////////////////////////////////////////////////////////////////////
// enum TerminateType
type TerminateType int32

const TerminateType_None	TerminateType = 0;
const TerminateType_Active	TerminateType = 1;
const TerminateType_Passive	TerminateType = 2;
const TerminateType_Error	TerminateType = 3;


/////////////////////////////////////////////////////////////////////////////////////////////
// other define in here.
const COMPRESS_MARKER  = uint32(0x40000000);   // (0100)
const CRYPTO_MARKER    = uint32(0x80000000);   // (1000)
const SIZE_MARKER      = uint32(0x3FFFFFFF);   // (0011)
const HEAD_SIZE		   = uint32( 8 );		   // uint32(head) + uint32(checksum)
const MSG_OPCODE	   = uint32( 2 );		   // uint8(first) + uint8(second)
const MAX_OPCODE_SIZE  = uint32( 255 );		   // 0xFF

/////////////////////////////////////////////////////////////////////////////////////////////
// golang packet.
type Packet struct {
//private:
	head		uint32
	checksum	uint32

//public:
	First		uint8
	Second		uint8
	Message		proto.Message
}

// member function of Packet
func (this Packet) CheckCompress() bool {

	return 0 != (this.head & COMPRESS_MARKER);
}

func (this Packet) CheckCrypto() bool {
	return 0 != (this.head & CRYPTO_MARKER);
}

func (this Packet) FetchSize() uint32 {
	return uint32( this.head & SIZE_MARKER );
}

////////////////////////////////////////////////////////////////////////////////////////////
type OpcodeHandle struct {
	pb interface{}
	handle func( proto.Message, IConnection )()
}

// global variable.
var gs_opcode_handlers [MAX_OPCODE_SIZE][MAX_OPCODE_SIZE]*OpcodeHandle

func init() {
	//
}

func REGISTER_OPCODE_HANDLER( pb interface{}, handle func( proto.Message, IConnection )() ) {

	first, second := GetOpcode( pb.(proto.Message) );

	if !(uint32(first) < MAX_OPCODE_SIZE && uint32(second) < MAX_OPCODE_SIZE) {
		panic( fmt.Errorf( "opcode out of range: %d %d", first, second ) );
	}

	if gs_opcode_handlers[first][second] != nil {
		panic( fmt.Errorf( "the opcode aready bind the handler: [%d:%d]", first, second ) );
	}

	gs_opcode_handlers[first][second] = &OpcodeHandle{ pb, handle }
}
////////////////////////////////////////////////////////////////////////////////////////////

// get opcode [first, second] from message.
func GetOpcode( pb proto.Message ) ( first, second byte ) {

	// check it.
	if _, ok := pb.(proto.Message); !ok {
		panic( "the [pb] nonprotobuf object." );
	}

	pb_mutable := reflect.ValueOf( pb );

	first_func := pb_mutable.MethodByName( "GetFIRST" )
	second_func := pb_mutable.MethodByName( "GetSECOND" )

	if !first_func.IsValid() || !second_func.IsValid() {
		panic( "invalid first/second method." );
	}

	// call void.
	params := make( []reflect.Value, 0 );

	// call func and fetch return value.
	first = byte( first_func.Call( params )[0].Int() );
	second = byte( second_func.Call( params )[0].Int() );

	return
}

// check message.
func ChecksumForMessage( checksum uint32, message []byte ) bool {
/*
	if Adler32( message ) != checksum {
		// ERROR package.
		return false;
	}
*/

	return 0 == checksum;
}

// decode packet from []byte
func ParsePacket( buf []byte ) ( consumeSize uint32, packets []*Packet, err error ) {

	consumeSize	= 0
	packets		= []*Packet{}
	err			= nil

	recvSize	:= uint32(len(buf))

	// check readable length.
	if recvSize < (HEAD_SIZE + MSG_OPCODE) {
		return
	}

	for {

		// check buffer length.
		if recvSize - consumeSize < (HEAD_SIZE + MSG_OPCODE) {
			break;
		}

		// fetch head checksum and opcode.
		head, checksum, first, second := uint32(0), uint32(0), uint8(0), uint8(0);

		binary.Read(bytes.NewReader(buf[consumeSize+0:consumeSize+4]),binary.LittleEndian, &head)
		binary.Read(bytes.NewReader(buf[consumeSize+4:consumeSize+8]),binary.LittleEndian, &checksum)
		binary.Read(bytes.NewReader(buf[consumeSize+8:consumeSize+9]),binary.LittleEndian, &first)
		binary.Read(bytes.NewReader(buf[consumeSize+9:consumeSize+10]),binary.LittleEndian, &second)

		// opcode + message
		data_len := uint32(head & SIZE_MARKER);

		// less a message body
		if recvSize - consumeSize < (data_len + HEAD_SIZE) {
			break;
		}

		data_buff := buf[consumeSize+HEAD_SIZE:consumeSize+HEAD_SIZE+data_len]

		// checksum
		if !ChecksumForMessage( checksum, data_buff ) {
			err = errors.New( "the package checksum to failed." );
			return
		}

		// uncrypto
		if 0 != head & CRYPTO_MARKER {
			//TODO: data_buff
		}

		// uncompress
		if 0 != head & COMPRESS_MARKER {
			//TODO: data_buff
		}

		message_buff := data_buff[MSG_OPCODE:];
		packet, perr := TransitionToPacket( head, checksum, first, second, message_buff );

		if perr != nil {
			err = fmt.Errorf( "message error. %s", perr.Error() );
			return
		}

		packets = append( packets, packet )

		consumeSize += data_len + HEAD_SIZE;
	}

	return
}

// transition packet from []byte
func TransitionToPacket( head, checksum uint32, first, second uint8, msg_buff []byte ) ( packet *Packet, err error ) {
	packet = nil;
	err = nil;

	// check opcode
	if !( uint32(first) < MAX_OPCODE_SIZE && uint32(second) < MAX_OPCODE_SIZE  ) {
		err = fmt.Errorf( "opcode: out of range: [%d:%d]", first, second );
		return
	}

	opcode_handle := gs_opcode_handlers[first][second];

	if opcode_handle == nil {
		return nil, fmt.Errorf( "unknown opcode: %d %d", first, second );
	}

	// new pb object from prototype.
	pb := reflect.New( reflect.TypeOf( opcode_handle.pb ).Elem() ).Interface();

	// deserialize pb from bytebuff
	if perr := proto.Unmarshal( msg_buff, pb.(proto.Message) ); perr != nil {
		err = fmt.Errorf( "Unmarshal Failed. %s", err.Error() );
		return
	}

	// check reality opcode.
	real_first, real_second := GetOpcode( pb.(proto.Message) );

	if first != real_first || second != real_second {
		err = fmt.Errorf( "the opcode [%d:%d] != [%d:%d]", first, second, real_first, real_second )
		return
	}

	return &Packet{ head, checksum, first, second, pb.(proto.Message) }, nil
}

// read uint32 from byte buffer.
// LittleEndian
func ReadUint32( buff []byte ) uint32 {

	val := uint32( buff[0] );
	val = val | uint32(buff[1]) << 8;
	val = val | uint32(buff[2]) << 16;
	val = val | uint32(buff[3]) << 24;

	return val
}

// write uint32 to byte buffer.
// LittleEndian
func WriteUint32( buff []byte, val uint32 ) {

	buff[0] = byte(val & 0x000000FF >> 0);
	buff[1] = byte(val & 0x0000FF00 >> 8);
	buff[2] = byte(val & 0x00FF0000 >> 16);
	buff[3] = byte(val & 0xFF000000 >> 24);
}

// encode the message to byte buffer.
func EncodeMessage( pb proto.Message, compress, crypto bool ) ( msg_buf []byte, err error ) {

	msg_buf = []byte{};

	buff, err := proto.Marshal( pb );

	if err != nil {
		err = fmt.Errorf( "encode message to falied. %s", err.Error() );
		return
	}

	// binary data.
	data_buf := make( []byte, uint32(len( buff )) + MSG_OPCODE )
	// opcode.
	data_buf[0], data_buf[1] = GetOpcode( pb );
	// message data.
	copy( data_buf[MSG_OPCODE:], buff );


	if compress {
		//TODO: compress data_buf
	}

	if crypto {
		//TODO: crypto data_buf
	}

	// packet data.
	msg_buf = make( []byte, uint32(len(data_buf)) + HEAD_SIZE );

	// write message head.
	WriteUint32( msg_buf[0:4], uint32(len( data_buf )) );
	WriteUint32( msg_buf[4:8], uint32(0) );	// Adler32( data_buf )

	// copy message data to packet.
	copy( msg_buf[8:], data_buf );

	return
}

// call message fucntion.
func CallOnMessage( packets []*Packet, connection IConnection ) {
	// call it.
	for _, packet := range( packets ) {
		// find the callback.
		opcode_handle := gs_opcode_handlers[packet.First][packet.Second];

		if opcode_handle.handle == nil {
			fmt.Println( "the opcode handle is nil: ", packet.First, ":", packet.Second );
			continue
		}

		opcode_handle.handle( packet.Message, connection );
	}
}
