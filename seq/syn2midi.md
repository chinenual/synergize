
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

The converter can handle this in multiple ways, depending on the preference of the user:

* ignore the voice annotation on the event and simply emit note on/off events to one MIDI track (i.e., one Synergy track is exactly one MIDI track).
* distribute voices to separate MIDI tracks.  In this case, for example, if the synergy track contains notes for 2 different voices, the MIDI file will contain two tracks - one for the first voice, one for the second.
* In the case of multi-voice patches, the user may wish to treat the multi-voice notes as a single thing (it's logically a "layered" patch of a single sound), so the converter will combine the separate voice events in the Synergy track to a single event on the MIDI track.

The second and third options may be combined such that the converter creates separate MIDI tracks for each voice, but in the case of a "multi-voice" event, combines that event into a single event, but on its own MIDI track.

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

A future release may support a virtual playback option that renders tracks as the Synergy would (repeating tracks when the track button is flashing, turning then on/off based on events recorded in the sequencer).

In this release, track button selection is ignored.   Tracks are rendered once to the MIDI file with the expectation that the user may copy/paste them in a DAW to effect a repeat and to otherwise orchestrate the interaction of each track.

### Transpose

In this release, transpose events are ignored.  See the discussion above regarding Track buttons. If a virtual playback mode is added to the converter, transpose will be supported.

### Inbound MIDI

Inbound MIDI events are emitted to a separate track since in the original recording, they were intended to control "other" instruments (not the Synergy itself).

### Tempo

**NOTE**:  The following may change: I am investigating SMPTE mode MIDI format - that may allow millisecond based timing directly in the MIDI file without requiring the user to select a tempo BPM.

The Synergy has no strict concept of "tempo" (there is no metronome or "click").  Events are recorded in terms of clock time (milliseconds since the last event).  The Sequencer Speed can change the rate that events are played back, but there is no inherent "beats per minute" recorded in the sequencer.  Any correspondence between the sequencer speed and a musical tempo is up to the user to make.

MIDI, on the other hand, has a definite notion of musical tempo.  Time in a MIDI stream is measured in pulses-per-quarter-note (PPQ).  In order for the converter to convert a Synergy event timestamp to MIDI PPQ, the user must supply a tempo (in beats per minute).  The accuracy of how the translated events line up to measures or beats in a DAW will depend on 1) how accurately the supplied tempo BPM matches the sequenced recording and 2) how close the original recording was to a strict click.
