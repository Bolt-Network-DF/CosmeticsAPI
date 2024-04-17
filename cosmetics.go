package cosmetics

import (
	"errors"
	"github.com/Bolt-Network-DF/CosmeticsAPI/utils"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/skin"
)

var capes = make(map[string]Cape)

// Cape ...
type Cape struct {
	name string
	cape skin.Cape
}

// Register ...
func Register(name string, path string) {
	capes[name] = Cape{
		name: name,
		cape: utils.Read("assets/capes/" + path),
	}
}

// SetCape ...
func SetCape(p *player.Player, name string) error {
	cape, ok := capes[name]

	if !ok {
		return errors.New("cape not found")
	}

	s := p.Skin()

	s.Cape = cape.cape

	return nil
}
