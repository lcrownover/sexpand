# sexpand-go

SLURM node name expander library and binary

The SLURM HPC scheduler allows specifying node lists using the following syntax:

```text
n1-n4               =>  n1,n2,n3,n4
n[01-04]            =>  n01,n02,n03,n04
n[02-03,09-11],n01  =>  n01,n02,n03,n09,n10,n11
```

## Usage

This library exposes a single function, `SExpand`, that takes a string and
returns a list of hostnames expanded from that string:

```go
hostnames, err := SExpand('n[02-03,09-11],n01')
if err != nil {
    panic("invalid hostname string")
}
fmt.Println(hostnames) // prints: [n01,n02,n03,n09,n10,n11]
```
