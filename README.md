# modular
A modular arithmetic library with an emphasis on speed
[![Build Status](https://travis-ci.org/stewi1014/modular.svg?branch=master)](https://travis-ci.org/stewi1014/modular)

Float64
[![GoDoc](https://godoc.org/github.com/stewi1014/modular/modular64?status.svg)](https://godoc.org/github.com/stewi1014/modular/modular64)

Float32
[![GoDoc](https://godoc.org/github.com/stewi1014/modular/modular32?status.svg)](https://godoc.org/github.com/stewi1014/modular/modular32)


Modular tries to leverage pre-computation as much as possible to allow direct computation in Congruent() and Index(), using [fastdiv] and pre-computed lookup tables. I can't test it on all hardware, but in principle should perform better than traditional modulo functions on all but the strangest of hardware.

For example, on my machine with a modulo of 1e-25;

**float64**

| Number | Math.Mod | Modulus.Congruent | Indexer.Index |
| ------ | ------ | ------ | ------ |
| 0 | 12.1 ns/op | 5.70 ns/op | 16.4 ns/op |
| 2.5e-25 | 25.8 ns/op | 21.2 ns/op | 19.5 ns/op |
| 1 | 569 ns/op | 56.1 ns/op | 49.4 ns/op |
| 1e300 | 9797 ns/op | 56.7 ns/op | 50.2 ns/op |


**float32**

| Number | Math.Mod | Modulus.Congruent | Indexer.Index |
| ------ | ------ | ------ | ------ |
| 0 | 12.4 ns/op | 5.09 ns/op | 12.9 ns/op |
| 2.5e-25 | 29.6 ns/op | 20.1 ns/op | 15.8 ns/op |
| 1 | 766 ns/op | 20.5 ns/op | 14.4 ns/op |
| 1e25 | 1240 ns/op | 22.2 ns/op | 16.5 ns/op |

***

I've use the name 'Congruent' as it's a more explicit definition of the function, and helps avoid confusion with other functions. It is the same as a euclidian 'modulo' function; that is, it finds the number satisfying '0 <= n < modulus' that is representative of the given number's congruency class, hence the name Congruent.


[fastdiv]: <https://github.com/bmkessler/fastdiv>
