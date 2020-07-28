---
layout: default
title: Voicing Theory
nav_order: 10
description: Voicing Theory
permalink: /docs/voicing-theory
---

# Synergy Voicing Theory


<p class="callout">
<b>NOTE:</b> This chapter is adapted almost verbatim from Tom Piggott's original
SYNHCS documentation (with permission). 
<br/><br/>
The section describing the descriptions of the Patch registers
(the Output, Addr and Freq DSR's) is adapted from the firmware
documentation.
<br/><br/>
The graphical display of the patch routing is new to Synergize.
</p>

Here are a series of comments, descriptions and hints for comprehension
of the GDS voicing system. These can be reviewed as often as necessary,
to strengthen your understanding of how the parameters you control on
the displays all interact to produce good sounds. Various books, that
have been published on acoustics, physics, sound synthesis and the like,
can all be used as reference towards building a solid understanding
of the components of sound. These will all contribute to your ability
to attain sounds you have wanted to synthesize for yourself.

## Oscillators

There are 32 oscillators in this system. They consist of either
SINE or modified TRIANGLE waveforms. If you use 2 oscillators
per key to make a voice, you will have 16 playable keys. If you
use 4 oscillators per voice, you will have 8 keys playable. If
you use 3 oscillators per voice, you will have 10 keys playable
with 2 left over.

The more complex the voice is, the more oscillators will be required.

Each oscillator can be used in three ways: to be heard, to be
summed with another one, or as a modulator to another. This
is shown in the patch display, which indicates how the oscillator
is being used in a particular voice patch.

## Patch Type

The "patch type" is selectable on the `Voice Tab`.
It contains various ways in which patching of oscillators is
accomplished. Synergize displays the factory patch types with both the
textual format used in SYNHCS and a graphical format similar to a
Yamaha's DX7 "algorithm" graphic.

### Display format

The display shows each oscillator in the voice.  A line between them
means the oscillator(s) above it frequency modulate the one below.
Oscillators on the very bottom row are "heard"; oscillators on rows
above are modulators for other oscillators. When more than one
oscillator modulates a

For example the standard patch 4:  Oscillator 1 modulates
Oscillator 2.  Oscillator 3 is summed with that and the sum modulates
oscillator 4.  Similarly for oscillators 5 through 8.  The output of
oscillators 4 and 8 are "heard":<br><img style="zoom:50%;" src="https://github.com/chinenual/synergize/raw/master/docs/screenshots/patch-type-4.png?raw=true" />

SYNHCHS used a textual representation of this modulation routing.  A
`~` indicates "modulation", a `+` indicates "additive" (summation):

```
(1~2+3)~4) + ((5~6+7)~8)
```

The actual routing is controlled by a set of three values on each
oscillator (the Freq, Adder and Output registers) -- described below.

### Standard Patch Routing


* PATCH #1, `1 + 2 + 3 + 4 + 5 + 6 + 7 + 8`,
represents additive synthesis. Each oscillator has values
which are "heard" as they occur. It is the most accurate form
of synthesis, but requires a larger number of oscillators to
attain results.<br><img style="zoom:50%" src="https://github.com/chinenual/synergize/raw/master/docs/screenshots/patch-type-1.png?raw=true" />

* PATCH #2, `(1~2) + (3~4) + (5~6) + (7~8)`, is used in sets of 2 oscillators at a time. Number 1 uses
its parameters to "modulate" number 2 which is listened to.
Number 3 modulates 4 which is listened to etc. This modulation
sets up a predictable array of sonic qualities , which helps to
provide rich sonorities , aiding in the conservation of oscillators.
This modulation is actually called PHASE MODULATION and CANCELLATION ,
not to be confused with Frequency Modulation as used in some
Synthesis techniques.<br><img style="zoom:50%" src="https://github.com/chinenual/synergize/raw/master/docs/screenshots/patch-type-2.png?raw=true" />

* PATCHES #3-10 are various combinations of one oscillator modulating
another, and in many cases mixing the modulation results with
some additive synthesis techniques. This makes the synthesis
process very unique and full of possibilities:

  * `((1+2+3)~4) + ((5+6+7)~8)`<br><img style="zoom:50%" src="https://github.com/chinenual/synergize/raw/master/docs/screenshots/patch-type-3.png?raw=true" />

  * `(1~2+3)~4) + ((5~6+7)~8)`<br><img style="zoom:50%" src="https://github.com/chinenual/synergize/raw/master/docs/screenshots/patch-type-4.png?raw=true" />

  * `(1~2) + 3 + 4 + (5~6) + 7 + 8`<br><img style="zoom:50%" src="https://github.com/chinenual/synergize/raw/master/docs/screenshots/patch-type-5.png?raw=true" />

  * `((1~2)~3) + ((1~2)~4) + ((5~6)~7) + ((5~6)~8)`<br><img style="zoom:50%" src="https://github.com/chinenual/synergize/raw/master/docs/screenshots/patch-type-6.png?raw=true" />

  * `(1~2) + (1~3) + (1~4) + (1~5) + (1~6) + (1~7) + (1~8)`<br><img style="zoom:50%" src="https://github.com/chinenual/synergize/raw/master/docs/screenshots/patch-type-7.png?raw=true" />

  * `(1~2~3) + (1~2~4) + (1~2~5) + (1~2~6) + (1~2~7) + (1~2~8)`<br><img style="zoom:50%" src="https://github.com/chinenual/synergize/raw/master/docs/screenshots/patch-type-8.png?raw=true" />

  * `(1~2~3~4) + (1~2~3~5) + (1~2~3~6) + (1~2~3~7) + (1~2~3~8)`<br><img style="zoom:50%" src="https://github.com/chinenual/synergize/raw/master/docs/screenshots/patch-type-9.png?raw=true" />

  * `((1~2+3)~4) + ((1~2+3)~5) + ((1~2+3)~6) + ((1~2+3)~7) + ((1~2+3)~8)`<br><img style="zoom:50%" src="https://github.com/chinenual/synergize/raw/master/docs/screenshots/patch-type-10.png?raw=true" />

These patches are changable from the Synergize `Voice Tab`
via the `Patch Type` selector, or by directly editing any of the
oscilator patch registers.  By selecting register combinations that
don't match the factory defaults, you can create unique and custom modulation schemes.

The conclusion and concept is that the GDS has the accuracy: of
additive synthesis, and the conservation of combinational forms
of synthesis as needed for certain sounds.

## Patch Registers

The synthesizer circuit is a time multiplexed microprocessor 
controlled  digital  oscillator which provides the user  with  32 
independent audio oscillator functions.  Each oscillator function 
also  includes  a frequency ramp  generator,  an  amplitude  ramp 
generator,  and a timer.  A patching network is provided allowing 
each  oscillator  to  receive additive  and/or  frequency  offset 
inputs from other oscillator's outputs.

The  oscillator  output  is a series of samples  or  numbers 
(32,000 per second) which represents the synthesized  sound.  All 
of  the synthesis,  amplitude control,  and sound mixing is  done 
while  the  sound  is  represented digitally  by  the  series  of 
samples.  The  changing  analog voltage needed to  represent  the 
audio  external  of  the  systhesizer is generated by  a  16  bit 
digital-to-analog  converter  (D/A) acting upon  this  series  of 
samples.

Actually,  only  one high speed oscillator circuit is  used, 
and is time multiplexed 32 ways in order to generate the multiple 
oscillator functions.  Each oscillator function uses 11 registers 
for  control.

The  oscillator  circuit performs all  of  the  calculations 
necessary to generate one sample of an oscillator output in about 
1 usec. A final output sample is required only every 32 usec. for 
each channel.  Each 1 usec  calculation period is called a "time 
slot",  and these are numbered 0,  1,  2,  ..., 31. The output of 
oscillator number 0 is calculated during time slot 0,  the output 
of  oscillator number 1 is calculated during time slot  1,  etc., 
until  the  output of oscillator number 31 is  calculated  during 
time  slot 31.  It is the completion of oscillator number 31 that 
initiates a D/A converter cycle.

Each  oscillator has two inputs that may receive  data  from 
the  outputs of other oscillators.  One of these is the frequency 
offset ("Freq") input used to produce vibrato or "voicing"  effects. 
The  second  is the "Adder" input which is used  to  combine  the 
output  of  some other oscillator with the output of the  current 
oscillator.  Four  temporary  data storage  registers  (DSR)  are 
provided  to  receive this combined output as well as supply  the 
data for the two types of oscillator inputs.

The  oscillators  may be interconnected using the  DSRs,  as 
specified by the control bytes. One oscillator output may be used 
to frequency modulate another oscillator,  and the output of this 
oscillator added to the output of another,  etc.  The calculation 
proceeds  sequentially (i.e.  0,  1,  ...,  31) with the  current 
oscillator's output available to any later calculated oscillator

Four  data  storage  registers (DSR) are available  to  pass 
oscillator  outputs to to the inputs of  other  oscillators.  The 
"PATCH"  control  byte contains three fields to specify which  of 
these  registers  is used for the  FM  input,  adder  input,  and 
output.

<p class="callout">
<b>NOTE:</b> SYNHCS allowed use of only 2 of the 4 registers. Synergize supports all 4.
</p>


The value of the Freq register is
used to shift the frequency of the oscillator. 

Since  each DSR may be used several times during each sample period
to  pass  data between  oscillators,
almost  any  interconnection between oscillators may be accomplished. However, care 
must be taken to assign the oscillators in the proper order.  For 
a  group  of  oscillators  interconnected,  the  lowest  numbered 
oscillator  should  be  assigned to  produce  the  signal  needed 
earliest in the calculation.

Note  that  the  output of oscillator 31 may be used  as  an 
input to oscillator 0 or a later calculated oscillator.

## Harmonics and Detune

HARMONIC values can be assigned to any oscillator, no matter
how it is used. In addition to the standard harmonic series
from 1 (the fundamental) through the 30th harmonic, "s"
harmonics are possible, 1s,2s,3s, etc. These represent the
Semi-tones that would exist between the fundamental and 2nd
harmonics. Access to them is especially useful in modulation
applications and very useful in helping to produce "unmusical"
qualities in sounds such as breath, bow scratch and lip buzz.

Harmonics do not necessarily have to be in tune. In most acoustic
instruments they are often not in tune, or are fluctuating under
different playing activities. Therefore DETUNING of harmonics
can be accomplished roughly 1/30hz at a time. Additionally,
degrees of random detuning fluctuations are possible per
oscillator. This is useful for many chorusing and imperfect
treatments.

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

### Envelopes Types

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

## Filters

* Af FILTERS are filters (one per voice, but it can effect several
oscillators). which are used to have the same values effect
several oscillators. The Af is used mostly for modulators. In
certain frequency ranges along the keyboard, the modulation
may be excessive in one range and not enough in another. This
filter is used to increase or decrease that effect by manipulating
the degree of modulation according to key number.

* BF FILTERS are filters that belong to one oscillator at a time.

You can have several B filters per voice, up to the number of
oscillators used in that voice. Again according to key number,
these values let you add or subtract from the amplitude values
of a particular oscillator in a particular range of the keyboard.
This is especially important for instruments whose timbral
characteristics change up and down the frequency range.

Filters can be used to make exaggerations of amplitude or modulation
up or down the frequency range, such as in the case where every
few keys produce a different instrumental timbre or a different
sound effect. This is how voice #24 of the internals, which has
different percussion sounds in different ranges of the keyboard,
is accomplished.

Filters can be used to shut off certain ranges of the keyboard,
such as extremes that are not needed, or at some voiced split
point.

Remember that the filters are SOFTWARE FILTERS , not digital filtering.

Use the FILTER GRAPH to help in achieving smooth contours of
filtering from one range to the next.

### Keyboard Equalization

KEYBOARD EQUALIZATION is the display used to raise or alter groups
of 4 keys at a time with respect to volume. In additive synthesis
the ear may hear certain qualities louder or softer than they should
be for a particular instrument. The ability to raise or lower the
volume levels according to frequency range will make a more
realistic balancing of the sound. It works for modulated sounds
in the same way, but effects only the additive (listened to)
oscillators. Alterations to the modulating oscillators would
alter timbre, not desireable in this display.

TRANPOSE should be used in the KEYBOARD EQUALIZATION display
to assure that extreme ranges are in balance with respect to
volume, especially should a user transpose a voice to a
different range during a performance patch. (This use of
transposition is also suggested when setting FILTER values
for the same reason).

### Key Proportionality

KEY PROPORTIONALITY is used to alter the decay rates of a sound ©
on a global basis or selected oscillators at a time. Those
effected have a (k) assigned in the VOICE DISPLAY. It is used
for PIANO (to make the upper part of the keyboard decay more
quickly than the lower) and other instruments effected by
a time/decay curve. It is especially useful to make synthesizer
Sounds, due to the unique nature of its activity. You could have
a long decaying drone synthesizer sound at the low end of the keyboard,
while the upper end was a sharp,plucked sound. Or you could
have a plucked bass at the lower part of the keyboard, while the
upper end has a Slowly decaying lead synthesizer sound.

A more subtle use of KEY PROPORTIONALITY is with respect to the
decay times of modulators, dependent upon frequency range. This
is especially useful in making brass sounds, where the lower
lip buzz sound is to be less exaggerated than at the top end.

### Summary

A study of different voices will give the best starting points
for learning how to voice the GDS system. Almost all parameters
can be altered and listened to.

During the voicing process, the sequencer can still be used to
hear how a voice responds and sounds during playback.

Each voice can be colored and stored with vibrato values, desired
defaults etc.

Voices can be tried with sustain pedal, portamento values and
pitchbend treatments as a test of "playability".

As a reference, it should be remembered that the Velocity can control
the following parameters, assigned in the voicing process: Volume,
timbre, attack time, decay time, pitch degree, pitch time, harmonic
entry, speed of loops, modulation degree, speed of repeats, speed of
musical tremelos, sustain lengths, a combination of these and more.
As a guide to voicing, decide what you want the velocity to accomplish
in a certain voice, then try to attain it.
