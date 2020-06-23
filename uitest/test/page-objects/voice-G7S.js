module.exports = {
  name: 'G7S',
  voicetab: {
    VNAME: 'G7S',
    nOsc: { value: '4' },
    keysPlayable: '8',
    vibType: 'Sine',
    VIBDEP: { value: '0' },
    VIBRAT: { value: '0' },
    VIBDEL: { value: '0' },
    APVIB: { value: '0' },
    VACENT: { value: '15' },
    VASENS: { value: '27' },
    VTCENT: { value: '20' },
    VTSENS: { value: '3' },
    VTRANS: { value: '0' },
    nFilter: '2',
    patchType: { value: '7' },

    "patchFOInputDSR[1]": '',
    "patchFOInputDSR[2]": '2',
    "patchFOInputDSR[3]": '2',
    "patchFOInputDSR[4]": '2',
    "patchFOInputDSR[5]": { exist: false },

    "patchAdderInDSR[1]": '',
    "patchAdderInDSR[2]": '',
    "patchAdderInDSR[3]": '1',
    "patchAdderInDSR[4]": '1',

    "patchOutputDSR[1]": '2',
    "patchOutputDSR[2]": '2',
    "patchOutputDSR[3]": '1',
    "patchOutputDSR[4]": '1',

    "OHARM[1]": { value: '3' },
    "OHARM[2]": { value: '1' },
    "OHARM[3]": { value: '1' },
    "OHARM[4]": { value: '2' },
    "OHARM[5]": { exist: false },

    "FDETUN[1]": { value: '15' },
    "FDETUN[2]": { value: '0' },
    "FDETUN[3]": { value: '0' },
    "FDETUN[4]": { value: '0' },

    "wkWAVE[1]": { value: 'Sin' },
    "wkWAVE[2]": { value: 'Sin' },
    "wkWAVE[3]": { value: 'Sin' },
    "wkWAVE[4]": { value: 'Sin' },

    "wkKEYPROP[1]": { selected: true },
    "wkKEYPROP[2]": { selected: false },
    "wkKEYPROP[3]": { selected: true },
    "wkKEYPROP[4]": { selected: true },

    "FILTER[1]": { value: '1' },
    "FILTER[2]": { value: '2' },
    "FILTER[3]": { value: '' },
    "FILTER[4]": { value: '' },
  },

  envelopestab: {
    selections: [
      {
        select: { value: '1', text: '1' },
        "envFreqLoop[1]": { value: '' },
        "envFreqLowVal[1]": { value: '0' },
        "envFreqUpVal[1]": { value: '1' },
        //"envFreqTime[1]": { value: '0' }, // times in first row are fixed
        //"envFreqUpTime[1]": { value: '0' },

        "envFreqLoop[2]": { visible: false },

        "envAmpLoop[1]": { value: 'S' },
        "envAmpLowVal[1]": { value: '70' },
        "envAmpUpVal[1]": { value: '70' },
        "envAmpLowTime[1]": { value: '0' },
        "envAmpUpTime[1]": { value: '0' },

        "envAmpLoop[2]": { value: '' },
        "envAmpLowVal[2]": { value: '69' },
        "envAmpUpVal[2]": { value: '69' },
        "envAmpLowTime[2]": { value: '6576' },
        "envAmpUpTime[2]": { value: '6576' },

        "envAmpLoop[3]": { visible: false },

        "accelFreqLow": { value: '30' },
        "accelFreqUp": { value: '30' },
        "accelAmpLow": { visible: false },
        "accelAmpUp": { visible: false },

      },
      {
        select: { value: '2', text: '2' },
        "envFreqLoop[1]": { value: '' },
        "envFreqLowVal[1]": { value: '0' },
        "envFreqUpVal[1]": { value: '0' },
        //"envFreqTime[1]": { value: '0' }, // times in first row are fixed
        //"envFreqUpTime[1]": { value: '0' },

        "envFreqLoop[2]": { visible: false },

        "envAmpLoop[1]": { value: '' },
        "envAmpLowVal[1]": { value: '52' },
        "envAmpUpVal[1]": { value: '60' },
        "envAmpLowTime[1]": { value: '0' },
        "envAmpUpTime[1]": { value: '0' },

        "envAmpLoop[2]": { value: '' },
        "envAmpLowVal[2]": { value: '37' },
        "envAmpUpVal[2]": { value: '37' },
        "envAmpLowTime[2]": { value: '6576' },
        "envAmpUpTime[2]": { value: '6576' },

        "envAmpLoop[3]": { value: '' },
        "envAmpLowVal[3]": { value: '0' },
        "envAmpUpVal[3]": { value: '60' },
        "envAmpLowTime[3]": { value: '6576' },
        "envAmpUpTime[3]": { value: '6576' },

        "envAmpLoop[4]": { visible: false },

        "accelFreqLow": { value: '30' },
        "accelFreqUp": { value: '30' },
        "accelAmpLow": { value: '0' },
        "accelAmpUp": { value: '0' },

      },
      {
        select: { value: '3', text: '3' },
        "envFreqLoop[1]": { value: '' },
        "envFreqLowVal[1]": { value: '0' },
        "envFreqUpVal[1]": { value: '2' },
        //"envFreqTime[1]": { value: '0' }, // times in first row are fixed
        //"envFreqUpTime[1]": { value: '0' },

        "envFreqLoop[2]": { visible: false },


        "envAmpLoop[1]": { value: '' },
        "envAmpLowVal[1]": { value: '72' },
        "envAmpUpVal[1]": { value: '64' },
        "envAmpLowTime[1]": { value: '26' },
        "envAmpUpTime[1]": { value: '28' },

        "envAmpLoop[2]": { value: '' },
        "envAmpLowVal[2]": { value: '69' },
        "envAmpUpVal[2]": { value: '57' },
        "envAmpLowTime[2]": { value: '326' },
        "envAmpUpTime[2]": { value: '91' },

        "envAmpLoop[3]": { value: '' },
        "envAmpLowVal[3]": { value: '54' },
        "envAmpUpVal[3]": { value: '49' },
        "envAmpLowTime[3]": { value: '2929' },
        "envAmpUpTime[3]": { value: '326' },

        "envAmpLoop[4]": { value: 'S' },
        "envAmpLowVal[4]": { value: '1' },
        "envAmpUpVal[4]": { value: '0' },
        "envAmpLowTime[4]": { value: '6576' },
        "envAmpUpTime[4]": { value: '652' },

        "envAmpLoop[5]": { value: '' },
        "envAmpLowVal[5]": { value: '1' },
        "envAmpUpVal[5]": { value: '0' },
        "envAmpLowTime[5]": { value: '64' },
        "envAmpUpTime[5]": { value: '57' },

        "envAmpLoop[6]": { visible: false },

        "accelFreqLow": { value: '30' },
        "accelFreqUp": { value: '30' },
        "accelAmpLow": { visible: false },
        "accelAmpUp": { visible: false },

      },
      {
        select: { value: '4', text: '4' },
        "envFreqLoop[1]": { value: '' },
        "envFreqLowVal[1]": { value: '0' },
        "envFreqUpVal[1]": { value: '4' },
        //"envFreqTime[1]": { value: '0' }, // times in first row are fixed
        //"envFreqUpTime[1]": { value: '0' },

        "envFreqLoop[2]": { visible: false },


        "envAmpLoop[1]": { value: '' },
        "envAmpLowVal[1]": { value: '72' },
        "envAmpUpVal[1]": { value: '50' },
        "envAmpLowTime[1]": { value: '25' },
        "envAmpUpTime[1]": { value: '36' },

        "envAmpLoop[2]": { value: '' },
        "envAmpLowVal[2]": { value: '68' },
        "envAmpUpVal[2]": { value: '39' },
        "envAmpLowTime[2]": { value: '72' },
        "envAmpUpTime[2]": { value: '91' },

        "envAmpLoop[3]": { value: '' },
        "envAmpLowVal[3]": { value: '53' },
        "envAmpUpVal[3]": { value: '24' },
        "envAmpLowTime[3]": { value: '2929' },
        "envAmpUpTime[3]": { value: '24' },

        "envAmpLoop[4]": { value: 'S' },
        "envAmpLowVal[4]": { value: '0' },
        "envAmpUpVal[4]": { value: '0' },
        "envAmpLowTime[4]": { value: '6576' },
        "envAmpUpTime[4]": { value: '1035' },

        "envAmpLoop[5]": { value: '' },
        "envAmpLowVal[5]": { value: '0' },
        "envAmpUpVal[5]": { value: '0' },
        "envAmpLowTime[5]": { value: '129' },
        "envAmpUpTime[5]": { value: '81' },

        "envFreqLoop[6]": { visible: false },

        "accelFreqLow": { value: '30' },
        "accelFreqUp": { value: '30' },
        "accelAmpLow": { visible: false },
        "accelAmpUp": { visible: false },

      },

    ],
  },
  filterstab: {
    select: { selector: '#filterSelect', value: '-1', text: 'All' },
    selections: [
      {
        select: { value: '1', text: 'Bf 1' },
        "flt[8]": { value: '-5' },
        "flt[16]": { value: '2' },
        "flt[24]": { value: '0' },
        "flt[32]": { value: '0' },
      },
      {
        select: { value: '2', text: 'Bf 2' },
        "flt[8]": { value: '0' },
        "flt[16]": { value: '2' },
        "flt[24]": { value: '-19' },
        "flt[32]": { value: '-19' },
      }],
  },
  keyeqtab: {
    "keyeq[6]": { value: '0' },
    "keyeq[12]": { value: '3' },
    "keyeq[18]": { value: '6' },
    "keyeq[24]": { value: '-4' },
  },
  keyproptab: {
    "keyprop[6]": { value: '4' },
    "keyprop[12]": { value: '9' },
    "keyprop[18]": { value: '24' },
    "keyprop[24]": { value: '32' },
  },
};
