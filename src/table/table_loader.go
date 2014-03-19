package table

import( "fmt" )

type ITableManager interface {
	Source() string
	Load( path string ) bool
}


var tableFolder = "./table"

func LoadTables( itable ITableManager ) bool {

	fmt.Print( "load table: ", itable.Source(), " -> " )

	ret := itable.Load( tableFolder )

	if ret {
		fmt.Println( "Ok." )
	} else {
		fmt.Println( "Failed." )
	}

	return ret
}
