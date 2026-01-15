/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)



// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gost-cli",
	Short: "gost-cli - консольная утилита для вычисления государственных криптографчиеских стандартов.",
	Long: `gost-cli - консольная утилита для вычисления государственных криптографчиеских стандартов. Реализованы:
1. Шифр Кузнечик (kuznyechik)
2. Шифр Магма (magma)
3. Хэш-функция Стрибог (streebog)
4. Схема ЭЦП ГОСТ 34.10-2012 (gost-sign и gost-verify)
5. ГПСЧ ГОСТ Р 1323565.1.006-2017 (rng)`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}


