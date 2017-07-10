package main

import (
	"fmt"

	proto "protot"
	proto_interface "protot_interface"
)

func main() {
	label := "Hello"
	t := int32(123)
	r := []int64{1,2,3}
	hello := proto.Test{
		Label: &label,
		Type: &t,
		Reps: r,
	}
	fmt.Println(hello.String())

	hello2 := proto.Test{
		Label: &label,
		Reps: r,
	}
	fmt.Println(hello2)
	fmt.Println(hello2.String())

	proto_interface.Print()
	proto_interface.PrintString()
	fmt.Println("hehe")
}
