/*
            DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE
                    Version 2, December 2004

 Copyright (C) 2004 Sam Hocevar <sam@hocevar.net>

 Everyone is permitted to copy and distribute verbatim or modified
 copies of this license document, and changing it is allowed as long
 as the name is changed.

            DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE
   TERMS AND CONDITIONS FOR COPYING, DISTRIBUTION AND MODIFICATION

  0. You just DO WHAT THE FUCK YOU WANT TO.
*/

package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

var (
	// hiragana
	hira1   []string // pure vowels, consonants without diacritics, moraic n
	hira2   []string // consonants with diacritics
	hira3   []string // palatal glides
	shira1  = []rune("あいうえおかきくけこさしすせそたちつてとなにぬねのはひふへほまみむめもやゆよらりるれろわをん")
	shira2  = []rune("がぎぐげござじずぜぞだぢづでどばびぶべぼぱぴぷぺぽ")
	shira3x = []rune("きぎしじちぢにひびぴみり")
	shira3y = []rune("ゃゅょ")

	// katakana
	kata1   []string
	kata2   []string
	kata3   []string
	skata1  = []rune("アイウエオカキクケコサシスセソタチツテトナニヌネノハヒフヘホマミムメモヤユヨラリルレロワヲン")
	skata2  = []rune("ガギグゲゴザジズゼゾダヂヅデドバビブベボパピプペポ")
	skata3x = []rune("キギシジチヂニヒビピミリ")
	skata3y = []rune("ャュョ")

	// kana romanizations
	nkana1 = []string{
		"a", "i", "u", "e", "o",
		"ka", "ki", "ku", "ke", "ko",
		"sa", "shi", "su", "se", "so",
		"ta", "chi", "tsu", "te", "to",
		"na", "ni", "nu", "ne", "no",
		"ha", "hi", "fu", "he", "ho",
		"ma", "mi", "mu", "me", "mo",
		"ya", "yu", "yo",
		"ra", "ri", "ru", "re", "ro",
		"wa", "wo",
		"n",
	}
	nkana2 = []string{
		"ga", "gi", "gu", "ge", "go",
		"za", "ji", "zu", "ze", "zo",
		"da", "ji", "zu", "de", "do",
		"ba", "bi", "bu", "be", "bo",
		"pa", "pi", "pu", "pe", "po",
	}
	nkana3 = []string{
		"kya", "kyu", "kyo", "gya", "gyu", "gyo",
		"sha", "shu", "sho", "ja", "ju", "jo",
		"cha", "chu", "cho", "ja", "ju", "jo",
		"nya", "nyu", "nyo",
		"hya", "hyu", "hyo", "bya", "byu", "byo", "pya", "pyu", "pyo",
		"mya", "myu", "myo",
		"rya", "ryu", "ryo",
	}

	// first year kanji
	kanji1a  []string
	kanji1b  []string
	kanji1c  []string
	kanji1d  []string
	kanji1e  []string
	kanji1f  []string
	kanji1g  []string
	kanji1h  []string
	kanji1i  []string
	skanji1a = "一二三四五六七八九十百千"
	skanji1b = "上下左右中大小月日年早"
	skanji1c = "木林山川土空田天生"
	skanji1d = "花草虫犬人名女男子"
	skanji1e = "目耳口手足見音"
	skanji1f = "力気円入出立休先"
	skanji1g = "夕本文字学校村町"
	skanji1h = "森正水火玉王石竹"
	skanji1i = "糸貝車金雨赤青白"

	mkanji = map[rune]kanji{
		// first year kanji, group A
		'一': kanji{'一', []string{"ichi", "itsu"}, []string{"hitotsu"}, "one"},
		'二': kanji{'二', []string{"ni"}, []string{"futatsu"}, "two"},
		'三': kanji{'三', []string{"san"}, []string{"mitsu"}, "three"},
		'四': kanji{'四', []string{"shi"}, []string{"yottsu", "yon"}, "four"},
		'五': kanji{'五', []string{"go"}, []string{"itsutsu"}, "five"},
		'六': kanji{'六', []string{"roku"}, []string{"muttsu"}, "six"},
		'七': kanji{'七', []string{"shichi"}, []string{"nanatsu"}, "seven"},
		'八': kanji{'八', []string{"hachi"}, []string{"yattsu"}, "eight"},
		'九': kanji{'九', []string{"ku", "kyuu"}, []string{"kokonotsu"}, "nine"},
		'十': kanji{'十', []string{"juu"}, []string{"too"}, "ten"},
		'百': kanji{'百', []string{"hyaku"}, []string{"momo"}, "hundred"},
		'千': kanji{'千', []string{"sen"}, []string{"chi"}, "thousand"},
		// group B
		'上': kanji{'上', []string{"jou"}, []string{"ue"}, "top, above"},
		'下': kanji{'下', []string{"ka", "ge"}, []string{"shita", "shimo", "moto"}, "bottom, below"},
		'左': kanji{'左', []string{"sa"}, []string{"hidari"}, "left"},
		'右': kanji{'右', []string{"u", "yuu"}, []string{"migi"}, "right"},
		'中': kanji{'中', []string{"chuu", "juu"}, []string{"naka"}, "center"},
		'大': kanji{'大', []string{"dai", "tai"}, []string{"ookii"}, "large"},
		'小': kanji{'小', []string{"shou"}, []string{"chiisai", "ko", "o"}, "small"},
		'月': kanji{'月', []string{"getsu", "gatsu"}, []string{"tsuki"}, "moon; month"},
		'日': kanji{'日', []string{"nichi", "jitsu"}, []string{"hi", "ka"}, "sun; day"},
		'年': kanji{'年', []string{"nen"}, []string{"toshi"}, "year"},
		'早': kanji{'早', []string{"sou", "satsu"}, []string{"hayai"}, "early"},
		// group C
		'木': kanji{'木', []string{"boku", "moku"}, []string{"ki"}, "tree"},
		'林': kanji{'林', []string{"rin"}, []string{"hayashi"}, "woods"},
		'山': kanji{'山', []string{"san", "zan"}, []string{"yama"}, "mountain"},
		'川': kanji{'川', []string{"sen"}, []string{"kawa"}, "river"},
		'土': kanji{'土', []string{"do", "to"}, []string{"tsuchi"}, "soil"},
		'空': kanji{'空', []string{"kuu"}, []string{"sora", "aku", "kara"}, "sky"},
		'田': kanji{'田', []string{"den"}, []string{"da", "ta"}, "rice paddy"},
		'天': kanji{'天', []string{"ten"}, []string{"ame", "ama"}, "heaven"},
		'生': kanji{'生', []string{"sei", "shou"}, []string{"ikiru", "umu", "nama"}, "life; bare; raw"},
		// group D
		'花': kanji{'花', []string{"ka"}, []string{"hana"}, "flower"},
		'草': kanji{'草', []string{"sou"}, []string{"kusa"}, "grass"},
		'虫': kanji{'虫', []string{"chuu"}, []string{"mushi"}, "insect"},
		'犬': kanji{'犬', []string{"ken"}, []string{"inu"}, "dog"},
		'人': kanji{'人', []string{"jin", "nin"}, []string{"hito"}, "person"},
		'名': kanji{'名', []string{"mei", "myou"}, []string{"na"}, "name"},
		'女': kanji{'女', []string{"jo", "nyo"}, []string{"onna"}, "female"},
		'男': kanji{'男', []string{"dan", "nan"}, []string{"otoko"}, "male"},
		'子': kanji{'子', []string{"shi", "su"}, []string{"ko"}, "child"},
		// group E
		'目': kanji{'目', []string{"moku"}, []string{"me"}, "eye"},
		'耳': kanji{'耳', []string{"ji", "ni"}, []string{"mimi"}, "ear"},
		'口': kanji{'口', []string{"kou"}, []string{"kuchi"}, "mouth"},
		'手': kanji{'手', []string{"shu"}, []string{"te"}, "hand"},
		'足': kanji{'足', []string{"soku"}, []string{"ashi"}, "foot"},
		'見': kanji{'見', []string{"ken", "gen"}, []string{"miru"}, "see"},
		'音': kanji{'音', []string{"on"}, []string{"ne", "oto"}, "sound"},
		// group F
		'力': kanji{'力', []string{"riki", "ryoku"}, []string{"chikara"}, "power"},
		'気': kanji{'気', []string{"ki", "ke"}, []string{"iki"}, "spirit; air"},
		'円': kanji{'円', []string{"en"}, []string{"maru"}, "yen; circle"},
		'入': kanji{'入', []string{"nyuu"}, []string{"hairu", "iru"}, "enter"},
		'出': kanji{'出', []string{"shutsu"}, []string{"deru"}, "exit"},
		'立': kanji{'立', []string{"ritsu"}, []string{"tatsu"}, "stand up"},
		'休': kanji{'休', []string{"kyuu"}, []string{"yasumu"}, "rest"},
		'先': kanji{'先', []string{"sen"}, []string{"saki"}, "previous"},
		// group G
		'夕': kanji{'夕', []string{"seki"}, []string{"yuu"}, "evening"},
		'本': kanji{'本', []string{"hon"}, []string{"moto"}, "book"},
		'文': kanji{'文', []string{"bun", "mon"}, []string{"fumi"}, "text"},
		'字': kanji{'字', []string{"ji"}, []string{"aza"}, "character, letter"},
		'学': kanji{'学', []string{"gaku"}, []string{"manabu"}, "study"},
		'校': kanji{'校', []string{"kou"}, []string{"kase"}, "school"},
		'村': kanji{'村', []string{"son"}, []string{"mura"}, "village"},
		'町': kanji{'町', []string{"chou"}, []string{"machi"}, "town"},
		// group H
		'森': kanji{'森', []string{"shin"}, []string{"mori"}, "forest"},
		'正': kanji{'正', []string{"sei", "shou"}, []string{"tadashii", "masa"}, "correct"},
		'水': kanji{'水', []string{"sui"}, []string{"mizu"}, "water"},
		'火': kanji{'火', []string{"ka"}, []string{"hi"}, "fire"},
		'玉': kanji{'玉', []string{"gyoku"}, []string{"tama"}, "gem"},
		'王': kanji{'王', []string{"ou"}, []string{"kimi"}, "king"},
		'石': kanji{'石', []string{"seki", "koku"}, []string{"ishi"}, "stone"},
		'竹': kanji{'竹', []string{"chiku"}, []string{"take"}, "bamboo"},
		// group I
		'糸': kanji{'糸', []string{"shi"}, []string{"ito"}, "thread"},
		'貝': kanji{'貝', []string{"bai"}, []string{"kai"}, "shellfish"},
		'車': kanji{'車', []string{"sha"}, []string{"kuruma"}, "vehicle"},
		'金': kanji{'金', []string{"kin"}, []string{"kane", "kana"}, "gold, money"},
		'雨': kanji{'雨', []string{"u"}, []string{"ame", "ama"}, "rain"},
		'赤': kanji{'赤', []string{"seki"}, []string{"aka"}, "red"},
		'青': kanji{'青', []string{"sei", "shou"}, []string{"ao"}, "blue"},
		'白': kanji{'白', []string{"haku"}, []string{"shiro", "shira"}, "white"},
	}

	syllab = map[string][]string{}

	chars    = map[string]interface{}{}
	cur      []string
	curonkun int
)

type kanji struct {
	rune    rune
	onyomi  []string
	kunyomi []string
	meaning string
}

func init() {
	for i, r := range shira1 {
		x := string(r)
		hira1 = append(hira1, x)
		chars[x] = nkana1[i]
	}
	for i, r := range shira2 {
		x := string(r)
		hira2 = append(hira2, x)
		chars[x] = nkana2[i]
	}
	for i, r := range shira3x {
		for j, a := range shira3y {
			x := string(r) + string(a)
			hira3 = append(hira3, x)
			chars[x] = nkana3[3*i+j]
		}
	}
	for i, r := range skata1 {
		x := string(r)
		kata1 = append(kata1, x)
		chars[x] = nkana1[i]
	}
	for i, r := range skata2 {
		x := string(r)
		kata2 = append(kata2, x)
		chars[x] = nkana2[i]
	}
	for i, r := range skata3x {
		for j, a := range skata3y {
			x := string(r) + string(a)
			kata3 = append(kata3, x)
			chars[x] = nkana3[3*i+j]
		}
	}
	for _, r := range skanji1a {
		kanji1a = append(kanji1a, string(r))
		chars[string(r)] = mkanji[r]
	}
	for _, r := range skanji1b {
		kanji1b = append(kanji1b, string(r))
		chars[string(r)] = mkanji[r]
	}
	for _, r := range skanji1c {
		kanji1c = append(kanji1c, string(r))
		chars[string(r)] = mkanji[r]
	}
	for _, r := range skanji1d {
		kanji1d = append(kanji1d, string(r))
		chars[string(r)] = mkanji[r]
	}
	for _, r := range skanji1e {
		kanji1e = append(kanji1e, string(r))
		chars[string(r)] = mkanji[r]
	}
	for _, r := range skanji1f {
		kanji1f = append(kanji1f, string(r))
		chars[string(r)] = mkanji[r]
	}
	for _, r := range skanji1g {
		kanji1g = append(kanji1g, string(r))
		chars[string(r)] = mkanji[r]
	}
	for _, r := range skanji1h {
		kanji1h = append(kanji1h, string(r))
		chars[string(r)] = mkanji[r]
	}
	for _, r := range skanji1i {
		kanji1i = append(kanji1i, string(r))
		chars[string(r)] = mkanji[r]
	}
	syllab["hira1"] = hira1
	syllab["hira2"] = hira2
	syllab["hira3"] = hira3
	syllab["hira"] = append(append(hira1, hira2...), hira3...)
	syllab["kata1"] = kata1
	syllab["kata2"] = kata2
	syllab["kata3"] = kata3
	syllab["kata"] = append(append(kata1, kata2...), kata3...)
	syllab["kana1"] = append(hira1, kata1...)
	syllab["kana"] = append(syllab["hira"], syllab["kata"]...)
	syllab["1-1"] = kanji1a
	syllab["1-2"] = kanji1b
	syllab["1-3"] = kanji1c
	syllab["1-4"] = kanji1d
	syllab["1-5"] = kanji1e
	syllab["1-6"] = kanji1f
	syllab["1-7"] = kanji1g
	syllab["1-8"] = kanji1h
	syllab["1-9"] = kanji1i
	syllab["1a"] = kanji1a
	syllab["1b"] = kanji1b
	syllab["1c"] = kanji1c
	syllab["1d"] = kanji1d
	syllab["1e"] = kanji1e
	syllab["1f"] = kanji1f
	syllab["1g"] = kanji1g
	syllab["1h"] = kanji1h
	syllab["1i"] = kanji1i
	syllab["kyouiku1"] = kanji1a
	syllab["kyouiku1"] = append(syllab["kyouiku1"], kanji1b...)
	syllab["kyouiku1"] = append(syllab["kyouiku1"], kanji1c...)
	syllab["kyouiku1"] = append(syllab["kyouiku1"], kanji1d...)
	syllab["kyouiku1"] = append(syllab["kyouiku1"], kanji1e...)
	syllab["kyouiku1"] = append(syllab["kyouiku1"], kanji1f...)
	syllab["kyouiku1"] = append(syllab["kyouiku1"], kanji1g...)
	syllab["kyouiku1"] = append(syllab["kyouiku1"], kanji1h...)
	syllab["kyouiku1"] = append(syllab["kyouiku1"], kanji1i...)
	syllab["kanji"] = syllab["kyouiku1"]
}

func mksyllab(syllabi ...string) {
	cur = cur[:0]
	for _, s := range syllabi {
		cur = append(cur, syllab[strings.ToLower(s)]...)
	}
	if len(cur) == 0 {
		cur = append(cur, syllab["hira1"]...)
	}
}

func prompt(i int) {
	t := chars[cur[i]]
	if k, iskanji := t.(kanji); iskanji {
		if curonkun == 0 {
			fmt.Print(string(k.rune), "の音読み ? ")
		} else {
			fmt.Print(string(k.rune), "の訓読み ? ")
		}
	} else {
		fmt.Print(cur[i], " ? ")
	}
}

func main() {
	rd := bufio.NewScanner(os.Stdin)
	mksyllab(os.Args[1:]...)
	for {
		i := rand.Intn(len(cur))
		curonkun := rand.Int() & 1
		t := chars[cur[i]]
		k, iskanji := t.(kanji)
		var onkun []string
		if iskanji {
			// Assume that no kanji has neither on'yomi nor kun'yomi.
			if curonkun == 0 && len(k.onyomi) == 0 {
				curonkun = 1
			} else if curonkun == 1 && len(k.kunyomi) == 0 {
				curonkun = 0
			}
			if curonkun == 0 {
				onkun = k.onyomi
			} else {
				onkun = k.kunyomi
			}
		}
		prompt(i)
	iloop:
		for {
			if rd.Scan() {
				g := strings.Fields(strings.ToLower(rd.Text()))
				if len(g) == 0 {
					prompt(i)
					continue
				}
				switch g[0] {
				case "q", "quit", "end":
					return
				case "?":
					if iskanji {
						if curonkun == 0 {
							fmt.Println("音読み:", strings.Join(k.onyomi, ", "), "=", k.meaning)
						} else {
							fmt.Println("訓読み:", strings.Join(k.kunyomi, ", "), "=", k.meaning)
						}
					} else {
						fmt.Println(t)
					}
					break iloop
				case "w", "what":
					if iskanji {
						fmt.Println(k.meaning)
					} else {
						fmt.Println("you call yourself a weeb?")
					}
				case "s", "syllab":
					mksyllab(g[1:]...)
					break iloop
				case "h", "help":
					fmt.Println("guess the kana or kanji. commands:")
					fmt.Println("q\tquit")
					fmt.Println("?\tgive up")
					fmt.Println("w\task for meaning, if kanji")
					fmt.Println("s\tchange syllabus, e.g. \"s kata1 kata2\" changes to katakana without glides")
					fmt.Println("available syllabi:")
					for s, l := range syllab {
						fmt.Printf("\t%s (%d characters)\n", s, len(l))
					}
				default:
					if iskanji {
						ok := false
						for _, v := range onkun {
							if g[0] == v {
								ok = true
								break
							}
						}
						if ok {
							fmt.Println("Correct!")
							if curonkun == 0 {
								fmt.Println("音読み:", strings.Join(onkun, ", "))
							} else {
								fmt.Println("訓読み:", strings.Join(onkun, ", "))
							}
							fmt.Println("Meaning:", k.meaning)
							break iloop
						} else {
							fmt.Println("Incorrect.")
						}
					} else {
						if g[0] == t.(string) {
							fmt.Println("Correct!")
							break iloop
						} else {
							fmt.Println("Incorrect.")
						}
					}
				}
				prompt(i)
			}
		}
	}
}
