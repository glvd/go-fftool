package fftool

var sf *StreamFormat

func init() {
	var e error
	sf, e = NewFFProbe().StreamFormat(`d:\video\女大学生的沙龙室.Room.Salon.College.Girls.2018.HD720P.X264.AAC.Korean.CHS.mp4`)
	if e != nil {
		panic(e)
	}
}
