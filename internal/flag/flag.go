package flag

import (
	"flag"
	"os"
)

var (
	FlagRunAddr  string
	Databaseflag string
	BonusAddress string
)

func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&Databaseflag, "d", "", "database path")
	flag.StringVar(&BonusAddress, "r", "", "bonus address")
	flag.Parse()
	if envRunAddr := os.Getenv("RUN_ADDRESS"); envRunAddr != "" {
		FlagRunAddr = envRunAddr
	}
	if envDatabaseflag := os.Getenv("DATABASE_URI"); envDatabaseflag != "" {
		Databaseflag = envDatabaseflag
	}
	if envBonusAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envBonusAddress != "" {
		BonusAddress = envBonusAddress
	}
}
