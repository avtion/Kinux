package translator

import (
	localesZh "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
)

var Trans ut.Translator

func init() {
	Trans = newZhTranslator()
}

func newZhTranslator() ut.Translator {
	zhTranslator := localesZh.New()
	uni := ut.New(zhTranslator, zhTranslator)
	trans, _ := uni.GetTranslator("zh")
	return trans
}
