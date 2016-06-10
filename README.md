# Backoff

[![GoDoc](https://godoc.org/github.com/efritz/backoff?status.svg)](https://godoc.org/github.com/efritz/backoff)
[![Build Status](https://secure.travis-ci.org/efritz/backoff.png)](http://travis-ci.org/efritz/backoff)
[![codecov.io](http://codecov.io/github/efritz/backoff/coverage.svg?branch=master)](http://codecov.io/github/efritz/backoff?branch=master)

Algorithms to generate intervals.

## Example

`BackOff` is an interface which implements the `Next` and the `Reset` methods.
The `Next` method will return a `time.Duration` according to the algorithm of
the back-off and the back-off's current state (generally, the number of times
the method has been called). The `Reset` method resets any such state.

```go
b := NewLinearBackOff(time.Second, time.Second, time.Minute)

b.Next() // time.Second * 1
b.Next() // time.Second * 2
b.Next() // time.Second * 3

// ...

b.Next() // time.Minute
b.Next() // time.Minute
b.Next() // time.Minute

// ...

b.Reset()
b.Next() // time.Second * 1
```

Four algorithms are provided. `ZeroBackOff` and `ConstantBackOff` return a
constant duration on each call to `Next`, and `Reset` is a no-op.

`LinearBackOff`, shown above, returns a linearly increasing duration according
to a minimum interval, and a maximum interval, and the interval to increase by
on each call.

`ExponentialBackOff` returns an exponentially increasing duration according to
a minimum interval, a maximum interval, a multiplier, and a random factor. The
multiplier dictates the *base* interval - e.g. (*min* * *multiplier* ^ *n*) on
the the *nth* attempt, and the random factor dictates the interval's *jitter*
such that the interval value *i* is randomized around *i* - *i* * *jitter* and
*i* + *i* * *jitter*.

## License

Copyright (c) 2016 Eric Fritz

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
