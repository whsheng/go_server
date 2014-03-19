package conn

import(
	"fmt"
	"net"
//	"runtime"
	"time"
	"container/list"
	"os"
	"runtime/pprof"
)

type IServer interface {
	Start( string, uint16 ) bool
	Stop()
	Run()
}

type Server struct {
	addr		*net.TCPAddr
	lintener	*net.TCPListener
	conn_list	*list.List
	rm_list		*list.List
}

// create server
func NewServer( ip string, port uint16 ) *Server {

	server := new(Server)
	server.conn_list = list.New()
	server.rm_list = list.New()

	if !server.Start( ip, port ) {
		fmt.Println( "server startup failed." );
		return nil
	}

	return server
}


// start listener
func (this *Server) Start( ip string, port uint16 ) bool {

	if this.lintener != nil {
		fmt.Println( "TCPListener exists" );
		return false
	}

	addr, err := net.ResolveTCPAddr( "tcp4", fmt.Sprintf( "%s:%d", ip, port ) )

	if err != nil {
		fmt.Println( "ResolveTCPAddr failed. ", err.Error() );
		return false
	}

	this.addr = addr;

	listener, err := net.ListenTCP( "tcp", this.addr );

	if err != nil {
		fmt.Println( "ListenTCP to failed. ", err.Error() )
		return false
	}

	this.lintener = listener

	fmt.Println( "server is listen on ", ip, port )
	go this.RunAccept()

	return true
}

// stop listener
func (this *Server) Stop() {

	if this.lintener != nil {
		tmp := this.lintener;
		this.lintener = nil

		tmp.Close()
	}
}

// accept handle.
func (this *Server) RunAccept() {

	defer this.Stop();

	for {
		c, err := this.lintener.Accept();

		if err != nil {
			if this.lintener != nil {
				fmt.Println( "accept to failed. ", err.Error() );
				break;
			}

			fmt.Println( "acceptor closed." )
			break;
		}

		conn := NewConnection( this, c )
		this.conn_list.PushBack( conn );

		fmt.Println( "accept connection: ", c.RemoteAddr() );
	}
}

// run
func (this *Server) Run() {

	const DELAY_TIME = time.Duration( time.Millisecond * time.Duration(35) )

	f, err := os.Create( "cpu_profile" );

	if err != nil {
		return
	}

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	for {
		for elm := this.conn_list.Front(); elm != nil; elm = elm.Next() {
			conn := elm.Value.(*Connection)

			conn.Process()

			if conn.CheckTerminate() {
				conn.Close()
				this.rm_list.PushBack( elm )
			}
		}

		// cleanup
		if this.rm_list.Len() > 0 {
			for elm := this.rm_list.Front(); elm != nil; elm = elm.Next() {
				this.conn_list.Remove( elm.Value.(*list.Element) )
			}
			this.rm_list.Init()
			break
		}

		time.Sleep( DELAY_TIME )
		fmt.Println( "current connection: ", this.conn_list.Len() )
	}
}
