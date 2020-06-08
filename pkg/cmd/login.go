package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

// loginCmd type for init command
type loginCmd struct {
	security.LoginManager
	formula.Loader
	prompt.InputText
	prompt.InputPassword
}

// NewLoginCmd creates new cmd instance
func NewLoginCmd(
	t prompt.InputText,
	p prompt.InputPassword,
	lm security.LoginManager,
	fm formula.Loader) *cobra.Command {
	l := loginCmd{
		LoginManager:  lm,
		Loader:        fm,
		InputText:     t,
		InputPassword: p,
	}
	return &cobra.Command{
		Use:   "login",
		Short: "User login",
		Long:  "Authenticates and creates a session for the user of the organization",
		RunE:  RunFuncE(l.runStdin(), l.runPrompt()),
	}
}

func (l loginCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		u, err := l.Text("Username: ", true)
		if err != nil {
			return err
		}
		p, err := l.Password("Password: ")
		if err != nil {
			return err
		}
		us := security.User{
			Username: u,
			Password: p,
		}
		if err = l.Login(us); err != nil {
			return err
		}
		fmt.Println("Login successfully!")
		return err
	}
}

func (l loginCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {

		u := security.User{}

		err := stdin.ReadJson(os.Stdin, &u)
		if err != nil {
			fmt.Println("The STDIN inputs weren't informed correctly. Check the JSON used to execute the command.")
			return err
		}

		if err = l.Login(u); err != nil {
			return err
		}
		fmt.Println("Login successfully!")
		return err
	}
}