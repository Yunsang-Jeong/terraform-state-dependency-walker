package main

import "github.com/Yunsang-Jeong/terraform-state-dependency-walker/cmd"

func main() {
	getAll := &cmd.GetAllCmd{}
	check := &cmd.CheckCmd{}

	cmd.RootCmd.AddCommand(getAll.Init())
	cmd.RootCmd.AddCommand(check.Init())
	cmd.Execute()
}
