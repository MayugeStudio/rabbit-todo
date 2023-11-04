package param

type Type int

const (
	STRING Type = iota
	INT
	BOOL
)

func ParameterTypeToString(paramType Type) string {
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
