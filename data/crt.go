package data

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/chinenual/synergize/logger"
	"github.com/orcaman/writerseeker"
	"github.com/pkg/errors"
)

const verboseParsing = false
const verboseWriting = false

type CRTHead struct {
	VRAM   [51]byte   // "free storage"
	VOITAB byte       // always zero
	VCHK   [5]byte    // check bytes - each should be 170
	TSTVEC uint16     // test vector
	IPDVEC uint16     // ret. from input w/ data
	IPNVEC uint16     // ret. from input w/out data
	VERSON byte       // version
	VOIPTR [24]uint16 // 24 pointers to voice data (offsets from VOITAB)
	BFILTR [24]byte   // b-filter # start for voices
	AFILTR [24]int8   // a-filter # start for voices
}

type CRT struct {
	Head   CRTHead
	Voices []*VCE
}

func ReadCrtFile(filename string) (crt CRT, err error) {
	// A CRT file is a long header containing filter info, followed by a list of CCE fragments (each voice missing the filter params since they are concatenated elsewhere in the file).

	var b []byte

	if b, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	buf := bytes.NewReader(b)

	if crt, err = ReadCrt(buf); err != nil {
		return
	}
	return
}

func ReadCrt(buf io.ReadSeeker) (crt CRT, err error) {
	// A CRT file is a long header containing filter info, followed by a list of CCE fragments (each voice missing the filter params since they are concatenated elsewhere in the file).

	if err = binary.Read(buf, binary.LittleEndian, &crt.Head); err != nil {
		logger.Error("binary.Read failed:", err)
		return
	}

	if verboseParsing {
		logger.Info(crtHeadToJson(crt.Head))
	}

	for i, offset := range crt.Head.VOIPTR {
		if offset != 0 {
			if verboseParsing {
				logger.Infof("READ voice #%d: seek to 0x%04x -> 0x%04x\n", i+1, offset, Off_VRAM_VOITAB+offset)
			}

			if _, err = buf.Seek(int64(Off_VRAM_VOITAB+offset), io.SeekStart); err != nil {
				err = errors.Wrapf(err, "failed to seek to voice #%d start", i)
				return
			}
			var vce VCE
			if vce, err = ReadVce(buf, true); err != nil {
				err = errors.Wrapf(err, "failed to read voice #%d start", i)
				return
			}

			if VceAFilterCount(vce) > 0 {
				offset = uint16(crt.Head.AFILTR[i]-1) * 32
				if verboseParsing {
					logger.Infof("READ voice #%d: A-Filter: seek to %d 0x%04x -> 0x%04x\n", i+1, crt.Head.AFILTR[i], offset, Off_VRAM_FILTAB+offset)
				}
				if _, err = buf.Seek(int64(Off_VRAM_FILTAB+offset), io.SeekStart); err != nil {
					err = errors.Wrapf(err, "failed to seek to voice #%d filter-b start", i)
					return
				}
				if err = VceReadAFilters(buf, &vce); err != nil {
					return
				}
			}
			if VceBFilterCount(vce) > 0 {
				offset = uint16(crt.Head.BFILTR[i]-1) * 32
				if verboseParsing {
					logger.Infof("READ voice #%d: B-Filters: seek to %d 0x%04x -> 0x%04x\n", i+1, crt.Head.BFILTR[i], offset, Off_VRAM_FILTAB+offset)
				}
				if _, err = buf.Seek(int64(Off_VRAM_FILTAB+offset), io.SeekStart); err != nil {
					err = errors.Wrapf(err, "failed to seek to voice #%d filter-b start", i)
					return
				}
				if err = VceReadBFilters(buf, &vce); err != nil {
					return
				}
			}

			crt.Voices = append(crt.Voices, &vce)
		}
	}

	return
}

// Filters.
// in the VRAM_AFILTR array, indexes into the VRAM_FILTAB for each voice's A-filters (if any) (index is the voice number; value is 0 if no filter,
// a negative number whose absolute value is the index into the FILTAB)
//
// in the VRAM_BFILTR array, indexes into the VRAM_FILTAB for each voice's B-filters (if any) (as above except positive numbers)
//
// in the VRAM_FILTAB, all of the filters
//
// in each Voice Head, the FILTER array has negative entry of the oscilator uses an A filter;
// or a positive entry if a B filter - or 0 if no filter.
// these are essentially indexes offset by the voice's entry in AFILTR or BFILTR into the FILTAB

func addVce(buf io.WriteSeeker, slot /*one-based*/ int, cursor *crtCursor, vce VCE) (err error) {
	// update the VOIPTR entry in the CRT header

	if _, err = buf.Seek(int64(Off_VRAM_VOIPTR+((slot-1)*2)), io.SeekStart); err != nil {
		return
	}
	if verboseWriting {
		logger.Infof(" WRITE voice #%d: VOIPTR  0x%04x\n", slot, uint16(cursor.VoiceOffset)-Off_VRAM_VOITAB)
	}

	// VOIPTR offsets are relative to the VOITAB field not the start of the file (sigh...):
	if err = binary.Write(buf, binary.LittleEndian, uint16(cursor.VoiceOffset)-Off_VRAM_VOITAB); err != nil {
		return
	}
	var filterindex = byte(0xff) // -1
	if VceAFilterCount(vce) > 0 {
		filterindex = byte(int8(cursor.AfilterIndex))
	}
	// update the AFILTR entry in the CRT header
	if _, err = buf.Seek(int64(Off_VRAM_AFILTR+(slot-1)), io.SeekStart); err != nil {
		return
	}
	if verboseWriting {
		logger.Infof(" WRITE voice #%d: a-filter index %d at 0x%04x\n", slot, filterindex, int64(Off_VRAM_AFILTR+(slot-1)))
	}
	if err = binary.Write(buf, binary.LittleEndian, filterindex); err != nil {
		return
	}

	// if there's an A-filter, write it:
	if _, err = buf.Seek(cursor.AfilterOffset, io.SeekStart); err != nil {
		return
	}
	if verboseWriting {
		logger.Infof(" WRITE voice #%d: %d a-filters at 0x%04x\n", slot, VceAFilterCount(vce), cursor.AfilterOffset)
	}
	if err = VceWriteAFilters(buf, vce); err != nil {
		return
	}
	// update cursor and index for next voice
	if cursor.AfilterOffset, err = buf.Seek(0, io.SeekCurrent); err != nil {
		return
	}
	cursor.AfilterIndex = cursor.AfilterIndex + byte(VceAFilterCount(vce))

	filterindex = byte(0)
	if VceBFilterCount(vce) > 0 {
		filterindex = cursor.BfilterIndex
	}
	if verboseWriting {
		logger.Infof(" WRITE voice #%d: b-filter index %d at 0x%04x\n", slot, filterindex, int64(Off_VRAM_BFILTR+(slot-1)))
	}
	// update the BFILTR entry in the CRT header
	if _, err = buf.Seek(int64(Off_VRAM_BFILTR+(slot-1)), io.SeekStart); err != nil {
		return
	}
	if err = binary.Write(buf, binary.LittleEndian, filterindex); err != nil {
		return
	}
	// if there are B-filters, write them:
	if _, err = buf.Seek(cursor.BfilterOffset, io.SeekStart); err != nil {
		return
	}
	if verboseWriting {
		logger.Infof(" WRITE voice #%d: %d b-filters at 0x%04x\n", slot, VceBFilterCount(vce), cursor.BfilterOffset)
	}
	if err = VceWriteBFilters(buf, vce); err != nil {
		return
	}
	var oldOffset = cursor.BfilterOffset
	// update cursor and index for next voice
	if cursor.BfilterOffset, err = buf.Seek(0, io.SeekCurrent); err != nil {
		return
	}
	cursor.BfilterIndex = cursor.BfilterIndex + byte(VceBFilterCount(vce))
	if verboseWriting {
		logger.Infof("         WRITE voice #%d: b-filters cursor advanced: 0x%04x (%d filters)\n", slot, cursor.BfilterOffset-oldOffset, (cursor.BfilterOffset-oldOffset)/VRAM_FILTR_length)
	}

	// Now the voice itself:
	if verboseWriting {
		logger.Infof(" WRITE voice #%d: at 0x%04x\n", slot, cursor.VoiceOffset)
	}

	if cursor.VoiceOffset > VRAM_Max_length {
		err = errors.Errorf("Error adding voice #%d (%s) to CRT - requires %d bytes, which exceeds maximum of %d",
			slot, VceName(vce.Head), cursor.VoiceOffset, VRAM_Max_length)
		return
	}

	if _, err = buf.Seek(cursor.VoiceOffset, io.SeekStart); err != nil {
		return
	}
	if err = WriteVce(buf, vce, VceName(vce.Head), true); err != nil {
		return
	}
	// update cursor for next voice
	if cursor.VoiceOffset, err = buf.Seek(0, io.SeekCurrent); err != nil {
		return
	}
	return
}

func WriteCrtFileFromVCEPaths(filename string, vcePaths []string) (err error) {
	var vces []*VCE
	for _, path := range vcePaths {
		var vce VCE
		if vce, err = ReadVceFile(path); err != nil {
			return
		}
		vces = append(vces, &vce)
	}
	if err = WriteCrtFileFromVCEArray(filename, vces); err != nil {
		return
	}
	return
}

func WriteCrtFileFromVCEArray(filename string, vces []*VCE) (err error) {
	var writebuf = writerseeker.WriterSeeker{}

	if err = WriteCrt(&writebuf, vces); err != nil {
		return
	}
	var write_bytes []byte
	if write_bytes, err = ioutil.ReadAll(writebuf.Reader()); err != nil {
		return
	}
	if err = ioutil.WriteFile(filename, write_bytes, 0644); err != nil {
		return
	}
	return
}

type crtCursor struct {
	AfilterOffset int64
	AfilterIndex  byte
	BfilterOffset int64
	BfilterIndex  byte
	VoiceOffset   int64
}

func WriteCrt(buf io.WriteSeeker, vces []*VCE) (err error) {
	if len(vces) < 1 || len(vces) > 24 {
		err = errors.Errorf("Must have at least 1 and no more than 24 voices")
		return
	}
	// write the CRT header - omit the filters and rest of the edata since the voicing image
	// has empty spaces that we won't use for a CRT file
	if err = binary.Write(buf, binary.LittleEndian, VRAM_EDATA[:Off_VRAM_FILTAB]); err != nil {
		return
	}

	// Write the header string (first 51 bytes are available for "copyright" info)
	{
		var headerString = "Generated by Synergize " + appVersion
		if len(headerString) > 51 {
			headerString = headerString[:50]
		}
		// if there are B-filters, write them:
		if _, err = buf.Seek(0, io.SeekStart); err != nil {
			return
		}
		if err = binary.Write(buf, binary.LittleEndian, []byte(headerString)); err != nil {
			return
		}

	}
	// advance to the start of the first voice

	var aFilterCount = 0
	var bFilterCount = 0

	for _, vce := range vces {
		if vce != nil {
			aFilterCount = aFilterCount + VceAFilterCount(*vce)
			bFilterCount = bFilterCount + VceBFilterCount(*vce)
		}
	}
	if verboseWriting {
		logger.Infof(" WRITE a filter count: %d, b filter count: %d\n", aFilterCount, bFilterCount)
	}
	var cursor crtCursor
	cursor.AfilterOffset = Off_VRAM_FILTAB
	cursor.AfilterIndex = 1 // one-based index
	cursor.BfilterOffset = cursor.AfilterOffset + int64(aFilterCount*VRAM_FILTR_length)
	cursor.BfilterIndex = byte(aFilterCount + 1) // one past the last A filter
	cursor.VoiceOffset = cursor.BfilterOffset + int64((bFilterCount)*VRAM_FILTR_length)

	if verboseWriting {
		logger.Infof(" WRITE cursors before first voice: a, b voice: 0x%04x 0x%04x 0x%04x\n", cursor.AfilterOffset, cursor.BfilterOffset, cursor.VoiceOffset)
	}
	for i, vce := range vces {
		if vce != nil {
			if err = addVce(buf, i+1, &cursor, *vce); err != nil {
				return
			}
		}
	}
	if verboseWriting {
		logger.Infof(" WRITE cursors after last voice: a, b voice: 0x%04x 0x%04x 0x%04x\n", cursor.AfilterOffset, cursor.BfilterOffset, cursor.VoiceOffset)
	}
	return
}

func CrtToJson(crt CRT) (result string) {
	b, _ := json.Marshal(crt)
	result = string(b)

	return
}
func crtHeadToJson(head CRTHead) (result string) {
	b, _ := json.Marshal(head)
	result = string(b)

	return
}
