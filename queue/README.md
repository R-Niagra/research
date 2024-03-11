Looks like lock-free queue implementation is no better in benchmarks. More benchmarking to follow...

I suspect GC cleanup overhead due to linked-list pointers in memory. Will memory profile to further diagnose it.

I'd also check overhead of compareAndSwap operation