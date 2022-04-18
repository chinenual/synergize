---
layout: default
title: SYN Sequencer to MIDI Converter
nav_order: 81
description: SYN Sequencer to MIDI Converter
permalink: /docs/syn2midi
---
# SYN Sequencer to MIDI Converter

## User Interface

Click the `Load/Save` button and select the `Convert SYN Sequencer to MIDI` menu item to
initiate the conversion. From the dialog, specify the SYN file to be converted and its [tempo](syn2midi.md#tempo).  

Choose `Virtual Playback mode`.  When checked, the converter will emulate the Synergy sequencer playback and repeat tracks and apply transpose settings in the same way that the original sequencer did.
If unchecked, each track is converted without repeats.  The intent of this mode is to produce raw MIDI data that the user can then slice and dice to reorchestrate in ways unconnected to the original repeat behavior.

If Virtual Playback mode is selected, you must specify `Max Time`.  In case of a sequence that contains looping (repeated) tracks, the converter will stop the "playback" after this aount of time has elapsed (after any running track reaches its end point).  

When Virtual Playback mode is selected, you can also alter the `Track Playback Modes`.  The values default to the settings stored in the SYN file; you can alter them here before initiating the conversion. 

The resulting MIDI file will be in the same folder as the original SYN file, with a `.mid` suffix appened to the SYN file name.

## Capabilities of the Synergy sequencer

The Synergy sequencer records the following events:

* key up / down
* pitchbend
* modulation
* pedal up/down
* transpose
* track switches (for turning other tracks on or off during playback)
* inbound MIDI

Limitations in the fidelity of the recorded events:

* While the Synergy's internal key velocity range is 1..32, the sequencer only records 3-bits of the velocity - so the values are 'compressed' to only seven distinct velocities; those are then converted back to the 1..32 range on playback.

## Other data stored in the SYN file

The SYN file also contains a snapshot of the "state" of the CMOS.  This includes all the voicing related parameters:

* sequencer speed
* portamento type and rate
* vibrato delay, rate and depth
* amplitude center and sensitivity
* timbre center and sensitivity
* multi-voice "program" configuration
* voice assign (rolling, 1st available, etc.)
* channel assign

NOTE values are snapshots of the values when the SYN file was downloaded from the instrument. Except for the events listed in the previous section, they are not tracked as events in the sequencer (so changed during a sequence recording are not reflected during playback - only the "final" value of each parameter  is stored in the SYN file.  Since none of these parameters change based on sequencer events, the translator does not emit any of them to the MIDI file.

## Conversion to MIDI

### Voices

A Synergy sequencer track can record more than just a single voice (if the player pressed a Voice# button while recording, the notes played after that are played on the new voice).  Similarly, when using multi-voice patches (program #1..#4), a keypress can trigger up to 4 notes at a time.

The converter creates a separate output MIDI track for each distinct voice used on the Synergy track.

In the case of multi-voice patches, the user may wish to treat the multi-voice notes as a single thing (it's logically a "layered" patch of a single key on/off event), and combine the separate voice events in the Synergy track to a single event on the MIDI track.  The converter, however does not attempt to recognize such multi-voice events and compress them to a single track.  A multi-voice event creates multiple output tracks.  The user can delete the redundant tracks from the resulting MIDI file in his/her DAW.

### Key Velocity

The SYN file records velocities in 3-bits, compressing the range to seven values.  The converter maps these seven velocities across the 1..127 values available in MIDI.

### Pitchbend

Pitchbend is recorded as -48..48 and mapped to the full range of MIDI pitchbend (-8192..8191).

### Modulation

Modulation is mapped to MIDI CC1 ("mod wheel"), with Synergy values 0..48 mapped to 0..127 MIDI CC values.

### Pedals

The right pedal is mapped to MIDI CC64 ("Sustain pedal") (0 when up, 127 when depressed).

The left pedal is mapped to MIDI CC65 ("Portamento switch").

### Track buttons

When using the Virtual Playback mode, tracks are rendered as the Synergy would play them (repeating tracks when the track button is flashing, turning then on/off based on events recorded in the sequencer).

When not using Virtual Playback mode, the track buttons are ignored.
Tracks are rendered once to the MIDI file with the expectation that the user may copy/paste them in a DAW to effect a repeat and to otherwise orchestrate the interaction of each track.

### Transpose

When using Virtual Playback mode, transpose events are handled in the same way the Synergy sequencer does during playback.

When not using Virtual Playback mode, transpose events are only applied to the track on which they occur.

### Inbound MIDI

Inbound MIDI events are emitted to a separate track since in the original recording, they were intended to control "other" instruments (not the Synergy itself).

### Even Mode

This version of the converter does not emulate the EVEN button.  Contact me if this is something you need and I'll bump its priority.
### Tempo

The Synergy has no strict concept of "tempo" (there is no metronome or "click").  Events are recorded in terms of clock time (milliseconds since the last event).  The Sequencer Speed can change the rate that events are played back, but there is no inherent "beats per minute" recorded in the sequencer.  Any correspondence between the sequencer speed and a musical tempo is up to the user to determine.

MIDI, on the other hand, has a definite notion of musical tempo.  Events written to the MIDI file will be accurate to the original speed recorded in the SYN, but for them to "line up" with MIDI's notion of measures, quarter / eight notes, etc., the user must supply the tempo (in beats per minute). The accuracy of how the translated events line up to measures or beats in a DAW will depend on 1) how accurately the supplied tempo BPM matches the sequenced recording and 2) how close the original recording was to a strict click. 

