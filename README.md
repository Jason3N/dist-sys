## 3/12

- currently scrapping around go-routines that utilize concurrency and seeing the effects of parralelism! specifically working on web-requests.. inspired by uc davis' website always pooping itself whenever registration happens, will want to test out a system in which can handle 8000 high-stress tests
- will want to create a website similar to that of signing up to davis registration, letting users sign in and queue up for a registration slot (i really want to learn more software architecture related to this, how to handle large spikes of activity.. will do more research later)
- goroutines take advantage of lightweight threads. these are threads that are language-based as go uses a scheduler to do these routines in concurrency (heavy-weight threads are more for embedded / OS / heavy tasks)


## 3/14

- added a waitgroup and used go routines to proc a function that runs way faster.. way faster than 1 second, need to figure that out because i moved a wait function inside of a go func.. will check out later but if not malicious will continue with this approach
