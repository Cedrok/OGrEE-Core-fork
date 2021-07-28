// Code generated by goyacc - DO NOT EDIT.

package main

import __yyfmt__ "fmt"

import (
	cmd "cli/controllers"
	"strings"
)

func resMap(x *string) map[string]interface{} {
	resarr := strings.Split(*x, "=")
	res := make(map[string]interface{})
	attrs := make(map[string]string)

	for i := 0; i+1 < len(resarr); {
		if i+1 < len(resarr) {
			switch resarr[i] {
			case "id", "name", "category", "parentID",
				"description", "domain", "parentid", "parentId":
				res[resarr[i]] = resarr[i+1]

			default:
				attrs[resarr[i]] = resarr[i+1]
			}
			i += 2
		}
	}
	res["attributes"] = attrs
	return res
}

type yySymType struct {
	yys int
	//n int
	s string
}

type yyXError struct {
	state, xsym int
}

const (
	yyDefault        = 57372
	yyEofCode        = 57344
	TOKEN_ATTR       = 57355
	TOKEN_BASHTYPE   = 57361
	TOKEN_BLDG       = 57349
	TOKEN_CD         = 57367
	TOKEN_CLR        = 57369
	TOKEN_CMDFLAG    = 57363
	TOKEN_CREATE     = 57356
	TOKEN_DELETE     = 57359
	TOKEN_DEVICE     = 57352
	TOKEN_DOC        = 57366
	TOKEN_EQUAL      = 57362
	TOKEN_EXIT       = 57365
	TOKEN_GET        = 57357
	TOKEN_GREP       = 57370
	TOKEN_LS         = 57371
	TOKEN_PWD        = 57368
	TOKEN_RACK       = 57351
	TOKEN_ROOM       = 57350
	TOKEN_SEARCH     = 57360
	TOKEN_SITE       = 57348
	TOKEN_SLASH      = 57364
	TOKEN_SUBDEVICE  = 57353
	TOKEN_SUBDEVICE1 = 57354
	TOKEN_TENANT     = 57347
	TOKEN_UPDATE     = 57358
	TOKEN_WORD       = 57346
	yyErrCode        = 57345

	yyMaxDepth = 200
	yyTabOfs   = -44
)

var (
	yyPrec = map[int]int{}

	yyXLAT = map[int]int{
		57344: 0,  // $end (36x)
		57346: 1,  // TOKEN_WORD (26x)
		57355: 2,  // TOKEN_ATTR (19x)
		57375: 3,  // F (8x)
		57381: 4,  // P (7x)
		57374: 5,  // E (4x)
		57349: 6,  // TOKEN_BLDG (4x)
		57352: 7,  // TOKEN_DEVICE (4x)
		57351: 8,  // TOKEN_RACK (4x)
		57350: 9,  // TOKEN_ROOM (4x)
		57348: 10, // TOKEN_SITE (4x)
		57353: 11, // TOKEN_SUBDEVICE (4x)
		57354: 12, // TOKEN_SUBDEVICE1 (4x)
		57347: 13, // TOKEN_TENANT (4x)
		57364: 14, // TOKEN_SLASH (3x)
		57363: 15, // TOKEN_CMDFLAG (2x)
		57373: 16, // BASH (1x)
		57376: 17, // K (1x)
		57377: 18, // NT_CREATE (1x)
		57378: 19, // NT_DEL (1x)
		57379: 20, // NT_GET (1x)
		57380: 21, // NT_UPDATE (1x)
		57382: 22, // Q (1x)
		57383: 23, // start (1x)
		57367: 24, // TOKEN_CD (1x)
		57369: 25, // TOKEN_CLR (1x)
		57356: 26, // TOKEN_CREATE (1x)
		57359: 27, // TOKEN_DELETE (1x)
		57366: 28, // TOKEN_DOC (1x)
		57362: 29, // TOKEN_EQUAL (1x)
		57365: 30, // TOKEN_EXIT (1x)
		57357: 31, // TOKEN_GET (1x)
		57370: 32, // TOKEN_GREP (1x)
		57371: 33, // TOKEN_LS (1x)
		57368: 34, // TOKEN_PWD (1x)
		57358: 35, // TOKEN_UPDATE (1x)
		57372: 36, // $default (0x)
		57345: 37, // error (0x)
		57361: 38, // TOKEN_BASHTYPE (0x)
		57360: 39, // TOKEN_SEARCH (0x)
	}

	yySymNames = []string{
		"$end",
		"TOKEN_WORD",
		"TOKEN_ATTR",
		"F",
		"P",
		"E",
		"TOKEN_BLDG",
		"TOKEN_DEVICE",
		"TOKEN_RACK",
		"TOKEN_ROOM",
		"TOKEN_SITE",
		"TOKEN_SUBDEVICE",
		"TOKEN_SUBDEVICE1",
		"TOKEN_TENANT",
		"TOKEN_SLASH",
		"TOKEN_CMDFLAG",
		"BASH",
		"K",
		"NT_CREATE",
		"NT_DEL",
		"NT_GET",
		"NT_UPDATE",
		"Q",
		"start",
		"TOKEN_CD",
		"TOKEN_CLR",
		"TOKEN_CREATE",
		"TOKEN_DELETE",
		"TOKEN_DOC",
		"TOKEN_EQUAL",
		"TOKEN_EXIT",
		"TOKEN_GET",
		"TOKEN_GREP",
		"TOKEN_LS",
		"TOKEN_PWD",
		"TOKEN_UPDATE",
		"$default",
		"error",
		"TOKEN_BASHTYPE",
		"TOKEN_SEARCH",
	}

	yyTokenLiteralStrings = map[int]string{}

	yyReductions = map[int]struct{ xsym, components int }{
		0:  {0, 1},
		1:  {23, 1},
		2:  {23, 1},
		3:  {17, 1},
		4:  {17, 1},
		5:  {17, 1},
		6:  {17, 1},
		7:  {18, 3},
		8:  {18, 4},
		9:  {20, 1},
		10: {20, 2},
		11: {20, 3},
		12: {21, 3},
		13: {21, 4},
		14: {19, 3},
		15: {19, 4},
		16: {5, 1},
		17: {5, 1},
		18: {5, 1},
		19: {5, 1},
		20: {5, 1},
		21: {5, 1},
		22: {5, 1},
		23: {5, 1},
		24: {3, 4},
		25: {3, 3},
		26: {4, 3},
		27: {4, 1},
		28: {4, 2},
		29: {22, 3},
		30: {22, 2},
		31: {22, 2},
		32: {22, 2},
		33: {22, 3},
		34: {22, 2},
		35: {22, 1},
		36: {16, 1},
		37: {16, 1},
		38: {16, 1},
		39: {16, 1},
		40: {16, 2},
		41: {16, 1},
		42: {16, 1},
		43: {16, 1},
	}

	yyXErrors = map[yyXError]string{}

	yyParseTab = [57][]uint8{
		// 0
		{16: 58, 46, 48, 51, 49, 50, 47, 45, 56, 59, 52, 55, 63, 30: 62, 53, 60, 57, 61, 54},
		{44},
		{43},
		{42},
		{41},
		// 5
		{40},
		{39},
		{38},
		{5: 97, 77, 80, 79, 78, 76, 81, 82, 75},
		{35, 69, 4: 94, 95, 77, 80, 79, 78, 76, 81, 82, 75},
		// 10
		{5: 90, 77, 80, 79, 78, 76, 81, 82, 75},
		{5: 74, 77, 80, 79, 78, 76, 81, 82, 75},
		{8, 71, 4: 72},
		{5, 65, 4: 66},
		{9, 64},
		// 15
		{7, 7},
		{6, 6},
		{3, 3},
		{2, 2},
		{1, 1},
		// 20
		{10},
		{17, 17, 14: 67, 68},
		{12, 4},
		{16, 69, 16, 4: 70},
		{11},
		// 25
		{17, 17, 17, 14: 67},
		{18, 18, 18},
		{17, 14: 67, 73},
		{13},
		{15},
		// 30
		{1: 69, 85, 83, 84},
		{1: 28, 28},
		{1: 27, 27},
		{1: 26, 26},
		{1: 25, 25},
		// 35
		{1: 24, 24},
		{1: 23, 23},
		{1: 22, 22},
		{1: 21, 21},
		{30},
		// 40
		{2: 85, 89},
		{29: 86},
		{1: 87},
		{19, 2: 85, 88},
		{20},
		// 45
		{29},
		{1: 69, 85, 91, 92},
		{32},
		{2: 85, 93},
		{31},
		// 50
		{34},
		{2: 85, 96},
		{33},
		{1: 69, 85, 98, 99},
		{37},
		// 55
		{2: 85, 100},
		{36},
	}
)

var yyDebug = 0

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyLexerEx interface {
	yyLexer
	Reduced(rule, state int, lval *yySymType) bool
}

func yySymName(c int) (s string) {
	x, ok := yyXLAT[c]
	if ok {
		return yySymNames[x]
	}

	if c < 0x7f {
		return __yyfmt__.Sprintf("%q", c)
	}

	return __yyfmt__.Sprintf("%d", c)
}

func yylex1(yylex yyLexer, lval *yySymType) (n int) {
	n = yylex.Lex(lval)
	if n <= 0 {
		n = yyEofCode
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("\nlex %s(%#x %d), lval: %+v\n", yySymName(n), n, n, lval)
	}
	return n
}

func yyParse(yylex yyLexer) int {
	const yyError = 37

	yyEx, _ := yylex.(yyLexerEx)
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	yyS := make([]yySymType, 200)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yyerrok := func() {
		if yyDebug >= 2 {
			__yyfmt__.Printf("yyerrok()\n")
		}
		Errflag = 0
	}
	_ = yyerrok
	yystate := 0
	yychar := -1
	var yyxchar int
	var yyshift int
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	if yychar < 0 {
		yylval.yys = yystate
		yychar = yylex1(yylex, &yylval)
		var ok bool
		if yyxchar, ok = yyXLAT[yychar]; !ok {
			yyxchar = len(yySymNames) // > tab width
		}
	}
	if yyDebug >= 4 {
		var a []int
		for _, v := range yyS[:yyp+1] {
			a = append(a, v.yys)
		}
		__yyfmt__.Printf("state stack %v\n", a)
	}
	row := yyParseTab[yystate]
	yyn = 0
	if yyxchar < len(row) {
		if yyn = int(row[yyxchar]); yyn != 0 {
			yyn += yyTabOfs
		}
	}
	switch {
	case yyn > 0: // shift
		yychar = -1
		yyVAL = yylval
		yystate = yyn
		yyshift = yyn
		if yyDebug >= 2 {
			__yyfmt__.Printf("shift, and goto state %d\n", yystate)
		}
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	case yyn < 0: // reduce
	case yystate == 1: // accept
		if yyDebug >= 2 {
			__yyfmt__.Println("accept")
		}
		goto ret0
	}

	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			if yyDebug >= 1 {
				__yyfmt__.Printf("no action for %s in state %d\n", yySymName(yychar), yystate)
			}
			msg, ok := yyXErrors[yyXError{yystate, yyxchar}]
			if !ok {
				msg, ok = yyXErrors[yyXError{yystate, -1}]
			}
			if !ok && yyshift != 0 {
				msg, ok = yyXErrors[yyXError{yyshift, yyxchar}]
			}
			if !ok {
				msg, ok = yyXErrors[yyXError{yyshift, -1}]
			}
			if yychar > 0 {
				ls := yyTokenLiteralStrings[yychar]
				if ls == "" {
					ls = yySymName(yychar)
				}
				if ls != "" {
					switch {
					case msg == "":
						msg = __yyfmt__.Sprintf("unexpected %s", ls)
					default:
						msg = __yyfmt__.Sprintf("unexpected %s, %s", ls, msg)
					}
				}
			}
			if msg == "" {
				msg = "syntax error"
			}
			println("OGREE: Unrecognised command!")
			/*yylex.Error(msg)*/
			Nerrs++
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				row := yyParseTab[yyS[yyp].yys]
				if yyError < len(row) {
					yyn = int(row[yyError]) + yyTabOfs
					if yyn > 0 { // hit
						if yyDebug >= 2 {
							__yyfmt__.Printf("error recovery found error shift in state %d\n", yyS[yyp].yys)
						}
						yystate = yyn /* simulate a shift of "error" */
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery failed\n")
			}
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yySymName(yychar))
			}
			if yychar == yyEofCode {
				goto ret1
			}

			yychar = -1
			goto yynewstate /* try again in the same state */
		}
	}

	r := -yyn
	x0 := yyReductions[r]
	x, n := x0.xsym, x0.components
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= n
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	exState := yystate
	yystate = int(yyParseTab[yyS[yyp].yys][x]) + yyTabOfs
	/* reduction by production r */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce using rule %v (%s), and goto state %d\n", r, yySymNames[x], yystate)
	}

	switch r {
	case 1:
		{
			println("@State start")
		}
	case 7:
		{
			cmd.PostObj(yyS[yypt-1].s, "", resMap(&yyS[yypt-0].s)) /*println("@State NT_CR");*/
		}
	case 8:
		{
			yyVAL.s = yyS[yypt-0].s /*println("Finally: "+$$);*/
			cmd.Disp(resMap(&yyS[yypt-0].s))
			cmd.PostObj(yyS[yypt-2].s, yyS[yypt-1].s, resMap(&yyS[yypt-0].s))
		}
	case 9:
		{
			println("@State NT_GET")
			cmd.GetObject("")
		}
	case 10:
		{
			cmd.GetObject(yyS[yypt-0].s)
		}
	case 11:
		{ /*cmd.Disp(resMap(&$4)); */
			cmd.SearchObjects(yyS[yypt-1].s, resMap(&yyS[yypt-0].s))
		}
	case 12:
		{
			println("@State NT_UPD")
		}
	case 13:
		{
			yyVAL.s = yyS[yypt-0].s
			cmd.Disp(resMap(&yyS[yypt-0].s))
			cmd.UpdateObj(yyS[yypt-2].s, resMap(&yyS[yypt-0].s))
		}
	case 14:
		{
			println("@State NT_DEL")
			cmd.DeleteObj(yyS[yypt-1].s, resMap(&yyS[yypt-0].s))
		}
	case 15:
		{
			yyVAL.s = yyS[yypt-0].s
			cmd.Disp(resMap(&yyS[yypt-0].s))
			cmd.DeleteObj(yyS[yypt-2].s, resMap(&yyS[yypt-0].s))
		}
	case 24:
		{
			yyVAL.s = string(yyS[yypt-3].s + "=" + yyS[yypt-1].s + "=" + yyS[yypt-0].s)
			println("So we got: ", yyVAL.s)
		}
	case 25:
		{
			yyVAL.s = yyS[yypt-2].s + "=" + yyS[yypt-0].s
			println("Taking the M")
			println("SUP DUDE: ", yyS[yypt-0].s)
		}
	case 26:
		{
			yyVAL.s = yyS[yypt-2].s + "/" + yyS[yypt-0].s
		}
	case 27:
		{
			yyVAL.s = yyS[yypt-0].s
		}
	case 28:
		{
			yyVAL.s = yyS[yypt-1].s
		}
	case 30:
		{
			cmd.CD(yyS[yypt-0].s)
		}
	case 31:
		{
			cmd.CD(yyS[yypt-0].s)
		}
	case 32:
		{
			cmd.LS(yyS[yypt-0].s)
		}
	case 35:
		{
			cmd.Execute()
		}
	case 36:
		{
			cmd.CD("")
		}
	case 39:
		{
			cmd.LS("")
		}
	case 40:
		{
			cmd.LS(yyS[yypt-0].s)
		}
	case 41:
		{
			cmd.PWD()
		}
	case 42:
		{
			cmd.Exit()
		}
	case 43:
		{
			cmd.Help()
		}

	}

	if yyEx != nil && yyEx.Reduced(r, exState, &yyVAL) {
		return -1
	}
	goto yystack /* stack new state and value */
}
