package money

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	ofx "github.com/aclindsa/ofxgo"
	sqlc "github.com/keeth/money/model/sqlc"
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

func ParseOfxTransaction(trans ofx.Transaction) sqlc.Tx {
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
	return sqlc.Tx{
		Date:   trans.DtPosted.Format("2006-01-02"),
		Desc:   cleanDesc(strings.Join(desc, " ")),
		Amount: amount,
		Xid:    trans.FiTID.String(),
	}
}

type ParsedResponse struct {
	Transactions []sqlc.Tx
	Kind         string
	ID           string
}

func ParseOfxResponse(file *os.File) (ParsedResponse, error) {
	resp, err := ofx.ParseResponse(file)
	if err != nil {
		return ParsedResponse{}, err
	}
	parsed := ParsedResponse{}
	xidParts := []string{}
	if len(resp.Bank) > 0 {
		if stmt, ok := resp.Bank[0].(*ofx.StatementResponse); ok {
			parsed.Transactions = make([]sqlc.Tx, len(stmt.BankTranList.Transactions))
			for i, tx := range stmt.BankTranList.Transactions {
				parsed.Transactions[i] = ParseOfxTransaction(tx)
			}
			if stmt.BankAcctFrom.BankID != "" {
				xidParts = append(xidParts, stmt.BankAcctFrom.BankID.String())
			}
			if stmt.BankAcctFrom.BranchID != "" {
				xidParts = append(xidParts, stmt.BankAcctFrom.BranchID.String())
			}
			if stmt.BankAcctFrom.AcctID != "" {
				xidParts = append(xidParts, stmt.BankAcctFrom.AcctID.String())
			}
		}
		parsed.Kind = "bank"
	} else if len(resp.CreditCard) > 0 {
		if stmt, ok := resp.CreditCard[0].(*ofx.CCStatementResponse); ok {
			parsed.Transactions = make([]sqlc.Tx, len(stmt.BankTranList.Transactions))
			for i, tx := range stmt.BankTranList.Transactions {
				parsed.Transactions[i] = ParseOfxTransaction(tx)
			}
			xidParts = append(xidParts, stmt.CCAcctFrom.AcctID.String())
		}
		parsed.Kind = "cc"
	} else {
		return parsed, fmt.Errorf("no information found in file")
	}
	parsed.ID = strings.Join(xidParts, " ")
	return parsed, nil
}
