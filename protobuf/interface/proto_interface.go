package protot_interface

import (
	"fmt"
	proto "protot"
)

var (
	label = "Hello"
	t = int32(111)
	p = proto.Test{
		Label: &label,
		Type: &t,
		Reps: []int64{1,2,3},
	}
	p3 = proto.Test3{
		Label: "Hello",
		Type: 111,
		Reps: []int64{1,2,3},
	}
)

func Print() {
	fmt.Println(p)
	fmt.Print("Proto3: ")
	fmt.Println(p3)
}

func PrintString() {
	fmt.Println(p.String())
	fmt.Print("Proto3: ")
	fmt.Println(p3.String())
}
