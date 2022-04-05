
## Capabilities of the Synergy sequencer

Records the following events:

* key up / down
* pitchbend
* modulation
* pedal up/down
* transpose
* track switches (for turning other tracks on or off during playback)
* inbound MIDI

Limitations in the fidelity of the recorded events:

* While the Synergy's internal key velocity range is 1..32, the sequencer only record 3-bits of the velocity - so the values are 'compressed' to only seven distinct velocities; those are then converted back to the 1..32 range on playback.

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

NOTE values are snapshots of the values when the SYN file was downloaded from the instrument. Except for the events listed above, they are not tracked as events in the sequencer (so changed during a sequence recording are not reflected during playback - only the "final" value of each parameter  is stored in the SYN file.  Since none of these parameters change based on sequencer events, the translator does not emit any of them to the MIDI file.

## Conversion to MIDI

### Key Velocity

The SYN file records velocities in 3-bits, compressing the range to seven values.  The converter maps these seven velocities across the 1..127 values available in MIDI.

### Pitchbend

Similarly, pitchbend is recorded as -127..127 - it is mapped to .
Depending on the user's preference, these are preserved as-is in the MIDI, or mapped to the full range of MIDI pitchbend (-8192..8191).

### Modulation

Modulation is mapped to MIDI CC1 ("mod wheel"), with Synergy values -127..127 mapped to 0..127 MIDI CC values.

### Pedals

The right pedal is mapped to MIDI CC64 ("Sustain pedal") (0 when up, 127 when depressed).

The left pedal is mapped to MIDI CC65 ("Portamento switch").

### Track buttons

A future release may support a virtual playback option that renders tracks as the Synergy would (repeating tracks when the track button is flashing, turning then on/off based on events recorded in the sequencer.

In this release, track button selection is ignored.   Tracks are rendered once to the MIDI file with the expectation that the user may copy/paste them in a DAW to effect a repeat and to otherwise orchestrate the interaction of each track.

### Transpose

In this release, transpose events are ignored.  See the discussion above regarding Track buttons. If a virtual playback mode is added to the converted, transpose will be supported.

###

### Tempo

NOTE:  This may change: I am investigating SMPTE mode MIDI format - that may allow millisecond based timing directly in the MIDI file without requiring the user to select a tempo BPM.

The Synergy has no strict concept of "tempo" (there is no metronome or "click").  Events are recorded in terms of clock time (milliseconds since the last event).  The Sequencer Speed can change the rate that events are played back, but there is no inherent "beats per minute" recorded in the sequencer.  Any correspondence between the sequencer speed and a musical tempo is up to the user to make.

MIDI, on the other hand, has a definite notion of musical tempo.  Time in a MIDI stream is measured in pulses-per-quarter-note (PPQ).  In order for the converter to convert a Synergy event timestamp to MIDI PPQ, the user must supply a tempo (in beats per minute).  The accuracy of how the translated events line up to measures or beats in a DAW will depend on 1) how accurately the supplied tempo BPM matches the sequenced recording and 2) how close the original recording was to a strict click.
