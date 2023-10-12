# modular
A math library with an emphasis on speed.

[![Status](https://github.com/stewi1014/modular/actions/workflows/go.yml/badge.svg)](https://github.com/stewi1014/modular/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/stewi1014/modular.svg)](https://pkg.go.dev/github.com/stewi1014/modular)


Modular tries to leverage pre-computation as much as possible to allow direct computation in Mod() and Index(), using [fastdiv] and pre-computed lookup tables. I can't test it on all hardware, but in principle should perform better than traditional modulo functions on all but the strangest of hardware.

For example, on my machine with a modulo of 1e-25;

**float64**

| Number | Math.Mod | Modulus.Mod | Indexer.Index |
| ------ | ------ | ------ | ------ |
| 0 | 14.1 ns/op | 6.07 ns/op | 12.4 ns/op |
| 2.5e-25 | 29.2 ns/op | 14.4 ns/op | 15.7 ns/op |
| 1 | 467 ns/op | 37.5 ns/op | 38.1 ns/op |
| 1e300 | 8063 ns/op | 36.3 ns/op | 39.4 ns/op |


**float32**

| Number | Math.Mod | Modulus.Mod | Indexer.Index |
| ------ | ------ | ------ | ------ |
| 0 | 15.6 ns/op | 5.49 ns/op | 11.3 ns/op |
| 2.5e-25 | 33.5 ns/op | 12.1 ns/op | 11.3 ns/op |
| 1 | 766 ns/op | 12.0 ns/op | 10.9 ns/op |
| 1e25 | 1259 ns/op | 12.0 ns/op | 11.0 ns/op |

***

Modular implements the euclidian modulo function; that is, it finds the number satisfying '0 <= n < modulus' that is representative of the given number's congruency class.


[fastdiv]: <https://github.com/bmkessler/fastdiv>
