package main

import (
	"github.com/NikosGour/chatter/internal"
	"github.com/NikosGour/chatter/internal/common"
	"github.com/NikosGour/chatter/internal/storage"
)

func main() {
	common.InitDotenv()

	db := storage.NewPostgreSQLStorage()

	api := internal.NewAPIServer(db)

	api.Start()
}
