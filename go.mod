module gtts

go 1.21

toolchain go1.21.5

require (
	github.com/bwmarrin/discordgo v0.27.1
	github.com/goccy/go-json v0.10.2
	github.com/joho/godotenv v1.5.1
	local.packages/dca v0.0.0-00010101000000-000000000000
	local.packages/disgolf v0.0.0-00010101000000-000000000000
)

require (
	github.com/FedorLap2006/disgolf v0.0.0-20221004200601-99cfc3d9a0e1 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/jonas747/ogg v0.0.0-20161220051205-b4f6f4cf3757 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
)

replace local.packages/dca => ./dca

replace local.packages/disgolf => ./disgolf
