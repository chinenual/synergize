---
layout: default
title: Voicing Mode
nav_order: 4
description: "Voicing mode"
permalink: /docs/voicing-mode
has_children: true
has_toc: false
---

# Voicing Mode

The design of each tab is closely modeled on screens from the original
SYNHCS software.

* [Voice](voicing-voice.md)
* [Envelopes](voicing-envs.md)
* [Filter](voicing-filters.md)
* [Key Equalization](voicing-keyeq.md)
* [Key Proportion](voicing-keyprop.md)

Those of you who know SYNHCS should have no trouble
using Synergize. 

A main difference is that SYNHCS relied on using buttons and knobs on
the Synergy front panel to adjust each value; in Synergize, all values
are edited directly in the application, so you use your keyboard and
touchpad or mouse rather than the Synergy knobs and buttons. 
     
Synergize allows full control of all editable parameters in the
underlying Synergy sound generator.  When the `Voicing Mode` button is
toggled on, each parameter displayed on the `Voice`, `Envelopes`,
`Filters`, etc. tabs becomes editable.  Changes made in the Synergize UI
are immediately transmitted through the serial port to the Synergy, so
the effect can be heard in real time.

## Default Voice

When Voicing Mode is initiated, Synergize loads an empty "default
voice" to the Synergy.  You can either start creating your new voice
from this default, or can load an existing VCE file and start from
there.

# Saving your edits

To save your new voice, select `Save/Load -> Save Voice (VCE)`.  If
you exit Voicing Mode without first saving your voice, your edits will
be discarded.
