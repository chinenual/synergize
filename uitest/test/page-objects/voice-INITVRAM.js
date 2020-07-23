module.exports = {
  name: 'INITVRAM',
  voicetab: {
    VNAME: { value: '' },
    nOsc: '1',
    keysPlayable: '32',
    vibType: 'Sine',
    VIBDEP: { value: '0' },
    VIBRAT: { value: '16' },
    VIBDEL: { value: '0' },
    APVIB: { value: '32' },
    VACENT: { value: '24' },
    VASENS: { value: '0' },
    VTCENT: { value: '0' },
    VTSENS: { value: '0' },
    VTRANS: { value: '0' },
    nFilter: '0',
    patchType: { value: '0' },

    "patchFOInputDSR[1]": { value: '' },
    "patchFOInputDSR[2]": { exist: false },

    "patchAdderInDSR[1]": { value: '1' },

    "patchOutputDSR[1]": {value: '1'},

    "OHARM[1]": { value: '1' },
    "OHARM[2]": { exist: false },

    "FDETUN[1]": { value: '0' },

    "wkWAVE[1]": { value: 'Sin' },

    "wkKEYPROP[1]": { selected: false },

    "FILTER[1]": { value: '' },
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

        "envAmpLoop[1]": { value: '' },
        "envAmpLowVal[1]": { value: '0' },
        "envAmpUpVal[1]": { value: '0' },
        "envAmpLowTime[1]": { value: '0' },
        "envAmpUpTime[1]": { value: '0' },

        "envAmpLoop[2]": { visible: false },

        "accelFreqLow": { value: '30' },
        "accelFreqUp": { value: '30' },
        "accelAmpLow": { value: '30' },
        "accelAmpUp": { value: '30' },

      },
    ],
  },


  filterstab: {
    select: { selector: '#filterSelect', value: '', text: '' },
    selections: [
    ],
  },
  keyeqtab: {
    "keyeq[6]": { value: '0' },
    "keyeq[12]": { value: '0' },
    "keyeq[18]": { value: '0' },
    "keyeq[24]": { value: '0' },
  },
  keyproptab: {
    "keyprop[6]": { value: '4' },
    "keyprop[12]": { value: '16' },
    "keyprop[18]": { value: '28' },
    "keyprop[24]": { value: '32' },
  },
};
