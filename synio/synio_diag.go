package synio

import (
	"log"
	"github.com/pkg/errors"
)

func DiagCOMTST() (err error) {

	var i int
	for i = 0; i < 256; i++ {
		b := byte(i)
		log.Printf("%d (%02x) ...\n", b, b)

		if err = serialWriteByte(TIMEOUT_MS, b, "write test byte"); err != nil {
			return errors.Wrapf(err, "failed to write byte %02x", b)
		}
		var read_b byte
		if read_b, err = serialReadByte(TIMEOUT_MS, "read test byte"); err != nil {
			return errors.Wrapf(err, "failed to read byte %02x", b)
		}
		if read_b != b {
			return errors.Errorf("read byte (%02x) does not match what we sent (%02x)", read_b, b)
		}
	}
	// errors will implicitly show  up in the log but we need to explicitly log success
	if synioVerbose {
		log.Printf("COMTST Success\n")
	}
	return nil
}

func DiagLOOPTST() (err error) {

	if synioVerbose {
		log.Printf("WARNING: LOOPTST causes Synergize to enter an infinte loop supporting the Synergy based test.  You must explicitly kill the Synergize process to stop the test.\n")
	}
	for {

		var b byte
		if b, err = serialReadByte(1000*60*5, "read test byte"); err != nil {
			return errors.Wrapf(err, "failed to read byte %02x", b)
		}

		log.Printf("%d (%02x) ...\n", b, b)

		if err = serialWriteByte(TIMEOUT_MS, b, "write test byte"); err != nil {
			return errors.Wrapf(err, "failed to write byte %02x", b)
		}
	}
}
