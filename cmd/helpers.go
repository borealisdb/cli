package cmd

func defaultTokenFunc() (string, error) {
	return cliConfig.GetToken()
}
