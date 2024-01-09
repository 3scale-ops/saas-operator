package ssh

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"golang.org/x/crypto/ssh"
)

type RemoteExecutor struct {
	Host       string
	User       string
	Port       uint32
	PrivateKey string
	Logger     logr.Logger
	CmdTimeout time.Duration
	Commands   []Runnable
}

func (re *RemoteExecutor) Run() error {

	key, err := ssh.ParsePrivateKey([]byte(re.PrivateKey))
	if err != nil {
		return err
	}
	// Authentication
	config := &ssh.ClientConfig{
		User:            re.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		Timeout: re.CmdTimeout,
	}
	client, err := ssh.Dial("tcp", net.JoinHostPort(re.Host, strconv.Itoa(int(re.Port))), config)
	if err != nil {
		return err
	}
	defer client.Close()

	for _, cmd := range re.Commands {
		re.Logger.V(1).Info(cmd.Info())
		output, err := cmd.Run(client)
		if output != "" {
			re.Logger.V(1).Info(fmt.Sprintf("remote ssh command output: %s", output))
		}
		if err != nil {
			re.Logger.V(1).Info(fmt.Sprintf("remote ssh command error: %s", err.Error()))
			return fmt.Errorf("remote ssh command failed: %w (%s)", err, output)
		}
	}

	// be silent on success
	return nil

}

type Runnable interface {
	Run(*ssh.Client) (string, error)
	Info() string
	WithSudo(bool) Runnable
}

type Command struct {
	value     string
	sensitive []string
	sudo      bool
}

var _ Runnable = &Command{}

func NewCommand(value string, sensitive ...string) *Command {
	return &Command{value: value, sensitive: sensitive}
}

func (c *Command) WithSudo(sudo bool) Runnable {
	c.sudo = sudo
	return c
}

func (c *Command) resolveValue() string {
	if c.sudo {
		return "sudo " + c.value
	}
	return c.value
}

func (c *Command) Info() string {
	return fmt.Sprintf("run command: %s", hideSensitive(c.resolveValue(), c.sensitive...))
}

func (c *Command) Run(client *ssh.Client) (string, error) {
	// Create a session. It is one session per command.
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput(c.resolveValue())
	if err != nil {
		return string(output), err
	}

	return "", nil
}

type Script struct {
	value       []byte
	interpreter string
	sensitive   []string
	sudo        bool
}

var _ Runnable = &Script{}

func NewScript(interpreter string, script string, sensitive ...string) *Script {
	return &Script{
		value:       []byte(script),
		interpreter: interpreter,
		sensitive:   sensitive,
	}
}

func (s *Script) WithSudo(sudo bool) Runnable {
	s.sudo = sudo
	return s
}

func (s *Script) Info() string {
	return fmt.Sprintf("run script with: '%s' \n'%s'",
		hideSensitive(s.resolveInterpreter(), s.sensitive...),
		hideSensitive(string(s.value), s.sensitive...),
	)
}

func (s *Script) resolveInterpreter() string {
	if s.sudo {
		return "sudo " + s.interpreter
	}
	return s.interpreter
}

func (s *Script) Run(client *ssh.Client) (string, error) {
	// Create a session. It is one session per command.
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		return "", err
	}
	defer stdin.Close()

	type response struct {
		output []byte
		err    error
	}

	chRsp := make(chan response)

	go func() {
		output, err := session.CombinedOutput(s.resolveInterpreter())
		if err != nil {
			chRsp <- response{output: output, err: err}
			return
		}
		chRsp <- response{output: []byte(""), err: nil}
	}()

	_, err = stdin.Write(s.value)
	if err != nil {
		return "", err
	}
	stdin.Close()

	rsp := <-chRsp

	return string(rsp.output), rsp.err
}

func hideSensitive(msg string, hide ...string) string {
	for _, ss := range hide {
		msg = strings.ReplaceAll(msg, ss, "*****")
	}
	return msg
}
