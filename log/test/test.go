package main

import mlog "mall/log"

func main() {
	mlog.SetName("test")
	mlog.Error("err in this")
	mlog.Info("info in this")
}
