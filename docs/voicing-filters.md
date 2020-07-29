---
layout: default
parent: Voicing Mode
title: Filter Tab
nav_order: 3
description: "Voicing mode - Filter Tab"
permalink: /docs/voicing-mode/filters
---

# Filter Tab

![screenshot](/synergize/docs/screenshots/viewVCE_filters_annotated.png)

1. Click the `Filters` tab to select this Filters editor
2. `Filter` Selector - use this to select which filter to view or
   edit.  If there are more than one filters, the default display will
   be "All" (and the table below the graph will be hidden). You must
   select one of the filters in order to be able to edit it.
3. `Copy From` Selector - use this to copy the filter values from
another filter into this one. See [below](#copying-filters).
4. The chart shows a graphical view of the filter.  X axis is the frequency, Y axis
   is the amount of amplitude change.
5. The table at the bottom are numeric values of each filter data point. Adjust these up or
   down to change the behavior of the filter.  Changes are transmitted in real time to
   the Synergy, so you can hear the affect of the change immediately.

## Filter Basics

A Synergy filter affects the amplitude of the oscillator it is paired
with. It is a table where each value affects a given frequency range,
and the value is a + / - amplitude adjustment to the oscillator result.

### A Filters and B Filters

The Synergy supports two sets of filters.   Each voice can have an "A
filter".  That A filter is shared by any oscillator that is configured
to use the A filter.The Af is used mostly for modulators. In
certain frequency ranges along the keyboard, the modulation
may be excessive in one range and not enough in another. This
filter is used to increase or decrease that effect by manipulating
the degree of modulation according to key number.

An oscillator can instead use a dedicated "B filter".  In this case,
the filter is used only by that oscillator.  If another oscillator
declares a B filter, an separate, independently configured filter is
used.

According to key number, these values let you add or subtract from the
amplitude values of a particular oscillator in a particular range of
the keyboard.  This is especially important for instruments whose
timbral characteristics change up and down the frequency range.

Filters can be used to make exaggerations of amplitude or modulation
up or down the frequency range, such as in the case where every
few keys produce a different instrumental timbre or a different
sound effect. This is how voice #24 of the internals, which has
different percussion sounds in different ranges of the keyboard,
is accomplished.

## Copying Filters

You can copy the entire set of filter values from another filter by
using the `Copy From` selector.  This copies the values from the
selected filter to the one currently being displayed.

## Adding or Removing Filters

To add or remove a filter, adjust the oscillator's `Filter` selector 
on the [Voice Tab](voicing-voice.md#adjust-filters).

