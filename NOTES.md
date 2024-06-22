# Notes


Recursion example:

```

n[0-2]
n0, n1, n2

n[1-3,4]
(n,[1-3])   n4
n1, n2, n3, n4

n[1-3,t[5-6]]
(n,[1-3])   (n,t[5-6])
(n1,n2,n3]) (n,     (t,[5-6]))
(n1,n2,n3]) (n,     (t5,t6))
(n,[1-3])   nt5, nt6
n1,n2,n3,nt5,nt6

```
