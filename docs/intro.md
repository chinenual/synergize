---
layout: default
title: Introduction
nav_order: 1
description: "Overview of the documentation"
permalink: /docs
has_children: true
has_toc: false
---
# Synergize User Manual

## The Synergy

![The Synergy](/synergize/docs/screenshots/synergy.jpg)

The Synergy was first released in 1981.  It features 32 digital
oscillators and a unique architecture that supports very musical and
expressive sounds, using both additive synthesis and phase
modulation. Both frequency and amplitude envelopes can be altered by
keyboard velocity. Even today there's still nothing quite like it.

Under the covers it has the same voicing architecture as the Crumar
GDS - a rare and very expensive early digital synth.

The Synergy was a favorite of Wendy Carlos and featured on her Digital
Moonscapes and Beauty and the Beast albums (she'd previously used a
GDS for the Tron soundtrack). 

However, without computer based software, the Synergy couldn't be
programmed.  It had 24 built in voices. DK produced a library of 9
carts (24 voices each), including 3 programmed by Wendy Carlos
herself.  For many users, these have been the only sounds their
instruments can produce. 

## SYNHCS

In 1984, DK (and later Mulogix) released SYNHCS - a CP/M application
that ran on the Kaypro and talked to the Synergy over a serial
cable. It opened up the full programability of the GDS on the Synergy.
In fact, the SYNHCS documentation refers to the SynergyII+ and SYNHCS
combo as "a better GDS".

## Synergize

That was almost 40 years ago. I don't know how many Synergy's are
still out there in various studios, but I do know that many of them
are running without a Kaypro and some don't even have any of the
original library carts.  Even those that have a Kaypro are dealing
with 30+ year old hardware that's not easy to integrate into a modern
network.

Synergize is a drop-in replacement for SYNHCS/Kaypro, but it runs on
modern hardware and operating systems (Mac, Windows, Linux).  No need
for a 30+ year old Kaypro. Except for a couple of minor features, Synergize
replicates everything you can do in SYNHCS.

With version 2.1.0, Synergize supports an tablet based "control
surface", providing faders akin to the original GDS user interface.
Make your Synergy a bit more like a GDS with Synergize.
See [Control Surface](csurface.md) for details.

## Features

### Cross Platform
* Runs on Mac, Windows, Linux

### GDS-like faders

* Control the editor with an external control surface (move sliders/faders to
  control harmonic frequency, envelope values, filter curves).

### Librarian (VCE and CRT files)
* View voice settings
* Load "carts" and play them. The original DK voice libraries are
  available online. See [Voice Library](voice-library.md) for
  links.
* Create your own "carts" of 24 voices

### Performance Parameter archive (SYN files)
* Persist Performance Parameters (sequencer, portamento, amplitude and timbre settings)

### Voicing Mode (Editor)
* Exposes full control of the Synergy/GDS hardware. As they said with the original SYNHCS, this turns your Synergy into a GDS
* Full control of all voice settings
  * Oscillator routing (phase modulation, additive synthesis and others)
  * Oscillator configuration (wave type, frequency, filter mapping, keyboard proportionality)
  * Default amplitude and timbre settings
  * Envelope control
  * Filter control
  * Keyboard equalization
  * Keyboard proportionality
