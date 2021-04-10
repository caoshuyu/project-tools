package tools

//将MySQL数据库类型与Go类型做转换
func DbDataTypeChange(dataType string) (newDataType string) {
	switch dataType {
	case "varchar":
		newDataType = "string"
	case "datetime":
		newDataType = "time.Time"
	case "date":
		newDataType = "time.Time"
	case "time":
		newDataType = "time.Time"
	case "char":
		newDataType = "string"
	case "text":
		newDataType = "string"
	case "float":
		newDataType = "float64"
	case "double":
		newDataType = "float64"
	case "int":
		newDataType = "int64"
	case "bigint":
		newDataType = "int64"
	case "tinyint":
		newDataType = "int64"
	case "longtext":
		newDataType = "string"
	case "mediumint":
		newDataType = "int64"
	case "smallint":
		newDataType = "int64"
	case "mediumtext":
		newDataType = "string"
	case "timestamp":
		newDataType = "time.Time"

	default:
		newDataType = dataType
	}
	return
}
