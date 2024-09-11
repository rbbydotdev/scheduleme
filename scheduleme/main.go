package main

import (
	"flag"
	"log"
	"os"
	"scheduleme/server"
	// "scheduleme/seeds"
)

func main() {
	// config.InitConfig()
	if len(os.Args) < 2 {
		log.Println("no command provided")
		os.Exit(2)
	}

	switch os.Args[1] {
	case "db": // handle 'flags' command
		flags := flag.NewFlagSet("db", flag.ExitOnError)

		// Subcommands
		migrate := flags.Bool("migrate", false, "Migrate the database")
		drop := flags.Bool("drop", false, "Drop the database")
		reset := flags.Bool("reset", false, "Reset the database")
		// seed := db.Bool("seed", false, "Seed the database")

		// Parse 'db' command
		flags.Parse(os.Args[2:])

		if *migrate {
			// cli.MigrateDB() // Replace with your function
		} else if *reset {
			// cli.ResetDB()

			// } else if *seed {
			// 	seeds.Seed()
		} else if *drop {
			// cli.DropDB() // Replace with your function
		} else {
			flags.PrintDefaults()
			os.Exit(1)
		}

	case "serve": // handle 'serve' command
		server.RunMain()

	default:
		log.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}
	os.Exit(0)
}
