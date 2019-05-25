# modular
A modular arithmetic library with an emphasis on speed
[![Build Status](https://travis-ci.org/stewi1014/modular.svg?branch=master)](https://travis-ci.org/stewi1014/modular)

Float64
[![GoDoc](https://godoc.org/github.com/stewi1014/modular/modular64?status.svg)](https://godoc.org/github.com/stewi1014/modular/modular64)

Float32
[![GoDoc](https://godoc.org/github.com/stewi1014/modular/modular32?status.svg)](https://godoc.org/github.com/stewi1014/modular/modular32)


Modular tries to leverage pre-computation as much as possible for calculating modulus and indexing floats, leveraging [fastdiv] and pre-computed lookup tables. I can't test it on all hardware, but in principle should perform better than traditional modulo functions on all but the strangest of hardware.

For example, on my machine with a modulo of 1e-20;

**float64**

| Number | Math.Mod | Modulus.Congruent | Indexer.Index |
| ------ | ------ | ------ | ------ |
| 0 | 14.2 ns/op | 7.68 ns/op | 11.2 ns/op |
| 20 | 638 ns/op | 57.0 ns/op | 15.9 ns/op |
| 1e300 | 8153 ns/op | 590 ns/op | 144 ns/op |


**float32**

| Number | Math.Mod | Modulus.Congruent | Indexer.Index |
| ------ | ------ | ------ | ------ |
| 0 | 13.8 ns/op | 8.78 ns/op | 11.3 ns/op |
| 20 | 481 ns/op | 17.9 ns/op | 13.6 ns/op |
| 1e20 | 1000 ns/op | 17.8 ns/op | 13.8 ns/op |

***

Modulus.Congruent is called Congruent as it is as accurate to euclidian division as possible, and finds the number within the range 0 <= n < modulus that is representative of the congruency class. While many implementations of the modulo function approximate this, many truncate, floor or otherwise don't perform division true to euclidian mathematics and hence often return negative numbers.

[fastdiv]: <https://github.com/bmkessler/fastdiv>
