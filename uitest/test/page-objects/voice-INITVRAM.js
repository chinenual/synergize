module.exports = {
  name: 'INITVRAM',
  voicetab: {
    name: '',
    nOsc: { value: '1' },
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

    "patchFOInputDSR[1]": '',
    "patchFOInputDSR[2]": { exist: false },

    "patchAdderInDSR[1]": '1',

    "patchOutputDSR[1]": '1',

    "OHARM[1]": { value: '1' },
    "OHARM[2]": { exist: false },

    "FDETUN[1]": { value: '0' },

    "wkWAVE[1]": { value: 'Sin' },

    "wkKEYPROP[1]": { selected: false },

    "FILTER[1]": { value: '' },
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
