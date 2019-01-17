# SQRL

An implementation of the [SQRL protocol](https://www.grc.com/sqrl/sqrl.htm), 
designed to be easy to integrate into a http server or SQRL client.

This is not production ready, please proceed with caution.

Simple Quick Reliable Login (SQRL) is a protocol designed and formalised by 
Steve Gibson of the [Gibson Research Corporation](https://www.grc.com). [Visit 
his site](https://www.grc.com/sqrl/sqrl.htm) for more information about the SQRL.

### SSP Example

The SQRL Service Provider (SSP) example is based on Steve's own example at 
[sqrl.grc.com](https://sqrl.grc.com/msa). To run the sample use the following;

```
$ cd ssp/example
$ go run *.go
```