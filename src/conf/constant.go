package conf

const SERVER_NAME = "init-project"

//明确枚举标记，不使用iota
const (
	DAO_ACTION_ADD         = 1 //添加insert
	DAO_ACTION_UPDATE      = 2 //更新
	DAO_ACTION_SELECT      = 3 //查询
	DAO_ACTION_SELECT_PAGE = 4 //分页查询
	DAO_ACTION_BATCH_ADD   = 5 //批量添加
)
