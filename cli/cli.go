package cli

type Func func(args []string, opts []string) (string, error)

type Command struct {
	arguments []string
	options   []string
	function  Func
}
