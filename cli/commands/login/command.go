package login

//Importing packages defined elsewhere in the application with the exception of urfave/cli; an external package
import (
	"github.com/taubyte/tau-cli/cli/common/options"
	"github.com/taubyte/tau-cli/flags"
	loginFlags "github.com/taubyte/tau-cli/flags/login"
	"github.com/taubyte/tau-cli/i18n"
	loginI18n "github.com/taubyte/tau-cli/i18n/login"
	loginLib "github.com/taubyte/tau-cli/lib/login"
	"github.com/taubyte/tau-cli/prompts"
	loginPrompts "github.com/taubyte/tau-cli/prompts/login"
	slices "github.com/taubyte/utils/slices/string"
	"github.com/urfave/cli/v2" //urfave/cli package for building the command line
)

//Defining a login command for the CLI
var Command = &cli.Command{
	Name: "login", //Name of the command
	Flags: flags.Combine( //Defining the command flags 
		flags.Name, //Name of the profile 
		loginFlags.Token, //Login token flag for authentication
		loginFlags.Provider, //Provider of the authentication service --> Github login
		loginFlags.New, //Flag for new profile creation
		loginFlags.SetDefault, //Sets profile to default selected on login
	),
	ArgsUsage: i18n.ArgsUsageName, //Setting the usage message for command args from internationalization package 
	Action:    Run, //Executes run function after all subcommands complete
	Before:    options.SetNameAsArgs0, //Executing SetNameAsArgs0 before subcommands
}

func Run(ctx *cli.Context) error {
	_default, options, err := loginLib.GetProfiles() //Grabbing profile variables where options is the list of profiles 
	//Checking for errors
	if err != nil { 
		return loginI18n.GetProfilesFailed(err) //Throwing error with the help of the GetProfilesFailed function that does the formatting
	}

	// New: if --new or no selectable profiles
	if ctx.Bool(loginFlags.New.Name) || len(options) == 0 { 
		//Checking that the new flag has a set name OR if there are no profiles 
		//initiates a function that creates a new profile for authentication
		return New(ctx, options)
	}

	// Selection
	var name string
	if ctx.IsSet(flags.Name.Name) { //Checking that the context has been set in the CLI
		name = ctx.String(flags.Name.Name) //setting the name var to the flags value

		if !slices.Contains(options, name) { //Checking if the list of profiles has no matches with the name var
			return loginI18n.DoesNotExistIn(name, options) //Throwing error formatted with DoesNotExistIn function
		}
	} else { //If there is a match
		/*
		the SelectInterface function below handles the logic 
		for creating a selection prompt based on the options 
		profiles list, and default profile name. If there
		is no default name it sets the default to names[0]. 
		It utilizes the survey package to help with creating 
		the selection prompt.		
		*/
		name, err = prompts.SelectInterface(options, loginPrompts.SelectAProfile, _default)
		if err != nil {
			return err //returns error if it exists
		}
	}
	/*
	Executes the select function in commands/login/select with the given parameters. 
	select function then runs Select located in /lib/login/select which handles updating
	the default profile and returns any relevant errors. The function then grabs the
	network information from selected profile by using the NetworkType and Network variables
	from the profile in context and sets them with the SetSelectedNetwork function and 
	SetNetworkUrl function.
	*/
	return Select(ctx, name, ctx.Bool(loginFlags.SetDefault.Name))
}
