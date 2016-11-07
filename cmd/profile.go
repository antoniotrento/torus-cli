package cmd

import (
	"context"
	"fmt"
	"sync"

	"github.com/manifoldco/torus-cli/api"
	"github.com/manifoldco/torus-cli/apitypes"
	"github.com/manifoldco/torus-cli/config"
	"github.com/manifoldco/torus-cli/errs"

	"github.com/urfave/cli"
)

func init() {
	profile := cli.Command{
		Name:     "profile",
		Usage:    "Manage your Torus account",
		Category: "ACCOUNT",
		Action:   chain(ensureDaemon, signupCmd),
		Subcommands: []cli.Command{
			{
				Name:  "update",
				Usage: "Update your profile",
				Action: chain(
					ensureDaemon, ensureSession, loadDirPrefs, loadPrefDefaults,
					setUserEnv, profileEdit,
				),
			},
		},
	}
	Cmds = append(Cmds, profile)
}

// profileEdit is used to update name and email for an account
func profileEdit(ctx *cli.Context) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg)
	c := context.Background()

	session, err := client.Session.Who(c)
	if err != nil {
		return errs.NewErrorExitError("Error fetching user details", err)
	}
	if session.Type() == apitypes.MachineSession {
		return errs.NewExitError("Machines do not have profiles")
	}

	ogName := session.Name()
	name, err := FullNamePrompt(ogName)
	if err != nil {
		return err
	}

	ogEmail := session.Email()
	email, err := EmailPrompt(ogEmail)
	if err != nil {
		return err
	}

	warning := "\nYou are about to update your profile to the values above."
	if email != ogEmail {
		warning = "\nYou will be required to re-verify your email address before taking any further actions within Torus."
	}

	if ogEmail == email && ogName == name {
		fmt.Println("\nNo changes made :)")
		return nil
	}

	err = ConfirmDialogue(ctx, nil, &warning)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(2)

	var eErr error
	go func() {
		if ogEmail != email {
			_, err := client.Users.UpdateEmail(c, email)
			eErr = err
		}
		wg.Done()
	}()

	var nErr error
	go func() {
		if ogName != name {
			_, err := client.Users.UpdateName(c, name)
			nErr = err
		}
		wg.Done()
	}()

	wg.Wait()
	msg := "Failed to update profile"
	if nErr != nil {
		return errs.NewErrorExitError(msg, nErr)
	}
	if eErr != nil {
		return errs.NewErrorExitError(msg, eErr)
	}

	err = client.Session.Refresh(c)
	if err != nil {
		return errs.NewExitError("Failed to refresh profile, please log out and log back in")
	}

	return nil
}
