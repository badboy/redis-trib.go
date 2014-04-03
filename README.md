redis-trib.go
=============

Create and administrate your [Redis Cluster][cluster-tutorial] from the Command Line.

Inspired heavily by [ruster][] and the original [redis-trib.rb][].

## Dependencies

* [redigo][]
* [cli][]

Dependencies are handled by [gpm][], simple install it and type `gpm` to fetch them.

## Install

~~~haskell
git clone https://github.com/badboy/redis-trib.go.git
cd redis-trib.go
make
~~~

## How to use

~~~haskell
# Execute a command on each cluster node
./redis-trib each 127.0.0.1:7001 info memory
# Check that the cluster is ok
./redis-trib check 127.0.0.1:7001
~~~

## Roadmap

This project needs a lot of work. I listed some things in no particular order:

* Implement more subcommands
  * `create`
  * `add-node`
  * `del-node`
  * `reshard`
  * `fix`
* Find a nice and easy way to implement these subcommands
* Write tests
* Better logging functions
* Documentation about the code
* Documentation about how to use
* License, contributions guideline, ...
* Release a proper version

## State of the code

This is the first Go code I ever wrote. It is not really good, the program is not anywhere near to be a complete replacement for either ruster or redis-trib.rb and it lacks a few basic things (proper error checks, cleaned-up code, tests, comments, ...)

I will try to improve it and I welcome all ideas, bug reports or code improvements.
Just open an [issue][], drop me a message on [twitter][] or write an email.

[cluster-tutorial]: http://redis.io/topics/cluster-tutorial
[ruster]: https://github.com/inkel/ruster
[redis-trib.rb]: https://github.com/antirez/redis/blob/unstable/src/redis-trib.rb
[issue]: https://github.com/badboy/redis-trib.go/issues
[twitter]: https://twitter.com/badboy_
[gpm]: https://github.com/pote/gpm
[redigo]: https://github.com/garyburd/redigo/
[cli]: https://github.com/codegangsta/cli
