// Generated by the proto compiler.  DO NOT EDIT!
// author: limpo1989@gmail.com
// source: item_proto.prot
// date  : 2014/02/20 11:53:39

package table

type ItemBase struct {
    Id                 uint32    // 物品ID
    Name               string    // 物品名称
    InitLv             int8      // 物品等级
    MainType           float64   // 物品主类型
    SubType            int8      // 物品子类型
    Rand               uint8     // 是否随机
    ModelPath          string    // 模型路径
    PicturePatch       string    // 图标路径
    DropModelPath      string    // 掉落模型路径
    ItemsMaterialPatch string    // 美术材质路径
    ItemsMaterial      uint8     // 物品材质
    SellOrNot          uint8     // 能否出售
    Overlap            uint16    // 叠加上限
    Discard            int32     // 能否丢弃
    Business           uint8     // 能否交易
    Binding            int8      // 绑定类型
    CancelBinding      uint16    // 能否解绑
    Lock               int16     // 能否锁定
    AutoPick           uint64    // 自动拾取
    Resolve            int64     // 能否分解
    Characterization   string    // 文字描述
    SellPrice          uint32    // 卖出价格
    Disappear          uint8     // 使用后消失
    Quality            uint8     // 物品品质
    NameListAa         []string  // 适用主职
    Auxiliary []struct {
        AuxiliaryLvId struct {
            AuxiliaryId uint8
            AuxiliaryLv uint8
        }
    } // 副职限制
    CanUse             uint8     // 能否主动使用
    GoldSell           uint8     // 是否金币买卖
    PicturePatchNull   string    // 图标路径无框
    User struct {
        UserId uint32
        Role   uint8
    } // 限制职业
    EntreatMainType    uint8     // 求购行主类型
    EntreatSubType     uint8     // 求购行子类型
    BusinessAstrict    uint8     // 交易限制
}

// -------------------------------------------------------------------
//member methos of ItemBase
func (this *ItemBase) Key() uint32 {
    return this.Id
}

func (this *ItemBase) SizeOf() int32 {
    return int32(105)
}

// ===================================================================
type ItemBaseManager struct {
    data []*ItemBase
}

//member methos of ItemBaseManager
func (this *ItemBaseManager) Source() string {
    return "item_proto.tbl"
}

func (this *ItemBaseManager) Size() int {
    return len(this.data)
}

func (this *ItemBaseManager) Get( index int ) *ItemBase {
    if index >= this.Size() {
    	panic( "out of range" )
    }

    return this.data[index]
}

func (this *ItemBaseManager) Load(path string) bool {

    if this.Size() > 0 {
    	return true
    }

    loader := &TableLoader{}

    path += "/"
    path += this.Source()
    if result, ok := loader.Load( &ItemBase{}, path ); ok {

        for _, v := range( result ) {
        	this.data = append( this.data, v.(*ItemBase) )
        }

        return true
    }

    return false
}

func (this *ItemBaseManager) Find( key uint32 ) *ItemBase {

    if this.Size() <= 0 {
    	return nil
    }

    start	:= 0
    stop	:= int(this.Size() - 1)
    middle	:= 0

    for ; start <= stop ; {

    	middle = int( (start + stop) / 2 )

    	tbl := this.data[middle]

    	if tbl.Key() == key {
    		return tbl
    	}

    	if tbl.Key() > key {
    		stop = middle - 1
    	} else {
    		start = middle + 1
    	}
    }

    return nil
}

// ===================================================================
var gs_ItemBaseInstance = &ItemBaseManager{}

func GetItemBaseManager() *ItemBaseManager {
	return gs_ItemBaseInstance
}

func init() {
	LoadTables( gs_ItemBaseInstance )
}
