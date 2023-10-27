package speech

const (
	ActiveParticiple     = "اسم فاعل"
	PassiveParticiple    = "اسم المفعول"
	VerbalNoun           = "المصدر"
	DemonstrativePronoun = "اسم اشارة"
	SuffixedPronoun      = "ضمبر متصل"
	NominativeNoun       = "اسم مرفوع"
	AccusativeNoun       = "اسم منصوب"
	GenetiveNoun         = "اسم مجرور"

	NegativeParticle        = "حرف نفي"
	CoordinatingConjunction = "حرف عطف"

	IndicativeVerb  = "فعل مرفوع"
	SubjunctiveVerb = "فعل منصوب"
)

// TODO(Amr Ojjeh): Convert to English constants
var GrammaticalTags = []string{
	ActiveParticiple,
	"اسم فعل الأمر",
	"اسم فعل الماضي",
	"اسم فعل المضارع",
	PassiveParticiple,
	VerbalNoun,
	"اسم المصدر",
	"الصفة المشبهة",
	"مثال المبالغة",
	DemonstrativePronoun,
	"موصول محختص",
	"موصول مسترك",
	SuffixedPronoun,
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
	NominativeNoun,
	AccusativeNoun,
	GenetiveNoun,
	NegativeParticle,
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
	CoordinatingConjunction,
	"حرف محذوف",
	"فعل تام",
	"أخوات كان",
	"أخوات كاد",
	"أفعال القلوب",
	"أفعال التحويل",
	IndicativeVerb,
	SubjunctiveVerb,
	"فعل مجزوم",
}
