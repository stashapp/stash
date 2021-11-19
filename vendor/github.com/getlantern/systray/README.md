systray is a cross-platform Go library to place an icon and menu in the notification area.

## Features

* Supported on Windows, macOS, and Linux
* Menu items can be checked and/or disabled
* Methods may be called from any Goroutine

## API

```go
func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("Awesome App")
	systray.SetTooltip("Pretty awesome超级棒")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	// Sets the icon of a menu item. Only available on Mac and Windows.
	mQuit.SetIcon(icon.Data)
}

func onExit() {
	// clean up here
}
```

See [full API](https://pkg.go.dev/github.com/getlantern/systray?tab=doc) as well as [CHANGELOG](https://github.com/getlantern/systray/tree/master/CHANGELOG.md).

## Try the example app!

Have go v1.12+ or higher installed? Here's an example to get started on macOS:

```sh
git clone https://github.com/getlantern/systray
cd example
env GO111MODULE=on go build
./example
```

On Windows, you should build like this:

```
env GO111MODULE=on go build -ldflags "-H=windowsgui"
```

The following text will then appear on the console:


```sh
go: finding github.com/skratchdot/open-golang latest
go: finding github.com/getlantern/systray latest
go: finding github.com/getlantern/golog latest
```

Now look for *Awesome App* in your menu bar!

![Awesome App screenshot](example/screenshot.png)

## The Webview example

The code under `webview_example` is to demostrate how it can co-exist with other UI elements. Note that the example doesn't work on macOS versions older than 10.15 Catalina.

## Platform notes

### Linux

* Building apps requires gcc as well as the `gtk3` and `libappindicator3` development headers to be installed. For Debian or Ubuntu, you may install these using:

```sh
sudo apt-get install gcc libgtk-3-dev libappindicator3-dev
```

On Linux Mint, `libxapp-dev` is also required .

To build `webview_example`, you also need to install `libwebkit2gtk-4.0-dev` and remove `webview_example/rsrc.syso` which is required on Windows.

### Windows

* To avoid opening a console at application startup, use these compile flags:

```sh
go build -ldflags -H=windowsgui
```

### macOS

On macOS, you will need to create an application bundle to wrap the binary; simply folders with the following minimal structure and assets:

```
SystrayApp.app/
  Contents/
    Info.plist
    MacOS/
      go-executable
    Resources/
      SystrayApp.icns
```

When running as an app bundle, you may want to add one or both of the following to your Info.plist:

```xml
<!-- avoid having a blurry icon and text -->
	<key>NSHighResolutionCapable</key>
	<string>True</string>

	<!-- avoid showing the app on the Dock -->
	<key>LSUIElement</key>
	<string>1</string>
```

Consult the [Official Apple Documentation here](https://developer.apple.com/library/archive/documentation/CoreFoundation/Conceptual/CFBundles/BundleTypes/BundleTypes.html#//apple_ref/doc/uid/10000123i-CH101-SW1).

## Credits

- https://github.com/xilp/systray
- https://github.com/cratonica/trayhost
