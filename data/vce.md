# VCE file format

## Header

| Type   | Count | Name   | Description 
|--------|-------|--------|-----
| byte  | 1     | VOITAB | index of last oscillator (osc count minus 1) 
| uint16 | 16    | OSCPTR | array of 16 offsets (one for each osc table (zero based)) 
| int8  | 1     | VTRANS | current transpose (+/- semitone) 
| byte  | 1     | VTCENT | envelope center (0 .. 32)
| byte  | 1     | VTSENS | envelope sensitivity
| byte  | 1     |        | \<unused - always 0\>
| int8  | 24    | VEQ    | array of voice amplitude for each key (one value per every 4 keys)
| byte  | 8     | VNAME  | voice name. left justified, blank padded
| byte  | 1     | VACENT | amplitude center
| byte  | 1     | VASENS | amplitude sensitivity
| byte  | 1     | VIBRAT | default vibrato rate (0..127)
| int8   | 1     | VIBDEL | default vibrato delay (0..127)
| int8   | 1     | VIBDEP | default vibrato depth (-127..127)
| byte  | 24    | KPROP  | key based proportion table
| byte  | 1     | APVIB  | depth control for random vibrato
| int8   | 16    | FILTER | filter assignment for each osc. - =AFILT, 0=NONE, + =BFILT

## Envelopes

One pair of Frequency and Amplitude envelopes for each Oscillator

### Frequency Envelope
| Type   | Count | Name   | Description 
|--------|-------|--------|-----
| byte  | 	1    | OPTCH | 
| int8   | 	1    | OHARM | 
| int8   | 	1    | FDETUN | 
| byte  | 	1    | FENVL | 
| byte  | 	1    | ENVTYPE | 
| byte  | 	1    | NPOINTS | 
| byte  | 	1    | SUSTAINPT | 
| byte  | 	1    | LOOPPT |
| byte  |  4    | \<elements\> | 4 bytes per point


### Amplitude Envelope
| Type   | Count | Name   | Description 
|--------|-------|--------|-----
| byte | 	1 | ENVTYPE | 
| byte | 	1 | NPOINTS | 
| byte | 	1 | SUSTAINPT | 
| byte | 	1 | LOOPPT | 
| byte  |  4    | \<elements\> | 4 bytes per point

## Filters

One filter for any non-zero entry in the header FILTER array

| Type   | Count | Name   | Description 
|--------|-------|--------|-----
| int8   | 32    | \<elements\>  | each filter table is 32 bytes long

