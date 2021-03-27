# ParallelCoin

## Plan 9 Hard Fork Specifications

### 1. Hash Function

DivHash is a hash function designed to never be made into an ASIC. It 
requires integer multiplication and division, exploiting the fact that long 
division is an inherently iterative process and cannot be parallelised.

To perform a DivHash - this is one cycle, which is repeated another 2 times:

1. cut block in half, reverse the first half
2. append reversed half to end of block
3. same thing for the other half
4. square first and second extended halves
5. multiply these two squared halves together
6. divide the result by the original block (or result of previous cycle)

After another two rounds, hash with Blake 3 256 bit hash function to reduce to 
the solution size.

Proof of Work hash functions need to have strong resistance to collisions, 
apart from this, the ways of shuffling the bits is mostly arbitrary. The 
described cycle will have an effect similar to photocopying a photocopy, and 
so on and at some point the text will become illegible. So the 
multiplication/division cycles 'rot' some of the entropy, but in a specific 
way that very likely won't be easily optimized, nor reversed.

DivHash also should have quite high memory-hardness, as its 40kb results 
will overflow most CPU's L1 caches meaning the calculation has a component 
of cache latency. I estimate that my old Ryzen 7 processor could move 14gb 
of cache back and forth at a time per second, so I suppose absolute 
structural limit of 350,000 cycles per second in addition to consuming about 
1/3 of the cpu's cycles nonstop.

### 2. Prime Based Parallel Block Schedule

Instead of having one single target interval, Plan 9 HardFork has 9, each 
interval created from the first 9 prime numbers (2, 3, 5, 7, 11, 13, 17, 19, 
23) and multiplied by a base of 9 seconds.

The formula to calculate the average of the 9 starts with the shortest (18 
second) block interval, divided by itself, to get 1, then repeat this with 
each longer interval

    18 / 18 => 1
    18 / 27 => 0.666666666666666
    18 / 45 => 0.4
    18 / 63 => 0.285714286
    18 / 99 => 0.181818182
    18 / 117 => 0.153846154
    18 / 171 => 0.105263158
    18 / 153 => 0.117647059
    18 / 207 => 0.086956522

The sums on the left each generate a value


1 + 0.666666666666666 + 0.4 + 0.285714286 + 0.181818182 + 0.153846154 + 0.105263158 + 0.117647059 + 0.086956522  
