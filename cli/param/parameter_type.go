package param

type ParameterType int

const (
	STRING ParameterType = iota
	INT
	BOOL
)

func ParameterTypeToString(paramType ParameterType) string {
	switch paramType {
	case BOOL:
		return "bool"
	case INT:
		return "int"
	case STRING:
		return "string"
	default:
		return "unknownType"
	}
}
