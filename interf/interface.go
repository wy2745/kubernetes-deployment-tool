package interf

import (
	abc "../type124"
	"fmt"
)

type podinface interface {
	GetName() string
	GetNodeName() string
}

func haha() {
	fmt.Print(abc.ClaimBound)
}