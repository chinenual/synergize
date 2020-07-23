module.exports = {
  name: 'GUITAR2A',
  voicetab: {
    VNAME: {value: 'GUITAR2A'},
    nOsc: '4',
    keysPlayable: '8',
    vibType: 'Sine',
    VIBDEP: { value: '0' },
    VIBRAT: { value: '17' },
    VIBDEL: { value: '33' },
    APVIB: { value: '6' },
    VACENT: { value: '16' },
    VASENS: { value: '15' },
    VTCENT: { value: '16' },
    VTSENS: { value: '8' },
    VTRANS: { value: '0' },
    nFilter: '3',
    patchType: { value: '3' },

    "patchFOInputDSR[1]": { value: '' },
    "patchFOInputDSR[2]": { value: '2' },
    "patchFOInputDSR[3]": { value: '' },
    "patchFOInputDSR[4]": { value: '2' },
    "patchFOInputDSR[5]": { exist: false },

    "patchAdderInDSR[1]": { value: '' },
    "patchAdderInDSR[2]": { value: '' },
    "patchAdderInDSR[3]": { value: '2' },
    "patchAdderInDSR[4]": { value: '1' },

    "patchOutputDSR[1]": { value: '2' },
    "patchOutputDSR[2]": { value: '2' },
    "patchOutputDSR[3]": { value: '2' },
    "patchOutputDSR[4]": { value: '1' },

    "OHARM[1]": { value: 's5' },
    "OHARM[2]": { value: '1' },
    "OHARM[3]": { value: '3' },
    "OHARM[4]": { value: '1' },
    "OHARM[5]": { exist: false },

    "FDETUN[1]": { value: '0' },
    "FDETUN[2]": { value: '0' },
    "FDETUN[3]": { value: '0' },
    "FDETUN[4]": { value: '57' },

    "wkWAVE[1]": { value: 'Sin' },
    "wkWAVE[2]": { value: 'Sin' },
    "wkWAVE[3]": { value: 'Sin' },
    "wkWAVE[4]": { value: 'Sin' },

    "wkKEYPROP[1]": { selected: false },
    "wkKEYPROP[2]": { selected: false },
    "wkKEYPROP[3]": { selected: false },
    "wkKEYPROP[4]": { selected: true },

    "FILTER[1]": { value: '1' },
    "FILTER[2]": { value: '2' },
    "FILTER[3]": { value: '3' },
    "FILTER[4]": { value: '' },
  },

  envelopestab: {
    selections: [
      {
        select: { value: '4', text: '4' },
        "envFreqLoop[1]": { value: '' },
        "envFreqLowVal[1]": { value: '0' },
        "envFreqUpVal[1]": { value: '0' },
        //"envFreqTime[1]": { value: '0' }, // times in first row are fixed
        //"envFreqUpTime[1]": { value: '0' },

        "envFreqLoop[2]": { visible: false },

        "envAmpLoop[1]": { value: '' },
        "envAmpLowVal[1]": { value: '72' },
        "envAmpUpVal[1]": { value: '72' },
        "envAmpLowTime[1]": { value: '16' },
        "envAmpUpTime[1]": { value: '16' },

        "envAmpLoop[2]": { visible: true },
        "envAmpLoop[3]": { visible: true },
        "envAmpLoop[4]": { visible: true },
        "envAmpLoop[5]": { visible: false },

        "accelFreqLow": { value: '30' },
        "accelFreqUp": { value: '30' },
        "accelAmpLow": { value: '42'  },
        "accelAmpUp": { value: '37'  },

      },
      // skip the rest
    ],
  },


  filterstab: {
    select: { selector: '#filterSelect', value: '-1', text: 'All' },
    selections: [
      {
        select: { value: '1', text: 'Bf 1' },
        "flt[8]": { value: '-18' },
        "flt[16]": { value: '-18' },
        "flt[24]": { value: '0' },
        "flt[32]": { value: '0' },
      },
      {
        select: { value: '2', text: 'Bf 2' },
        "flt[8]": { value: '-4' },
        "flt[16]": { value: '0' },
        "flt[24]": { value: '0' },
        "flt[32]": { value: '0' },
      },
      {
        select: { value: '3', text: 'Bf 3' },
        "flt[8]": { value: '-20' },
        "flt[16]": { value: '-5' },
        "flt[24]": { value: '0' },
        "flt[32]": { value: '-19' },
      }],
  },
  keyeqtab: {
    "keyeq[6]": { value: '0' },
    "keyeq[12]": { value: '0' },
    "keyeq[18]": { value: '0' },
    "keyeq[24]": { value: '0' },
  },
  keyproptab: {
    "keyprop[6]": { value: '16' },
    "keyprop[12]": { value: '2' },
    "keyprop[18]": { value: '28' },
    "keyprop[24]": { value: '32' },
  },
};
