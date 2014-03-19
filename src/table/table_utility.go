package table

import(
	"os"
	"io"
	"bytes"
	"fmt"
	"reflect"
	"encoding/binary"
)

/////////////////////////////////////////////////////////////////////////////////////////
/*
struct TableHeader
{
	ALL_SIZE_FUNC( size, data )
	DATA_SIZE_FUNC( size, data )

	uint32	version;
	uint16	size;
	char	data[0];
};

struct DataHeader
{
	int	data_size;
	int	string_size;
	int	offset;
};

struct RepeatData
{
	int	size;
	int	offset;
};

*/

type ITableFace interface {
	SizeOf() int32
}

const (
	SIZEOF_TABLE_HEADER	= 6		// uint32 + uint16
	SIZEOF_DATA_HEADER	= 12	// int * 3
	SIZEOF_REPEAT_DATA	= 8		// int * 2
)

/////////////////////////////////////////////////////////////////////////////////////////
// table header info.
type TableHeader struct {
	Version uint32
	Size	uint16
	Data	[]byte
}

func (this *TableHeader) AllSize() int {

	return SIZEOF_TABLE_HEADER + len( this.Data )
}

/////////////////////////////////////////////////////////////////////////////////////////
// data header
type DataHeader struct {
	DataSize	int32
	StringSize	int32
	Offset		int32
}

/////////////////////////////////////////////////////////////////////////////////////////
// repeat data
type RepeatData struct {
	Size	int32
	Offset	int32
}

/////////////////////////////////////////////////////////////////////////////////////////
// data pool
type DataPool struct {
	data []byte
}

func (this *DataPool) Data( offset, size int32 ) []byte {
	// check
	if !(int32(len(this.data)) > (offset + size)) {
		panic( "out of range." );
	}

	return this.data[offset:offset+size]
}

func (this *DataPool) DataEx( offset int32 ) []byte {
	if !(int32(len( this.data )) > offset) {
		panic( "out of range." )
	}

	return this.data[offset:]
}

func (this *DataPool) Len() int32 {
	return int32( len( this.data ) )
}

/////////////////////////////////////////////////////////////////////////////////////////
// string pool
type StringPool struct {
	data []byte
}

func ( this *StringPool ) String( offset int32 ) string {
	// check.
	if offset >= int32(len( this.data )) {
		panic( "offset out of range." );
	}

	str, err := bytes.NewBuffer(this.data[offset:]).ReadString( 0 );

	if err != nil {
		panic( err )
	}

	return str
}

func (this *StringPool) Len() int32 {
	return int32( len( this.data ) )
}

//######################################################################################
// load file to bytes.
func LoadFileToBytes( fpath string ) ( file_buf []byte, err error ) {

	file_buf = []byte{}
	err = nil

	// open tbl file. 
	file, e := os.OpenFile( fpath, os.O_RDONLY, 0666 );

	if e != nil {
		err = fmt.Errorf( e.Error() );
		return
	}

	defer file.Close();

	file_stat, e := file.Stat();

	if e != nil {
		err = e
		return
	}

	file_size := file_stat.Size()
	//fmt.Println( "file size: ", file_size, " bytes" );

	file_buf  = make( []byte, file_size );

	readSize := int64(0);

	for {
		if readSize >= file_size {
			break;
		}

		read, err := file.Read( file_buf[:] );

		if err != nil {
			panic( err )
		}

		readSize += int64(read);
	}

	return
}

//######################################################################################
// 

type TableLoader struct {
	file_buffer		[]byte

	table_header	TableHeader
	data_header		DataHeader
	data_pool		DataPool
	string_pool		StringPool
	object_array	RepeatData
}

// load
func (this *TableLoader) Load( proto interface{}, file string ) ( result []interface{}, ret bool) {
	// load file to buffer
	result = []interface{}{}
	ret = false

	if _, ok := proto.(ITableFace); !ok {
		panic( "that proto is not impl ITableFace!" )
		return
	}

	tf := proto.(ITableFace)

	var err error
	if this.file_buffer, err = LoadFileToBytes( file ); err != nil {
		panic( err )
		return
	}

	// preprocess 
	if !this.Prepare() {
		panic( "prepare tbl buffer to failed." )
		return
	}

	ret = true

	/*
	fmt.Println( "####################################################" );
	fmt.Println( "table_header: ", this.table_header );
	fmt.Println( "data_header: ", this.data_header );
	fmt.Println( "data_pool_len: ", len( this.data_pool.data ) );
	fmt.Println( "string_pool_len: ", len( this.string_pool.data ) );
	fmt.Println( "string test: ", this.string_pool.String( 1 ) );
	fmt.Println( "object_array: ", this.object_array )
	fmt.Println( "####################################################" );
	*/

	for i := int32(0); i < this.object_array.Size; i++ {
		obj_buf := this.data_pool.Data( i * tf.SizeOf(), tf.SizeOf())

		new_tbl := reflect.New( reflect.TypeOf( proto ).Elem() ).Interface()

		if pos, ok := this.Unmarshal( new_tbl, bytes.NewBuffer( obj_buf ) ); ok {
			//fmt.Println( "pos=> ", pos, "new => ", new_tbl )

			if pos != len( obj_buf ) {
			//	panic( fmt.Errorf( "pos != len(buf), [%d] %d != %d", i, pos, len( obj_buf ) ) )
			}

			result = append( result, new_tbl )
		}
	}


	return
}

// initialize buffer.
func (this *TableLoader) Prepare() bool {

	reader := bytes.NewBuffer( this.file_buffer[:] );

	// TableHeader
	binary.Read( reader, binary.LittleEndian, &this.table_header.Version );
	binary.Read( reader, binary.LittleEndian, &this.table_header.Size );
	//fmt.Println( this.table_header.Version, "   ", this.table_header.Size )
	this.table_header.Data = this.file_buffer[SIZEOF_TABLE_HEADER:this.table_header.Size+SIZEOF_TABLE_HEADER]

	// DataHeader Pointer.
	reader = bytes.NewBuffer( this.file_buffer[this.table_header.AllSize():] )

	// DataHeader
	binary.Read( reader, binary.LittleEndian, &this.data_header.DataSize );
	binary.Read( reader, binary.LittleEndian, &this.data_header.StringSize );
	binary.Read( reader, binary.LittleEndian, &this.data_header.Offset );

	data := this.file_buffer[ this.table_header.AllSize() + SIZEOF_DATA_HEADER: ];

	this.data_pool.data = data[:this.data_header.DataSize]
	this.string_pool.data = data[this.data_header.DataSize: this.data_header.DataSize + this.data_header.StringSize]

	// RepeatData
	reader = bytes.NewBuffer( this.data_pool.data[this.data_header.Offset:] )

	binary.Read( reader, binary.LittleEndian, &this.object_array.Size );
	binary.Read( reader, binary.LittleEndian, &this.object_array.Offset );

	return true
}

// 
func (this *TableLoader) Fill( t reflect.Type, v reflect.Value, obj_reader io.Reader ) int {

	pos := 0;

	switch v.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		{
			switch t.Bits() {
			case 8:
				{
					val := int8(0)
					binary.Read( obj_reader, binary.LittleEndian, &val );

					v.SetInt( int64(val) )
				}
			case 16:
				{
					val := int16(0)
					binary.Read( obj_reader, binary.LittleEndian, &val );

					v.SetInt( int64(val) )
				}
			case 32:
				{
					val := int32(0)
					binary.Read( obj_reader, binary.LittleEndian, &val );

					v.SetInt( int64(val) )
				}
			case 64:
				{
					val := int64(0)
					binary.Read( obj_reader, binary.LittleEndian, &val );

					v.SetInt( val )
				}
			default:
				panic( "unsupported!!" )
			}

			pos += t.Bits() / 8;
		}
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		{
			switch t.Bits() {
			case 8:
				{
					val := uint8(0)
					binary.Read( obj_reader, binary.LittleEndian, &val );

					v.SetUint( uint64(val) )
				}
			case 16:
				{
					val := uint16(0)
					binary.Read( obj_reader, binary.LittleEndian, &val );

					v.SetUint( uint64(val) )
				}
			case 32:
				{
					val := uint32(0)
					binary.Read( obj_reader, binary.LittleEndian, &val );

					v.SetUint( uint64(val) )
				}
			case 64:
				{
					val := uint64(0)
					binary.Read( obj_reader, binary.LittleEndian, &val );

					v.SetUint( val )
				}
			default:
				panic( "unsupported!!" )
			}

			pos += t.Bits() / 8;
		}
	case reflect.Float32:
		{
			val := float32(0)
			binary.Read( obj_reader, binary.LittleEndian, &val );

			pos += 4;
			v.SetFloat( float64(val) )
		}
	case reflect.Float64:
		{
			val := float64(0)
			binary.Read( obj_reader, binary.LittleEndian, &val );

			pos += 8;
			v.SetFloat( (val) )
		}
	case reflect.String:
		{
			val := int32(0)
			binary.Read( obj_reader, binary.LittleEndian, &val );

			pos += 4
			v.SetString( this.string_pool.String( int32(val) ) )
		}
	default:
		panic( "unsupported!!" )
	}

	return pos
}

//
func (this *TableLoader) Unmarshal( result interface{}, obj_reader io.Reader ) (int, bool) {
//	fmt.Println( "struct => ", reflect.TypeOf( result ) )
	typ := reflect.TypeOf( result )
	valp := reflect.ValueOf( result )

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if valp.Kind() == reflect.Ptr {
		valp = valp.Elem()
	}

	if typ.Kind() != reflect.Struct {
		panic( "must is struct" )
	}

	pos := 0

	for index := 0; index < typ.NumField(); index++ {

		//filed_type := typ.Field( index )
		filed_value := valp.Field( index )

		switch filed_value.Kind() {
		case reflect.Int8, reflect.Uint8, reflect.Int16, reflect.Uint16, reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.String:
			{
				off := this.Fill( filed_value.Type(), filed_value, obj_reader )

				if off <= 0 {
					panic( "asssert" )
				}

				pos += off;
			}
		case reflect.Slice:
			{
				//fmt.Println( "******[", filed_type.Name, " ", filed_value.Kind(), pos, "]*****" )

				repeat := &RepeatData{}
				binary.Read( obj_reader, binary.LittleEndian, &repeat.Size );
				binary.Read( obj_reader, binary.LittleEndian, &repeat.Offset );

				//fmt.Println( "size: ", repeat.Size, " offset: ", repeat.Offset )

				if 0 == repeat.Size {
					break
				}

				new_slice := reflect.MakeSlice( filed_value.Type(), int(repeat.Size), int(repeat.Size) )
				//fmt.Println( "slice type: ", reflect.TypeOf( new_slice.Index(0).Interface() ).Kind()  )

				slice_reader := bytes.NewBuffer( this.data_pool.DataEx( repeat.Offset ) )

				for index := 0; index < int(repeat.Size); index ++ {
					slice_t := new_slice.Index(index).Type()
					slice_v := new_slice.Index(index)

					switch reflect.TypeOf( new_slice.Index(index).Interface() ).Kind() {
					case reflect.Int8, reflect.Uint8, reflect.Int16, reflect.Uint16, reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.String:
						{
							off := this.Fill( slice_t, slice_v, slice_reader )

							if off <= 0 {
								panic( "assert" )
							}

							//fmt.Println( "slice read :", off )
						}
					case reflect.Slice:
						panic( "unsupported!!" )
					case reflect.Struct:
						{
							if _, ok := this.Unmarshal( slice_v.Addr().Interface(), slice_reader ); !ok {
								panic( "error." )
							}
						}
					default:
						panic( "unsupported!!" )
					}
				}

				filed_value.Set( new_slice )

				pos += SIZEOF_REPEAT_DATA;
			}
		case reflect.Struct:
			{
				off, ok := this.Unmarshal( filed_value.Addr().Interface(), obj_reader );

				if !ok {
					return 0, false;
				}

				pos += off
			}
		default:
			fmt.Println( "unsupported: ", filed_value.Kind() )
			panic( "unsupported!!!" )
		}
	}

	return pos, true
}
