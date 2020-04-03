# VCE file format

## Header

| Type   | Count | Name   | Description 
|--------|-------|--------|-----
| uint8  | 1     | VOITAB | index of last oscillator (osc count minus 1) 
| uint16 | 16    | OSCPTR | array of 16 offsets (one for each osc table (zero based)) 
| uint8  | 1     | VTRANS | current transpose (+/- semitone) 
| uint8  | 1     | VTCENT | envelope center (0 .. 32)
| uint8  | 1     | VTSENS | envelope sensitivity
| uint8  | 1     |        | <unused - always 0>
| uint8  | 24    | VEQ    | array of voice amplitude for each key (one value per every 4 keys)
| uint8  | 8     | VNAME  | voice name. left justified, blank padded
| uint8  | 1     | VACENT | amplitude center
| uint8  | 1     | VASENS | amplitude sensitivity
| uint8  | 1     | VIBRAT | default vibrato rate (0..127)
| int8   | 1     | VIBDEL | default vibrato depth (-128..127)
| uint8  | 24    | KPROP  | key based proportion table
| uint8  | 1     | APVIB  | depth control for random vibrato
| int8   | 16    |        | filter assignment for each osc. - =AFILT, 0=NONE, + =BFILT

## Envelopes

One pair of Frequency and Amplitude envelopes for each Oscillator

### Frequency Envelope
| Type   | Count | Name   | Description 
|--------|-------|--------|-----
| uint8  | 	1    | OPTCH | 
| int8   | 	1    | OHARM | 
| int8   | 	1    | FDETUN | 
| uint8  | 	1    | FENVL | 
| uint8  | 	1    | ENVTYPE | 
| uint8  | 	1    | NPOINTS | 
| uint8  | 	1    | SUSTAINPT | 
| uint8  | 	1    | LOOPPT |
| uint8  |  4    | <elements> | 4 bytes per point


### Amplitude Envelope
| Type   | Count | Name   | Description 
|--------|-------|--------|-----
| uint8 | 	1 | ENVTYPE | 
| uint8 | 	1 | NPOINTS | 
| uint8 | 	1 | SUSTAINPT | 
| uint8 | 	1 | LOOPPT | 
| uint8  |  4    | <elements> | 4 bytes per point

## Filters

One filter for any non-zero entry in the header FILTER array

| Type   | Count | Name   | Description 
|--------|-------|--------|-----
| uint8  | 32    | <elements>  | each filter table is 32 bytes long

