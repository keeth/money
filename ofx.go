package money

import (
	"regexp"
	"strings"

	ofx "github.com/aclindsa/ofxgo"
)

var (
	descRe1 = regexp.MustCompile(`(\w)\*+(\w)`)
	descRe2 = regexp.MustCompile(`(\b\*+|\*+\b)`)
	memoRe  = regexp.MustCompile(`[,;(].*$`)
)

type ParsedTransaction struct {
	Date   string
	Desc   string
	Amount float64
	ID     string
}

func cleanDesc(s string) string {
	s = strings.ToLower(s)
	s = descRe1.ReplaceAllString(s, `$1-$2`)
	s = descRe2.ReplaceAllString(s, "")
	s = strings.TrimSpace(s)
	return s
}

func cleanMemo(s string) string {
	return memoRe.ReplaceAllString(s, "")
}

func ParseOfxTransaction(trans *ofx.Transaction) ParsedTransaction {
	amount, _ := trans.TrnAmt.Float64()
	var desc []string
	if trans.Name != "" {
		desc = append(desc, trans.Name.String())
	} else if trans.Payee != nil {
		desc = append(desc, trans.Payee.Name.String())
	}
	if trans.Memo != "" {
		desc = append(desc, cleanMemo(trans.Memo.String()))
	}
	return ParsedTransaction{
		Date:   trans.DtPosted.Format("2006-01-02"),
		Desc:   cleanDesc(strings.Join(desc, " ")),
		Amount: amount,
		ID:     trans.FiTID.String(),
	}
}
