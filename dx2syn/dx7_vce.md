# DX7 VCE file format as used in 32 voice cart

| Count  | Name   | Description 
|-------|--------|-----
| 0     | Status Byte | 11110000 
| 1     |  ID #        | 67 = Yamaha 
| 2     | SubStatus   | 0sssnnnn    s=0   n=0 ; ch=1
| 3     | Format #    | 0fffffff    f=9;32 voices   
| 4     | Byte Cnt MS | 0bbbbbbb    b=4096 ; 32 voices
| 5     | Byte Cnt LS | 0bbbbbbb
| 6     | Data        | 0ddddddd   32 voices 
| 4102  | Data        | 0ddddddd   32 voices
| 4103  | CheckSum    | 0bbbbbbb   add 4096th byte and make the 2's complement
| 4104  | EOX | 11110111

##  indivual voice format    17 bytes per OP, 
##  32 consecutive 128 byte chunks.
##  KLS Key Level Scaling  

| Count  | Name   | Description 
|-------|-------|--------
|0   | OP6 EG Rate 1  | 0-99      
|1   | OP6 EG Rate 2  | 0-99   
|2   | OP6 EG Rate 3  | 0-99   
|3   | OP6 EG Rate 4  | 0-99   
|4   | OP6 EG Level 1  | 0-99    x .727
|5   | OP6 EG Level 2  | 0-99    x .727
|6   | OP6 EG Level 3  | 0-99    x .727
|7   | OP6 EG Level 4  | 0-99    x .727
|8   | KLS Break Point  | 0-99     $27 = C3
|9   | KLS Left Depth  | 0-99 
|10  | KLS Right Bepth  | 0-99 
|11  | KLS Left Curve 0-3, KLS Right Curve 0-3| x000LCrc  0=-LIN, 1=-EXP, 2=+EXP, 3-+LIN
|12  | OSC DETUNE, 0-14  OSC RATE SCALE 0-7| xDETUors    Detune 1:1  *3      
|13  | Key Velocity Sensitivity 0-7, Amplitude Mod Sensitivity 0-3| x00KVSam   temp ignore, will add later
|14  | OP6 Operator Output Level  | 0-99  
|15  | OSC Freq Coarse 0-31, Oscillator Mode 0-1|  xxCOARSm  Coarse=1:1 *,  Mode: 0=Ratio  1=Fixed
|16  | OP6 OSC Freq Fine  | 0-99  will be added to Harmonic...

## 17  thru 33 = OP5
## 34  thru 50 = OP4
## 51  thru 67 = OP3
## 68  thru 84 = OP2
## 85  theu 101 = OP1

| Count  | Name   | Description 
|-------|-------|--------
|102 | Pitch EG Rate 1  | 0-99 = 
|103 | Pitch EG Rate 2  | 0-99 = 
|104 | Pitch EG Rate 3  | 0-99 = 
|105 | Pitch EG Rate 4  | 0-99 = 
|106 | Pitch EG Level 1  | 0-99 = 0-72   x .727
|107 | Pitch EG Level 2  | 0-99 = 0-72   x .727
|108 | Pitch EG Level 3  | 0-99 = 0-72   x .727
|109 | Pitch EG Level 4  | 0-99 = 0-72   x .727
|110 | Algorithm Select  | 0-31     Temp use algo = 2 on Syn
|111 | OSC Key Sync 0-1,  Feedback  | 0-7 |  xxxxSfdb  ignore Feedback, temp ignore sync 
|112 | LFO Speed  | 0-99 = 0-127    x 1.28 
|113 | LFO Delay  | 0-99 = 0-127    x 1.28
|114 | LFO Pitch Mod Depth  | 0-99 = 0-127    x 1.28
|115 | LFO Amplitude Mod Depth  | 0-99 = 0-127    x 1.28
|116 | LFO Pitch MOD Sensitivity 0-7, LFO WAVE 0-5, LFO SYNC 0-1 |  xTPMwavS 
|117 | Transpose  | 0-48      1:1  may change
|118 to 127 | 10 Char Voice Name ASCII   |  Use  1st 8 letters
|128 | ????? OP ON/OFF 0=OFF 1=ON  |  xx123456 *2 

# OSC order will depend on algorithm  temp order usin Syn algo 2 is 5 6 3 4 1 2
# *1 If Coarse = .5 or .25, then .25 or .5 = 1 in Syn, and reset transpose
# *2 This not shown in all carts ????
# *3 will have to compare, since vales extremes are so different.  Not used in test voice 

