package speech

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

const (
	Shadda = string(rune(0x0651))

	Sukoon   = string(rune(0x0652))
	Damma    = string(rune(0x064F))
	Fatha    = string(rune(0x064E))
	Kasra    = string(rune(0x0650))
	Dammatan = string(rune(0x064C))
	Fathatan = string(rune(0x064B))
	Kasratan = string(rune(0x064D))

	Placeholder = string(rune(0x25CC))

	SuperscriptAlef = string(rune(0x670))
)

var GrammaticalTags = []string{
	"اسم فاعل",
	"اسم فعل الأمر",
	"اسم فعل الماضي",
	"اسم فعل المضارع",
	"اسم المفعول",
	"المصدر",
	"اسم المصدر",
	"الصفة المشبهة",
	"مثال المبالغة",
	"اسم اشارة",
	"موصول محختص",
	"موصول مسترك",
	"ضمبر متصل",
	"ضمير منفصل",
	"مستتر جوازا",
	"مستتر وجوبا",
	"علم جنسي",
	"علم شخصي",
	"معرف بالاضافة المحضة",
	"معرف بالاضافة غير المحضة",
	"المعرف بأل الجنسية",
	"المعرف بأل العهدية",
	"نكرة منصوبة",
	"نكرة مرفوعة",
	"اسم استفهام",
	"اسم مرفوع",
	"اسم منصوب",
	"اسم مجرور",
	"حرف نفي",
	"حرف جر أصلي",
	"حرف جر زائد",
	"حرف جر شبه زائد",
	"حرف لنداء القريب",
	"حرف لنداء البعيد",
	"أخوات إن",
	"لا النافية للجنس",
	"ناصب المضارع",
	"حرف جزم",
	"حرف استفهام",
	"حرف التفسير",
	"موصول حرفي",
	"حرف عطف",
	"حرف محذوف",
	"فعل تام",
	"أخوات كان",
	"أخوات كاد",
	"أفعال القلوب",
	"أفعال التحويل",
	"فعل مرفوع",
	"فعل منصوب",
	"فعل مجزوم",
	"مبني",

	// Second column
	"نائب الفاعل",
	"اسم أخوات كاد",
	"اسم أخوات كان",
	"خبر أخوات إن",
	"اسم حرف نفي",
	"فاعل",
	"اعرابه",
	"خبر مبتدأ",

	"خبر أخوات كان",
	"خبر حرف نفي",
	"مفعول معه",
	"مفعول لأجله",
	"مفعول فيه",
	"مفعول به ثان",
	"مفعول به ثالث",
	"مفعول لأجله",
	"مفعول مطلق",
	"مفعول به",
	"مفعول لأجله",
	"مفعول مطلق",
	"منادى",
	"مفعول معه",
	"مفعول به",
	"اسم أخوات إن",

	"مجرور بالاضافة",
	"مجرور بحرف",

	"منصوب بحرف",
	"مجزوم بحرف",
	"مجزوم بالطلب",

	"صلة الموصول",
	"صلة",
}

func IsWhitespace(letter rune) bool {
	return letter == ' '
}

func IsWordPunctuation(word string) bool {
	for _, l := range word {
		if IsPunctuation(l) {
			return true
		}
	}
	return false
}

var Punctuation = regexp.MustCompile("[\\.:«»،\"—]")

func IsPunctuation(letter rune) bool {
	return Punctuation.MatchString(string(letter))
}

// isArabicLetter does not include tashkeel
func IsArabicLetter(letter rune) bool {
	if letter >= 0x0621 && letter <= 0x063A {
		return true
	}
	if letter >= 0x0641 && letter <= 0x064A {
		return true
	}
	return false
}

func IsVowel(letter rune) bool {
	sl := string(letter)
	return sl == Sukoon || sl == Damma || sl == Fatha || sl == Kasra ||
		sl == Dammatan || sl == Fathatan || sl == Kasratan
}

func IsShadda(letter rune) bool {
	return string(letter) == Shadda
}

func CleanContent(content string) (string, error) {
	for _, c := range content {
		if !(IsArabicLetter(c) || IsWhitespace(c) || IsPunctuation(c)) {
			return "", errors.New(fmt.Sprintf("%v is an invalid letter", c))
		}
	}

	// Remove double spaces
	r, _ := regexp.Compile(" +")
	content = r.ReplaceAllString(content, " ")

	// Trim sentence
	content = strings.TrimFunc(content, unicode.IsSpace)
	return content, nil
}
