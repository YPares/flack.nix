package cli

import (
	"fmt"
	"os"
	"slices"

	"github.com/YPares/flack/internal/flake"
	"github.com/YPares/flack/internal/nix"
	"github.com/YPares/flack/internal/profile"
	"github.com/YPares/flack/internal/tui"
	"github.com/spf13/cobra"
)

var (
	profilePath  string
	flakePath    string
	priorityStep int
	jsonOutput   bool
	refreshFlag  bool
	allFlag      bool
	fromFlag     string
)

var rootCmd = &cobra.Command{
	Use:   "flack",
	Short: "Manage your nix profile as a stack and update flake inputs",
	Long:  "flack wraps nix profile and nix flake to provide stack-based package management and interactive input selection.",
	RunE:  runShow,
}

func init() {
	p := defaultProfile()
	if v := os.Getenv("FLACK_PROFILE"); v != "" {
		p = v
	}
	rootCmd.PersistentFlags().StringVarP(&profilePath, "profile", "p", p, "nix profile path ($FLACK_PROFILE)")
	rootCmd.PersistentFlags().IntVar(&priorityStep, "priority-step", 10, "priority step between stack layers ($FLACK_PRIORITY_STEP)")
	rootCmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "j", false, "output as JSON")

	if v := os.Getenv("FLACK_PRIORITY_STEP"); v != "" {
		if n, err := fmt.Sscanf(v, "%d", &priorityStep); n == 1 && err == nil {
			rootCmd.PersistentFlags().Lookup("priority-step").Value.Set(v)
		}
	}

	pushCmd.Flags().BoolVar(&refreshFlag, "refresh", false, "pass --refresh to nix")
	pushCmd.Flags().StringVarP(&fromFlag, "from", "f", "", "select interactively from a flake's packages")
	rootCmd.AddCommand(pushCmd)

	rootCmd.AddCommand(popCmd)

	upgradeCmd.Flags().BoolVar(&refreshFlag, "refresh", false, "pass --refresh to nix profile upgrade")
	upgradeCmd.Flags().BoolVar(&allFlag, "all", false, "upgrade all packages without selection")
	rootCmd.AddCommand(upgradeCmd)

	f := "."
	if v := os.Getenv("FLACK_FLAKE"); v != "" {
		f = v
	}
	updateCmd.Flags().StringVarP(&flakePath, "flake", "f", f, "flake path ($FLACK_FLAKE)")
	updateCmd.Flags().BoolVar(&refreshFlag, "refresh", false, "pass --refresh to nix flake update")
	updateCmd.Flags().BoolVar(&allFlag, "all", false, "update all inputs without selection")
	rootCmd.AddCommand(updateCmd)

	inputsCmd.Flags().StringVarP(&flakePath, "flake", "f", f, "flake path ($FLACK_FLAKE)")
	rootCmd.AddCommand(inputsCmd)
}

func defaultProfile() string {
	home, _ := os.UserHomeDir()
	if home != "" {
		return home + "/.nix-profile"
	}
	return "/nix/var/nix/profiles/default"
}

var pushCmd = &cobra.Command{
	Use:   "push [<flake-ref>]",
	Short: "Install a package with highest priority in the stack",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if fromFlag != "" {
			if len(args) > 0 {
				return fmt.Errorf("cannot specify both --from and a flake-ref argument")
			}
			names, err := nix.FlakePackageNames(fromFlag)
			if err != nil {
				return err
			}
			if len(names) == 0 {
				return fmt.Errorf("no packages found for current system in flake %s", fromFlag)
			}
			slices.Sort(names)
			items := make([]tui.Selectable, len(names))
			for i, name := range names {
				items[i] = tui.Selectable{
					ID:    name,
					Label: name,
				}
			}
			selected, err := tui.SingleSelect("Select a package to push", items)
			if err != nil {
				return err
			}
			return profile.Push(profilePath, fromFlag+"#"+selected, priorityStep)
		}
		if len(args) == 0 {
			return fmt.Errorf("requires a flake-ref argument or --from flag")
		}
		return profile.Push(profilePath, args[0], priorityStep)
	},
}

var popCmd = &cobra.Command{
	Use:   "pop",
	Short: "Remove the highest-priority package from the stack",
	RunE: func(cmd *cobra.Command, args []string) error {
		return profile.Pop(profilePath)
	},
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Interactively select packages to upgrade",
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := profile.List(profilePath)
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			return fmt.Errorf("profile is empty")
		}
		if allFlag {
			var names []string
			for _, e := range entries {
				names = append(names, e.Name)
			}
			return profile.Upgrade(profilePath, names, refreshFlag)
		}
		items := make([]tui.Selectable, len(entries))
		for i, e := range entries {
			items[i] = tui.Selectable{
				ID:    e.Name,
				Label: fmt.Sprintf("%s (%s)", e.Name, e.OriginalURL),
			}
		}
		selected, err := tui.MultiSelect("Select packages to upgrade", items)
		if err != nil {
			return err
		}
		if len(selected) == 0 {
			return fmt.Errorf("no package selected")
		}
		return profile.Upgrade(profilePath, selected, refreshFlag)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Interactively select flake inputs to update",
	RunE: func(cmd *cobra.Command, args []string) error {
		inputs, err := flake.Inputs(flakePath)
		if err != nil {
			return err
		}
		if len(inputs) == 0 {
			return fmt.Errorf("no inputs found")
		}
		if allFlag {
			var names []string
			for _, inp := range inputs {
				names = append(names, inp.Name)
			}
			return flake.Update(flakePath, names, refreshFlag)
		}
		items := make([]tui.Selectable, len(inputs))
		for i, inp := range inputs {
			items[i] = tui.Selectable{
				ID:    inp.Name,
				Label: fmt.Sprintf("%s (%s)", inp.Name, inp.Original),
			}
		}
		selected, err := tui.MultiSelect("Select inputs to update", items)
		if err != nil {
			return err
		}
		if len(selected) == 0 {
			return fmt.Errorf("no input selected")
		}
		return flake.Update(flakePath, selected, refreshFlag)
	},
}

var inputsCmd = &cobra.Command{
	Use:   "inputs",
	Short: "Show flake lockfile inputs",
	RunE: func(cmd *cobra.Command, args []string) error {
		inputs, err := flake.Inputs(flakePath)
		if err != nil {
			return err
		}
		if jsonOutput {
			out, err := flake.InputsToJSON(inputs)
			if err != nil {
				return err
			}
			fmt.Println(out)
			return nil
		}
		fmt.Println(flake.FormatInputs(inputs))
		return nil
	},
}

func runShow(cmd *cobra.Command, args []string) error {
	entries, err := profile.List(profilePath)
	if err != nil {
		return err
	}
	if jsonOutput {
		out, err := profile.EntriesToJSON(entries)
		if err != nil {
			return err
		}
		fmt.Println(out)
		return nil
	}
	fmt.Println(profile.FormatEntries(entries))
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}