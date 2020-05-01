include VERSION

SRCS=*.md LICENSE [c-z]*.go resources/app/*html resources/app/static/css/*.css  resources/app/static/js/*js

all: TAGS output 

output: $(SRCS)
	astilectron-bundler

package: packageMac packageWindows

# uses create-dmg (installed via "brew install create-dmg"):
packageMac: release/Synergize-Installer-$(VERSION).dmg
release/Synergize-Installer-$(VERSION).dmg : output
	mkdir release
	create-dmg \
		--volname "Synergize Installer" \
		--volicon resources/icon.icns \
		--icon-size 100 \
		--window-size 450 400 \
		--icon "Synergize.app" 100 120 \
		--app-drop-link 300 120 \
		"release/Synergize-Installer-$(VERSION).dmg" \
		output/darwin-amd64

packageWindows: output

version.go : VERSION
	echo package main > version.go
	echo const Version = \"$(VERSION)\" >> version.go

TAGS: $(SRCS)
	etags $(SRCS)

clean:
	rm -rf release output bind_*.go *.log
