# (R)ead-(C)opy-Update
read copy update map for golang 1.18+

## How it works
This is a simple generic implementation for https://en.wikipedia.org/wiki/Read-copy-update

Inserts are slow(er) while reading is fast. We always return a read-copy of the map to the consumers, which only requires atomic access to the map at the moment of changing. 

If an update to the map happens, we then update the Read Model and the next time it is requested, it will yield the updated copy.

## Caveats of the Read-Copy-Update model
- The read-copy of the model might contain outdated values. This is suitable to return to clients when they request current values
- it is not suitable if you need the consistent underlying data. In that case a `sync.Map` is what you need.
- The current read-model is shared. You should not modify data in it as this will still not be thread-safe and will modify the map for all current users. 

## How to get it 

`go get github.com/mier85/rcu`



