package block

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

type Plugin struct {
	Entrypoint string
	cmd        *exec.Cmd
}

func NewPlugin(entrypoint string) Plugin {
	return Plugin{Entrypoint: entrypoint}
}

func (p *Plugin) Start(configPath string) error {
	p.cmd = exec.Command("sh", "-c", p.Entrypoint)
	p.cmd.Env = os.Environ()
	p.cmd.Env = append(p.cmd.Env, fmt.Sprintf("BLOCK_CONFIG=%s", configPath))
	p.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	p.cmd.Stdout = os.Stdout
	p.cmd.Stderr = os.Stderr
	return p.cmd.Start()
}

func (p *Plugin) Kill() error {
	if p.cmd.Process == nil {
		return nil
	}
	pgid, err := syscall.Getpgid(p.cmd.Process.Pid)
	if err != nil {
		return err
	}
	return syscall.Kill(-pgid, 2)
}

func (p *Plugin) Wait() error {
	return p.cmd.Wait()
}
