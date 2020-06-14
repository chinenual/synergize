module.exports = {
  name: 'G7S',
  voicetab: {
    name: 'G7S',
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
  /*
  envelopestab: {
    select: { selector: '#envOscSelect', value: '1', text: '1' },
    selections: {
      select: { value: '1', text: 'Bf1' },
      "flt[8]": { value: '-5' },
      "flt[16]": { value: '2' },
      "flt[24]": { value: '0' },
      "flt[32]": { value: '0' },
    },
  },*/
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
