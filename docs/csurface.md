---
layout: default
title: MIDI Control Surface
nav_order: 5
description: "Using a MIDI Control Surface"
permalink: /docs/control-surface
---
# MIDI Control Surface

Make your Synergy a little bit more like a GDS. Synergize includes support for an external MIDI control surface to provide GDS-like analog controls (sliders and knobs) to the Synergize. You can use both the Synergize UI and the MIDI control surface at the same time - making a change in one place is reflected in the other. 


## TouchOSC setup

At the moment, the only control surface supported is one written especially for Synergize.  It is based on the excellent TouchOSC tablet based control surface application (which is available for both iPad and Android tablets).

Upload the supplied touchosc configuration file, [Synergize.touchosc](https://github.com/chinenual/synergize/raw/master/midi/touchosc/Synergize.touchosc)
to your tablet following the instructions at [the TouchOSC home site](https://hexler.net/docs/touchosc).

Synergize-TouchOSC is hardwired to run on MIDI channel 16 (a limitation in the TouchOSC-side configuration).

I run this via a networked MIDI ("Network Session 1") connection - I've not needed to use the TouchOSC bridge software. Your milage may vary.

## Synergize side configuration

In Synergize, enable the control surface, by enabling the checkbox in the Help->Preferences page.  You must select the MIDI interface on which to communicate with the control surface.


## Operation

The Synergize Control Surface is organized in a similar way to the Synergize UI. Each editor view has a separate "tab" on the control surface.  These closely correspond to the tabs in the Synergize UI. In some cases, in order to de-clutter the control surface, a Synergize UI tab is split into more than one tab on the control surface.

Touching the page/tab on the control surface will change its view to the corresponding editor view. It will also change the Synergize UI to focus on the related editor tab (and vice versa).

### Voice Tab

The voice table offers controls for everything editable on the UI's Voice Tab.

Each row corresponds to an oscillator.  The yellow lights indicate which oscillators are active for the voice. Inactive oscillator buttons and sliders are not hidden, but their values are ignored.

To add or delete oscillators, use the Synergize UI `+` and `-` buttons.

### Voice Freqs Tab

This duplicates the Harmonics and Detune sliders from the voice tab, but in a larger format giving each slide a longer "throw".

### Freq Envelopes Tab

The Synergize Envelopes tab is split into two on the control surface due to the sheer number of things that can be controlled.  The Frequency and Amplitude envelopes have their own pages on the control surface.  When the Synergize UI changes to the Envelopes tab, the control surface switches to the Frequency Envelope. You must manually switch to the Amp envelope if needed.

NOTE: as for the oscillator harminics and detune, unused envelope points still have visible sliders on the control surface - their values are ignored.

### Amp Envelopes Tab

Like the Freq tab, but for the Amp envelopes.

### Filters Tab

You must select the filter you want to edit via the Synergize UI.

### Key Equalization Tab

Allows direct editing of the Key Equalization curve.

### Key Proportion Tab

Allows direct editing of the Key Proportion curve.

