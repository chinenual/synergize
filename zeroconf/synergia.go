package zeroconf

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/chinenual/synergize/logger"
)

const vstViaSharedFile = true

func getSynergiaState() (result []Service) {
	path, _ := os.UserConfigDir()
	path = path + "/Synergia/state.dat"
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		// not state file, so no VSTs to find
		logger.Infof("No Synergia state file %s.  %v", path, err)
		return
	}
	var f *os.File
	if f, err = os.Open(path); err != nil {
		logger.Errorf("Error opening Synergia state from %s.  %v", path, err)
		return
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		space := strings.Index(line, " ")
		if space >= 0 {
			portstr := line[0:space]
			var port int
			if port, err = strconv.Atoi(portstr); err != nil {
				logger.Errorf("Error parsing Synergia state from %s. Invalid port %s %v", path, portstr, err)
				continue
			}
			name := line[space+1:]
			svc := Service{
				InstanceName: name,
				Port:         uint(port),
				HostName:     "localhost",
			}
			result = append(result, svc)
			logger.Infof("ZEROCONF: found Synergia via state: %s:%d (%s)", svc.HostName, svc.Port, svc.InstanceName)
		}
	}
	if err := scanner.Err(); err != nil {
		logger.Errorf("Error reading Synergia state from %s.  %v", path, err)
		return
	}
	return
}
