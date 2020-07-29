---
layout: default
title: Getting Started
nav_order: 1
description: "Getting Started"
parent: Introduction
permalink: /docs/getting-started
---

# Getting Started

## Requirements

### Synergy II+

Synergize connects to a Synergy II+ via its RS232 serial port.  It is tested with a Synergy running the latest 3.22 firmware but is probably compatible with earlier versions.

You'll also need a null modem serial cable and possibly a USB serial device depending on your computer's capablities. See [HARDWARE](hardware.md) for some details.

### Operating Systems

Synergize has been tested on:

* MacOS 10.14.6 (Mojave)
* Windows 10
* Ubuntu Linux 18.04
* Ubuntu Linux 16.04


## Install Synergize

### Mac

The Mac release is packaged as an installable DMG file. Open the DMG and drag Synergize the Applications folder.

### Windows

The Windows release is as a standard Windows .MSI installer. Run the installer.

### Linux

The Linux release is packaged as a simple tarball.  Untar it to a location of your choosing.

## Download the Voice Library

If you don't already have a copy of the factory CRT and VCE files, get a copy as described [here](voice-library.md).

## Settings

Use the `Help -> Preferences` menu to configure Synergize to your local setup.

### Serial Port

Set your Serial Port device here. 

### Baud Rate

Set the baud rate of your Synergy here. See [Serial Port Configuration](hardware.md) for details.

### Voice Library Location

This tells Synergize where to find your voice library (Synergize defaults your Home directory).  This will be the default location displayed in the left-pane file browser when you start Synergize.

## Test your connection to the Synergy

You can test the connection by selecting the `Connect->Connect to Synergy` menu.  If successful, Synergize will report the firmware version of the connected Synergy in the upper left pane of the display.

It is not necessary to explicitly connect in this way. Synergize will connect the first time you invoke  a command that needs to communicate with the Synergy.

See [Troubleshooting](hardware-troubleshooting.md) for things to check if this does not "just work".
