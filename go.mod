module twinQuasarAppV2

go 1.15

require (
	github.com/Jeffail/gabs v1.4.0
	github.com/filecoin-project/go-address v0.0.5
	github.com/filecoin-project/go-jsonrpc v0.1.4-0.20210217175800-45ea43ac2bec
	github.com/filecoin-project/go-state-types v0.1.1-0.20210506134452-99b279731c48
	github.com/filecoin-project/lotus v1.9.0
	github.com/filecoin-project/specs-actors/v5 v5.0.1
	github.com/gopherjs/gopherjs v0.0.0-20200217142428-fce0ec30dd00 // indirect
	github.com/ilyakaznacheev/cleanenv v1.2.5
	github.com/ipfs/go-cid v0.0.7
	github.com/jackc/pgx/v4 v4.11.0
	github.com/lib/pq v1.7.0
	github.com/slack-go/slack v0.9.1
	github.com/smartystreets/assertions v1.1.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

replace github.com/filecoin-project/filecoin-ffi => ./extern/lotus/extern/filecoin-ffi

replace github.com/filecoin-project/lotus => ./extern/lotus
