# retry

Implements a retry queue with an input channel and internal buffered channel. The retry queue will continually attempt to retry a given action on the elements of that queue, until both no more input has been recieved for a number of seconds and the length of the queue is unchanged for at least one cycle through the queue. 

Because this queue is type specific, it should probably not be imported, and instead simply implemented in a project that uses this. Copy the implementation and retype as needed. For this reason, the struct is implemented as private and typed as integer.