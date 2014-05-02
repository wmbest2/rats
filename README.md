RATS
===========
Remote Android Test Server
-----

![Imgur](http://i.imgur.com/7skL9wq.jpg)

What and Why?
----
RATS is a easy to run Test Runner Service that can be hosted anywhere a device can be connected.  It was designed to thwart device hogs and provide a level playing field for shops who can't afford more expensive testing services.

### Who?
So far mostly me, but not without support from @vokalinteractive


###Installation


Requirements:

`go, adb (at $ANDROID_HOME), bzr, git`

1. `go get github.com/wmbest2/rats_server`
2. `rats_server`


###Capabilities

* Concurrent Test Running on all devices connect
* Filter Devices for a particular run
 * Strict mode for cases where you must run a particular subset
* Easy to use Api for extra info
* Easy to curl test execution
  * `curl -X POST myserver.local:3000/api/run -F apk=@myapk.apk -F test-apk=@myapk-test.apk`
  * Seriously thats it

###Whats Missing

* Serious Go know-how (this was a learning opportunity for me)
* Log capture from tests
* Screenshots
* Better Filtering/Scripting
* Project Grouping

### Let's see more!

![Imgur](http://i.imgur.com/MrQX3Gz.jpg)

![Imgur](http://i.imgur.com/73hH7Qd.jpg)

