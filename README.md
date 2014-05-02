RATS: Remote Android Test Service
-----

![Imgur](http://i.imgur.com/7skL9wq.jpg)

What and Why?
----
RATS is an easy to run Test Runner Service that can be hosted anywhere a device can be connected.  It was designed to thwart device hogs and provide a level playing field for shops who can't afford more expensive testing services.

### Who?
So far mostly me, but not without support from [VOKAL Interactive](http://www.vokalinteractive.com)


###Installation

Binaries:

Go-Pros:

1. `go get github.com/wmbest2/rats_server`
2. `rats_server`

### Run Your Tests

  `curl -X POST myserver.local:3000/api/run -F apk=@myapk.apk -F test-apk=@myapk-test.apk`

  _Seriously, that's it!_


###Capabilities

* Concurrent Test Running on all devices connected
* Filter Devices for a particular run
 * Device count filter so you don't piss off your coworkers 
 * Strict mode and Serial Specification for cases where you must run a particular subset of devices
* Easy to use Api for extra info

###Whats Missing / What's to come

* Proper encoding for Accept headers (e.g. junit for xml)
* Serious Go know-how (this was a learning opportunity for me)
* Log capture from tests
* Screenshots
* Better Filtering/Scripting
* Project Grouping

### Let's see more!

![Imgur](http://i.imgur.com/MrQX3Gz.jpg)

![Imgur](http://i.imgur.com/73hH7Qd.jpg)

### Extras

* CLI for running tests against all devices and spitting out the results

