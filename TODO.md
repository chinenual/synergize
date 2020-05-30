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
* [ ] [ ] Edit VRAM image
* [x] [x] Disable VRAM
* [x] [x] Test communications
* [ ] [ ] <s>"Extended programmer" Performance Program controller (sequence of CRT/SYN files for live performance)</s>

* Voicing System Main screen
    * [x] [x] Load a VCE directly into VRAM
    * [ ] [ ] Save a VCE from VRAM
    * [ ] [ ] set # oscillators
    * [ ] [ ] set default patch
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
    * [ ] [ ] select sine vs. triangle
    * [ ] [ ] assign A-filter, B-filter or no-filter to each oscillator
    * [x] [x] set timbre/amp sensitivity and center
    * [ ] [ ] set portamento
    * [x] [x] set vibrato
    * [ ] [ ] toggle keyboard proportionalization per oscillator
    * [ ] [ ] adjust the "patch" routing (adder, frequency, output registers)

* Envelopes page (AKA Oscillator page)
    * [ ] [ ] Change number of points in an envelope
    * [ ] [ ] set Sustain point
    * [ ] [ ] set Loop point
    * [ ] [ ] change envelope type
    * [ ] [ ] scale up/down
    * [ ] [ ] change freq, amp or time values for both upper and lower bounds
    * [ ] [ ] copy upper <-> lower

* Filters
    * [ ] [ ] set each point in the curve
    * [ ] [ ] scale display
    * [ ] [ ] copy other oscillator's filter

* Keyboard Equalization
    * [ ] [x] set each point in the curve

* Keyboard Proportion
    * [ ] [x] set each point in the curve
