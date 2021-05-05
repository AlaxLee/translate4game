package translate4game

import (
	"encoding/json"
	"fmt"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	"io/ioutil"
	"os"
	"time"
)

func GetPosition(name string) (p *Position) {
	fmt.Println("Now you can mouse left for ", name)
	mleft := robotgo.AddEvent("mleft")
	if mleft {
		p = GetCurrentPosition()
		fmt.Println("You got "+name+": ", p.X, p.Y)
	}
	return
}

type Translate struct {
	LeftUP     *Position `json:"left_up"`
	RightDown  *Position `json:"right_down"`
	Image      *Position `json:"image"`
	Game       *Position `json:"game"`
	Baidufanyi *Position `json:"baidufanyi"`
}

func (tr *Translate) Init() {
	configName := "tr.conf"
	if err := tr.loadFromConfig(configName); err == nil {
		fmt.Println("load conf from " + configName + " successful: ")
		fmt.Println("LeftUpPositionForScreen:", tr.LeftUP.X, tr.LeftUP.Y)
		fmt.Println("RightDownPositionForScreen:", tr.RightDown.X, tr.RightDown.Y)
		fmt.Println("CatchedImagePosition:", tr.Image.X, tr.Image.Y)
		fmt.Println("GamePosition:", tr.Game.X, tr.Game.Y)
		fmt.Println("BaidufanyiPosition:", tr.Baidufanyi.X, tr.Baidufanyi.Y)
		return
	} else {
		fmt.Println("load conf from "+configName+" failed: ", err)
	}

	tr.LeftUP = GetPosition("LeftUpPositionForScreen")
	time.Sleep(time.Second)
	tr.RightDown = GetPosition("RightDownPositionForScreen")
	time.Sleep(time.Second)
	tr.Image = GetPosition("CatchedImagePosition")
	time.Sleep(time.Second)
	tr.Game = GetPosition("GamePosition")
	time.Sleep(time.Second)
	tr.Baidufanyi = GetPosition("BaidufanyiPosition")
	time.Sleep(time.Second)

	if err := tr.writeToConfig(configName); err != nil {
		fmt.Println("write conf file "+configName+"failed: ", err)
	}
}

func (tr *Translate) loadFromConfig(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	d := json.NewDecoder(f)
	err = d.Decode(tr)
	if err != nil {
		return err
	}
	return nil
}

func (tr *Translate) writeToConfig(filename string) error {
	b, err := json.Marshal(tr)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filename, b, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (tr *Translate) Start() {
	fmt.Println("---- Now you can translate by press key 'g' ----")
	robotgo.EventHook(hook.KeyDown, []string{"g"}, func(e hook.Event) {
		startPosition := GetCurrentPosition()
		tr.CaptureScreen()
		time.Sleep(time.Millisecond * 500)
		tr.CopyImage()
		time.Sleep(time.Millisecond * 500)
		tr.TranslateByBaidufanyi()
		time.Sleep(time.Millisecond * 100)
		tr.ReturnToGame(startPosition)
	})

	s := robotgo.EventStart()
	<-robotgo.EventProcess(s)
}

func (tr *Translate) CaptureScreen() {
	bitmap := robotgo.CaptureScreen(tr.LeftUP.X, tr.LeftUP.Y, tr.RightDown.X-tr.LeftUP.X, tr.RightDown.Y-tr.LeftUP.Y)
	// use `defer robotgo.FreeBitmap(bit)` to free the bitmap
	defer robotgo.FreeBitmap(bitmap)

	//fmt.Println("...", bitmap)

	fx, fy := robotgo.FindBitmap(bitmap)
	fmt.Println("FindBitmap------ ", fx, fy)

	robotgo.SaveBitmap(bitmap, "test.png")
}

func (tr *Translate) CopyImage() {
	robotgo.MoveMouseSmooth(tr.Image.X, tr.Image.Y, 0.0, 0.0, 100)
	robotgo.MouseClick("left", true)
	time.Sleep(time.Millisecond * 500)
	robotgo.KeyTap("c", "control")
}

func (tr *Translate) TranslateByBaidufanyi() {
	robotgo.MoveMouseSmooth(tr.Baidufanyi.X, tr.Baidufanyi.Y, 0.0, 0.0, 100)
	robotgo.MouseClick("left", true)
	time.Sleep(time.Millisecond * 100)
	robotgo.KeyTap("v", "control")
}

func (tr *Translate) ReturnToGame(startPosition *Position) {
	robotgo.MoveMouseSmooth(tr.Game.X, tr.Game.Y, 0.0, 0.0, 100)
	robotgo.MouseClick("left", false)
	robotgo.MoveMouseSmooth(startPosition.X, startPosition.Y, 0.0, 0.0, 100)
}
