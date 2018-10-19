## Collatz Sequence Homework
### Brendan Horng

### Cache vs. Uncached

The uncached version of my program ran considerably faster than the cached version of my program.
The reason for this is probably due to the fact that the goroutines must be being halted while waiting for access to the lock.
This results in more context switches and waiting around for the lock.

This tradeoff is sometimes worth it especially when you consider very long operations.
Waiting for that cached value might be able to save a significant amount of time
However, in our case, that wait time might have been better served just actually recalculating the value.
To make it worse, we might wait for the lock and then realize that it isn't even cached which means we have to calculate it anyway.

This idea of caching results might be better suited for operations that are more expensive than our Collatz Length calculation.

### Buffered vs. Unbuffered Channels

In my code, I used buffered channels because I had 5 processes reading from a single channel and writing to a single channel.
While the rate that the collector thread was consuming values was faster than the rate at which any one of the worker threads were putting in `CollatzPair`s, there are still 5 workers and only one collector.
So to minimize blocking, I buffered the output to give the collector some time to consume and allow the workers to continue to the next task.

For a similar reason on the input, I only had one thread generating values, but five of them consuming it.
To reduce the wait time for a value, I buffered the input thread so, while the workers were calculating Collatz Lengths, the generator could place in more values.

### Other Designs
