package company

type CompanyType string

func (ct CompanyType) IsValid() bool {
	if ct == Corporations || ct == NonProfit || ct == Cooperative || ct == SoleProprietorship {
		return true
	}
	return false
}

const (
	Corporations       CompanyType = "Corporations"
	NonProfit          CompanyType = "NonProfit"
	Cooperative        CompanyType = "Cooperative"
	SoleProprietorship CompanyType = "Sole Proprietorship"
)
