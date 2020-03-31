rm -rf synergize.app/Contents/MacOS
mkdir synergize.app/Contents/MacOS
cp -R public synergize.app/Contents/MacOS
go build -o synergize.app/Contents/MacOS/synergize

