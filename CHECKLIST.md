# Release Checklist

## Windows installer GUIDs

See https://www.firegiant.com/wix/tutorial/upgrades-and-modularization/
and https://www.firegiant.com/wix/tutorial/upgrades-and-modularization/checking-for-oldies/

If contents of the install has changed (more or less files), the product ID GUID in the wxs needs to change.

When version is bumped (minor release), the Package GUID has to change

When version is bumped (major release), both Product ID and Package GUID has to change

