package engine

import (
	"bytes"
	"embed"
	"fmt"
	"game/layout"
	"game/maps"
	"game/role"
	"game/tools"
	"math"
	"runtime"
	"strconv"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//config
const (
	SCREENWIDTH         int   = 490
	SCREENHEIGHT        int   = 300
	OFFSETX             int   = -30
	OFFSETY             int   = -30
	LAYOUTX             int   = 790
	LAYOUTY             int   = 480
	PLAYERCENTERX       int64 = 388 //LAYOUTX/2
	PLAYERCENTERY       int64 = 242 //LAYOUTY/2
	WEOFFSETX           int   = 127
	WEOFFSETY           int   = 14
	GAMESCENELOGIN      int   = 1
	GAMESCENESELECTROLE int   = 2
	GAMESCENEOPENDOOR   int   = 3
	GAMESCENESTART      int   = 4
)

type Game struct {
	count             int
	player            *role.Player
	maps              *maps.MapBase
	ui                *layout.UI
	currentGameScence int
	audioContext      *audio.Player
	//monster *role.Monster
}

var (
	counts          int  = 0
	frameNums       int  = 4
	flg             bool = false
	frameSpeed      int  = 5
	changeScenceFlg bool = false
	doorCountFlg    bool = false
	loadingFlg      bool = false
	musicIsPlay     bool = false
)
var op, opS, opMouse, opWea, opSkill *ebiten.DrawImageOptions

var mouseX, mouseY int

var mouseIcon *ebiten.Image

var images *embed.FS

var gameSceneType int = 0
var cont *audio.Context

//factory
func NewGame(img *embed.FS) *Game {
	images = img
	cont = audio.NewContext(48000)
	//map
	m := maps.NewMap(img)
	//palayer
	r := role.NewPlayer(float64(LAYOUTX/2), float64(LAYOUTY/2), tools.IDLE, 0, 0, 0, img, m)
	//UI
	u := layout.NewUI(img)

	gameEngine := &Game{
		count:             0,
		player:            r,
		maps:              m,
		ui:                u,
		currentGameScence: GAMESCENELOGIN,
		audioContext:      nil,
	}
	return gameEngine
}

func (g *Game) StartEngine() {
	//Hidden Mouse Icon
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	opMouse = &ebiten.DrawImageOptions{}
	s, _ := images.ReadFile("resource/UI/mouse.png")
	mouseIcon = tools.GetEbitenImage(s)
	opS = &ebiten.DrawImageOptions{}
	op = &ebiten.DrawImageOptions{}
	//
	w := sync.WaitGroup{}
	w.Add(1)
	//UI Init
	go func() {
		g.ui.LoadGameLoginImages()
		runtime.GC()
		w.Done()
	}()
	w.Wait()
	go func() {
		runtime.GC()
	}()
}

//Change Load
func (g *Game) ChangeScene(name string) {
	if name == "select" {
		g.currentGameScence = GAMESCENESELECTROLE
		g.ui.LoadGameCharaSelectImages()
		runtime.GC()

	} else if name == "loading" {
		g.currentGameScence = GAMESCENEOPENDOOR
		//GC
		g.ui.ClearSlice(0)
		g.ui.ClearGlobalVariable()

	} else if name == "game" {
		g.currentGameScence = GAMESCENESTART
		w := sync.WaitGroup{}
		w.Add(3)
		//Palyer Init
		go func() {
			g.player.LoadImages()
			runtime.GC()
			w.Done()
		}()
		//Ui
		go func() {
			g.ui.LoadGameImages()
			runtime.GC()
			w.Done()
		}()
		//Map Init
		go func() {
			g.maps.LoadMap()
			w.Done()
		}()
		w.Wait()
		go func() {
			runtime.GC()
		}()
	}

}
func (g *Game) Update() error {
	mouseX, mouseY = ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && changeScenceFlg == false {
		if g.currentGameScence == GAMESCENELOGIN {
			if mouseX > 286 && mouseX < 503 && mouseY > 150 && mouseY < 218 {
				g.currentGameScence = GAMESCENESELECTROLE
				changeScenceFlg = true
				w := sync.WaitGroup{}
				w.Add(1)
				go func() {
					g.ChangeScene("select")
					w.Done()
				}()
				w.Wait()
				changeScenceFlg = false
			}

		} else if g.currentGameScence == GAMESCENESELECTROLE {
			if mouseX > 684 && mouseX < 759 && mouseY > 390 && mouseY < 470 {
				changeScenceFlg = true
				g.currentGameScence = GAMESCENEOPENDOOR
				changeScenceFlg = false
			}
		}
	}

	//check
	if changeScenceFlg == false {
		if g.currentGameScence == GAMESCENESTART {
			g.changeScenceGameUpdate()
		} else if g.currentGameScence == GAMESCENEOPENDOOR {
			g.ChangeScenceOpenDoorUpdate()
		} else {
			g.ChangeScenceLoginUpdate()
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	//Is Change Scence ?
	if changeScenceFlg == false {
		if g.currentGameScence == GAMESCENESTART {
			g.ChangeScenceGameDraw(screen)
		} else if g.currentGameScence == GAMESCENESELECTROLE {
			g.ChangeScenceSelectDraw(screen)
		} else if g.currentGameScence == GAMESCENEOPENDOOR {
			g.ChangeScenceOpenDoorDraw(screen)
		} else {
			g.ChangeScenceLoginDraw(screen)
		}
	}
	//Draw Mouse Icon
	g.DrawMouseIcon(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return LAYOUTX, LAYOUTY
}

//Draw Mouse Icon
func (g *Game) DrawMouseIcon(screen *ebiten.Image) {
	opMouse.GeoM.Reset()
	opMouse.GeoM.Rotate(-0.5)
	opMouse.Filter = ebiten.FilterLinear
	opMouse.GeoM.Translate(float64(mouseX), float64(mouseY))
	screen.DrawImage(mouseIcon, opMouse)
}

//Draw Game Update
func (g *Game) changeScenceGameUpdate() {
	if musicIsPlay == false {
		g.PlayMusic("Bar_act2_complete_tombs.wav", "wav")
	}
	g.count++
	if g.player.State != tools.ATTACK {
		g.player.State = tools.IDLE
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.player.MouseX = mouseX
		g.player.MouseY = mouseY
		//防止点击UI界面也移动
		if mouseY < 436 {
			flg = true
		}
		//for test
		if mouseX > 201 && mouseX < 228 && mouseY > 446 && mouseY < 468 {
			g.ui.SetDisplay()
		}

	}

	//Calculate direction
	dir := tools.CaluteDir(PLAYERCENTERX, PLAYERCENTERY, int64(g.player.MouseX), int64(g.player.MouseY))

	//EVENT Listen
	g.ui.EventLoop()

	//attack
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		g.player.SkillName = "liehuo"
		if g.player.State != tools.ATTACK {
			counts = 0
		}
		flg = false
		if g.player.Direction != dir || g.player.State != tools.ATTACK {
			g.player.SetPlayerState(tools.ATTACK, dir)

		}
	}

	//mouse controll
	flg = g.player.GetMouseController(dir, flg)
	//states
	if g.player.State == tools.IDLE {
		frameNums = 16
		frameSpeed = 5
		g.player.SetPlayerState(tools.IDLE, dir)

	} else if g.player.State == tools.ATTACK {
		frameNums = 16
		frameSpeed = 1
	} else {
		frameNums = 8
		frameSpeed = 5
	}
}

//Draw Game Scence
func (g *Game) ChangeScenceGameDraw(screen *ebiten.Image) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("has error is :", r)
		}
	}()
	//Draw Background
	screen.DrawImage(g.maps.BgImage, g.maps.OpBg)
	//
	var name string
	//nameSkill := ""
	switch g.player.State {
	case tools.ATTACK:
		name = strconv.Itoa(g.player.Direction) + "_attack_" + strconv.Itoa(counts) + ".png"
		//nameSkill = strconv.Itoa(g.player.Direction) + "_skill_" + strconv.Itoa(counts) + ".png"
	case tools.IDLE:
		name = strconv.Itoa(g.player.Direction) + "_stand_" + strconv.Itoa(counts) + ".png"
	default:
		if counts >= 8 {
			counts = 0
		}
		name = strconv.Itoa(g.player.Direction) + "_run_" + strconv.Itoa(counts) + ".png"
	}
	imagess, x, y := g.player.GetAnimator("man", name)
	//Idel -> Walk Offset
	if g.player.State == tools.RUN {
		x += -4
		y += -3
	}

	//Idel -> Walk -> Attack Offset
	if g.player.State == tools.ATTACK {
		x += -50
		y += -30
	}
	//Draw Shadow
	opS.GeoM.Reset()
	opS.GeoM.Translate(float64(LAYOUTX/2+x-350), float64(LAYOUTY/2+y+365))
	opS.Filter = ebiten.FilterLinear
	opS.GeoM.Rotate(-0.5)
	opS.GeoM.Scale(1, 0.5)
	opS.ColorM.Scale(0, 0, 0, 1)
	opS.ColorM.Translate(0, 0, 0, 0)
	screen.DrawImage(imagess, opS)
	//Draw Player
	op.GeoM.Reset()
	op.GeoM.Translate(float64(LAYOUTX/2+OFFSETX+x), float64(LAYOUTY/2+OFFSETY+y))
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(imagess, op)

	//Draw UI
	g.ui.DrawUI(screen)

	//Draw Skill
	// if g.player.State == ATTACK {
	// 	imagey, x, y := g.player.GetAnimator("skill", nameSkill)
	// 	//skill option
	// 	opSkill = &ebiten.DrawImageOptions{}
	// 	opSkill.GeoM.Translate(float64(SCREENWIDTH/2+x), float64(SCREENHEIGHT/2+y))
	// 	opSkill.CompositeMode = ebiten.CompositeModeLighter
	// 	opSkill.GeoM.Scale(1.5, 1.5)
	// 	opSkill.Filter = ebiten.FilterLinear
	// 	screen.DrawImage(imagey, opSkill)
	// }

	//Draw Debug
	len := tools.Distance(PLAYERCENTERX, PLAYERCENTERY, int64(mouseX), int64(mouseY))
	re := tools.Angle(math.Abs(float64(int64(mouseY)-PLAYERCENTERY)), len)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS %d\nplayer position %d,%d\nmouse position %d,%d\ndir %d\nAngle %f",
		int64(ebiten.CurrentFPS()), int64(g.player.X), int64(g.player.Y), g.player.MouseX, g.player.MouseY, tools.CaluteDir(PLAYERCENTERX, PLAYERCENTERY, int64(g.player.MouseX), int64(g.player.MouseY)), re))

	//Change Frame
	if g.count > frameSpeed {
		counts++
		g.count = 0
		if counts >= frameNums {
			counts = 0
		}
	}
}

//Draw Login Update
func (g *Game) ChangeScenceLoginUpdate() {
	if musicIsPlay == false {
		g.PlayMusic("Act0-Intro.mp3", "mp3")
	}
	frameSpeed_clone := 0
	if g.currentGameScence == GAMESCENELOGIN {
		frameSpeed_clone = 2
	} else {
		frameSpeed_clone = 5
	}
	g.count++
	//Change Frame
	if g.count > frameSpeed_clone {
		counts++
		g.count = 0
		if counts >= 30 {
			counts = 0
		}
	}
}

//Draw Login Scence
func (g *Game) ChangeScenceLoginDraw(screen *ebiten.Image) {
	//Draw UI
	g.ui.DrawUI(screen)

	//Draw Logo Left
	name := "logoFireLeft_" + strconv.Itoa(counts) + ".png"
	left, _, _ := g.ui.GetAnimator("logo", name)
	opLo := &ebiten.DrawImageOptions{}
	opLo.Filter = ebiten.FilterLinear
	opLo.GeoM.Translate(220, 0)
	opLo.CompositeMode = ebiten.CompositeModeLighter
	opLo.GeoM.Scale(1, 0.7)
	screen.DrawImage(left, opLo)
	//Draw Logo Right
	name = "logoFireRight_" + strconv.Itoa(counts) + ".png"
	right, _, _ := g.ui.GetAnimator("logo", name)
	opRo := &ebiten.DrawImageOptions{}
	opRo.Filter = ebiten.FilterLinear
	opRo.GeoM.Translate(float64(220+right.Bounds().Max.X), 0)
	opRo.CompositeMode = ebiten.CompositeModeLighter
	opRo.GeoM.Scale(1, 0.7)
	screen.DrawImage(right, opRo)
	//Draw Debug
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS %d\nmouse position %d,%d",
		int64(ebiten.CurrentFPS()), mouseX, mouseY))

	//Change Frame
	if g.count > frameSpeed {
		counts++
		g.count = 0
		if counts >= 30 {
			counts = 0
		}
	}
}

//Draw Select Scence
func (g *Game) ChangeScenceSelectDraw(screen *ebiten.Image) {
	//Draw UI
	g.ui.DrawUI(screen)

	//Draw Fire
	name := "fire_" + strconv.Itoa(counts) + ".png"
	fire, _, _ := g.ui.GetAnimator("logo", name)
	opf := &ebiten.DrawImageOptions{}
	opf.Filter = ebiten.FilterLinear
	opf.GeoM.Translate(350, 375)
	opf.CompositeMode = ebiten.CompositeModeLighter
	opf.GeoM.Scale(1, 0.7)
	screen.DrawImage(fire, opf)

	//Draw Role Ba
	name = "ba_" + strconv.Itoa(counts) + ".png"
	ba, _, _ := g.ui.GetAnimator("role", name)
	opBa := &ebiten.DrawImageOptions{}
	opBa.Filter = ebiten.FilterLinear
	opBa.GeoM.Translate(356, 120)
	screen.DrawImage(ba, opBa)

	//Draw Debug
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS %d\nmouse position %d,%d", int64(ebiten.CurrentFPS()), mouseX, mouseY))
}

//Draw OpenDoor Scence
func (g *Game) ChangeScenceOpenDoorDraw(screen *ebiten.Image) {
	//Draw Open Door
	name := "loading_" + strconv.Itoa(counts) + ".png"
	door, _, _ := g.ui.GetAnimator("logo", name)
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	op.GeoM.Translate(268, 120)
	screen.DrawImage(door, op)

	//Draw Debug
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS %d\nmouse position %d,%d", int64(ebiten.CurrentFPS()), mouseX, mouseY))
}

//Draw OpenDoor Update
func (g *Game) ChangeScenceOpenDoorUpdate() {
	if doorCountFlg == false {
		counts = 0
		doorCountFlg = true
	}
	g.count++
	//Change Frame
	if g.count > 10 && counts != 9 {
		counts++
		g.count = 0
	}

	// Change Scence
	if counts == 9 && loadingFlg == false {
		loadingFlg = true
		g.currentGameScence = GAMESCENESTART
		w := sync.WaitGroup{}
		w.Add(1)
		go func() {
			//close music
			g.CloseMusic()
			g.ChangeScene("game")
			w.Done()
		}()
		w.Wait()
		changeScenceFlg = false
	}
}

//Play Music
func (g *Game) PlayMusic(name, ty string) {
	musicIsPlay = true
	switch ty {
	case "mp3":
		go func() {
			bgm, _ := images.ReadFile("resource/BGM/" + name)
			ss, _ := mp3.Decode(cont, bytes.NewReader(bgm))
			g.audioContext = nil
			g.audioContext, _ = cont.NewPlayer(ss)
			g.audioContext.Play()
		}()
	case "wav":
		go func() {
			bgm, _ := images.ReadFile("resource/BGM/" + name)
			ss, _ := wav.Decode(cont, bytes.NewReader(bgm))
			g.audioContext = nil
			g.audioContext, _ = cont.NewPlayer(ss)
			g.audioContext.Play()
		}()
	}
}

//Close Music
func (g *Game) CloseMusic() {
	musicIsPlay = false
	go func() {
		g.audioContext.Close()
	}()
}
