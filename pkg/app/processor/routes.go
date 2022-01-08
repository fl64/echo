package processor

import (
	"echo-http/pkg/app/models"
	"encoding/json"
	"os/exec"
)

func (p *Processor) GetRoutes() (*[]models.Route, error) {
	routes := make([]models.Route, 1)
	out, err := exec.Command("ip", "-j", "r").Output()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(out, &routes)
	if err != nil {
		return nil, err
	}
	return &routes, nil
}
