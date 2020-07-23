module.exports = {
  name: 'CATHERG',
  voicetab: {
    VNAME: { value: 'CATHERG'},
    nOsc: '4' ,
    keysPlayable: '8',
    vibType: 'Sine',
    VIBDEP: { value: '0' },
    VIBRAT: { value: '16' },
    VIBDEL: { value: '1' },
    APVIB: { value: '0' },
    VACENT: { value: '28' },
    VASENS: { value: '0' },
    VTCENT: { value: '0' },
    VTSENS: { value: '31' },
    VTRANS: { value: '-12' },
    nFilter: '1',
    patchType: { value: '0' },

    "patchFOInputDSR[1]": { value: '' },
    "patchFOInputDSR[2]": { value: '' },
    "patchFOInputDSR[3]": { value: '' },
    "patchFOInputDSR[4]": { value: '' },
    "patchFOInputDSR[5]": { exist: false },

    "patchAdderInDSR[1]": { value: '1' },
    "patchAdderInDSR[2]": { value: '1' },
    "patchAdderInDSR[3]": { value: '1' },
    "patchAdderInDSR[4]": { value: '1' },

    "patchOutputDSR[1]": { value: '1' },
    "patchOutputDSR[2]": { value: '1' },
    "patchOutputDSR[3]": { value: '1' },
    "patchOutputDSR[4]": { value: '1' },

    "OHARM[1]": { value: '4' },
    "OHARM[2]": { value: '2' },
    "OHARM[3]": { value: '8' },
    "OHARM[4]": { value: '1' },
    "OHARM[5]": { exist: false },

    "FDETUN[1]": { value: '-114' },
    "FDETUN[2]": { value: '-114' },
    "FDETUN[3]": { value: '-114' },
    "FDETUN[4]": { value: '-114' },

    "wkWAVE[1]": { value: 'Tri' },
    "wkWAVE[2]": { value: 'Tri' },
    "wkWAVE[3]": { value: 'Tri' },
    "wkWAVE[4]": { value: 'Tri' },

    "wkKEYPROP[1]": { selected: false },
    "wkKEYPROP[2]": { selected: false },
    "wkKEYPROP[3]": { selected: false },
    "wkKEYPROP[4]": { selected: false },

    "FILTER[1]": { value: '' },
    "FILTER[2]": { value: '' },
    "FILTER[3]": { value: '3' },
    "FILTER[4]": { value: '' },
  },

  envelopestab: {
    selections: [
      {
        select: { value: '1', text: '1' },
        "envFreqLoop[1]": { value: '' },
        "envFreqLowVal[1]": { value: '0' },
        "envFreqUpVal[1]": { value: '0' },
        //"envFreqTime[1]": { value: '0' }, // times in first row are fixed
        //"envFreqUpTime[1]": { value: '0' },

        "envFreqLoop[2]": { visible: false },

        "envAmpLoop[1]": { value: 'S' },
        "envAmpLowVal[1]": { value: '0' },
        "envAmpUpVal[1]": { value: '58' },
        "envAmpLowTime[1]": { value: '20' },
        "envAmpUpTime[1]": { value: '40' },

        "envAmpLoop[2]": { visible: true },
        "envAmpLoop[3]": { visible: false },

        "accelFreqLow": { value: '30' },
        "accelFreqUp": { value: '30' },
        "accelAmpLow": { visible: false },
        "accelAmpUp": { visible: false },

      },
      // skip the rest
    ],
  },

  filterstab: {
    select: { selector: '#filterSelect', value: '3', text: 'Bf 1' },
    selections: [
      {
        select: { value: '3', text: 'Bf 1' },
        "flt[8]": { value: '0' },
        "flt[16]": { value: '0' },
        "flt[24]": { value: '0' },
        "flt[31]": { value: '-64' },
      }],
  },
  keyeqtab: {
    "keyeq[6]": { value: '0' },
    "keyeq[12]": { value: '0' },
    "keyeq[18]": { value: '-6' },
    "keyeq[24]": { value: '-10' },
  },
  keyproptab: {
    "keyprop[6]": { value: '0' },
    "keyprop[12]": { value: '0' },
    "keyprop[18]": { value: '0' },
    "keyprop[24]": { value: '0' },
  },
};
