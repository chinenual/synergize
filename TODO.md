# SYNHCS functionality and corresponding Synergize implementation

The following is organized by the SYNCHS screens each function appears on:

First check is full UI support; second check is lower-level IO and tests

* [ ] [ ] Create CRT file from existing VCE files
* [ ] [ ] <s>Create a DOC file</s>
* [ ] [ ] <s>View a DOC file</s>
* [x] [x] Load Synergy Machine State
* [x] [x] Save Synergy Machine State
* [x] [x] Load VRAM image
* [ ] [x] Save VRAM image
* [-] [-] Edit VRAM image
* [x] [x] Disable VRAM
* [x] [x] Test communications
* [ ] [ ] <s>"Extended programmer" Performance Program controller (sequence of CRT/SYN files for live performance)</s>

* Voicing System Main screen
    * [x] [x] Load a VCE directly into VRAM
    * [ ] [ ] Save a VCE from VRAM
    * [x] [x] set # oscillators
    * [x] [x] set default patch
    * [ ] [ ] clear filters
    * [ ] [ ] import filters from other voice
    * [ ] [ ] clear voicing system
    * [ ] [ ] <s>print voicing parameters</s>

* Voice page
    * [x] [x] "ensemble"/"group"/"solo" oscillators
    * [ ] [ ] copy an oscillator
    * [x] [x] change harmonic
    * [x] [x] detune
    * [x] [x] aperiodic detune
    * [x] [x] select sine vs. triangle
    * [x] [x] assign A-filter, B-filter or no-filter to each oscillator
    * [x] [x] set timbre/amp sensitivity and center
    * [x] [x] <s>set portamento</s>
    * [x] [x] set vibrato
    * [x] [x] toggle keyboard proportionalization per oscillator
    * [x] [x] adjust the "patch" routing (adder, frequency, output registers)

* Envelopes page (AKA Oscillator page)
    * [ ] [ ] Change number of points in an envelope
    * [ ] [ ] set Sustain point
    * [ ] [ ] set Loop point
    * [ ] [ ] change envelope type
    * [ ] [ ] <s>scale up/down</s>
    * [ ] [ ] change freq, amp or time values for both upper and lower bounds
    * [ ] [ ] copy upper <-> lower

* Filters
    * [x] [x] set each point in the curve
    * [ ] [ ] <s>scale display</s>
    * [ ] [ ] copy other oscillator's filter

* Keyboard Equalization
    * [x] [x] set each point in the curve

* Keyboard Proportion
    * [x] [x] set each point in the curve
