package synio

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/chinenual/synergize/logger"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh/terminal"
)

func DiagCOMTST() (err error) {

	var i int
	for i = 0; i < 256; i++ {
		b := byte(i)
		logger.Infof("%d (%02x) ...\n", b, b)

		if err = c.conn.WriteByteWithTimeout(TIMEOUT_MS, b, "write test byte"); err != nil {
			return errors.Wrapf(err, "failed to write byte %02x", b)
		}
		var read_b byte
		if read_b, err = c.conn.ReadByteWithTimeout(TIMEOUT_MS, "read test byte"); err != nil {
			return errors.Wrapf(err, "failed to read byte %02x", b)
		}
		if read_b != b {
			return errors.Errorf("read byte (%02x) does not match what we sent (%02x)", read_b, b)
		}
	}
	// errors will implicitly show  up in the log but we need to explicitly log success
	if synioVerbose {
		logger.Infof("SYNIO: COMTST Success\n")
	}
	return nil
}

func DiagLOOPTST() (err error) {

	if synioVerbose {
		logger.Warnf("LOOPTST causes Synergize to enter an infinte loop supporting the Synergy based test.  You must explicitly kill the Synergize process to stop the test.\n")
	}
	for {

		var b byte
		if b, err = c.conn.ReadByteWithTimeout(1000*60*5, "read test byte"); err != nil {
			return errors.Wrapf(err, "failed to read byte %02x", b)
		}

		logger.Infof("%d (%02x) ...\n", b, b)

		if err = c.conn.WriteByteWithTimeout(TIMEOUT_MS, b, "write test byte"); err != nil {
			return errors.Wrapf(err, "failed to write byte %02x", b)
		}
	}
}

func DiagLINKTST() (err error) {

	state, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalln("setting stdin to raw:", err)
		return
	}
	defer func() {
		fmt.Println("Exiting...\n\r")
		if err := terminal.Restore(int(os.Stdin.Fd()), state); err != nil {
			logger.Warn("failed to restore terminal:", err)
			return
		}
	}()

	in := bufio.NewReader(os.Stdin)
	for {
		r, _, err := in.ReadRune()
		if err != nil {
			logger.Error("stdin:", err)
			break
		}
		if r == '\x03' {
			break
		}
		var b = byte(r)
		if err = c.conn.WriteByteWithTimeout(TIMEOUT_MS, b, "write test byte"); err != nil {
			logger.Infof("failed to write byte %02x - %v\n", b, err)
			break
		}
		if b, err = c.conn.ReadByteWithTimeout(1000*60*5, "read test byte"); err != nil {
			logger.Infof("failed to read byte %02x - %v\n", b, err)
			break
		}
		fmt.Printf(" sent '%q' (0x%02x) ... received 0x%02x (control-C to quit)\n\r", r, r, b)
	}
	return
}
