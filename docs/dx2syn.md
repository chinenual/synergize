---
layout: default
title: DX7 to Synergy Converter
nav_order: 80
description: DX7 to Synergy Converter
permalink: /docs/dx2syn
---

# DX7 to Synergy Converter

## SYX to VCE and CRT

For each SYSX file, the converter converts each voice in the SYSX and
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
this sometimes involves compressing the 12-character DX name to the
available 8-characters on the Synergy. The converter attempts to
preserve 'significant' characters and prefers to eliminate spaces,
punctuation and vowels. If more than one voice have the same name
after such conversion, the converter appends a numeric suffix to
ensure that each name is unique.  For example,

| *DX name*      | *Synergy Name* |
| <tt>"---woohoo---"</tt>  |  <tt>"--woohoo"</tt>
| <tt>"Brass&#9251; &#9251; &#9251; &#9251; &#9251;01"</tt>| <tt>"Brass&#9251;01"</tt> |


Since a DX7 SYSX can contain up to 32 voices and a Synergy CRT can
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
every SYSX it finds.

## Conversion Approach

The DX7 and Synergy synthesis architectures are similar, but not the
same.  In some cases there is a one-to-one correspondence of a
parameter in a DX7 patch to a corresponding parameter in a Synergy
voice.  However, there are some DX7 parameters that have no
equivalance on the Synergy.

The Synergy carrier and modulating patches are, for the most part,
identical to the original DX7 algorithms, subject to the following constraints:

### Differences

Due to the order that oscillators are evaluated by the Synergy osc board,
oscillator numbering is reversed in Synergy vs. DX -- DX oscillator 1
converts to Synergy oscillator 6, DX 2 converts to Synergy 5, etc.

The Synergy has no way to feed back a signal to an oscillator
modulation chain.  This is, however, common on the DX7.  Oscillators
with feedback use the Triangle wave as an approximation of

DX7 pitch and amplitude envelopes are directly converted to Synergy
envelopes.

DX key level scaling are converted to Synergy filters.

DX algorithm 16 and 17 cannot be implemented directly on the
Synergy. Instead, we use a special 5-oscillator patch which produces
similar sounding voices.

Modulation levels for carriers have been scaled to emulate the sounds
of the original DX voices.  Some algorithm-specific adjustments are
made automatically during the conversion.