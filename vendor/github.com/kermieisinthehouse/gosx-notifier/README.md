gosx-notifier
===========================
A [Go](http://golang.org) lib for sending desktop notifications to OSX Mountain Lion's (10.8 or higher REQUIRED)
[Notification Center](http://www.macworld.com/article/1165411/mountain_lion_hands_on_with_notification_center.html).

[![GoDoc](http://godoc.org/github.com/deckarep/gosx-notifier?status.png)](http://godoc.org/github.com/deckarep/gosx-notifier)

Update 4/3/2014
------
On OSX 10.9 and above gosx-notifier now supports images and icons.
![Now with custom icon support](../master/example.png?raw=true)

Synopsis
--------
OSX Mountain Lion comes packaged with a built-in notification center. For whatever reason, [Apple sandboxed the
notification center API](http://forums.macrumors.com/showthread.php?t=1403807) to apps hosted in its App Store. The end
result? A potentially useful API shackled to Apple's ecosystem.

Thankfully, [Eloy Durán](https://github.com/alloy) put together [an osx app](https://github.com/alloy/terminal-notifier) that allows terminal access to the sandboxed API. **gosx-notifier** embeds this app with a simple interface to the closed API.

It's not perfect, and the implementor will quickly notice its limitations. However, it's a start and any pull requests are accepted and encouraged!

Dependencies:
-------------
There are none! If you utilize this package and create a binary executable it will auto-magically install the terminal-notifier component into a temp directory of the server.  This is possible because in this latest version the terminal-notifier binary is now statically embedded into the Go source files.


Installation and Requirements
-----------------------------
The following command will install the notification api for Go along with the binaries.  Also, utilizing this lib requires OSX 10.8 or higher. It will simply not work on lower versions of OSX.

```sh
go get github.com/deckarep/gosx-notifier
```

Using the Command Line
-------------
```Go
notify "Wow! A notification!!!"
```

useful for knowing when long running commands finish

```Go
longRunningCommand && notify done!
```

Using the Code
------------------
It's a pretty straightforward API:

```Go
package main

import (
    "github.com/deckarep/gosx-notifier"
    "log"
)

func main() {
    //At a minimum specifiy a message to display to end-user.
    note := gosxnotifier.NewNotification("Check your Apple Stock!")

    //Optionally, set a title
    note.Title = "It's money making time 💰"

    //Optionally, set a subtitle
    note.Subtitle = "My subtitle"

    //Optionally, set a sound from a predefined set.
    note.Sound = gosxnotifier.Basso

    //Optionally, set a group which ensures only one notification is ever shown replacing previous notification of same group id.
    note.Group = "com.unique.yourapp.identifier"

    //Optionally, set a sender (Notification will now use the Safari icon)
    note.Sender = "com.apple.Safari"

    //Optionally, specifiy a url or bundleid to open should the notification be
    //clicked.
    note.Link = "http://www.yahoo.com" //or BundleID like: com.apple.Terminal

    //Optionally, an app icon (10.9+ ONLY)
    note.AppIcon = "gopher.png"

    //Optionally, a content image (10.9+ ONLY)
    note.ContentImage = "gopher.png"

    //Then, push the notification
    err := note.Push()

    //If necessary, check error
    if err != nil {
        log.Println("Uh oh!")
    }
}
```

Sample App: Desktop Pinger Notification - monitors your websites and will notifiy you when a website is down.
```Go
package main

import (
	"github.com/deckarep/gosx-notifier"
	"net/http"
	"strings"
	"time"
)

//a slice of string sites that you are interested in watching
var sites []string = []string{
	"http://www.yahoo.com",
	"http://www.google.com",
	"http://www.bing.com"}

func main() {
	ch := make(chan string)

	for _, s := range sites {
		go pinger(ch, s)
	}

	for {
		select {
		case result := <-ch:
			if strings.HasPrefix(result, "-") {
				s := strings.Trim(result, "-")
				showNotification("Urgent, can't ping website: " + s)
			}
		}
	}
}

func showNotification(message string) {

	note := gosxnotifier.NewNotification(message)
	note.Title = "Site Down"
	note.Sound = gosxnotifier.Default

	note.Push()
}

//Prefixing a site with a + means it's up, while - means it's down
func pinger(ch chan string, site string) {
	for {
		res, err := http.Get(site)

		if err != nil {
			ch <- "-" + site
		} else {
			if res.StatusCode != 200 {
				ch <- "-" + site
			} else {
				ch <- "+" + site
			}
			res.Body.Close()
		}
		time.Sleep(30 * time.Second)
	}
}
```

Usage Ideas
-----------
* Monitor your awesome server cluster and push notifications when something goes haywire (we've all been there)
* Scrape Hacker News looking for articles of certain keywords and push a notification
* Monitor your stock performance, push a notification, before you lose all your money
* Hook it up to ifttt.com and push a notification when your motion-sensor at home goes off

Coming Soon
-----------
* Remove ID

Licence
-------
This project is dual licensed under [any licensing defined by the underlying apps](https://github.com/alloy/terminal-notifier) and MIT licensed for this version written in Go.


[![Bitdeli Badge](https://d2weczhvl823v0.cloudfront.net/deckarep/gosx-notifier/trend.png)](https://bitdeli.com/free "Bitdeli Badge")
