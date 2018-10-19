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

So, unless I implemented the design incorrectly or missed something, I believe that this design has the following problem:

Each calculation of a Collatz Length, while not instantaneous, is fast enough that advanced techniques like parallelization and caching result in an overhead that actually makes the operation slower.

A potential other solution is to recognize that if the non-parallel, iterative solution was faster, we can decompose our search space into smaller problems.

The following solution I tested with an alternative version (for my own amusement and learning) of the Collatz code and I did the following to achieve faster results than any of my other solutions:

1. Write a generic function, `maxCollatz`, that given a set of numbers, for me this was described with a `start`, `end`, and `step` parameters, determines the maximum Collatz Length in the set.
2. Split up your larger search space into sub search spaces.
3. Run `maxCollatz` on each of the subspaces in separate goroutines and output their result to a channel.
4. Determine which is the maximum of these subspaces much like the collector from our other architecture.

#### Problems

This method requires you to have a pretty reasonable way to split up your search space.
The reason for this is if the search spaces aren't split up well, we could potentially have a situation where one of the search spaces is running for a lot longer than the others.
The more discrepancy in how long each sub search takes there is, the worse this method is.
Having high differences in the completion times essentially makes this method similar to just the iterative solution with some additional overhead.

The method I used to split up the search spaces was to evenly distribute the subspaces such that they all had integers of varying size.
I achieved this by offsetting their start points and having them step in uniform step sizes.
Example: a two process run from 1 to 10 would split up the sets in the following way
```
[1, 3, 5, 7, 9]
[2, 4, 6, 8, 10]
```

I did this because I made the assumption that larger numbers have larger Collatz lengths.
There are probably better ways to split the search space, and finding a good way to split up the search space is crucial for my proposed method.
