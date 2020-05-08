include VERSION

DOCS=*.md LICENSE
SRCS=bundler.json [c-z]*.go */*.go resources/app/*html resources/app/static/css/*.css  resources/app/static/js/*js 
EXES=$(EXE_MAC) $(EXE_WINDOWS) $(EXE_LINUX) ${EXE_WINDOWS_TEST}
EXE_MAC=output/darwin-amd64/Synergize.app/Contents/MacOS/Synergize 
EXE_WINDOWS=output/windows-386/Synergize.exe
EXE_WINDOWS_TEST=output/windows-386-cmd/Synergize-cmd.exe
EXE_LINUX=output/linux-amd64/Synergize

# NOTE: must build the exes before we can run the test since some variables 
# used in main.go are generated as side-effects of the astielectron-bundler
all: TAGS $(EXES)

$(EXE_MAC) $(EXE_WINDOWS) $(EXE_LINUX): $(SRCS)
	astilectron-bundler

# command line friendly variant for running batch serial comms tests
$(EXE_WINDOWS_TEST): $(SRCS)
	mkdir -p output/windows-386-cmd
	GOOS=windows GOARCH=386 go build -o $(EXE_WINDOWS_TEST)

mac: $(EXE_MAC)
windows : $(EXE_WINDOWS) $(EXE_WINDOWS_TEST)
linux: $(EXE_LINUX)

package: test packageMac packageWindows packageLinux


# uses create-dmg (installed via "brew install create-dmg"):
packageMac: packages/Synergize-Installer-$(VERSION).dmg
packages/Synergize-Installer-$(VERSION).dmg : $(EXE_MAC) 
	mkdir -p packages
	rm -f packages/Synergize-Installer-$(VERSION).dmg
	create-dmg \
		--volname "Synergize Installer" \
		--volicon resources/icon.icns \
		--icon-size 100 \
		--window-size 450 400 \
		--icon "Synergize.app" 100 120 \
		--app-drop-link 300 120 \
		"packages/Synergize-Installer-$(VERSION).dmg" \
		output/darwin-amd64

# uses msitools (installed via "brew install msitools"):
packageWindows: packages/Synergize-Installer-$(VERSION).msi $(EXE_WINDOWS)
packages/Synergize-Installer-$(VERSION).msi : windows-installer.wxs $(EXE_WINDOWS)
	mkdir -p packages
	rm -f packages/Synergize-Installer-$(VERSION).msi
	wixl -v \
		-a x86 \
		-D VERSION=$(VERSION) \
		-D SourceDir=output/windows-386/ \
		-o packages/Synergize-Installer-$(VERSION).msi \
		windows-installer.wxs

packageLinux: packages/Synergize-linux-amd64-$(VERSION).tar.gz
packages/Synergize-linux-amd64-$(VERSION).tar.gz: $(EXE_LINUX)
	mkdir -p packages
	cd output/linux-amd64 && tar czvf ../../packages/Synergize-linux-amd64-$(VERSION).tar.gz Synergize

test:
	cd data && go test
	go test

version.go : VERSION
	echo package main > version.go
	echo const Version = \"$(VERSION)\" >> version.go

tags TAGS: $(SRCS) $(DOCS)
	etags $(SRCS) $(DOCS)

updateAstilectron:
	go get -u github.com/asticode/go-astilectron
	go get -u github.com/asticode/go-astilectron-bundler
	go get -u github.com/asticode/go-astilectron-bootstrap
	go install github.com/asticode/go-astilectron
	go install github.com/asticode/go-astilectron-bundler
	go install github.com/asticode/go-astilectron-bootstrap

clean:
	rm -rf packages output bind_*.go *.log
