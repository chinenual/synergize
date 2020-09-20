package main

import (
	"fmt"

	"github.com/chinenual/synergize/io"
	"github.com/chinenual/synergize/logger"
	"github.com/chinenual/synergize/osc"
	"github.com/chinenual/synergize/synio"
	"github.com/chinenual/synergize/zeroconf"

	"github.com/pkg/errors"
)

// Synergy and Control Surface connection can be configured from the command line,
// hardcoded in the preferences file, or auto-configured via zeroconf.
// Orchestrate the initializations from here...
//
// Connection/Zeroconf Lifecycle:
//  OSC server (re)started when:
//       voicing mode starts
//  OSC client (re)started when:
//       voicing mode starts
//
//  VST client started when:
//       first time IO requiring synergy connection
//          user explictly connects
//          load CRT for editing
//          load CRT
//          load SYN
//          save SYN
//          disable VRAM
//          run COMTest
//          toggle voicing mode
//
//  zeroconf browses when:
//       at program startup
//       user explicitly asks for a re-scan
//
//  zeroconf publishes OSC server address when:
//       at program startup
//       whenever server restarted
//
// ----------
// Connections are initiated from Javascript due to inability to do a simple synchronous
// dialog initiated from Go. If Go could call javascript and prompt with "choose one of
// these devcies", all would be much more straightforward.  But since Javascript requires
// async callbacks from any sort of modal dialog, we have to structure it thusly:
//
//    javascript asks go what devices are available
//    if zeroconf is not enabled for the device, go returns "already configured" (in which case javascript stops)
//    else go can respond with "device is already configured" if we've already been through a zeroconf selection cycle
//       (in which case javascript stops)
//       or go returns a list of choices
//    javascript modal dialog shows the choices. user can choose "rescan" in which case we do
//       another zeroconf browse  - javascript starts over and reinitates the cycle
//       or, user can cancel - (in which case javascript stops)
//       or user can choose one of the choices.  It then passes that back to go to "connect"
//       to the chosen device

var firmwareVersion string

func DisconnectSynergy() (err error) {
	firmwareVersion = ""
	if err = synio.Close(); err != nil {
		return
	}
	return
}

func DisconnectControlSurface() (err error) {
	if err = osc.Quit(); err != nil {
		return
	}
	return
}

func ConnectSynergy(zeroconfConfig zeroconf.Service) (err error) {
	if zeroconfConfig.InstanceName == "serial-port" {
		logger.Infof("ZEROCONF: using Synergy preferences config %s at %d\n", prefsUserPreferences.SerialPort, prefsUserPreferences.SerialBaud)
		if err = synio.SetSynergySerialPort(prefsUserPreferences.SerialPort, prefsUserPreferences.SerialBaud,
			true, *serialVerboseFlag, *mockSynio); err != nil {
			return
		}
	} else {
		logger.Infof("ZEROCONF: using Synergy zeroconf config %s (%s:%d)\n", zeroconfConfig.InstanceName, zeroconfConfig.HostName, zeroconfConfig.Port)
		if err = synio.SetSynergyVst(zeroconfConfig.InstanceName, zeroconfConfig.HostName, zeroconfConfig.Port,
			true, *serialVerboseFlag, *mockSynio); err != nil {
			return
		}
	}
	return
}

func GetSynergyConfig() (hasDevice bool, alreadyConfigured bool, name string, choices *[]zeroconf.Service, err error) {
	if !io.SynergyConfigured() {
		if prefsUserPreferences.VstAutoConfig {
			vstServices := zeroconf.GetVstServices()
			if false && len(vstServices) == 1 && prefsUserPreferences.SerialPort == "" {
				logger.Infof("ZEROCONF: auto config VST: %#v\n", vstServices[0])
				firmwareVersion = ""
				if err = synio.SetSynergyVst(vstServices[0].InstanceName, vstServices[0].HostName, vstServices[0].Port,
					true, *serialVerboseFlag, *mockSynio); err != nil {
					return
				}
			} else {
				logger.Infof("ZEROCONF: zero or more than one VST: %#v\n", vstServices)
				if prefsUserPreferences.SerialPort != "" {
					var pseudoEntry zeroconf.Service
					pseudoEntry.InstanceName = "serial-port"
					list := append([]zeroconf.Service{pseudoEntry}, vstServices...)
					choices = &list
				} else {
					choices = &vstServices
				}
			}
		} else {
			firmwareVersion = ""
			if *vst != 0 {
				// VST command line option
				logger.Infof("ZEROCONF: VST zeroconf disabled - using -VST command line config %d\n", *vst)
				if err = synio.SetSynergyVst(fmt.Sprintf("VST localhost:%d", *vst), "localhost", *vst,
					true, *serialVerboseFlag, *mockSynio); err != nil {
					return
				}
			} else {
				logger.Infof("ZEROCONF: VST zeroconf disabled - using preferences config %s at %d\n", prefsUserPreferences.SerialPort, prefsUserPreferences.SerialBaud)
				if err = synio.SetSynergySerialPort(prefsUserPreferences.SerialPort, prefsUserPreferences.SerialBaud,
					true, *serialVerboseFlag, *mockSynio); err != nil {
					return
				}
			}
		}
	}
	hasDevice = true // we assume user has a Synergy or VST
	name = io.SynergyName()
	alreadyConfigured = io.SynergyConfigured()
	return
}

func GetControlSurfaceConfig() (hasDevice bool, alreadyConfigured bool, name string, choices *[]zeroconf.Service, err error) {
	if !osc.ControlSurfaceConfigured() && prefsUserPreferences.UseOsc {
		if prefsUserPreferences.OscAutoConfig {
			oscServices := zeroconf.GetOscServices()

			if false && len(oscServices) == 1 {
				logger.Infof("ZEROCONF: auto config Control Surface: %#v\n", oscServices[0])
				osc.SetControlSurface(oscServices[0].InstanceName, oscServices[0].HostName, oscServices[0].Port)
			} else {
				logger.Infof("ZEROCONF: zero or more than one Control Surface: %#v\n", oscServices)
				choices = &oscServices
			}
		} else {
			logger.Infof("ZEROCONF: Control Surface zeroconf disabled - using preferences config %s:%d\n", prefsUserPreferences.OscCSurfaceAddress, prefsUserPreferences.OscCSurfacePort)
			osc.SetControlSurface("", prefsUserPreferences.OscCSurfaceAddress, prefsUserPreferences.OscCSurfacePort)
		}
	}
	hasDevice = prefsUserPreferences.UseOsc
	alreadyConfigured = osc.ControlSurfaceConfigured()
	name = osc.ControlSurfaceName()
	return
}

func GetFirmwareVersion() (id string, err error) {
	if firmwareVersion == "" {
		if !io.SynergyConfigured() {
			err = errors.New("Not connected to a Synergy; can't check firmware version")
			return
		} else {
			var bytes [2]byte
			bytes, err = synio.GetID()
			if err != nil {
				err = errors.Wrap(err, "Cannot get firmware version")
				logger.Errorf(err.Error())
				return
			}
			firmwareVersion = fmt.Sprintf("%d.%d", bytes[0], bytes[1])
			logger.Infof("Connected to Synergy, firmware version: %s\n", firmwareVersion)
		}
	}
	return
}

func ConnectToSynergy(choice *zeroconf.Service) (err error) {
	if !io.SynergyConfigured() {
		firmwareVersion = ""
		if choice == nil || choice.InstanceName == "serial-port" {
			if err = synio.SetSynergySerialPort(prefsUserPreferences.SerialPort,
				prefsUserPreferences.SerialBaud, true, *serialVerboseFlag, *mockSynio); err != nil {
				err = errors.Wrapf(err, "Cannot connect to synergy on port %s at %d baud\n",
					prefsUserPreferences.SerialPort,
					prefsUserPreferences.SerialBaud)
				logger.Errorf(err.Error())
				return
			}
		} else {
			// VST
			if err = synio.SetSynergyVst(choice.InstanceName, choice.HostName,
				choice.Port, true, *serialVerboseFlag, *mockSynio); err != nil {
				err = errors.Wrapf(err, "Cannot connect to synergy VST %s at %s:%d\n",
					choice.InstanceName, choice.HostName,
					choice.Port)
				logger.Errorf(err.Error())
				return
			}
		}
		logger.Infof("Connected to Synergy %s\n", io.SynergyName())
	}
	return
}
