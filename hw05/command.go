package myCobra

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"
)

type Command struct {
	Use string

	Short string

	Long string

	Run func(cmd *Command, args []string)

	RunE func(cmd *Command, args []string) error

	helpFunc func(*Command, []string)

	helpCommand *Command

	args []string

	parent *Command

	commands []*Command

	pflags *flag.FlagSet
}

func (cmd *Command) Execute() error {
	cmd.InitDefaultHelpCmd()
	args := cmd.args
	if cmd.args == nil && filepath.Base(os.Args[0]) != "cobra.test" {
		args = os.Args[1:]
	}
	targetCmd, flags, err := cmd.Find(args)
	if err != nil {
		return err
	}
	err = targetCmd.execute(flags)
	if err != nil {
		if err == flag.ErrHelp {
			targetCmd.HelpFunc()(targetCmd, args)
			return nil
		}
	}
	return err
}

func (cmd *Command) execute(a []string) (err error) {
	for _, v := range a {
		if v == "-h" || v == "--help" {
			fmt.Println("Congratulation!")
			cmd.Print_help()
			return
		}
	}
	if cmd == nil {
		return fmt.Errorf("Called Execute() on a nil Command")
	}
	if !cmd.Runnable() {
		return flag.ErrHelp
	}
	if cmd.RunE != nil {
		err := cmd.RunE(cmd, a)
		return err
	}
	cmd.Run(cmd, a)
	return nil
}

func (cmd *Command) Runnable() bool {
	return cmd.Run != nil || cmd.RunE != nil
}

func ParseArgs(c *Command, args []string) {
	if len(args) < 1 {
		return
	}
	for _, v := range c.commands {
		if v.Use == args[0] {
			c.args = args[:1]
            c.AddCommand(v)
			v.parent = c
			ParseArgs(v, args[1:])
			return
		}
	}
	c.args = args
	c.PersistentFlags().Parse(c.args)
}

func (cmd *Command) Find(args []string) (*Command, []string, error) {
	var innerfind func(*Command, []string) (*Command, []string)
	innerfind = func(cmd *Command, innerArgs []string) (*Command, []string) {
		argsWOflags := innerArgs
		if len(argsWOflags) == 0 {
			return cmd, innerArgs
		}
		nextSubCmd := argsWOflags[0]
		targetCmd := cmd.findNext(nextSubCmd)
		if targetCmd != nil {
			return innerfind(targetCmd, argsMinusFirstX(innerArgs, nextSubCmd))
		}
		return cmd, innerArgs
	}
	commandFound, flags := innerfind(cmd, args)
	return commandFound, flags, nil
}


func (cmd *Command) findNext(next string) *Command {
	for _, cmd := range cmd.commands {
		if cmd.Name() == next {
			return cmd
		}
	}
	return nil
}

func argsMinusFirstX(args []string, x string) []string {
	for i, y := range args {
		if x == y {
			ret := []string{}
			ret = append(ret, args[:i]...)
			ret = append(ret, args[i+1:]...)
			return ret
		}
	}
	return args
}

func (cmd *Command) AddCommand(subCmds ...*Command) {
	for i, subCmd := range subCmds {
		if subCmds[i] == cmd {
			panic("Command can't be a child of itself")
		}
		subCmds[i].parent = cmd
		cmd.commands = append(cmd.commands, subCmd)
	}
}

func (cmd *Command) RemoveCommand(rmCmds ...*Command) {
	commands := []*Command{}
main:
	for _, command := range cmd.commands {
		for _, rmCmd := range rmCmds {
			if command == rmCmd {
				command.parent = nil
				continue main
			}
		}
		commands = append(commands, command)
	}
	cmd.commands = commands
}

func (c *Command) Print_help() {
	fmt.Printf("%s\n\n", c.Long)
	fmt.Printf("Usage:\n")
	fmt.Printf("\t%s [flags]\n", c.Name())
	if (len(c.commands) > 0) {
		fmt.Printf("\t%s [command]\n\n", c.Name())
		fmt.Printf("Available Commands:\n")
		for _, v := range c.commands {
			fmt.Printf("\t%-10s%s\n", v.Name(), v.Short)
		}
	}

	fmt.Printf("\nFlags:\n")

	c.PersistentFlags().VisitAll(func (flag *flag.Flag) {
		fmt.Printf("\t-%1s, --%-6s %-12s%s (default \"%s\")\n", flag.Shorthand, flag.Name,  flag.Value.Type(), flag.Usage, flag.DefValue)
	})
	fmt.Printf("\t-%1s, --%-19s%s%s\n", "h", "help", "help for ", c.Name())
	fmt.Println()
	if len(c.commands) > 0 {
		fmt.Printf("Use \"%s [command] --help\" for more information about a command.\n", c.Name())
	}
	fmt.Println()
}

func (c *Command) PersistentFlags() *flag.FlagSet {
	if c.pflags == nil {
		c.pflags = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	}
	return c.pflags
}

func (cmd *Command) Name() string {
	name := cmd.Use
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func (cmd *Command) Root() *Command {
	if cmd.parent != nil {
		return cmd.parent.Root()
	}
	return cmd
}

func (cmd *Command) HelpFunc() func(*Command, []string) {
	if cmd.helpFunc != nil {
		return cmd.helpFunc
	}
	if cmd.parent != nil {
		return cmd.parent.HelpFunc()
	}
	return func(cmd *Command, a []string) {
		if cmd.Long != "" {
			fmt.Println(cmd.Long)
		}
		fmt.Println(" ")
	}
}

func (cmd *Command) SetHelpFunc(f func(*Command, []string)) {
	cmd.helpFunc = f
}

func (cmd *Command) InitDefaultHelpCmd() {
	if len(cmd.commands) == 0 {
		return
	}

	if cmd.helpCommand == nil {
		cmd.helpCommand = &Command{
			Use:   "help [command]",
			Short: "Help about any command",
			Long: `Help provides help for any command in the application.
Simply type ` + cmd.Name() + ` help [path to command] for full details.`,
			Run: func(c *Command, args []string) {
				targetCmd, _, e := c.Root().Find(args)
				if targetCmd == nil || e != nil {
					fmt.Printf("Unknown help topic %#q\n", args)
					c.HelpFunc()
				} else {
					targetCmd.HelpFunc()(targetCmd, []string{})
				}
			},
		}
	}
	cmd.RemoveCommand(cmd.helpCommand)
	cmd.AddCommand(cmd.helpCommand)
}

func (cmd *Command) SetHelpCommand(c *Command) {
	cmd.helpCommand = c
}
