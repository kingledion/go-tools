# queue

A single ended queue with push and pop functionality implemented with channels for built in thread safety. 

Because this queue is type specific, it should probably not be imported, and instead simply implemented in a project that uses this. Copy the implementation and retype as needed. For this reason, the struct is implemented as private.

There are optional add-ons for this implementation such as a queue counter to prevent the queue from being cycled through excessively. An example implementation is included in `run.go`.