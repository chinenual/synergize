rm -rf synergize.app/Contents/MacOS
mkdir synergize.app/Contents/MacOS
cp -R content/static synergize.app/Contents/MacOS
cp -R content/template synergize.app/Contents/MacOS
go build -o synergize.app/Contents/MacOS/synergize

