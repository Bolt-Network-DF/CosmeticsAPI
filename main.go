package cosmetic

import (
	"fmt"
	"github.com/Bolt-Network-DF/CosmeticsAPI/utils"
	"github.com/anthonynsimon/bild/blend"
	"github.com/anthonynsimon/bild/transform"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/skin"
	"image"
	"image/draw"
	"sync"
)

var originalSkins sync.Map

// CosmeticsManager handles player cosmetics.
type CosmeticsManager struct {
	capes map[string]skin.Cape
	wings map[string]struct {
		Texture  image.Image
		Geometry []byte
	}
	hats map[string]struct {
		Texture  image.Image
		Geometry []byte
	}
	legs map[string]struct {
		Texture  image.Image
		Geometry []byte
	}
}

// NewCosmeticsManager creates a new CosmeticsManager instance.
func NewCosmeticsManager() *CosmeticsManager {
	return &CosmeticsManager{
		capes: make(map[string]skin.Cape),
		wings: make(map[string]struct {
			Texture  image.Image
			Geometry []byte
		}),
		hats: make(map[string]struct {
			Texture  image.Image
			Geometry []byte
		}),
		legs: make(map[string]struct {
			Texture  image.Image
			Geometry []byte
		}),
	}
}

// GetCapes ...
func (cm *CosmeticsManager) GetCapes() map[string]skin.Cape {
	return cm.capes
}

// GetWings ...
func (cm *CosmeticsManager) GetWings() map[string]struct {
	return cm.wings
}

// GetHats ...
func (cm *CosmeticsManager) GetHats() map[string]struct {
	return cm.hats
}

// GetLegs ...
func (cm *CosmeticsManager) GetLegs() map[string]struct {
	return cm.legs
}

// RegisterCape registers a cape with a given name and path.
func (cm *CosmeticsManager) RegisterCape(name, path string) error {
	if _, ok := cm.capes[name]; !ok {
		cm.capes[name] = utils.ReadCapeData(path)
		return nil
	}
	
	return fmt.Errorf("cape %s already registered", name)
}

// RegisterWings registers wings with a given name and path.
func (cm *CosmeticsManager) RegisterWings(name, path string) error {
	return cm.registerCosmetic(cm.wings, name, path)
}

// RegisterHats registers hats with a given name and path.
func (cm *CosmeticsManager) RegisterHats(name, path string) error {
	return cm.registerCosmetic(cm.hats, name, path)
}

// RegisterLegs registers legs with a given name and path.
func (cm *CosmeticsManager) RegisterLegs(name, path string) error {
	return cm.registerCosmetic(cm.legs, name, path)
}

// SetCape sets the cape for the player.
func (cm *CosmeticsManager) SetCape(p *player.Player, name string) error {
	cm.RemoveCosmetics(p)

	cape, ok := cm.capes[name]
	if !ok {
		return fmt.Errorf("cape %v not registered", name)
	}

	cm.storeOriginalSkin(p)

	s := p.Skin()
	s.Cape = cape
	p.SetSkin(s)

	return nil
}

// SetWings sets wings for the player.
func (cm *CosmeticsManager) SetWings(p *player.Player, name string) error {
	cm.RemoveCosmetics(p)
	cm.storeOriginalSkin(p)

	return cm.setCosmetic(p, cm.wings, name)
}

// SetHats sets hats for the player.
func (cm *CosmeticsManager) SetHats(p *player.Player, name string) error {
	cm.RemoveCosmetics(p)
	cm.storeOriginalSkin(p)

	return cm.setCosmetic(p, cm.hats, name)
}

// SetLegs sets legs for the player.
func (cm *CosmeticsManager) SetLegs(p *player.Player, name string) error {
	cm.RemoveCosmetics(p)
	cm.storeOriginalSkin(p)

	return cm.setCosmetic(p, cm.legs, name)
}

// RemoveCosmetics removes any applied cosmetics and sets the player's skin back to the original.
func (cm *CosmeticsManager) RemoveCosmetics(p *player.Player) {
	xuid := p.XUID()
	if originalSkin, ok := originalSkins.Load(xuid); ok {
		p.SetSkin(originalSkin.(skin.Skin))
		originalSkins.Delete(xuid)
	} else {
		return
	}
}

// registerCosmetic registers a cosmetic with the given name and path.
func (cm *CosmeticsManager) registerCosmetic(cosmeticMap map[string]struct {
	Texture  image.Image
	Geometry []byte
}, name, path string) error {
	if _, ok := cosmeticMap[name]; !ok {
		texture, geometry := utils.ReadCosmeticData(path)
		cosmeticMap[name] = struct {
			Texture  image.Image
			Geometry []byte
		}{
			Texture:  texture,
			Geometry: geometry,
		}
		return nil
	}
	return fmt.Errorf("%s %v already registered", path, name)
}

// setCosmetic sets the cosmetic for the player.
func (cm *CosmeticsManager) setCosmetic(p *player.Player, cosmeticMap map[string]struct {
	Texture  image.Image
	Geometry []byte
}, name string) error {
	cosmetic, ok := cosmeticMap[name]
	if !ok {
		return fmt.Errorf("%s %v not registered", name)
	}

	s := p.Skin()
	img := blend.Add(transform.Resize(utils.GetImageFromSkin(s), 128, 128, transform.NearestNeighbor), cosmetic.Texture)
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, img.Bounds(), img, img.Bounds().Min, draw.Src)
	s.Pix = rgba.Pix
	s.Model = cosmetic.Geometry
	p.SetSkin(s)

	return nil
}

// storeOriginalSkin stores the original skin of the player before applying cosmetics.
func (cm *CosmeticsManager) storeOriginalSkin(p *player.Player) {
	xuid := p.XUID()

	if _, ok := originalSkins.Load(xuid); !ok {
		originalSkins.Store(xuid, p.Skin())
	}
}
