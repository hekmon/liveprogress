package colors

import (
	"github.com/hekmon/liveprogress/v2"
	"github.com/muesli/termenv"
)

/*
	Helpers for easy color styles usage in liveprogress
*/

func init() {
	// oportunistic init (default liveprogress.Output value)
	Generate()
}

var (
	NoColor termenv.Style
)

// Generate (re)generates all the styles with the current terminal profile.
// Call this function after liveprogress.Start() if you have changed the default liveprogress.Output value, otherwise no need to call it.
func Generate() {
	NoColor = liveprogress.BaseStyle()
	generateANSIBasic()
	generateANSIExtended()
	generateANSIExtendedGreyscale()
}

/*
	Basic ANSI colors 0 - 15
*/

var (
	ANSIBasicBlack         termenv.Style
	ANSIBasicRed           termenv.Style
	ANSIBasicGreen         termenv.Style
	ANSIBasicYellow        termenv.Style
	ANSIBasicBlue          termenv.Style
	ANSIBasicMagenta       termenv.Style
	ANSIBasicCyan          termenv.Style
	ANSIBasicWhite         termenv.Style
	ANSIBasicBrightBlack   termenv.Style
	ANSIBasicBrightRed     termenv.Style
	ANSIBasicBrightGreen   termenv.Style
	ANSIBasicBrightYellow  termenv.Style
	ANSIBasicBrightBlue    termenv.Style
	ANSIBasicBrightMagenta termenv.Style
	ANSIBasicBrightCyan    termenv.Style
	ANSIBasicBrightWhite   termenv.Style
)

func generateANSIBasic() {
	ANSIBasicBlack = liveprogress.BaseStyle().Foreground(termenv.ANSIBlack)
	ANSIBasicRed = liveprogress.BaseStyle().Foreground(termenv.ANSIRed)
	ANSIBasicGreen = liveprogress.BaseStyle().Foreground(termenv.ANSIGreen)
	ANSIBasicYellow = liveprogress.BaseStyle().Foreground(termenv.ANSIYellow)
	ANSIBasicBlue = liveprogress.BaseStyle().Foreground(termenv.ANSIBlue)
	ANSIBasicMagenta = liveprogress.BaseStyle().Foreground(termenv.ANSIMagenta)
	ANSIBasicCyan = liveprogress.BaseStyle().Foreground(termenv.ANSICyan)
	ANSIBasicWhite = liveprogress.BaseStyle().Foreground(termenv.ANSIWhite)
	ANSIBasicBrightBlack = liveprogress.BaseStyle().Foreground(termenv.ANSIBrightBlack)
	ANSIBasicBrightRed = liveprogress.BaseStyle().Foreground(termenv.ANSIBrightRed)
	ANSIBasicBrightGreen = liveprogress.BaseStyle().Foreground(termenv.ANSIBrightGreen)
	ANSIBasicBrightYellow = liveprogress.BaseStyle().Foreground(termenv.ANSIBrightYellow)
	ANSIBasicBrightBlue = liveprogress.BaseStyle().Foreground(termenv.ANSIBrightBlue)
	ANSIBasicBrightMagenta = liveprogress.BaseStyle().Foreground(termenv.ANSIBrightMagenta)
	ANSIBasicBrightCyan = liveprogress.BaseStyle().Foreground(termenv.ANSIBrightCyan)
	ANSIBasicBrightWhite = liveprogress.BaseStyle().Foreground(termenv.ANSIBrightWhite)
}

/*
	Extended ANSI colors 16-231
*/

var (
	// See termenv chart here : https://github.com/muesli/termenv?tab=readme-ov-file#color-chart
	// Or open the file in vscode with the following extension installed: https://marketplace.visualstudio.com/items?itemName=naumovs.color-highlight
	ANSIExtended16  termenv.Style // #000000
	ANSIExtended17  termenv.Style // #00005f
	ANSIExtended18  termenv.Style // #000087
	ANSIExtended19  termenv.Style // #0000af
	ANSIExtended20  termenv.Style // #0000d7
	ANSIExtended21  termenv.Style // #0000ff
	ANSIExtended22  termenv.Style // #005f00
	ANSIExtended23  termenv.Style // #005f5f
	ANSIExtended24  termenv.Style // #005f87
	ANSIExtended25  termenv.Style // #005faf
	ANSIExtended26  termenv.Style // #005fd7
	ANSIExtended27  termenv.Style // #005fff
	ANSIExtended28  termenv.Style // #008700
	ANSIExtended29  termenv.Style // #00875f
	ANSIExtended30  termenv.Style // #008787
	ANSIExtended31  termenv.Style // #0087af
	ANSIExtended32  termenv.Style // #0087d7
	ANSIExtended33  termenv.Style // #0087ff
	ANSIExtended34  termenv.Style // #00af00
	ANSIExtended35  termenv.Style // #00af5f
	ANSIExtended36  termenv.Style // #00af87
	ANSIExtended37  termenv.Style // #00afaf
	ANSIExtended38  termenv.Style // #00afd7
	ANSIExtended39  termenv.Style // #00afff
	ANSIExtended40  termenv.Style // #00d700
	ANSIExtended41  termenv.Style // #00d75f
	ANSIExtended42  termenv.Style // #00d787
	ANSIExtended43  termenv.Style // #00d7af
	ANSIExtended44  termenv.Style // #00d7d7
	ANSIExtended45  termenv.Style // #00d7ff
	ANSIExtended46  termenv.Style // #00ff00
	ANSIExtended47  termenv.Style // #00ff5f
	ANSIExtended48  termenv.Style // #00ff87
	ANSIExtended49  termenv.Style // #00ffaf
	ANSIExtended50  termenv.Style // #00ffd7
	ANSIExtended51  termenv.Style // #00ffff
	ANSIExtended52  termenv.Style // #5f0000
	ANSIExtended53  termenv.Style // #5f005f
	ANSIExtended54  termenv.Style // #5f0087
	ANSIExtended55  termenv.Style // #5f00af
	ANSIExtended56  termenv.Style // #5f00d7
	ANSIExtended57  termenv.Style // #5f00ff
	ANSIExtended58  termenv.Style // #5f5f00
	ANSIExtended59  termenv.Style // #5f5f5f
	ANSIExtended60  termenv.Style // #5f5f87
	ANSIExtended61  termenv.Style // #5f5faf
	ANSIExtended62  termenv.Style // #5f5fd7
	ANSIExtended63  termenv.Style // #5f5fff
	ANSIExtended64  termenv.Style // #5f8700
	ANSIExtended65  termenv.Style // #5f875f
	ANSIExtended66  termenv.Style // #5f8787
	ANSIExtended67  termenv.Style // #5f87af
	ANSIExtended68  termenv.Style // #5f87d7
	ANSIExtended69  termenv.Style // #5f87ff
	ANSIExtended70  termenv.Style // #5faf00
	ANSIExtended71  termenv.Style // #5faf5f
	ANSIExtended72  termenv.Style // #5faf87
	ANSIExtended73  termenv.Style // #5fafaf
	ANSIExtended74  termenv.Style // #5fafd7
	ANSIExtended75  termenv.Style // #5fafff
	ANSIExtended76  termenv.Style // #5fd700
	ANSIExtended77  termenv.Style // #5fd75f
	ANSIExtended78  termenv.Style // #5fd787
	ANSIExtended79  termenv.Style // #5fd7af
	ANSIExtended80  termenv.Style // #5fd7d7
	ANSIExtended81  termenv.Style // #5fd7ff
	ANSIExtended82  termenv.Style // #5fff00
	ANSIExtended83  termenv.Style // #5fff5f
	ANSIExtended84  termenv.Style // #5fff87
	ANSIExtended85  termenv.Style // #5fffaf
	ANSIExtended86  termenv.Style // #5fffd7
	ANSIExtended87  termenv.Style // #5fffff
	ANSIExtended88  termenv.Style // #870000
	ANSIExtended89  termenv.Style // #87005f
	ANSIExtended90  termenv.Style // #870087
	ANSIExtended91  termenv.Style // #8700af
	ANSIExtended92  termenv.Style // #8700d7
	ANSIExtended93  termenv.Style // #8700ff
	ANSIExtended94  termenv.Style // #875f00
	ANSIExtended95  termenv.Style // #875f5f
	ANSIExtended96  termenv.Style // #875f87
	ANSIExtended97  termenv.Style // #875faf
	ANSIExtended98  termenv.Style // #875fd7
	ANSIExtended99  termenv.Style // #875fff
	ANSIExtended100 termenv.Style // #878700
	ANSIExtended101 termenv.Style // #87875f
	ANSIExtended102 termenv.Style // #878787
	ANSIExtended103 termenv.Style // #8787af
	ANSIExtended104 termenv.Style // #8787d7
	ANSIExtended105 termenv.Style // #8787ff
	ANSIExtended106 termenv.Style // #87af00
	ANSIExtended107 termenv.Style // #87af5f
	ANSIExtended108 termenv.Style // #87af87
	ANSIExtended109 termenv.Style // #87afaf
	ANSIExtended110 termenv.Style // #87afd7
	ANSIExtended111 termenv.Style // #87afff
	ANSIExtended112 termenv.Style // #87d700
	ANSIExtended113 termenv.Style // #87d75f
	ANSIExtended114 termenv.Style // #87d787
	ANSIExtended115 termenv.Style // #87d7af
	ANSIExtended116 termenv.Style // #87d7d7
	ANSIExtended117 termenv.Style // #87d7ff
	ANSIExtended118 termenv.Style // #87ff00
	ANSIExtended119 termenv.Style // #87ff5f
	ANSIExtended120 termenv.Style // #87ff87
	ANSIExtended121 termenv.Style // #87ffaf
	ANSIExtended122 termenv.Style // #87ffd7
	ANSIExtended123 termenv.Style // #87ffff
	ANSIExtended124 termenv.Style // #af0000
	ANSIExtended125 termenv.Style // #af005f
	ANSIExtended126 termenv.Style // #af0087
	ANSIExtended127 termenv.Style // #af00af
	ANSIExtended128 termenv.Style // #af00d7
	ANSIExtended129 termenv.Style // #af00ff
	ANSIExtended130 termenv.Style // #af5f00
	ANSIExtended131 termenv.Style // #af5f5f
	ANSIExtended132 termenv.Style // #af5f87
	ANSIExtended133 termenv.Style // #af5faf
	ANSIExtended134 termenv.Style // #af5fd7
	ANSIExtended135 termenv.Style // #af5fff
	ANSIExtended136 termenv.Style // #af8700
	ANSIExtended137 termenv.Style // #af875f
	ANSIExtended138 termenv.Style // #af8787
	ANSIExtended139 termenv.Style // #af87af
	ANSIExtended140 termenv.Style // #af87d7
	ANSIExtended141 termenv.Style // #af87ff
	ANSIExtended142 termenv.Style // #afaf00
	ANSIExtended143 termenv.Style // #afaf5f
	ANSIExtended144 termenv.Style // #afaf87
	ANSIExtended145 termenv.Style // #afafaf
	ANSIExtended146 termenv.Style // #afafd7
	ANSIExtended147 termenv.Style // #afafff
	ANSIExtended148 termenv.Style // #afd700
	ANSIExtended149 termenv.Style // #afd75f
	ANSIExtended150 termenv.Style // #afd787
	ANSIExtended151 termenv.Style // #afd7af
	ANSIExtended152 termenv.Style // #afd7d7
	ANSIExtended153 termenv.Style // #afd7ff
	ANSIExtended154 termenv.Style // #afff00
	ANSIExtended155 termenv.Style // #afff5f
	ANSIExtended156 termenv.Style // #afff87
	ANSIExtended157 termenv.Style // #afffaf
	ANSIExtended158 termenv.Style // #afffd7
	ANSIExtended159 termenv.Style // #afffff
	ANSIExtended160 termenv.Style // #d70000
	ANSIExtended161 termenv.Style // #d7005f
	ANSIExtended162 termenv.Style // #d70087
	ANSIExtended163 termenv.Style // #d700af
	ANSIExtended164 termenv.Style // #d700d7
	ANSIExtended165 termenv.Style // #d700ff
	ANSIExtended166 termenv.Style // #d75f00
	ANSIExtended167 termenv.Style // #d75f5f
	ANSIExtended168 termenv.Style // #d75f87
	ANSIExtended169 termenv.Style // #d75faf
	ANSIExtended170 termenv.Style // #d75fd7
	ANSIExtended171 termenv.Style // #d75fff
	ANSIExtended172 termenv.Style // #d78700
	ANSIExtended173 termenv.Style // #d7875f
	ANSIExtended174 termenv.Style // #d78787
	ANSIExtended175 termenv.Style // #d787af
	ANSIExtended176 termenv.Style // #d787d7
	ANSIExtended177 termenv.Style // #d787ff
	ANSIExtended178 termenv.Style // #d7af00
	ANSIExtended179 termenv.Style // #d7af5f
	ANSIExtended180 termenv.Style // #d7af87
	ANSIExtended181 termenv.Style // #d7afaf
	ANSIExtended182 termenv.Style // #d7afd7
	ANSIExtended183 termenv.Style // #d7afff
	ANSIExtended184 termenv.Style // #d7d700
	ANSIExtended185 termenv.Style // #d7d75f
	ANSIExtended186 termenv.Style // #d7d787
	ANSIExtended187 termenv.Style // #d7d7af
	ANSIExtended188 termenv.Style // #d7d7d7
	ANSIExtended189 termenv.Style // #d7d7ff
	ANSIExtended190 termenv.Style // #d7ff00
	ANSIExtended191 termenv.Style // #d7ff5f
	ANSIExtended192 termenv.Style // #d7ff87
	ANSIExtended193 termenv.Style // #d7ffaf
	ANSIExtended194 termenv.Style // #d7ffd7
	ANSIExtended195 termenv.Style // #d7ffff
	ANSIExtended196 termenv.Style // #ff0000
	ANSIExtended197 termenv.Style // #ff005f
	ANSIExtended198 termenv.Style // #ff0087
	ANSIExtended199 termenv.Style // #ff00af
	ANSIExtended200 termenv.Style // #ff00d7
	ANSIExtended201 termenv.Style // #ff00ff
	ANSIExtended202 termenv.Style // #ff5f00
	ANSIExtended203 termenv.Style // #ff5f5f
	ANSIExtended204 termenv.Style // #ff5f87
	ANSIExtended205 termenv.Style // #ff5faf
	ANSIExtended206 termenv.Style // #ff5fd7
	ANSIExtended207 termenv.Style // #ff5fff
	ANSIExtended208 termenv.Style // #ff8700
	ANSIExtended209 termenv.Style // #ff875f
	ANSIExtended210 termenv.Style // #ff8787
	ANSIExtended211 termenv.Style // #ff87af
	ANSIExtended212 termenv.Style // #ff87d7
	ANSIExtended213 termenv.Style // #ff87ff
	ANSIExtended214 termenv.Style // #ffaf00
	ANSIExtended215 termenv.Style // #ffaf5f
	ANSIExtended216 termenv.Style // #ffaf87
	ANSIExtended217 termenv.Style // #ffafaf
	ANSIExtended218 termenv.Style // #ffafd7
	ANSIExtended219 termenv.Style // #ffafff
	ANSIExtended220 termenv.Style // #ffd700
	ANSIExtended221 termenv.Style // #ffd75f
	ANSIExtended222 termenv.Style // #ffd787
	ANSIExtended223 termenv.Style // #ffd7af
	ANSIExtended224 termenv.Style // #ffd7d7
	ANSIExtended225 termenv.Style // #ffd7ff
	ANSIExtended226 termenv.Style // #ffff00
	ANSIExtended227 termenv.Style // #ffff5f
	ANSIExtended228 termenv.Style // #ffff87
	ANSIExtended229 termenv.Style // #ffffaf
	ANSIExtended230 termenv.Style // #ffffd7
	ANSIExtended231 termenv.Style // #ffffff
)

func generateANSIExtended() {
	ANSIExtended16 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(16))
	ANSIExtended17 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(17))
	ANSIExtended18 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(18))
	ANSIExtended19 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(19))
	ANSIExtended20 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(20))
	ANSIExtended21 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(21))
	ANSIExtended22 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(22))
	ANSIExtended23 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(23))
	ANSIExtended24 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(24))
	ANSIExtended25 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(25))
	ANSIExtended26 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(26))
	ANSIExtended27 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(27))
	ANSIExtended28 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(28))
	ANSIExtended29 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(29))
	ANSIExtended30 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(30))
	ANSIExtended31 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(31))
	ANSIExtended32 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(32))
	ANSIExtended33 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(33))
	ANSIExtended34 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(34))
	ANSIExtended35 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(35))
	ANSIExtended36 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(36))
	ANSIExtended37 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(37))
	ANSIExtended38 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(38))
	ANSIExtended39 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(39))
	ANSIExtended40 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(40))
	ANSIExtended41 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(41))
	ANSIExtended42 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(42))
	ANSIExtended43 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(43))
	ANSIExtended44 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(44))
	ANSIExtended45 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(45))
	ANSIExtended46 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(46))
	ANSIExtended47 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(47))
	ANSIExtended48 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(48))
	ANSIExtended49 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(49))
	ANSIExtended50 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(50))
	ANSIExtended51 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(51))
	ANSIExtended52 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(52))
	ANSIExtended53 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(53))
	ANSIExtended54 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(54))
	ANSIExtended55 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(55))
	ANSIExtended56 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(56))
	ANSIExtended57 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(57))
	ANSIExtended58 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(58))
	ANSIExtended59 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(59))
	ANSIExtended60 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(60))
	ANSIExtended61 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(61))
	ANSIExtended62 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(62))
	ANSIExtended63 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(63))
	ANSIExtended64 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(64))
	ANSIExtended65 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(65))
	ANSIExtended66 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(66))
	ANSIExtended67 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(67))
	ANSIExtended68 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(68))
	ANSIExtended69 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(69))
	ANSIExtended70 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(70))
	ANSIExtended71 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(71))
	ANSIExtended72 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(72))
	ANSIExtended73 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(73))
	ANSIExtended74 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(74))
	ANSIExtended75 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(75))
	ANSIExtended76 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(76))
	ANSIExtended77 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(77))
	ANSIExtended78 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(78))
	ANSIExtended79 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(79))
	ANSIExtended80 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(80))
	ANSIExtended81 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(81))
	ANSIExtended82 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(82))
	ANSIExtended83 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(83))
	ANSIExtended84 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(84))
	ANSIExtended85 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(85))
	ANSIExtended86 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(86))
	ANSIExtended87 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(87))
	ANSIExtended88 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(88))
	ANSIExtended89 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(89))
	ANSIExtended90 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(90))
	ANSIExtended91 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(91))
	ANSIExtended92 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(92))
	ANSIExtended93 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(93))
	ANSIExtended94 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(94))
	ANSIExtended95 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(95))
	ANSIExtended96 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(96))
	ANSIExtended97 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(97))
	ANSIExtended98 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(98))
	ANSIExtended99 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(99))
	ANSIExtended100 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(100))
	ANSIExtended101 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(101))
	ANSIExtended102 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(102))
	ANSIExtended103 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(103))
	ANSIExtended104 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(104))
	ANSIExtended105 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(105))
	ANSIExtended106 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(106))
	ANSIExtended107 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(107))
	ANSIExtended108 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(108))
	ANSIExtended109 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(109))
	ANSIExtended110 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(110))
	ANSIExtended111 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(111))
	ANSIExtended112 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(112))
	ANSIExtended113 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(113))
	ANSIExtended114 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(114))
	ANSIExtended115 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(115))
	ANSIExtended116 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(116))
	ANSIExtended117 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(117))
	ANSIExtended118 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(118))
	ANSIExtended119 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(119))
	ANSIExtended120 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(120))
	ANSIExtended121 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(121))
	ANSIExtended122 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(122))
	ANSIExtended123 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(123))
	ANSIExtended124 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(124))
	ANSIExtended125 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(125))
	ANSIExtended126 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(126))
	ANSIExtended127 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(127))
	ANSIExtended128 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(128))
	ANSIExtended129 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(129))
	ANSIExtended130 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(130))
	ANSIExtended131 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(131))
	ANSIExtended132 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(132))
	ANSIExtended133 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(133))
	ANSIExtended134 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(134))
	ANSIExtended135 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(135))
	ANSIExtended136 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(136))
	ANSIExtended137 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(137))
	ANSIExtended138 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(138))
	ANSIExtended139 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(139))
	ANSIExtended140 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(140))
	ANSIExtended141 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(141))
	ANSIExtended142 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(142))
	ANSIExtended143 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(143))
	ANSIExtended144 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(144))
	ANSIExtended145 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(145))
	ANSIExtended146 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(146))
	ANSIExtended147 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(147))
	ANSIExtended148 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(148))
	ANSIExtended149 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(149))
	ANSIExtended150 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(150))
	ANSIExtended151 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(151))
	ANSIExtended152 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(152))
	ANSIExtended153 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(153))
	ANSIExtended154 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(154))
	ANSIExtended155 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(155))
	ANSIExtended156 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(156))
	ANSIExtended157 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(157))
	ANSIExtended158 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(158))
	ANSIExtended159 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(159))
	ANSIExtended160 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(160))
	ANSIExtended161 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(161))
	ANSIExtended162 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(162))
	ANSIExtended163 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(163))
	ANSIExtended164 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(164))
	ANSIExtended165 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(165))
	ANSIExtended166 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(166))
	ANSIExtended167 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(167))
	ANSIExtended168 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(168))
	ANSIExtended169 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(169))
	ANSIExtended170 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(170))
	ANSIExtended171 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(171))
	ANSIExtended172 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(172))
	ANSIExtended173 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(173))
	ANSIExtended174 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(174))
	ANSIExtended175 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(175))
	ANSIExtended176 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(176))
	ANSIExtended177 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(177))
	ANSIExtended178 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(178))
	ANSIExtended179 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(179))
	ANSIExtended180 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(180))
	ANSIExtended181 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(181))
	ANSIExtended182 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(182))
	ANSIExtended183 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(183))
	ANSIExtended184 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(184))
	ANSIExtended185 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(185))
	ANSIExtended186 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(186))
	ANSIExtended187 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(187))
	ANSIExtended188 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(188))
	ANSIExtended189 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(189))
	ANSIExtended190 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(190))
	ANSIExtended191 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(191))
	ANSIExtended192 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(192))
	ANSIExtended193 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(193))
	ANSIExtended194 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(194))
	ANSIExtended195 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(195))
	ANSIExtended196 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(196))
	ANSIExtended197 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(197))
	ANSIExtended198 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(198))
	ANSIExtended199 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(199))
	ANSIExtended200 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(200))
	ANSIExtended201 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(201))
	ANSIExtended202 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(202))
	ANSIExtended203 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(203))
	ANSIExtended204 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(204))
	ANSIExtended205 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(205))
	ANSIExtended206 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(206))
	ANSIExtended207 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(207))
	ANSIExtended208 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(208))
	ANSIExtended209 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(209))
	ANSIExtended210 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(210))
	ANSIExtended211 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(211))
	ANSIExtended212 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(212))
	ANSIExtended213 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(213))
	ANSIExtended214 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(214))
	ANSIExtended215 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(215))
	ANSIExtended216 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(216))
	ANSIExtended217 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(217))
	ANSIExtended218 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(218))
	ANSIExtended219 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(219))
	ANSIExtended220 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(220))
	ANSIExtended221 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(221))
	ANSIExtended222 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(222))
	ANSIExtended223 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(223))
	ANSIExtended224 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(224))
	ANSIExtended225 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(225))
	ANSIExtended226 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(226))
	ANSIExtended227 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(227))
	ANSIExtended228 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(228))
	ANSIExtended229 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(229))
	ANSIExtended230 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(230))
	ANSIExtended231 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(231))
}

/*
	Grayscale ANSI colors 232-255
*/

var (
	// See termenv chart here : https://github.com/muesli/termenv?tab=readme-ov-file#color-chart
	// Or open the file in vscode with the following extension installed: https://marketplace.visualstudio.com/items?itemName=naumovs.color-highlight
	ANSIExtendedGrayscale232 termenv.Style // #080808
	ANSIExtendedGrayscale233 termenv.Style // #121212
	ANSIExtendedGrayscale234 termenv.Style // #1c1c1c
	ANSIExtendedGrayscale235 termenv.Style // #262626
	ANSIExtendedGrayscale236 termenv.Style // #303030
	ANSIExtendedGrayscale237 termenv.Style // #3a3a3a
	ANSIExtendedGrayscale238 termenv.Style // #444444
	ANSIExtendedGrayscale239 termenv.Style // #4e4e4e
	ANSIExtendedGrayscale240 termenv.Style // #585858
	ANSIExtendedGrayscale241 termenv.Style // #626262
	ANSIExtendedGrayscale242 termenv.Style // #6c6c6c
	ANSIExtendedGrayscale243 termenv.Style // #767676
	ANSIExtendedGrayscale244 termenv.Style // #808080
	ANSIExtendedGrayscale245 termenv.Style // #8a8a8a
	ANSIExtendedGrayscale246 termenv.Style // #949494
	ANSIExtendedGrayscale247 termenv.Style // #9e9e9e
	ANSIExtendedGrayscale248 termenv.Style // #a8a8a8
	ANSIExtendedGrayscale249 termenv.Style // #b2b2b2
	ANSIExtendedGrayscale250 termenv.Style // #bcbcbc
	ANSIExtendedGrayscale251 termenv.Style // #c6c6c6
	ANSIExtendedGrayscale252 termenv.Style // #d0d0d0
	ANSIExtendedGrayscale253 termenv.Style // #dadada
	ANSIExtendedGrayscale254 termenv.Style // #e4e4e4
	ANSIExtendedGrayscale255 termenv.Style // #eeeeee
)

func generateANSIExtendedGreyscale() {
	ANSIExtendedGrayscale232 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(232))
	ANSIExtendedGrayscale233 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(233))
	ANSIExtendedGrayscale234 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(234))
	ANSIExtendedGrayscale235 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(235))
	ANSIExtendedGrayscale236 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(236))
	ANSIExtendedGrayscale237 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(237))
	ANSIExtendedGrayscale238 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(238))
	ANSIExtendedGrayscale239 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(239))
	ANSIExtendedGrayscale240 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(240))
	ANSIExtendedGrayscale241 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(241))
	ANSIExtendedGrayscale242 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(242))
	ANSIExtendedGrayscale243 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(243))
	ANSIExtendedGrayscale244 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(244))
	ANSIExtendedGrayscale245 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(245))
	ANSIExtendedGrayscale246 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(246))
	ANSIExtendedGrayscale247 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(247))
	ANSIExtendedGrayscale248 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(248))
	ANSIExtendedGrayscale249 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(249))
	ANSIExtendedGrayscale250 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(250))
	ANSIExtendedGrayscale251 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(251))
	ANSIExtendedGrayscale252 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(252))
	ANSIExtendedGrayscale253 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(253))
	ANSIExtendedGrayscale254 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(254))
	ANSIExtendedGrayscale255 = liveprogress.BaseStyle().Foreground(termenv.ANSI256Color(255))
}

/*
	TrueColor (RGB)
*/

// RGB generates a style with a foreground color set to rgb. rgb is a hex-encoded color, e.g. "#abcdef".
// You can call this function before liveprogress.Start() if you did not changed liveprogress.Output, otherwise you should call it after liveprogress.Start().
func RGB(rgb string) termenv.Style {
	return liveprogress.BaseStyle().Foreground(termenv.RGBColor(rgb))
}
