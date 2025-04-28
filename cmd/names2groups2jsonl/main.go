package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	ng "github.com/takanoriyanagitani/go-names2grouped"
	. "github.com/takanoriyanagitani/go-names2grouped/util"
)

var envValByKey func(string) IO[string] = Lift(
	func(key string) (string, error) {
		val, found := os.LookupEnv(key)
		switch found {
		case true:
			return val, nil
		default:
			return "", fmt.Errorf("env var %s missing", key)
		}
	},
)

var names ng.Names = ng.StdinToNames()

var commonPrefixByteLen IO[int] = Bind(
	envValByKey("ENV_COMMON_PREFIX_LEN"),
	Lift(strconv.Atoi),
)

var groupGen IO[ng.IsSameGroup] = Bind(
	commonPrefixByteLen,
	Lift(func(i int) (ng.IsSameGroup, error) {
		return ng.FirstBytesCheckNew(uint16(i)), nil
	}),
)

var grouped IO[ng.GroupedNames] = Bind(
	groupGen,
	Lift(func(g ng.IsSameGroup) (ng.GroupedNames, error) {
		return g.NamesToGrouped(names), nil
	}),
)

var grouped2stdout IO[Void] = Bind(
	grouped,
	Lift(func(n ng.GroupedNames) (Void, error) {
		return Empty, n.ToJsonToWriter(os.Stdout)
	}),
)

func main() {
	_, e := grouped2stdout(context.Background())
	if nil != e {
		log.Printf("%v\n", e)
	}
}
