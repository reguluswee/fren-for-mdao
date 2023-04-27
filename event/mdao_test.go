package event

import (
	"fmt"
	"testing"
)

func TestMdao(t *testing.T) {
	v1, v2 := BatchMdaoIssue()

	fmt.Println(v1, v2)
}
