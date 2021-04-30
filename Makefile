include VERSION

DOCS=*.md LICENSE
SRCS=bundler.json cmd/*/*.go [c-z]*.go */*.go resources/app/*html resources/app/static/css/*.css  resources/app/static/js/*js 
EXES=$(EXE_MAC) $(EXE_WINDOWS) $(EXE_LINUX_AMD64) $(EXE_LINUX_386) $(EXE_LINUX_ARM) ${EXE_WINDOWS_TEST}
EXE_MAC=output/darwin-amd64/Synergize.app/Contents/MacOS/Synergize 
EXE_WINDOWS=output/windows-386/Synergize.exe
EXE_WINDOWS_TEST=output/windows-386/Synergize-cmd.exe
EXE_LINUX_AMD64=output/linux-amd64/Synergize
EXE_LINUX_386=output/linux-386/Synergize
EXE_LINUX_ARM=output/linux-arm/Synergize

DX2SYNS=$(DX2SYN_MAC) $(DX2SYN_WINDOWS) $(DX2SYN_LINUX_AMD64) $(DX2SYN_LINUX_386) $(DX2SYN_LINUX_ARM) 
DX2SYN_MAC=output/darwin-amd64/Synergize.app/Contents/MacOS/dx2syn 
DX2SYN_WINDOWS=output/windows-386/dx2syn.exe
DX2SYN_LINUX_AMD64=output/linux-amd64/dx2syn
DX2SYN_LINUX_386=output/linux-386/dx2syn
DX2SYN_LINUX_ARM=output/linux-arm/dx2syn


# NOTE: must build the exes before we can run the test since some variables 
# used in main.go are generated as side-effects of the astielectron-bundler
.PHONY: all
all: $(DX2SYNS) $(EXES)

$(DX2SYN_MAC): $(SRCS)
	mkdir -p output/darwin-amd64/Synergize.app/Contents/MacOS
	cd cmd/dx2syn && go build -o ../../output/darwin-amd64/Synergize.app/Contents/MacOS

$(DX2SYN_WINDOWS): $(SRCS)
	mkdir -p output/windows-386
	cd cmd/dx2syn && GOOS=windows GOARCH=386 go build -o ../../output/windows-386

$(DX2SYN_LINUX_AMD64):
	mkdir -p output/linux-amd64
	cd cmd/dx2syn && GOOS=linux GOARCH=amd64 go build -o ../../output/windows-386

$(DX2SYN_LINUX_386):
	mkdir -p output/linux-386
	cd cmd/dx2syn && GOOS=linux GOARCH=386 go build -o ../../output/windows-386

$(DX2SYN_LINUX_ARM):
	mkdir -p output/linux-arm
	cd cmd/dx2syn && GOOS=linux GOARCH=arm go build -o ../../output/windows-386


$(EXE_MAC) $(EXE_WINDOWS) $(EXE_LINUX_AMD64) $(EXE_LINUX_386) $(EXE_LINUX_ARM) : $(SRCS)
	rm -f windows.syso # delete temporary file side effect of windows build - linux-arm chokes on it.
	astilectron-bundler

# command line friendly variant for running batch serial comms tests
$(EXE_WINDOWS_TEST): $(SRCS)
	mkdir -p output/windows-386
	GOOS=windows GOARCH=386 go build -o $(EXE_WINDOWS_TEST)

.PHONY: mac
mac: version.go
	rm -rf output/windows* output/linux* 
	astilectron-bundler -c bundler-mac-only.json

.PHONY: cibuild
cibuild: version.go
	astilectron-bundler -c bundler-ci.json

.PHONY: package
package: test packageMac packageWindows packageLinux
	cp osc/touchosc/Synergize-v*.touchosc packages


# uses create-dmg (installed via "brew install create-dmg"):
.PHONY: packageMac
packageMac: packages/Synergize-Installer-mac-$(VERSION).dmg
packages/Synergize-Installer-mac-$(VERSION).dmg : $(EXE_MAC) 
	mkdir -p packages
	rm -f packages/Synergize-Installer-mac-$(VERSION).dmg
	create-dmg \
		--volname "Synergize Installer" \
		--volicon resources/icon.icns \
		--icon-size 100 \
		--window-size 450 400 \
		--icon "Synergize.app" 100 120 \
		--app-drop-link 300 120 \
		"packages/Synergize-Installer-mac-$(VERSION).dmg" \
		output/darwin-amd64

# uses msitools (installed via "brew install msitools"):
.PHONY: packageWindows
packageWindows: packages/Synergize-Installer-windows-$(VERSION).msi $(EXE_WINDOWS) $(EXE_WINDOWS_TEST)
packages/Synergize-Installer-windows-$(VERSION).msi : windows-installer.wxs $(EXE_WINDOWS) $(EXE_WINDOWS_TEST)
	mkdir -p packages
	rm -f packages/Synergize-Installer-windows-$(VERSION).msi
	wixl -v \
		-a x86 \
		-D VERSION=$(VERSION) \
		-D SourceDir=output/windows-386/ \
		-o packages/Synergize-Installer-windows-$(VERSION).msi \
		windows-installer.wxs

.PHONY: packageLinux
packageLinux: \
  packages/Synergize-linux-amd64-$(VERSION).tar.gz \
  packages/Synergize-linux-386-$(VERSION).tar.gz \
  packages/Synergize-linux-arm-$(VERSION).tar.gz

packages/Synergize-linux-amd64-$(VERSION).tar.gz: $(EXE_LINUX_AMD64)
	mkdir -p packages
	cd output/linux-amd64 && tar czvf ../../packages/Synergize-linux-amd64-$(VERSION).tar.gz Synergize

packages/Synergize-linux-386-$(VERSION).tar.gz: $(EXE_LINUX_386)
	mkdir -p packages
	cd output/linux-386 && tar czvf ../../packages/Synergize-linux-386-$(VERSION).tar.gz Synergize

packages/Synergize-linux-arm-$(VERSION).tar.gz: $(EXE_LINUX_ARM)
	mkdir -p packages
	cd output/linux-arm && tar czvf ../../packages/Synergize-linux-arm-$(VERSION).tar.gz Synergize

.PHONY: test
test:
	cd data && go test
	cd osc && go test
	cd synio && go test
	cd dx2syn && go test

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
SCREENSHOT_ARCH=darwin
PORT=/dev/tty.usbserial-AL05OC8S
endif
ifeq ($(UNAME_S),Linux)
SCREENSHOT_ARCH=linux
PORT=/dev/ttyS0
endif
ifeq ($(PORT),'')
SCREENSHOT_ARCH=win32
PORT=COM1
endif

.PHONY: itest
itest:
	cd osc && go test -v -oscio
	cd synio && go test -v -synio -port $(PORT)

.PHONY: uitest
uitest:
	/bin/rm -f uitest/test/screenshots/$(SCREENSHOT_ARCH)/DEBUG*.png
	/bin/rm -f uitest/test/screenshots/$(SCREENSHOT_ARCH)/*failed.png
	-cd uitest && npm test | tee uitest-$(SCREENSHOT_ARCH).log
	ls uitest/test/screenshots/$(SCREENSHOT_ARCH)/*failed.png

.PHONY: uitest-diff
uitest-diff:
	-open uitest/test/screenshots/$(SCREENSHOT_ARCH)/AFTERHOOK*.failed.png
	for f in `/bin/ls -1 uitest/test/screenshots/$(SCREENSHOT_ARCH)/*.failed.png | grep -v AFTERHOOK`; do \
		echo f: $$f; \
		s="`basename "$$f" .failed.png`"; \
		echo s: $$s; \
		compare "uitest/test/screenshots/$(SCREENSHOT_ARCH)/$${s}.png" "$$f" /tmp/diff.png; \
		open /tmp/diff.png ;\
		echo press RETURN for next image; \
		read;\
	done

version.go : VERSION
	echo package main > version.go
	echo >> version.go
	echo const Version = \"$(VERSION)\" >> version.go

.PHONY: tags
tags: $(SRCS) $(DOCS)
	etags $(SRCS) $(DOCS)

.PHONY: installDependencies
installDependencies:
	go get -u github.com/asticode/go-astilectron
	go get -u github.com/asticode/go-astilectron-bundler
	go get -u github.com/asticode/go-astilectron-bootstrap
	go get -v -t -d ./...
	go install github.com/asticode/go-astilectron
	go install github.com/asticode/go-astilectron-bundler/astilectron-bundler
	go install github.com/asticode/go-astilectron-bootstrap

.PHONY: clean
clean:
	rm -rf packages output bind_*.go *.log
