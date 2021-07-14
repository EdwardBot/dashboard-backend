package utils

import (
	"log"
	"strconv"
	"strings"
)

func ArrayToPSQL(in []int64) string {
	tmp := "{"
	for i := range in {
		tmp = tmp + strconv.FormatInt(in[i], 10)
		if i != len(in)-1 {
			tmp = tmp + ","
		}
	}
	log.Println(in)
	log.Println(tmp)
	return tmp + "}"
}

func PSQLToArray(in string) []int64 {
	a := strings.ReplaceAll(strings.ReplaceAll(in, "{", ""), "}", "")
	aa := strings.Split(a, ",")
	c := make([]int64, len(aa))
	for b := range aa {
		bb := strings.ReplaceAll(aa[b], "'", "")
		cc, _ := strconv.ParseInt(bb, 10, 64)
		c[b] = cc
	}
	log.Println(in)
	log.Println(c)
	return c
}

func IToSArray(in []int64) []string {
	a := make([]string, len(in))
	for k := range in {
		a[k] = strconv.FormatInt(in[k], 10)
	}
	return a
}
