include VERSION

DOCS=*.md LICENSE
SRCS=bundler.json [c-z]*.go */*.go resources/app/*html resources/app/static/css/*.css  resources/app/static/js/*js 
EXES=$(EXE_MAC) $(EXE_WINDOWS) $(EXE_LINUX_AMD64) $(EXE_LINUX_386) $(EXE_LINUX_ARM) ${EXE_WINDOWS_TEST}
EXE_MAC=output/darwin-amd64/Synergize.app/Contents/MacOS/Synergize 
EXE_WINDOWS=output/windows-386/Synergize.exe
EXE_WINDOWS_TEST=output/windows-386-cmd/Synergize-cmd.exe
EXE_LINUX_AMD64=output/linux-amd64/Synergize
EXE_LINUX_386=output/linux-386/Synergize
EXE_LINUX_ARM=output/linux-arm/Synergize

# NOTE: must build the exes before we can run the test since some variables 
# used in main.go are generated as side-effects of the astielectron-bundler
.PHONY: all
all: TAGS $(EXES)


$(EXE_MAC) $(EXE_WINDOWS) $(EXE_LINUX_AMD64) $(EXE_LINUX_386) $(EXE_LINUX_ARM) : $(SRCS)
	rm windows.syso # delete temporary file side effect of windows build - linux-arm chokes on it.
	astilectron-bundler

# command line friendly variant for running batch serial comms tests
$(EXE_WINDOWS_TEST): $(SRCS)
	mkdir -p output/windows-386-cmd
	GOOS=windows GOARCH=386 go build -o $(EXE_WINDOWS_TEST)

.PHONY: mac
mac: $(EXE_MAC)
.PHONY: windows
windows : $(EXE_WINDOWS) $(EXE_WINDOWS_TEST)
.PHONY: linux
linux: $(EXE_LINUX_AMD64) $(EXE_LINUX_386) $(EXE_LINUX_ARM)

.PHONY: package
package: test packageMac packageWindows packageLinux


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
packageWindows: packages/Synergize-Installer-windows-$(VERSION).msi $(EXE_WINDOWS)
packages/Synergize-Installer-windows-$(VERSION).msi : windows-installer.wxs $(EXE_WINDOWS)
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
	go test

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
PORT=/dev/tty.usbserial-AL05OC8S
endif
ifeq ($(UNAME_S),Linux)
PORT=/dev/ttyS0
endif
ifeq ($(PORT),'')
PORT=COM1
endif

.PHONY: itest
itest:
	cd synio && go test -v -synio -port $(PORT)

version.go : VERSION
	echo package main > version.go
	echo const Version = \"$(VERSION)\" >> version.go

.PHONY: tags
tags: $(SRCS) $(DOCS)
	etags $(SRCS) $(DOCS)

.PHONY: updateAstilectron
updateAstilectron:
	go get -u github.com/asticode/go-astilectron
	go get -u github.com/asticode/go-astilectron-bundler
	go get -u github.com/asticode/go-astilectron-bootstrap
	go install github.com/asticode/go-astilectron
	go install github.com/asticode/go-astilectron-bundler
	go install github.com/asticode/go-astilectron-bootstrap

.PHONY: clean
clean:
	rm -rf packages output bind_*.go *.log
