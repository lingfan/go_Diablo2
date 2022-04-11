package layout

import (
	"embed"
	"fmt"
	"game/maps"
	"game/status"
	"game/tools"
	"image"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fzipp/texturepacker"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	plist_png     *image.NRGBA
	plist_R_png   *image.NRGBA
	plist_sheet   *texturepacker.SpriteSheet
	plist_R_sheet *texturepacker.SpriteSheet
	isClick       bool = false
	mouseIcon     *ebiten.Image
	mouseIconCopy ebiten.Image
	opMouse       *ebiten.DrawImageOptions
	mouseRoate    float64 = -0.5
)

//UI类
type UI struct {
	image             *embed.FS
	Compents          []*Sprite            //普通UI存放集合
	HiddenCompents    []*Sprite            //可以被隐藏的UI组件集合
	MiniPanelCompents []*Sprite            //MINI板的UI集合
	ItemsCompents     []*SpriteItems       //Items的UI集合
	status            *status.StatusManage //状态管理器
	maps              *maps.MapBase        //地图
	BagLayout         [5][10]string        //4*10 背包 1*10 装备栏
	tempBag           [1]*SpriteItems      //临时Items存放
}

func NewUI(images *embed.FS, s *status.StatusManage, m *maps.MapBase) *UI {
	//初始化背包 数据
	itemsLayout := [5][10]string{
		{"HP0", "HP0", "HP0", "HP0", "book_0,4", "dun-5_0,5", "dun-5_0,5", "", "dun_0,8", "dun_0,8"},
		{"body-3_1,0", "body-3_1,0", "", "", "book_0,4", "dun-5_0,5", "dun-5_0,5", "", "dun_0,8", "dun_0,8"},
		{"body-3_1,0", "body-3_1,0", "hand_2,2", "hand_2,2", "book_2,4", "dun-5_0,5", "dun-5_0,5", "", "head-3_2,8", "head-3_2,8"},
		{"body-3_1,0", "body-3_1,0", "hand_2,2", "hand_2,2", "book_2,4", "HP0", "", "", "head-3_2,8", "head-3_2,8"},
		{"", "", "", "", "", "", "", "", "", ""},
		//头盔526,8  左手武器412,54 右手武器644,54 项链599,36 铠甲526,80 手套413,182 左戒指485,181 腰带527,181 右戒指599,183 靴子644,183
	}
	ui := &UI{
		image:             images,
		Compents:          make([]*Sprite, 0, 12),
		HiddenCompents:    make([]*Sprite, 0, 6),
		MiniPanelCompents: make([]*Sprite, 0, 6),
		ItemsCompents:     make([]*SpriteItems, 0, 10),
		status:            s,
		maps:              m,
		BagLayout:         itemsLayout,
	}
	//鼠标Icon设置
	opMouse = &ebiten.DrawImageOptions{}
	ss, _ := ui.image.ReadFile("resource/UI/mouse.png")
	mouseIcon = tools.GetEbitenImage(ss)
	return ui
}

//加载进入游戏UI
func (u *UI) LoadGameImages() {
	u.ClearSlice(10)
	var len float64 = 0
	// go func() {
	// 	plist, _ := u.image.ReadFile("resource/UI/0000.png")
	// 	plist_json, _ := u.image.ReadFile("resource/man/warrior/ba.json")
	// 	plist_sheet, plist_png = tools.GetImageFromPlistPaletted(plist, plist_json)
	// 	runtime.GC()
	// }()
	s, _ := u.image.ReadFile("resource/UI/0000.png")
	mgUI := tools.GetEbitenImage(s)
	u.AddComponent(QuickCreate(len, 480-float64(mgUI.Bounds().Max.Y), mgUI, 0, nil), tools.ISNORCOM)

	s, _ = u.image.ReadFile("resource/UI/HP.png")
	mgUI = tools.GetEbitenImage(s)
	u.AddComponent(QuickCreate(28, 480-float64(mgUI.Bounds().Max.Y+13), mgUI, 0, nil), tools.ISNORCOM)

	len += 115

	s, _ = u.image.ReadFile("resource/UI/chisha.png")
	mgUI = tools.GetEbitenImage(s)

	u.AddComponent(QuickCreate(len, 480-float64(mgUI.Bounds().Max.Y), mgUI, 0, nil), tools.ISNORCOM)

	len += float64(mgUI.Bounds().Max.X)

	s, _ = u.image.ReadFile("resource/UI/0001.png")
	mgUI = tools.GetEbitenImage(s)

	u.AddComponent(QuickCreate(len, 480-float64(mgUI.Bounds().Max.Y), mgUI, 0, nil), tools.ISNORCOM)

	len += float64(mgUI.Bounds().Max.X)

	s, _ = u.image.ReadFile("resource/UI/0002.png")
	mgUI = tools.GetEbitenImage(s)

	u.AddComponent(QuickCreate(len, 480-float64(mgUI.Bounds().Max.Y), mgUI, 0, nil), tools.ISNORCOM)

	len += float64(mgUI.Bounds().Max.X)

	s, _ = u.image.ReadFile("resource/UI/0003.png")
	mgUI = tools.GetEbitenImage(s)

	u.AddComponent(QuickCreate(len, 480-float64(mgUI.Bounds().Max.Y), mgUI, 0, nil), tools.ISNORCOM)

	len += float64(mgUI.Bounds().Max.X)

	s, _ = u.image.ReadFile("resource/UI/0004.png")
	mgUI = tools.GetEbitenImage(s)

	u.AddComponent(QuickCreate(len, 480-float64(mgUI.Bounds().Max.Y), mgUI, 0, nil), tools.ISNORCOM)

	s, _ = u.image.ReadFile("resource/UI/liehuo.png")
	mgUI1 := tools.GetEbitenImage(s)
	u.AddComponent(QuickCreate(627, 480-float64(mgUI1.Bounds().Max.Y), mgUI1, 0, nil), tools.ISNORCOM)

	len += float64(mgUI.Bounds().Max.X)

	s, _ = u.image.ReadFile("resource/UI/0005.png")
	mgUI = tools.GetEbitenImage(s)

	u.AddComponent(QuickCreate(len, 480-float64(mgUI.Bounds().Max.Y), mgUI, 0, nil), tools.ISNORCOM)

	s, _ = u.image.ReadFile("resource/UI/MP.png")
	mgUI = tools.GetEbitenImage(s)

	u.AddComponent(QuickCreate(684, 480-float64(mgUI.Bounds().Max.Y+13), mgUI, 1, nil), tools.ISNORCOM)

	s, _ = u.image.ReadFile("resource/UI/skill_btn.png")
	mgUI = tools.GetEbitenImage(s)
	u.AddComponent(QuickCreate(204, 441, mgUI, 0, func(i spriteInterface) {
		if isClick == false {
			isClick = true
			go func() {
				on := *i.(*Sprite).images
				s, _ = u.image.ReadFile("resource/UI/skill_btn_down.png")
				mgUI = tools.GetEbitenImage(s)
				i.(*Sprite).images = mgUI
				time.Sleep(tools.CLOSEBTNSLEEP)
				i.(*Sprite).images = &on
				runtime.GC()
				isClick = false
			}()
		}
	}, true), tools.ISNORCOM)
	u.AddComponent(QuickCreate(562, 441, mgUI, 0, func(i spriteInterface) {
		if isClick == false {
			isClick = true
			go func() {
				on := *i.(*Sprite).images
				s, _ = u.image.ReadFile("resource/UI/skill_btn_down.png")
				mgUI = tools.GetEbitenImage(s)
				i.(*Sprite).images = mgUI
				time.Sleep(tools.CLOSEBTNSLEEP)
				i.(*Sprite).images = &on
				runtime.GC()
				isClick = false
			}()
		}
	}, true), tools.ISNORCOM)

	//描画装备栏和包裹UI
	s, _ = u.image.ReadFile("resource/UI/eq_0.png")
	mgUI = tools.GetEbitenImage(s)
	u.AddComponent(QuickCreate(395, 0, mgUI, 0, nil), tools.ISHIDDEN)

	len += float64(mgUI.Bounds().Max.X)

	s, _ = u.image.ReadFile("resource/UI/eq_1.png")
	mgUI = tools.GetEbitenImage(s)
	u.AddComponent(QuickCreate(395+256, 0, mgUI, 0, nil), tools.ISHIDDEN)

	len = 395
	s, _ = u.image.ReadFile("resource/UI/bag_0.png")
	mgUI = tools.GetEbitenImage(s)
	u.AddComponent(QuickCreate(395, 176, mgUI, 0, nil), tools.ISHIDDEN)

	len += float64(mgUI.Bounds().Max.X)
	s, _ = u.image.ReadFile("resource/UI/bag_1.png")
	mgUI = tools.GetEbitenImage(s)
	u.AddComponent(QuickCreate(395+256, 176, mgUI, 0, nil), tools.ISHIDDEN)

	//关闭装备栏按钮
	s, _ = u.image.ReadFile("resource/UI/close_btn_on.png")
	mgUI = tools.GetEbitenImage(s)
	u.AddComponent(QuickCreate(414, 384, mgUI, 0, func(i spriteInterface) {
		if isClick == false {
			isClick = true
			go func() {
				on := *i.(*Sprite).images
				s, _ = u.image.ReadFile("resource/UI/close_btn_down.png")
				mgUI = tools.GetEbitenImage(s)
				i.(*Sprite).images = mgUI
				time.Sleep(tools.CLOSEBTNSLEEP)
				i.(*Sprite).images = &on
				u.setHidden(tools.ISHIDDEN)
				go func() {
					for _, v := range u.MiniPanelCompents {
						v.SetPosition(100, 0)

						if v.clickMaxX != 0 {
							v.clickMinX += 100
							v.clickMaxX += 100
						}
					}
				}()
				runtime.GC()
				//恢复因打开包裹导致的人物偏移
				u.status.UIOFFSETX = 0
				//恢复影子偏移
				u.status.ShadowOffsetX = -350
				u.status.ShadowOffsetY = 365
				//恢复玩家中心位置
				u.status.PLAYERCENTERX = 388
				//恢复地图偏移
				u.maps.ChangeMapTranslate(200, 0)
				isClick = false
			}()
		}
	}, true), tools.ISHIDDEN)

	//背包物品LOOP Start
	//临时Map
	TempArray := make(map[string]int, 10)
	items := u.BagLayout
	for i := 0; i < 4; i++ {
		for j := 0; j < 10; j++ {
			if strings.Contains(items[i][j], "_") {
				if _, ok := TempArray[items[i][j]]; !ok {
					TempArray[items[i][j]] = 0
					t := strings.Split(items[i][j], "_")
					s, _ = u.image.ReadFile("resource/UI/" + t[0] + ".png")
					mgUI = tools.GetEbitenImage(s)
					x := 413 + j*29
					y := 254 + i*29
					u.AddComponent(QuickCreateItems(float64(x), float64(y), t[0], mgUI, 1, u.ItemsEvent(), 1, true), 0)
				}
			} else if items[i][j] != "" {
				s, _ = u.image.ReadFile("resource/UI/" + items[i][j] + ".png")
				mgUI = tools.GetEbitenImage(s)
				x := 413 + j*29
				y := 254 + i*29
				u.AddComponent(QuickCreateItems(float64(x), float64(y), items[i][j], mgUI, 1, u.ItemsEvent(), 1, true), 0)
			}
		}
	}
	//手动销毁临时Map
	TempArray = nil
	//背包物品LOOP END

	//注册mini板打开按钮
	s, _ = u.image.ReadFile("resource/UI/open_minipanel_btn.png")
	mgUI = tools.GetEbitenImage(s)
	u.AddComponent(QuickCreate(390, 443, mgUI, 0, func(i spriteInterface) {
		if isClick == false {
			isClick = true
			go func() {
				if u.status.OpenMiniPanel {
					u.setHidden(tools.ISMINICOM)
					s, _ = u.image.ReadFile("resource/UI/open_minipanel_down.png")
					mgUI = tools.GetEbitenImage(s)
					i.(*Sprite).images = mgUI
					time.Sleep(tools.CLOSEBTNSLEEP)
					s, _ = u.image.ReadFile("resource/UI/close_minipanel_btn.png")
					mgUI = tools.GetEbitenImage(s)
					i.(*Sprite).images = mgUI
				} else {
					u.SetDisplay(tools.ISMINICOM)
					s, _ = u.image.ReadFile("resource/UI/close_minipanel_down.png")
					mgUI = tools.GetEbitenImage(s)
					i.(*Sprite).images = mgUI
					time.Sleep(tools.CLOSEBTNSLEEP)
					s, _ = u.image.ReadFile("resource/UI/open_minipanel_btn.png")
					mgUI = tools.GetEbitenImage(s)
					i.(*Sprite).images = mgUI
				}
				runtime.GC()
				isClick = false
			}()
		}
	}, true), tools.ISNORCOM)
	//注册mini板
	s, _ = u.image.ReadFile("resource/UI/miniPanel.png")
	mgUI = tools.GetEbitenImage(s)
	baseX := float64(tools.LAYOUTX/2 - mgUI.Bounds().Max.X/2)
	u.AddComponent(QuickCreate(baseX, 406, mgUI, 0, nil), tools.ISMINICOM)
	baseX += 4
	//
	s, _ = u.image.ReadFile("resource/UI/mini_menu_man.png")
	mgUI = tools.GetEbitenImage(s)
	u.AddComponent(QuickCreate(baseX, 410, mgUI, 0, nil), tools.ISMINICOM)
	//
	baseX += float64(mgUI.Bounds().Max.X) + 4
	s, _ = u.image.ReadFile("resource/UI/mini_menu_wea.png")
	mgUI = tools.GetEbitenImage(s)
	u.AddComponent(QuickCreate(baseX, 410, mgUI, 0, func(i spriteInterface) {
		if isClick == false {
			isClick = true
			go func() {
				on := *i.(*Sprite).images
				s, _ = u.image.ReadFile("resource/UI/mini_menu_wea_down.png")
				mgUI = tools.GetEbitenImage(s)
				if !u.status.OpenBag {
					//判断MINI板子的最左端坐标是否超过最大极限
					if x, _ := u.MiniPanelCompents[0].GetPosition(); x > 209 {
						go func() {
							//设置因打开包裹导致的人物偏移
							u.status.UIOFFSETX = -200
							//修改地图偏移
							u.maps.ChangeMapTranslate(-200, 0)
							//修改玩家中心位置
							u.status.PLAYERCENTERX -= 200
							//修改人物影子偏移
							u.status.ShadowOffsetX = u.status.ShadowOffsetX + 14
							u.status.ShadowOffsetY = u.status.ShadowOffsetY - 79
							for _, v := range u.MiniPanelCompents {
								v.SetPosition(-100, 0)
								if v.clickMaxX != 0 {
									v.clickMinX -= 100
									v.clickMaxX -= 100
								}
							}
						}()
					}
				}
				i.(*Sprite).images = mgUI
				time.Sleep(tools.CLOSEBTNSLEEP)
				i.(*Sprite).images = &on
				u.SetDisplay(tools.ISHIDDEN)
				runtime.GC()
				isClick = false
			}()
		}
	}, true), tools.ISMINICOM)
	baseX += float64(mgUI.Bounds().Max.X) + 4
	s, _ = u.image.ReadFile("resource/UI/mini_menu_j.png")
	mgUI = tools.GetEbitenImage(s)
	u.AddComponent(QuickCreate(baseX, 410, mgUI, 0, nil), tools.ISMINICOM)
	baseX += float64(mgUI.Bounds().Max.X) + 4
	s, _ = u.image.ReadFile("resource/UI/mini_menu_m.png")
	mgUI = tools.GetEbitenImage(s)
	u.AddComponent(QuickCreate(baseX, 410, mgUI, 0, nil), tools.ISMINICOM)
	baseX += float64(mgUI.Bounds().Max.X) + 4
	s, _ = u.image.ReadFile("resource/UI/mini_menu_mess.png")
	mgUI = tools.GetEbitenImage(s)
	u.AddComponent(QuickCreate(baseX, 410, mgUI, 0, nil), tools.ISMINICOM)
	baseX += float64(mgUI.Bounds().Max.X) + 4
	s, _ = u.image.ReadFile("resource/UI/mini_menu_s.png")
	mgUI = tools.GetEbitenImage(s)
	u.AddComponent(QuickCreate(baseX, 410, mgUI, 0, nil), tools.ISMINICOM)
	baseX += float64(mgUI.Bounds().Max.X) + 4
	s, _ = u.image.ReadFile("resource/UI/mini_menu_st.png")
	mgUI = tools.GetEbitenImage(s)
	u.AddComponent(QuickCreate(baseX, 410, mgUI, 0, nil), tools.ISMINICOM)

	u.setHidden(tools.ISHIDDEN)
	u.setHidden(tools.ISMINICOM)

}

//加载登录游戏UI
func (u *UI) LoadGameLoginImages() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("has error is :", r)
		}
	}()
	var len float64 = 0
	var scales float64 = 0.8
	s, _ := u.image.ReadFile("resource/UI/login0.png")
	mgUI := tools.GetEbitenImage(s)
	op := newSprite()
	op.SetPosition(len, 0)
	op.op.GeoM.Scale(1, scales)
	op.addImage(mgUI)
	u.AddComponent(op, tools.ISNORCOM)

	len += float64(mgUI.Bounds().Max.X)

	s, _ = u.image.ReadFile("resource/UI/login1.png")
	mgUI = tools.GetEbitenImage(s)
	op = newSprite()
	op.SetPosition(len, 0)
	op.op.GeoM.Scale(1, scales)
	op.addImage(mgUI)
	u.AddComponent(op, tools.ISNORCOM)

	len += float64(mgUI.Bounds().Max.X)

	s, _ = u.image.ReadFile("resource/UI/login2.png")
	mgUI = tools.GetEbitenImage(s)
	op = newSprite()
	op.SetPosition(len, 0)
	op.op.GeoM.Scale(1, scales)
	op.addImage(mgUI)
	u.AddComponent(op, tools.ISNORCOM)

	len += float64(mgUI.Bounds().Max.X)

	s, _ = u.image.ReadFile("resource/UI/login3.png")
	mgUI = tools.GetEbitenImage(s)
	op = newSprite()
	op.SetPosition(len, 0)
	op.op.GeoM.Scale(1, scales)
	op.addImage(mgUI)
	u.AddComponent(op, tools.ISNORCOM)

	len = 0
	var offset float64 = 340
	s, _ = u.image.ReadFile("resource/UI/login8.png")
	mgUI = tools.GetEbitenImage(s)
	op = newSprite()
	op.SetPosition(len, offset)
	op.op.GeoM.Scale(1, scales)
	op.addImage(mgUI)
	u.AddComponent(op, tools.ISNORCOM)

	len += float64(mgUI.Bounds().Max.X)

	s, _ = u.image.ReadFile("resource/UI/login9.png")
	mgUI = tools.GetEbitenImage(s)
	op = newSprite()
	op.SetPosition(len, offset)
	op.op.GeoM.Scale(1, scales)
	op.addImage(mgUI)
	u.AddComponent(op, tools.ISNORCOM)

	len += float64(mgUI.Bounds().Max.X)

	s, _ = u.image.ReadFile("resource/UI/login10.png")
	mgUI = tools.GetEbitenImage(s)
	op = newSprite()
	op.SetPosition(len, offset)
	op.op.GeoM.Scale(1, scales)
	op.addImage(mgUI)
	u.AddComponent(op, tools.ISNORCOM)

	len += float64(mgUI.Bounds().Max.X)

	s, _ = u.image.ReadFile("resource/UI/login11.png")
	mgUI = tools.GetEbitenImage(s)
	op = newSprite()
	op.SetPosition(len, offset)
	op.op.GeoM.Scale(1, scales)
	op.addImage(mgUI)
	u.AddComponent(op, tools.ISNORCOM)

	len = 0

	s, _ = u.image.ReadFile("resource/UI/login4.png")
	mgUI = tools.GetEbitenImage(s)
	op = newSprite()
	op.SetPosition(len, float64(mgUI.Bounds().Max.Y))
	op.op.GeoM.Scale(1, scales)
	op.addImage(mgUI)
	u.AddComponent(op, tools.ISNORCOM)

	len += float64(mgUI.Bounds().Max.X)

	s, _ = u.image.ReadFile("resource/UI/login5.png")
	mgUI = tools.GetEbitenImage(s)
	op = newSprite()
	op.SetPosition(len, float64(mgUI.Bounds().Max.Y))
	op.op.GeoM.Scale(1, scales)
	op.addImage(mgUI)
	u.AddComponent(op, tools.ISNORCOM)

	len += float64(mgUI.Bounds().Max.X)

	s, _ = u.image.ReadFile("resource/UI/login6.png")
	mgUI = tools.GetEbitenImage(s)
	op = newSprite()
	op.SetPosition(len, float64(mgUI.Bounds().Max.Y))
	op.op.GeoM.Scale(1, scales)
	op.addImage(mgUI)
	u.AddComponent(op, tools.ISNORCOM)

	len += float64(mgUI.Bounds().Max.X)

	s, _ = u.image.ReadFile("resource/UI/login7.png")
	mgUI = tools.GetEbitenImage(s)
	op = newSprite()
	op.SetPosition(len, float64(mgUI.Bounds().Max.Y))
	op.op.GeoM.Scale(1, scales)
	op.addImage(mgUI)
	u.AddComponent(op, tools.ISNORCOM)

	go func() {
		plist, _ := u.image.ReadFile("resource/UI/logo.png")
		plist_json, _ := u.image.ReadFile("resource/UI/logo.json")
		plist_sheet, plist_png = tools.GetImageFromPlist(plist, plist_json)
		runtime.GC()
	}()

}

//加载游戏选择角色UI
func (u *UI) LoadGameCharaSelectImages() {
	u.ClearSlice(1)
	s, _ := u.image.ReadFile("resource/UI/charactSelect.png")
	mgUI := tools.GetEbitenImage(s)
	op := newSprite()
	op.SetPosition(0, 0)
	op.op.GeoM.Scale(1, 0.8)
	op.addImage(mgUI)
	u.AddComponent(op, tools.ISNORCOM)
	w := &sync.WaitGroup{}
	w.Add(2)
	go func() {
		plist, _ := u.image.ReadFile("resource/UI/logo.png")
		plist_json, _ := u.image.ReadFile("resource/UI/logo.json")
		plist_sheet, plist_png = tools.GetImageFromPlist(plist, plist_json)
		w.Done()
	}()
	go func() {
		plist, _ := u.image.ReadFile("resource/UI/selectRoles.png")
		plist_json, _ := u.image.ReadFile("resource/UI/selectRoles.json")
		plist_R_sheet, plist_R_png = tools.GetImageFromPlist(plist, plist_json)
		w.Done()
	}()
	w.Wait()
	go func() {
		runtime.GC()
	}()

}

//图集获取图片
func (u *UI) GetAnimator(flg, name string) (*ebiten.Image, int, int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	if flg == "role" {

		return ebiten.NewImageFromImage(plist_R_png.SubImage(plist_R_sheet.Sprites[name].Frame)), plist_R_sheet.Sprites[name].SpriteSourceSize.Min.X, plist_R_sheet.Sprites[name].SpriteSourceSize.Min.Y

	} else {

		return ebiten.NewImageFromImage(plist_png.SubImage(plist_sheet.Sprites[name].Frame)), plist_sheet.Sprites[name].SpriteSourceSize.Min.X, plist_sheet.Sprites[name].SpriteSourceSize.Min.Y
	}
}

//组件注册
func (u *UI) AddComponent(s spriteInterface, ImageType uint8) {
	if ImageType == tools.ISHIDDEN {
		//将UI压入通用集合
		u.Compents = append(u.Compents, s.(*Sprite))
		//将UI压入可隐藏集合
		u.HiddenCompents = append(u.HiddenCompents, s.(*Sprite))
	} else if ImageType == tools.ISMINICOM {
		//将UI压入通用集合
		u.Compents = append(u.Compents, s.(*Sprite))
		//将UI压入MINI板集合
		u.MiniPanelCompents = append(u.MiniPanelCompents, s.(*Sprite))
	} else if ImageType == tools.ISITEMS {
		//将UI压入Items集合
		u.ItemsCompents = append(u.ItemsCompents, s.(*SpriteItems))
	} else {
		//将UI压入通用集合
		u.Compents = append(u.Compents, s.(*Sprite))
	}
}

//显示UI
func (u *UI) SetDisplay(ImageType uint8) {
	if ImageType == tools.ISHIDDEN {
		u.status.OpenBag = true
		for _, v := range u.HiddenCompents {
			v.isDisplay = true
		}
	} else {
		u.status.OpenMiniPanel = true
		for _, v := range u.MiniPanelCompents {
			v.isDisplay = true
		}
	}

}

//隐藏UI
func (u *UI) setHidden(ImageType uint8) {
	if ImageType == tools.ISHIDDEN {
		u.status.OpenBag = false
		for _, v := range u.HiddenCompents {
			v.isDisplay = false
		}
	} else {
		u.status.OpenMiniPanel = false
		for _, v := range u.MiniPanelCompents {
			v.isDisplay = false
		}
	}

}

//清除切片
func (u *UI) ClearSlice(cap int) {
	u.Compents = make([]*Sprite, 0, cap)
	u.HiddenCompents = make([]*Sprite, 0, cap/2)
	u.MiniPanelCompents = make([]*Sprite, 0, cap/2)
	u.ItemsCompents = make([]*SpriteItems, 0, 10)
}

//渲染UI
func (u *UI) DrawUI(screen *ebiten.Image) {
	//渲染UI
	for _, v := range u.Compents {
		if v.layer == 0 && v.isDisplay {
			screen.DrawImage(v.images, v.op)
		}
	}
	//渲染层级为1的UI
	for _, v := range u.Compents {
		if v.layer == 1 && v.isDisplay {
			screen.DrawImage(v.images, v.op)
		}
	}
	//当包裹打开的时候，渲染包裹内物品和装备 TODO
	if u.status.OpenBag {
		for _, v := range u.ItemsCompents {
			if v.bgIsDisplay {
				screen.DrawImage(v.imageBg, v.opBg)
			}
			screen.DrawImage(v.images, v.op)
		}

	}
}

//事件轮询
func (u *UI) EventLoop(mouseX, mouseY int) {

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		//普通UI事件轮询
		for _, v := range u.Compents {
			if v.hasEvent == 1 && v.isDisplay {
				if mouseX >= v.clickMinX && mouseX <= v.clickMaxX && mouseY >= v.clickMinY && mouseY <= v.clickMaxY {
					//实行回调函数
					v.f(v)
				}
			}
		}
		//包裹打开的情况下监听
		if u.status.OpenBag {
			//items UI事件轮询
			for _, v := range u.ItemsCompents {
				if v.hasEvent == 1 {
					if mouseX > v.clickMinX && mouseX < v.clickMaxX && mouseY > v.clickMinY && mouseY < v.clickMaxY {
						v.clickEvnet(v, mouseX, mouseY)
					}
				}
			}
		}

		//点击包裹区域并且在包裹坐标范围内
		if u.status.OpenBag && mouseX >= 408 && mouseY >= 6 && mouseX <= 698 && mouseY <= 372 && u.tempBag[0] != nil && u.status.IsTakeItem {

			s := u.tempBag[0]
			//给鼠标加一个假偏移，防止双击
			if u.AddItemToBag(mouseX+u.status.Mouseoffset, mouseY+u.status.Mouseoffset, s.itemName) {
				//鼠标还原
				mouseIcon = &mouseIconCopy
				//清理临时区
				u.tempBag[0] = nil
				mouseRoate = -0.5
				//恢复防止双击的鼠标偏移量
				u.status.Mouseoffset = 500
				//拿起物品flag设置
				u.status.IsTakeItem = false
				go func() {
					time.Sleep(tools.CLOSEBTNSLEEP)
					isClick = false
				}()

			}
		}

	}
	//包裹打开的情况下监听
	if u.status.OpenBag {
		//items UI事件轮询
		for _, v := range u.ItemsCompents {
			if v.hasEvent == 1 {
				v.touchEvnet(v, mouseX, mouseY)
			}
		}
	}

}

//GC 清理变量
func (u *UI) ClearGlobalVariable() {
	plist_R_sheet = nil
	plist_R_png = nil
}

//添加物品到包裹 or 装备栏
func (u *UI) AddItemToBag(mousex, mousey int, itemName string) bool {
	//屏幕坐标转换成包裹坐标
	x := int(mousey-254) / 29
	y := int(mousex-413) / 29
	sizeX, sizeY := tools.GetItemsCellSize(itemName)
	if sizeX != 0 && sizeY != 0 {
		//x y这个单元格有位置是否
		if x >= 0 && x <= 3 && y >= 0 && y <= 9 && u.BagLayout[x][y] == "" {
			//是否相同size的时候
			if sizeX == 1 && sizeY == 1 {
				u.BagLayout[x][y] = itemName
				s, _ := u.image.ReadFile("resource/UI/" + itemName + ".png")
				mgUI := tools.GetEbitenImage(s)
				layoutX := 413 + y*29
				layoutY := 254 + x*29
				u.AddComponent(QuickCreateItems(float64(layoutX), float64(layoutY), itemName, mgUI, 1, u.ItemsEvent(), 1, true), 0)
				return true
			} else {
				//循环判断是否可以放下
				for i := 0; i < sizeX; i++ {
					for j := 0; j < sizeY; j++ {
						if x+j > 3 || y+i > 9 || u.BagLayout[x+j][y+i] != "" {
							return false
						}
					}
				}
				name := strconv.Itoa(x) + "," + strconv.Itoa(y)
				for i := 0; i < sizeX; i++ {
					for j := 0; j < sizeY; j++ {
						u.BagLayout[x+j][y+i] = itemName + "_" + name
					}
				}
				s, _ := u.image.ReadFile("resource/UI/" + itemName + ".png")
				mgUI := tools.GetEbitenImage(s)
				layoutX := 413 + y*29
				layoutY := 254 + x*29
				u.AddComponent(QuickCreateItems(float64(layoutX), float64(layoutY), itemName, mgUI, 1, u.ItemsEvent(), 1, true), 0)
				return true
			}
		} else if mousex >= 397 && mousey >= 5 && mousex <= 705 && mousey <= 247 {
			//判断是否放入装备栏
			return u.JudgeCanToEquip(mousex, mousey, itemName)
		} else {
			return false
		}
	} else {
		return false
	}
}

//从包裹删除物品
func (u *UI) DelItemFromBag(imageX, imageY int) {
	//屏幕坐标转换成包裹坐标
	x := int(imageY-254) / 29
	y := int(imageX-413) / 29
	if x >= 0 && x <= 3 && y >= 0 && y <= 9 && u.BagLayout[x][y] != "" {
		if strings.Contains(u.BagLayout[x][y], "_") {
			itemName := u.BagLayout[x][y]
			for i := 0; i < 4; i++ {
				for j := 0; j < 10; j++ {
					if u.BagLayout[i][j] == itemName {
						u.BagLayout[i][j] = ""
					}
				}
			}
		} else {
			u.BagLayout[x][y] = ""
		}
		layoutX := 413 + y*29
		layoutY := 254 + x*29
		for k, v := range u.ItemsCompents {
			//根据具体的图片坐标删除 支持唯一性
			if v.imagex == float64(layoutX) && v.imagey == float64(layoutY) {
				if k != len(u.ItemsCompents)-1 {
					u.ItemsCompents = append(u.ItemsCompents[0:k], u.ItemsCompents[k+1:]...)
				} else {
					u.ItemsCompents = u.ItemsCompents[0:k]
				}
			}
		}
	} else if xx, _, key := u.JudgeIsEquipArea(imageX, imageY); xx != 0 {
		//删除装备栏
		for k, v := range u.ItemsCompents {
			//根据具体的图片坐标删除 支持唯一性
			if v.imagex == float64(imageX) && v.imagey == float64(imageY) {
				if k != len(u.ItemsCompents)-1 {
					u.ItemsCompents = append(u.ItemsCompents[0:k], u.ItemsCompents[k+1:]...)
				} else {
					u.ItemsCompents = u.ItemsCompents[0:k]
				}
				u.BagLayout[4][key] = ""
				return
			}
		}
	}
}

//重新绘制鼠标ICON
func (u *UI) DrawMouseIcon(screen *ebiten.Image, mouseX, mouseY int) {
	opMouse.GeoM.Reset()
	opMouse.GeoM.Rotate(mouseRoate)
	opMouse.Filter = ebiten.FilterLinear
	opMouse.GeoM.Translate(float64(mouseX), float64(mouseY))
	screen.DrawImage(mouseIcon, opMouse)
}

//判断是否可以放入装备栏
func (u *UI) JudgeCanToEquip(mousex, mousey int, itemName string) bool {
	x, y, key := u.JudgeIsEquipArea(mousex, mousey)
	if x != 0 && u.BagLayout[4][key] == "" {
		s, _ := u.image.ReadFile("resource/UI/" + itemName + ".png")
		mgUI := tools.GetEbitenImage(s)
		u.BagLayout[4][key] = itemName
		u.AddComponent(QuickCreateItems(float64(x), float64(y), itemName, mgUI, 1, u.ItemsEvent(), 0, true), 0)
		return true
	} else {
		return false
	}
}

//物品事件
func (u *UI) ItemsEvent() func(i spriteInterface, x, y int) {
	//注册监听
	item_event := func(i spriteInterface, x, y int) {
		if isClick == false {
			isClick = true
			go func() {
				if !u.status.IsTakeItem {
					//拿起物品flag设置
					u.status.IsTakeItem = true
					s := i.(*SpriteItems)
					go func() {
						time.Sleep(tools.CLOSEBTNSLEEP)
						u.status.Mouseoffset = 0
					}()
					//将拿起的物品放入临时区
					u.tempBag[0] = s
					mouseIconCopy = *mouseIcon
					mouseIcon = s.images
					mouseRoate = 0
					//拿起物品，从包裹中删除物品
					u.DelItemFromBag(int(s.imagex), int(s.imagey))
				}
			}()
		}
	}
	return item_event
}

//判断鼠标是否位于装备区
func (u *UI) JudgeIsEquipArea(mousex, mousey int) (int, int, uint8) {
	if mousex >= 526 && mousey >= 7 && mousex <= 580 && mousey <= 58 {
		//判断是否可以放入头盔
		return 526, 7, 0
	} else if mousex >= 412 && mousey >= 52 && mousex <= 465 && mousey <= 158 {
		//判断是否可以放入左武器
		return 412, 52, 1

	} else if mousex >= 526 && mousey >= 78 && mousex <= 579 && mousey <= 158 {
		//判断是否可以放入铠甲
		return 526, 78, 4

	} else if mousex >= 643 && mousey >= 51 && mousex <= 695 && mousey <= 155 {
		//判断是否可以放入右武器
		return 643, 51, 2

	} else if mousex >= 410 && mousey >= 181 && mousex <= 464 && mousey <= 234 {
		//判断是否可以放入手套
		return 410, 181, 5

	} else if mousex >= 642 && mousey >= 181 && mousex <= 695 && mousey <= 233 {
		//判断是否可以放入鞋
		return 642, 181, 9
	} else {
		return 0, 0, 0
	}
}
