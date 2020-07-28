---
layout: default
parent: Voicing Mode
title: Filter Tab
nav_order: 3
description: "Voicing mode - Filter Tab"
permalink: /docs/voicing-mode/filters
---

# Filter Tab

<img title="Voice Filters view" src="https://github.com/chinenual/synergize/raw/master/docs/screenshots/viewVCE_filters.png?raw=true" width="100%"/>

1. Click the `Filters` tab to select this Filters editor
2. `Filter` Selector - use this to select which filter to view or edit.
3. `Copy From` Selector - use this to copy the filter values from
another filter into this one. See [below](#copying-filters).
4. The chart shows a graphical view of the filter.  X axis is the frequency, Y axis
   is the amount of amplitude change.
5. The table at the bottom are numeric values of each filter data point. Adjust these up or
   down to change the behavior of the filter.  Changes are transmitted in real time to
   the Synergy, so you can hear the affect of the change immediately.

## Filter Basics

A Synergy filter affects the amplitude of the oscilator it is paired
with. It is a table where each value affects a given frequency range,
and the value is a + / - amplitude adjustment to the oscilator result.

### A Filters and B Filters

The Synergy supports two sets of filters.   Each voice can have an "A
filter".  That A filter is shared by any oscillator that is configured
to use the A filter.

An oscillator can instead use a dedicated "B filter".  In this case,
the filter is used only by that oscillator.  If another oscilator
declares a B filter, an separate, independently configured filter is
used.

## Copying Filters

You can copy the entire set of filter values from another filter by
using the `Copy From` selector.  This copies the values from the
selected filter to the one currently being displayed.

## Adding or Removing Filters

To add or remove a filter, adjust the oscillator's `Filter` selector 
on the [Voice Tab](voicing-voice.md#adjust-filters).
