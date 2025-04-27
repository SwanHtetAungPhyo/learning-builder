package main

import (
	"fmt"
	"unsafe"
)

// Type	Size (bytes)	Notes
// bool	1	true or false
// byte	1	alias for uint8
// rune	4	alias for int32, Unicode code point
// int8	1	signed 8-bit integer
// uint8	1	unsigned 8-bit integer
// int16	2	signed 16-bit integer
// uint16	2	unsigned 16-bit integer
// int32	4	signed 32-bit integer
// uint32	4	unsigned 32-bit integer
// int64	8	signed 64-bit integer
// uint64	8	unsigned 64-bit integer
// int	4 or 8	depends on architecture (32-bit or 64-bit)
// uint	4 or 8	depends on architecture (32-bit or 64-bit)
// uintptr	4 or 8	enough to hold pointer values
// float32	4	32-bit floating point number
// float64	8	64-bit floating point number
// complex64	8	two float32 (real and imaginary parts)
// complex128	16	two float64 (real and imaginary parts)
// string	16	2 words: pointer + length (on 64-bit syst
//

//Memory:
//0xA010: 00000001  (a, lowest byte)
//0xA011: 00000000
//0xA012: 00000000
//0xA013: 00000000
//0xA014: 00000001  (b, true)
//

//  HardWare
// CPU
// Register
// RAM  (random access memory )
// ROM ( read only memory)

// Kernal
// OS

func main() {
	var a int8 = 1
	var b int16 = 2
	var c int32 = 3
	var d int64 = 4

	fmt.Println("Variable Info:")
	fmt.Println("-------------------------------")
	fmt.Printf("a: value = %d, address = %p, size = %d bytes\n", a, &a, unsafe.Sizeof(a))
	fmt.Printf("b: value = %d, address = %p, size = %d bytes\n", b, &b, unsafe.Sizeof(b))
	fmt.Printf("c: value = %d, address = %p, size = %d bytes\n", c, &c, unsafe.Sizeof(c))
	fmt.Printf("d: value = %d, address = %p, size = %d bytes\n", d, &d, unsafe.Sizeof(d))
	fmt.Println("-------------------------------")
	var a int = 1
}
// a = 1
// []
// 10 = 2,

13 /2 = remainder = 1 ,
component = 6 / 2 = remiander = 0
component = 3/ 2 = remiander = 1
component = 1/2 = reminder = 1

13 base 10
base 2 1101
2^0 = 1
2^2 = 4
2^3 = 8

111 binary 7
11011

1 = 2^0 = 1
1 = 2^1 = 2
0

1 = 2^3 = 8
1 = 2^4 = 16
base 10  =