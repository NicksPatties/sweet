var Cmd = &cobra.Command{
	Use:   "sweet [file|-]",
	Short: "The Software Engineer Exercise for Typing.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ex, err := fromArgs(cmd, args)
		if err != nil {
			log.Fatal(err)
		}
		Run(ex)
	},
}
