RATS: Remote Android Test Service
-----

![Imgur](http://i.imgur.com/s9Dl1ih.png)

What and Why?
----
RATS is an easy to run Test Runner Service that can be hosted anywhere a device can be connected.  It was designed to thwart device hogs and provide a level playing field for shops who can't afford more expensive testing services.

### Who?
So far mostly me[(wmbest2)](http://www.github.com/wmbest2), but not without support from [VOKAL Interactive](http://www.vokalinteractive.com)

###Installation

Binaries:

[Temporary Dropbox](https://www.dropbox.com/sh/z0spjt91pqufyh1/sOZ7cJ34-i)

Go-Pros:

1. `go get github.com/wmbest2/rats_server`
2. `rats_server -port 8080 -db mongodb://somehost`

### Run Your Tests

  `curl -X POST myserver.local:8080/api/run -F apk=@myapk.apk -F test-apk=@myapk-test.apk`

  _Seriously, that's it!_

  Parameters:

  * `apk` Main apk
  * `test-apk` Test package apk
  * `count` Number of devices
  * `serials` Comma separated list of device serials
  * `strict` Strict mode (will run forever if devices dont match)

###Capabilities

* Concurrent Test Running on all devices requested
 * Devices are held only as long as the test runs freeing them up sooner for more testing
 * Active Device monitoring lets you see which devices are in use
* Filter Devices for a particular run
 * Device count filter so you don't piss off your coworkers 
 * Automatic SDK filtering based on manifest parameters
 * Strict filtering mode and Serial filters for cases where you must run a particular subset of devices
* Easy to use Api for extra info

###Gradle Plugin

Coming soon. Plugin is working and nearly complete.

###Whats Missing / What's to come

* Proper encoding for Accept headers (e.g. xml for proper junit output)
* Serious Go know-how (this was a learning opportunity for me)
* Log capture from tests
* Screenshots
* Project Grouping
* Realtime test results
* Ironically: Tests

### Let's see more!

![Imgur](http://i.imgur.com/zEnBWu9.png)

![Imgur](http://i.imgur.com/oZsFNNG.png)
