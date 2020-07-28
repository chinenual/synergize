---
layout: default
parent: Voicing Mode
title: Envelopes Tab
nav_order: 2
description: "Voicing mode - Envelopes Tab"
permalink: /docs/voicing-mode/envelopes
---

# Envelopes Tab

![screenshot](/synergize/docs/screenshots/viewVCE_envs_annotated.png)

1. Click the `Envelopes` tab to select this Envelopes editor
2. `Oscillator` Selector - use this to select which oscillator's envelope to view or edit.
3. `Envelopes` selector - use this to limit the graphical view of the
   dialog to just one of the four envelopes (in case the overall range of time
   or frequency makes one of the envelopes hard to see).
4. `Copy From` Selector - use this to copy the envelope values from
another oscillator into this one. See [below](#copying-envelopes).
5. The chart shows a graphical view of the filter.  X axis is time in milliseconds, Y axis
   is the amount of frequency (in Hz) or amplitude (in dB) change.
6. The table at the bottom are numeric values of each data point. Adjust these up or
   down to change the behavior of the envelopes.  Changes are transmitted in real time to
   the Synergy, so you can hear the affect of the change immediately.
  * 6a - Loop point - use this to set/clear Loop/Sustain/Repeat points (see [below](#envelope-types)).
  * 6b - Amplitude Value (low and up). The Amplitide value for this point
  in each of the low and up envelopes.
  * 6c - Amplitude Time (low and up).  The time when this point occurs.
  * 6d - Frequency envelope controls.  Similar to the Amplitude controls,
     but controls the modulation of the frequency of the oscilator
     instead of altering the amplitude.
7. For TYPE1 envelopes, two acceleration values are used to control
   the envelope after the key is released (see
   [below](#envelope-types)).  Note: in this case, the Frequency
   envelope is Type1 (no loops), but the Amplitude envelope is Type2,
   so no acceleration values are shown.
8. The `+` and `-` buttons are used to add or remove points from the envelope. 
   
## Envelopes

The GDS is an expressive instrument because any musical parameter
can be programmed to be controlled by the VELOCITY of the keyboard.

Each oscillator is voiced twice. One set of numbers (called the
LOWER BOUNDS) react: either when the CENTER/SENSITIVITY graph
is at the full lower position on the display, or when the keys
are played lightly (with a light touch or slow velocity).

The second set of numbers (called the UPPER BOUNDS) react
when either the CENTER/SENSITIVITY graph is at the full upper
position on the display, or when the keys are played with a
heavy touch (or a faster velocity). When the two sets of
numbers between the LOWER and UPPER BOUNDS are significantly
different, positioning of the CENTER control will allow averaging
or INTERPOLATION between the two sets of values. The SENSITIVITY
control allows the key velocity to alter the interpolation values
and give "expression". Expressive nuance is actually the velocity
control of the values of two sets of numbers.

Each oscillator is voiced twice for amplitude values, specified
in up to 16 points of envelope for the lower bounds and up to 16
points of envelope for the upper bounds.

Each oscillator can be voiced twice for frequency values, specified
in up to 16 points of envelope for the lower bounds and up to 16
points of envelope for the upper bounds.

The VELOCITY controls interpolation characteristics of the envelope
data separately from overall volume expression. TIMBRE center &
sensitivity relate to the numerical data of the envelopes. Amplitude
center & sensitivity control volume response only.

The 16 point envelopes for each of amplitude and frequency,
are independent per oscillator and per voice. However, one
can "copy" one set of envelopes to other oscillators to
speed up the voicing process considerably, then making slight
alterations if desired.

The extensive envelope routines make it possible to duplicate
the amplitude and frequency paths of fundamentals and harmonics,
as does occur in the documented analysis of acoustical sounds.

The number of envelope points are used to aid in duplicating the
multiple changes that take place in natural sounds. This makes
the ADSR concepts of analog synthesizers outdated and impractical
for certain applications, since some sounds have multiple envelope
points in the attack stage alone, and different between one
harmonic to the next.

### Envelope Types

There are several "TYPES of ENVELOPES"

* TYPE 1: A series of attacks and decays. Basically the sound begins
with an attack, and decays while the key is still held down.
If the key is released before the decay is completed, the
envelope finishes at a given rate of speed. If it is
desired to alter how fast or slow the envelope completes
its path, when a key is released, it can be accomplished
by the ACCELLERATION RATE for both amplitude and frequency
and for both bounds.

* TYPE 2: An envelope which has a SUSTAIN (S) point after a
series of attack points. Decay points begin after release
of the key.

* TYPE 3: An envelope with a LOOP (L) point. The envelope begins,
goes to the SUSTAIN point and cycles back to the LOOP
point, and continues between the SUSTAIN and LOOP points
until the key is released. Oscillators with LOOP points
signified with the L are independently looping.

* TYPE 4: The same as TYPE 3, but the "L" is replaced with an "R"
which stands for repeat. In this envelope all oscillators
having the "R" wait for each other so as to be "together".

TYPE 1 is used for pianos, bass, harpsichords, plucked sounds,
percussion sounds etc.

TYPE 2 is used for sustaining sounds, brass, some strings, woodwinds,
sustained synthesizer sounds etc.

TYPE 3 is most useful for chorusing sounds, agitated sounds and active
sound effects.

TYPE 4 is used for rolling percussion, arpeggiation, rotary speakers
of organs and repeating mallet effects.

## Velocity Sensitivity

It should be remembered that the Velocity can control
the following parameters, assigned in the voicing process: Volume,
timbre, attack time, decay time, pitch degree, pitch time, harmonic
entry, speed of loops, modulation degree, speed of repeats, speed of
musical tremelos, sustain lengths, a combination of these and more.
As a guide to voicing, decide what you want the velocity to accomplish
in a certain voice, then try to attain it.
