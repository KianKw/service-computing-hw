package myCobra

import (
	"fmt"
	"reflect"
	"testing"
)

func TestExecute(t *testing.T) {
	cases := []struct {
		cmd *Command
		err error
	}{
		{&Command{}, nil},
		{&Command{
			RunE: func(*Command, []string) error {
				return fmt.Errorf("error")
			},
		}, fmt.Errorf("error")},
		{&Command{Run: func(*Command, []string) {}}, nil},
	}

	for _, c := range cases {
		err := c.cmd.Execute()
		if err == nil && c.err == nil {
			continue
		}
		if err == nil || c.err == nil || err.Error() != c.err.Error() {
			t.Errorf("want error: %v, got error: %v", c.err, err)
		}
	}

}

func TestRunnable(t *testing.T) {
	var rootCmd = &Command{
		Run:  nil,
		RunE: nil,
	}
	var subCmd = &Command{
		Run:  func(cmd *Command, args []string) {},
		RunE: nil,
	}
	var sub2Cmd = &Command{
		Run:  nil,
		RunE: func(cmd *Command, args []string) error { return nil },
	}

	cases := []struct {
		cmd      *Command
		runnable bool
	}{
		{rootCmd, false},
		{subCmd, true},
		{sub2Cmd, true},
	}
	for _, c := range cases {
		if c.cmd.Runnable() != c.runnable {
			t.Errorf("want runnable: %v, got runnable: %v", c.runnable, c.cmd.Runnable())
		}
	}
}

func TestAddCommand(t *testing.T) {
	var help = &Command{
		commands: []*Command{},
	}
	var sub2Cmd = &Command{
		commands: []*Command{},
	}
	var subCmd = &Command{
		commands: []*Command{},
	}
	var rootCmd = &Command{
		commands: []*Command{},
	}

	rootCmd.AddCommand(help)
	rootCmd.AddCommand(subCmd)
	subCmd.AddCommand(sub2Cmd)

	cases := []struct {
		cmd      *Command
		commands []*Command
	}{
		{rootCmd, []*Command{help, subCmd}},
		{subCmd, []*Command{sub2Cmd}},
		{sub2Cmd, []*Command{}},
	}

	for _, c := range cases {
		if !reflect.DeepEqual(c.cmd.commands, c.commands) {
			t.Errorf("want commands list: %v, got  commands list: %v", c.commands, c.cmd.commands)
		}
	}

}

func TestRemoveCommand(t *testing.T) {

	var help = &Command{
		commands: []*Command{},
	}
	var sub2Cmd = &Command{
		commands: []*Command{},
	}
	var subCmd = &Command{
		commands: []*Command{sub2Cmd},
	}
	var rootCmd = &Command{
		commands: []*Command{help, subCmd},
	}

	cases := []struct {
		cmd      *Command
		commands []*Command
	}{
		{rootCmd, []*Command{help}},
		{subCmd, []*Command{sub2Cmd}},
		{sub2Cmd, []*Command{}},
	}

	for _, c := range cases {
		c.cmd.RemoveCommand(subCmd)
		if !reflect.DeepEqual(c.cmd.commands, c.commands) {
			t.Errorf("want commands list: %v, got  commands list: %v", c.commands, c.cmd.commands)
		}
	}

}

func TestName(t *testing.T) {
	var rootCmd = &Command{
		Use: "root",
	}
	var subCmd = &Command{
		Use: "sub -t",
	}
	var sub2Cmd = &Command{
		Use: "sub2 -s -t",
	}

	cases := []struct {
		cmd  *Command
		name string
	}{
		{rootCmd, "root"},
		{subCmd, "sub"},
		{sub2Cmd, "sub2"},
	}

	for _, c := range cases {
		if c.cmd.Name() != c.name {
			t.Errorf("want name: %v, got name: %v", c.name, c.cmd.Name())
		}
	}
}

func TestRoot(t *testing.T) {
	var rootCmd = &Command{
		parent: nil,
	}
	var subCmd = &Command{
		parent: rootCmd,
	}
	var sub2Cmd = &Command{
		parent: subCmd,
	}

	cases := []struct {
		cmd  *Command
		root *Command
	}{
		{rootCmd, rootCmd},
		{subCmd, rootCmd},
		{sub2Cmd, rootCmd},
	}

	for _, c := range cases {
		if c.cmd.Root() != c.root {
			t.Errorf("want root: %v, got command: %v", c.root.Name(), c.cmd.Root().Name())
		}
	}

}


func TestHelpFunc(t *testing.T) {
	var rootCmd = &Command{
		parent:   nil,
		helpFunc: nil,
	}
	var subCmd = &Command{
		parent:   rootCmd,
		helpFunc: nil,
	}

	defaultHelp := rootCmd.HelpFunc()

	if defaultHelp == nil {
		t.Errorf("Default help function missing")
	}

	help := func(*Command, []string) {}
	emptyHelp := reflect.ValueOf(help).Pointer()

	rootCmd.helpFunc = help

	if reflect.ValueOf(rootCmd.HelpFunc()).Pointer() != emptyHelp {
		t.Errorf("Help function missing")
	}

	if reflect.ValueOf(subCmd.HelpFunc()).Pointer() != emptyHelp {
		t.Errorf("Parent help function missing")
	}

}

func TestSetHelpFunc(t *testing.T) {
	var help = func(cmd *Command, a []string) {}

	var rootCmd = &Command{
		helpFunc: nil,
	}
	var subCmd = &Command{
		helpFunc: nil,
	}
	var sub2Cmd = &Command{
		helpFunc: nil,
	}

	cases := []struct {
		cmd      *Command
		helpFunc func(*Command, []string)
	}{
		{rootCmd, help},
		{subCmd, help},
		{sub2Cmd, help},
	}

	for _, c := range cases {
		c.cmd.SetHelpFunc(c.helpFunc)
		got := reflect.ValueOf(c.cmd.helpFunc).Pointer()
		want := reflect.ValueOf(c.helpFunc).Pointer()
		if got != want {
			t.Errorf("want help function: %v, got help function: %v", want, got)
		}
	}

}

func TestInitDefaultHelpCmd(t *testing.T) {
	var rootCmd = &Command{
		commands:    []*Command{},
		helpCommand: nil,
	}
	var subCmd = &Command{
		commands:    []*Command{},
		helpCommand: nil,
	}

	rootCmd.InitDefaultHelpCmd()

	if rootCmd.helpCommand != nil {
		t.Errorf("Redundant help command")
	}

	rootCmd.AddCommand(subCmd)

	rootCmd.InitDefaultHelpCmd()

	if rootCmd.helpCommand == nil {
		t.Errorf("Default help command missing")
	}

	var emptyHelp = &Command{}

	rootCmd.helpCommand = emptyHelp

	rootCmd.InitDefaultHelpCmd()

	if rootCmd.helpCommand != emptyHelp {
		t.Errorf("User help command missing")
	}

}

func TestSetHelpCommand(t *testing.T) {

	var help = &Command{
		Use: "help",
	}

	var rootCmd = &Command{
		helpCommand: nil,
	}
	var subCmd = &Command{
		helpCommand: nil,
	}
	var sub2Cmd = &Command{
		helpCommand: nil,
	}

	cases := []struct {
		cmd         *Command
		helpCommand *Command
	}{
		{rootCmd, help},
		{subCmd, help},
		{sub2Cmd, help},
	}

	for _, c := range cases {
		c.cmd.SetHelpCommand(c.helpCommand)
		if c.cmd.helpCommand != c.helpCommand {
			t.Errorf("want help command: %v, got help command: %v", c.cmd.helpCommand.Name(), c.helpCommand.Name())
		}
	}

}

func TestFind(t *testing.T) {

	var sub2Cmd = &Command{
		Use:      "sub2",
		commands: []*Command{},
	}
	var subCmd = &Command{
		Use:      "sub",
		commands: []*Command{sub2Cmd},
	}
	var rootCmd = &Command{
		Use:      "root",
		commands: []*Command{subCmd},
	}

	cases := []struct {
		cmd          *Command
		args         []string
		commandFound *Command
		flags        []string
		err          error
	}{
		{rootCmd, []string{"sub", "sub2"}, sub2Cmd, []string{}, nil},
		{subCmd, []string{"sub2"}, sub2Cmd, []string{}, nil},
		{sub2Cmd, []string{}, sub2Cmd, []string{}, nil},
	}

	for _, c := range cases {
		targetCmd, flags, err := c.cmd.Find(c.args)
		if targetCmd != c.commandFound {
			t.Errorf("want target command: %v, got target command: %v", c.commandFound.Name(), targetCmd.Name())
		}
		if !reflect.DeepEqual(flags, c.flags) {
			t.Errorf("want target flags: %v, got target flags: %v", c.flags, flags)
		}
		if err != c.err {
			t.Errorf("want error: %v, got error: %v", c.err, err)
		}
	}

}
