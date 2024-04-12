package script

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"

	"bharvest.io/scheduled-voter/utils"
)

func Vote(scriptPath string, chain string, propsId string, option string) (*TxResponse, error) {
	utils.Info(fmt.Sprintf("Vote for proposal %s with %s", propsId, option))

	commands := []string{"bash", scriptPath, chain, option, propsId}
	utils.Debug(fmt.Sprintf("Command: %s", commands))

	output, err := exec.Command(commands[0], commands[1:]...).Output()
	if err != nil {
		return nil, err
	}

	cliOutput := TxResponse{}
	err = json.Unmarshal(output, &cliOutput)
	if err != nil {
		return nil, err
	}

	if cliOutput.Code != 0 {
		err := errors.New(fmt.Sprintf("Tx failed: %s", cliOutput.RawLog))
		return nil, err
	}

	utils.Debug(cliOutput)
	return &cliOutput, nil
}
