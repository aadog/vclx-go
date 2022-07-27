package vclx

import (
	"fmt"
	"github.com/aadog/dict-go"
	"github.com/praveen001/ds/list/arraylist"
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
	"time"
)

var _AutoListViewExclude map[uintptr][]int = map[uintptr][]int{}
var _SetupListViewCheckBoxsOnCb map[uintptr]func(sender vcl.IObject, item *vcl.TListItem) = map[uintptr]func(sender vcl.IObject, item *vcl.TListItem){}

// 安装list自动宽度
func ListViewSetupAutoWidth(lv *vcl.TListView, excludes ...int) {
	_AutoListViewExclude[lv.Instance()] = excludes
	lv.SetAutoWidthLastColumn(false)
	lv.SetOnResize(func(sender vcl.IObject) {
		listv := vcl.AsListView(sender)
		exclude := _AutoListViewExclude[listv.Instance()]
		lvcount := listv.Columns().Count()
		nowidth := int32(0)
		for _, it := range exclude {
			nowidth += listv.Columns().Items(int32(it)).Width()
		}
		isnoin := func(idx int32) bool {
			for _, it := range exclude {
				if int32(it) == idx {
					return true
				}
			}
			return false
		}
		countwidth := (listv.Width() - nowidth) - 20
		itcount := lvcount - int32(len(exclude))
		for i := int32(0); i < lvcount; i++ {
			if isnoin(i) == false {
				col := listv.Columns().Items(i)
				col.SetWidth(countwidth / itcount)
			}
		}
	})
}
func ListViewCheckedAutoCallBack(item *vcl.TListItem, ls *arraylist.ArrayList) {
	iv, ok := ls.Get(int(item.Index()))
	d := iv.(*dict.Dict)
	if ok == true {
		d.Set(":checked", !d.GetBool(":checked"))
	}
}
func ListViewSetupAutoChecked(lv *vcl.TListView, witdh int, onchecked func(sender vcl.IObject, item *vcl.TListItem)) {
	_SetupListViewCheckBoxsOnCb[lv.Instance()] = onchecked
	var checkboxbytes = []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x10\x00\x00\x00 \b\x02\x00\x00\x00\x94\xebo\x9b\x00\x00\x00\x06tRNS\x00\x00\x00\x00\x00\x00n\xa6\a\x91\x00\x00\x00yIDATx\x9cc`\xa0\x13\xf0Y\xf0T\xb4\xf9.~\x04T\x83\xd0\x00\xe4\u05ee\xbd\x82\x1f\x01ՠk\xf8\x8f\x1b\x8cj\x18\xcc\x1aH\x88i\x92\xd3\xd2 \x04\x15\x9b\xa3\xfd\xe6h\xe0G@5\b\r\x04UC\x10\xb1\x1a\x80\xf1@\x82\x06x\xdc\x11\xa5\x01Y5a\rh\xaa\xb1h@\x96\xc3T\x8d\xae\x01Y\x05V\xd5\xd8m@\x06D\x05+\x1e\xd58=\x8dK5\t\x11\x87E\x03\xc9i\x89x\x00\x00y\xe7w\x8dM;\xd9\xc9\x00\x00\x00\x00IEND\xaeB`\x82")
	checkboxsteam := vcl.NewMemoryStreamFromBytes(checkboxbytes)
	defer checkboxsteam.Free()
	checkboxsteam.SetPosition(0)
	pic := vcl.NewPicture()
	defer pic.Free()
	pic.LoadFromStream(checkboxsteam)
	stateImages := vcl.NewImageList(lv)
	stateImages.AddSliced(pic.Bitmap(), 1, 2)
	lv.SetStateImages(stateImages)
	lv.SetStateImagesWidth(int32(witdh))
	lv.SetOnItemChecked(func(sender vcl.IObject, item *vcl.TListItem) {
		onchecked = _SetupListViewCheckBoxsOnCb[sender.Instance()]
		if onchecked != nil {
			//item.SetChecked(!item.Checked())
			onchecked(sender, item)
		}
	})
	lv.SetOnMouseDown(func(sender vcl.IObject, button types.TMouseButton, shift types.TShiftState, x, y int32) {
		listview := vcl.AsListView(sender)
		onchecked = _SetupListViewCheckBoxsOnCb[sender.Instance()]
		//if listview.Checkboxes() && x <= int32(witdh) { //16= f.stateImages.Width
		if x <= int32(witdh) { //16= f.stateImages.Width
			item := listview.GetItemAt(x, y)
			if item != nil {
				idx := item.Index()
				r := item.DisplayRect(types.DrIcon)
				if y >= r.Top && y <= r.Bottom {
					if onchecked != nil {
						//item.SetChecked(!item.Checked())
						onchecked(sender, item)
						listview.Refresh()
					}
					// 不知道为啥idx=0时要repaint，但Repaint效率不如Invalidate
					if idx == 0 {
						listview.Repaint()
					} else {
						listview.Invalidate()
					}
				}
			}
		}
	})
}

// 返回当前运行时间
func RunTimeS(st time.Time) string {
	etm := time.Since(st)
	return fmt.Sprintf("%02d时 %02d分 %02d秒", int(etm.Hours())%60, int(etm.Minutes())%60, int(etm.Seconds())%60)
}
