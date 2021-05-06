---
layout: default
title: Alternate Tunings
nav_order: 60
description: Alternate Tunings
permalink: /docs/tunings
---

# Alternate Tunings

Many Synergy enthusiasts discovered the instrument due to Wendy Carlos's use of the Synergy on her Digital Moonscapes and Beauty and the Beast albums.   The later featured a variety of alternate tunings, including several invented by Carlos. Unfortunately, the Synergies she used had custom firmware (known by the community as "the WENDY firmware") and factory firmware does not support altering the instrument's tuning tables.

## Synergy firmware supporting alternate tunings

The WENDY firmware can be assembled from [Z80 assembly sources](https://chinenual.github.io/synergize/docs/acknowledgements#original-z80-firmware-synhcs-source-code-and-voice-libraries) by specifying a special macro variable (coincidentally named "WENDY").  This repurposes 256 bytes of the sequencer memory to support runtime-editable frequency tables.   Jari Kleimola has created ROMable images for hardware Synergies, available at XXXXX.   His virtual Synergy emulator, Synergia, also supports runtime tuning changes.

Synergize can send tuning changes to both Synergia and hardware Synergies running the WENDY firmware.

<p class="callout">
<b>WARNING:</b> Attempting to send tuning changes to a "factory" hardware Synergy may result in
  corrupted sequencer memory.  If you accidentally send a tuning change to a Synergy that doesn't
  support tuning  changes, you can restore your CMOS settings via a "factory reset" by pressing
  <code class="highlighter-rouge">RESTORE</code> then <code class="highlighter-rouge">SAVE</code>
  on the Synergy's front panel.
</p>

## Loading an alternate tuning

Synergize allows you to specify tunings via [Scala](http://www.huygens-fokker.org/scala/) [SCM](http://www.huygens-fokker.org/scala/scl_format.html) and [KBM](http://www.huygens-fokker.org/scala/help.htm#mappings) files. 

To load an alternate tuning into your Synergy, select `Load/Save -> Load Alternate Tuning`. Specify the tuning via the Alternate Tunings parameter form:

* `Use Standard Tuning` - when selected, the tuning will be standard 12 note Equal Temperament; the same as the factory tuning.

   * Even if Standard Tuning is selected, you can control how the tunings are applied to MIDI notes:
       * `Middle Note` - the MIDI note of first element of the scale (the "root" of the scale)
       * `Reference Frequency` - the frequency of a specific note on the keyboard
       * `Reference Note` - the MIDI note that the reference frequency specifies

       For factory tuning, specify `Middle Note = 60` (C4), `Reference Note = 69` (A4)  and
       `Reference      Frequency = 440` (i.e., A4 == 440Hz). 

* If `Use Standard Tuning` is unselected, the you must specify an alternate tuning via an SCM file.

   * If your scale has 12 notes, it is possible to use the standard keyboard layout and specify the
     reference note and frequency as above. Check `Use Standard Keyboard Mapping`.

   * If your scale has more or fewer than 12 notes, you will probably want to specify how the scale
     is mapped the the keyboard with a KBM file (you can skip notes in the scale and/or map some
     notes to more than one key).  Uncheck `Use Standard Keyboard Mapping` and
     select a KBM file. 

* To view the frequencies computed by the selected tuning parameter, press `Show Frequency Table`.  

* To send the new tuning to your Synergy, press `Send to Synergy`.



