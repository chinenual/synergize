---
layout: default
title: DX7 to Synergy Converter
nav_order: 80
description: DX7 to Synergy Converter
permalink: /docs/dx2syn
---

# DX7 to Synergy Converter

## Motivation

As one of the most popular digital synthesizers of all time, there
are thousands of DX7 voices available online.  The DX7 and
Synergy share some common characteristics and sounds developed for the
DX7 can be translated to similar sounds on the Synergy.  The converter
hopes to make some of those DX7 voices available on the Synergy, or
more realistically, make some of those voices an inspiration and
starting point for your own unique Synergy voices.

## SYX to VCE and CRT

For each SYX file, the converter converts each voice in the sysex and
produces a VCE file and then creates one or more CRT files containing
all the converted Synergy voices.  VCEs, DOCs and CRTs are created in
a folder named via the SYX base filename.  For example, given a SYX
named `foo/bar.SYX`, the converter would produce:

```
foo/bar.SYX
foo/bar/bar.CRT
foo/bar/bar-2.CRT
foo/bar/voice1.DOC
foo/bar/voice1.VCE
foo/bar/voice2.DOC
foo/bar/voice2.VCE
...
foo/bar/voice32.DOC
foo/bar/voice32.VCE
```

Synergy voice names attempt to preserve the original DX name, however
this sometimes involves compressing the 10-character DX name to the
available 8-characters on the Synergy. The converter attempts to
preserve 'significant' characters and prefers to eliminate spaces,
punctuation and vowels. If more than one voice have the same name
after such conversion, the converter appends a numeric suffix to
ensure that each name is unique.  For example,

| *DX name*      | *Synergy Name* |
| <tt>"--woohoo--"</tt>  |  <tt>"--woohoo"</tt>
| <tt>"FabVoice12"</tt>  |  <tt>"FabVoc12"</tt>
| <tt>"Brass&#9251; &#9251; &#9251;01"</tt>| <tt>"Brass&#9251;01"</tt> |


Since a DX7 sysex can contain up to 32 voices and a Synergy CRT can
support only 24 (and depending on their complexity can sometimes 
fit less into a CRT image), the converter puts the first set of VCE's into the
`bar.CRT` and then creates a `bar-2.CRT` for any voices that could not
fit in the first one.

## User Interface

Click the `Load/Save` button and select a `Convert DX7` menu item to
initiate the conversion.  On Mac, there's a single option since the
Mac file selector allows either folder or file selection.  On
WIndows and Linux there are separate menu items for selecting a folder
or selecting an individual SYX file.

When selecting an individual SYX file, only the voices in that file
are converted.    When selecting a folder, the converter recursively
walks that folder's contents (including sub-folders) and converts
every sysex it finds.

## Conversion Approach

The DX7 and Synergy synthesis architectures share some
characteristics, but are not identical.  In some cases there is a
one-to-one correspondence of a parameter in a DX7 patch to a
corresponding parameter in a Synergy voice.  However, there are some
DX7 parameters that have no equivalance on the Synergy.

The Synergy carrier and modulating patches are, for the most part,
identical to the original DX7 algorithms, subject to the following constraints:

* Due to the order that oscillators are evaluated by the Synergy osc board,
oscillator numbering is reversed in Synergy vs. DX -- DX oscillator 1
converts to Synergy oscillator 6, DX 2 converts to Synergy 5, etc.

* The Synergy has no way to feed back a signal to an oscillator
modulation chain.  Oscillators with feedback use the Triangle wave as
an approximation of of the DX sound. On the DX, the waveform produced by 
full feedback level (7) produces a Sawtooth wave, so the Triangle wave is 
closer than the standard Sine wave. 

* DX pitch and amplitude envelopes are directly converted to Synergy
envelopes. The Key Acceleration value on the DX is translated to the LOW (Yellow) 
Envelope, to give the sound change from soft to hard key strike.  You may need to
scale the Lower or Upper Envelope gain to your playing. 

* DX key level scaling are converted to Synergy filters.

* DX algorithm 16 and 17 cannot be implemented directly on the
Synergy. Instead, we use a special 5-oscillator patch which produces
similar sounding voices.

* Modulation levels for carriers have been scaled to emulate the sounds
of the original DX voices.  Some algorithm-specific adjustments are
made automatically during the conversion. One of these is due to the fact that
3 and 4 Oscillator 'towers', as they are decribed in the DX7 world, are not
handled the same in the Synergy. The 3rd and 4th upper Oscillators are lowered 
to about 25% of the original level. Since this isn't a one size fits all fix, 
some voices may be helped by raising the gain on these Oscillators. 

* One of the features on the DX7 that is not on the Synergy is being able
to set the Harmonic for the Oscillator at 0.5.  This is emulated by setting the Synergy 
Oscillator at 1, and raising the other Oscillators by 1 octave, then setting the 
Transpose down 1 octave so it matches the DX7 voice. 

## Recommendations

While some DX voices convert very accurately to Synergy, do not expect
a 100% accurate conversion in general. Although both are considered to be 
FM Synthesizers, there are multiple 'flavors' of FM Synthesis. The Carriers 
sound identical in additive mode, but when modulated, the DX7 and the Synergy
do sound different.  After creating the Synergy voice you may need to tweak it 
before it sounds right.

Things to try:

* Adjust oscillator gain levels to change frequency modulation and
  oscillator balance
* Adjust envelope rates. Add extra envelope points.
* Tweak filters 

In some cases, the initial Synergy voice may not sound much like the
DX voice at all -- but may inspire you to create a completely new
sound for the Synergy. 
