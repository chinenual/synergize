---
layout: default
title: Control Surface
nav_order: 5
description: "Using the Control Surface"
permalink: /docs/control-surface
---
# OSC-based Control Surface

Make your Synergy a little bit more like a GDS. Synergize includes support for an external control surface to provide GDS-like analog controls (sliders and knobs) to the Synergize. You can use both the Synergize UI and the control surface at the same time - making a change in one place is reflected in the other. 
![screenshot](/synergize/docs/screenshots/cs-animated.gif)

## TouchOSC setup

At the moment, the only control surface supported is one written especially for Synergize.  It is based on the excellent TouchOSC tablet based control surface application (which is available for both iPad and Android tablets).

Upload the supplied touchosc configuration file, [Synergize.touchosc](https://github.com/chinenual/synergize/raw/master/osc/touchosc/Synergize.touchosc)
to your tablet following the instructions at [the TouchOSC home site](https://hexler.net/docs/touchosc).

Synergize-TouchOSC uses an OSC connection.  Ensure both your tablet
and the computer running Synergize have good network connections and
are on the same local network. 

## TouchOSC side configuration

In TouchOSC, configure OSC:

* `Enabled` - must be turned on.

* `Host` - the IP address of the computer running Synergize

* `Port (outgoing)` - the network port Synergize is listening to.
  Must match the `OSC Listening Port` configuration in Synergize.

* `Port (incoming)` - the network port Synergize will use to send data
  to TouchOSC.
  Must match the `Control Surface Port` configuration in Synergize.

## Synergize side configuration

In Synergize, enable the control surface, and configure the
connection to the control surface:

* `Enable OSC Control Surface` - unless checked, Synergize will not
use the external control surface.

* `OSC Listening Port` - the network port Synergize will use to
recieve data from the control surface.  Must match the `Port
(outgoing)` setting in TouchOSC.

* `Control Surface Address` - the IP address of the control surface.
TouchOSC reports this as `Local IP address`.

* `Control Surface Port` - the network port Synergize will use to
send data to the control surface.  Must match the `Port
(incoming)` setting in TouchOSC.


## Operation

The Synergize Control Surface is organized in a similar way to the Synergize UI. Each editor view has a separate "tab" on the control surface.  These closely correspond to the tabs in the Synergize UI. In some cases, in order to de-clutter the control surface, a Synergize UI tab is split into more than one tab on the control surface.

Touching the page/tab on the control surface will change its view to
the corresponding editor view. It will also change the Synergize UI to
focus on the related editor tab (and vice versa).

Changing a control on the control surface (e.g. a fader or button)
sends the value to Synergize.  The value will be displayed on the
Synergize UI and sent to the connected Synergy.  The last value sent
to Synergize is displayed on the bottom left corner of the control
surface in yellow text.

### Voice Tab

![screenshot](/synergize/docs/screenshots/cs-voice-tab.png)

The voice table offers controls for everything (except the patch type
and 
routing registers) editable on the UI's Voice Tab.

Each row corresponds to an oscillator.  Inactive oscillator buttons and sliders are hidden.

### Voice Freqs Tab

![screenshot](/synergize/docs/screenshots/cs-voice-freqs-tab.png)

This duplicates the Harmonics and Detune sliders from the voice tab, but in a larger format giving each slide a longer "throw".

### Freq Envelopes Tab

![screenshot](/synergize/docs/screenshots/cs-freq-envs-tab.png)

The Synergize Envelopes tab is split into two on the control surface
due to the sheer number of things that can be controlled.  The
Frequency and Amplitude envelopes have their own pages on the control
surface.  When the Synergize UI changes to the Envelopes tab, the
control surface switches to the Frequency Envelope. You must manually
switch to the Amp envelope tab to control the Amplitude Envelope.

Loop / Sustain and Repeat points cannot be controlled from the control
surface. To change those, use the Synergize UI.

NOTE: as for the oscillator harmonics and detune, unused envelope
points are hidden.

### Amp Envelopes Tab

![screenshot](/synergize/docs/screenshots/cs-amp-envs-tab.png)

Like the Freq tab, but for the Amp envelopes.

### Filters Tab

![screenshot](/synergize/docs/screenshots/cs-filters-tab.png)

You must select the filter you want to edit via the Synergize UI.

### Key Equalization Tab

![screenshot](/synergize/docs/screenshots/cs-keyeq-tab.png)

Allows direct editing of the Key Equalization curve.

### Key Proportion Tab

![screenshot](/synergize/docs/screenshots/cs-keyprop-tab.png)

Allows direct editing of the Key Proportion curve.

## Troubleshooting

A common problem is a poor network connection.  The OSC protocol is
"fire and forget" so if your tablet or the computer running Synergize
are on on a weak Wifi connection or separated by multiple hops,
packets can be lost and updates may not be delivered.
* Ensure both devices are on the same network.
* If using Wifi, make sure you are using a strong connection.
* If you have more than one Wifi router/access point near your studio,
  ensure each device is using the closest/best one. 
