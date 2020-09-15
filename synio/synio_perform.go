package synio

import (
	"github.com/pkg/errors"
)

func SelectVoiceMapping(v1, v2, v3, v4 byte) (err error) {
	if err = command(OP_SELECT, "OP_SELECT"); err != nil {
		return errors.Wrapf(err, "failed to OP_SELECT")
	}
	if err = c.conn.WriteByteWithTimeout(RT_TIMEOUT_MS, v1, "voice1"); err != nil {
		return errors.Wrapf(err, "failed to voice1 mapping")
	}
	if err = c.conn.WriteByteWithTimeout(RT_TIMEOUT_MS, v2, "voice2"); err != nil {
		return errors.Wrapf(err, "failed to voice2 mapping")
	}
	if err = c.conn.WriteByteWithTimeout(RT_TIMEOUT_MS, v3, "voice3"); err != nil {
		return errors.Wrapf(err, "failed to voice3 mapping")
	}
	if err = c.conn.WriteByteWithTimeout(RT_TIMEOUT_MS, v4, "voice4"); err != nil {
		return errors.Wrapf(err, "failed to voice4 mapping")
	}
	return
}

// voice      1..4
// key        0..73
// velocity   0..32
func KeyDown(voice, key, velocity byte) (err error) {
	if err = command(OP_ASSIGNED_KEY, "OP_ASSIGNED_KEY"); err != nil {
		return errors.Wrapf(err, "failed to OP_ASSIGNED_KEY")
	}
	//	if err = c.conn.WriteByteWithTimeout(RT_TIMEOUT_MS, OP_KEYDWN, "OP_KEYDWN"); err != nil {
	//		return errors.Wrapf(err, "failed to OP_KEYDWN")
	//	}
	if err = c.conn.WriteByteWithTimeout(RT_TIMEOUT_MS, voice, "voice"); err != nil {
		return errors.Wrapf(err, "failed to send notedown voice")
	}
	if err = c.conn.WriteByteWithTimeout(RT_TIMEOUT_MS, key, "key"); err != nil {
		return errors.Wrapf(err, "failed to send notedown key")
	}
	if err = c.conn.WriteByteWithTimeout(RT_TIMEOUT_MS, velocity, "velocity"); err != nil {
		return errors.Wrapf(err, "failed to send notedown velocity")
	}
	return
}

// Synergy can't turn off voice-specific key - we're in rolling voice assign mode
// key        0..73
// velocity   0..32
func KeyUp(key, velocity byte) (err error) {
	if err = command(OP_KEYUP, "OP_KEYUP"); err != nil {
		return errors.Wrapf(err, "failed to OP_KEYUP")
	}
	if err = c.conn.WriteByteWithTimeout(RT_TIMEOUT_MS, key, "key"); err != nil {
		return errors.Wrapf(err, "failed to send noteup key")
	}
	if err = c.conn.WriteByteWithTimeout(RT_TIMEOUT_MS, velocity, "velocity"); err != nil {
		return errors.Wrapf(err, "failed to send noteup velocity")
	}
	return
}

func Pedal(up bool) (err error) {
	const OPERAND_PEDAL_SUSTAIN = byte(64)
	//const OPERAND_PEDAL_LATCH = byte(65)

	if err = command(OP_POT, "OP_POT"); err != nil {
		return errors.Wrapf(err, "failed to send pedal OP")
	}
	if err = c.conn.WriteByteWithTimeout(RT_TIMEOUT_MS, OPERAND_PEDAL_SUSTAIN, "OPERAND_PEDAL_SUSTAIN"); err != nil {
		return errors.Wrapf(err, "failed to send pedal SUSTAIN operand")
	}
	var value = byte(0) // down
	if up {
		value = 127
	}
	if err = c.conn.WriteByteWithTimeout(RT_TIMEOUT_MS, value, "pedal value"); err != nil {
		return errors.Wrapf(err, "failed to send pedal value")
	}
	return
}
