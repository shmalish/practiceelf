package oned

import (
	"testing"
)

func testLookupCountryIdentifier(t testing.TB, productCode, expect string) {
	t.Helper()
	r := eanManufacturerOrgSupportLookupCountryIdentifier(productCode)
	if r != expect {
		t.Fatalf("eanManufacturerOrgSupportList(%v) = \"%v\", expect \"%v\"", productCode, r, expect)
	}
}

func TestLookupCountryIdentifier(t *testing.T) {
	testLookupCountryIdentifier(t, "10", "")
	testLookupCountryIdentifier(t, "abcd", "")
	testLookupCountryIdentifier(t, "0001234", "US/CA")
	testLookupCountryIdentifier(t, "0191234", "US/CA")
	testLookupCountryIdentifier(t, "0201234", "")
	testLookupCountryIdentifier(t, "0291234", "")
	testLookupCountryIdentifier(t, "0301234", "US")
	testLookupCountryIdentifier(t, "0391234", "US")
	testLookupCountryIdentifier(t, "0401234", "")
	testLookupCountryIdentifier(t, "0591234", "")
	testLookupCountryIdentifier(t, "0601234", "US/CA")
	testLookupCountryIdentifier(t, "1391234", "US/CA")
	testLookupCountryIdentifier(t, "1401234", "")
	testLookupCountryIdentifier(t, "3101234", "FR")
	testLookupCountryIdentifier(t, "3801234", "BG")
	testLookupCountryIdentifier(t, "3831234", "SI")
	testLookupCountryIdentifier(t, "3851234", "HR")
	testLookupCountryIdentifier(t, "3871234", "BA")
	testLookupCountryIdentifier(t, "4201234", "DE")
	testLookupCountryIdentifier(t, "4551234", "JP")
	testLookupCountryIdentifier(t, "4651234", "RU")
	testLookupCountryIdentifier(t, "4711234", "TW")
	testLookupCountryIdentifier(t, "4741234", "EE")
	testLookupCountryIdentifier(t, "4751234", "LV")
	testLookupCountryIdentifier(t, "4761234", "AZ")
	testLookupCountryIdentifier(t, "4771234", "LT")
	testLookupCountryIdentifier(t, "4781234", "UZ")
	testLookupCountryIdentifier(t, "4791234", "LK")
	testLookupCountryIdentifier(t, "4801234", "PH")
	testLookupCountryIdentifier(t, "4811234", "BY")
	testLookupCountryIdentifier(t, "4821234", "UA")
	testLookupCountryIdentifier(t, "4841234", "MD")
	testLookupCountryIdentifier(t, "4851234", "AM")
	testLookupCountryIdentifier(t, "4861234", "GE")
	testLookupCountryIdentifier(t, "4871234", "KZ")
	testLookupCountryIdentifier(t, "4891234", "HK")
	testLookupCountryIdentifier(t, "4901234", "JP")
	testLookupCountryIdentifier(t, "5001234", "GB")
	testLookupCountryIdentifier(t, "5201234", "GR")
	testLookupCountryIdentifier(t, "5281234", "LB")
	testLookupCountryIdentifier(t, "5291234", "CY")
	testLookupCountryIdentifier(t, "5311234", "MK")
	testLookupCountryIdentifier(t, "5351234", "MT")
	testLookupCountryIdentifier(t, "5391234", "IE")
	testLookupCountryIdentifier(t, "5401234", "BE/LU")
	testLookupCountryIdentifier(t, "5601234", "PT")
	testLookupCountryIdentifier(t, "5691234", "IS")
	testLookupCountryIdentifier(t, "5701234", "DK")
	testLookupCountryIdentifier(t, "5901234", "PL")
	testLookupCountryIdentifier(t, "5941234", "RO")
	testLookupCountryIdentifier(t, "5991234", "HU")
	testLookupCountryIdentifier(t, "6001234", "ZA")
	testLookupCountryIdentifier(t, "6031234", "GH")
	testLookupCountryIdentifier(t, "6081234", "BH")
	testLookupCountryIdentifier(t, "6091234", "MU")
	testLookupCountryIdentifier(t, "6111234", "MA")
	testLookupCountryIdentifier(t, "6131234", "DZ")
	testLookupCountryIdentifier(t, "6161234", "KE")
	testLookupCountryIdentifier(t, "6181234", "CI")
	testLookupCountryIdentifier(t, "6191234", "TN")
	testLookupCountryIdentifier(t, "6211234", "SY")
	testLookupCountryIdentifier(t, "6221234", "EG")
	testLookupCountryIdentifier(t, "6241234", "LY")
	testLookupCountryIdentifier(t, "6251234", "JO")
	testLookupCountryIdentifier(t, "6261234", "IR")
	testLookupCountryIdentifier(t, "6271234", "KW")
	testLookupCountryIdentifier(t, "6281234", "SA")
	testLookupCountryIdentifier(t, "6291234", "AE")
	testLookupCountryIdentifier(t, "6401234", "FI")
	testLookupCountryIdentifier(t, "6901234", "CN")
	testLookupCountryIdentifier(t, "7001234", "NO")
	testLookupCountryIdentifier(t, "7291234", "IL")
	testLookupCountryIdentifier(t, "7301234", "SE")
	testLookupCountryIdentifier(t, "7401234", "GT")
	testLookupCountryIdentifier(t, "7411234", "SV")
	testLookupCountryIdentifier(t, "7421234", "HN")
	testLookupCountryIdentifier(t, "7431234", "NI")
	testLookupCountryIdentifier(t, "7441234", "CR")
	testLookupCountryIdentifier(t, "7451234", "PA")
	testLookupCountryIdentifier(t, "7461234", "DO")
	testLookupCountryIdentifier(t, "7501234", "MX")
	testLookupCountryIdentifier(t, "7541234", "CA")
	testLookupCountryIdentifier(t, "7591234", "VE")
	testLookupCountryIdentifier(t, "7601234", "CH")
	testLookupCountryIdentifier(t, "7701234", "CO")
	testLookupCountryIdentifier(t, "7731234", "UY")
	testLookupCountryIdentifier(t, "7751234", "PE")
	testLookupCountryIdentifier(t, "7771234", "BO")
	testLookupCountryIdentifier(t, "7791234", "AR")
	testLookupCountryIdentifier(t, "7801234", "CL")
	testLookupCountryIdentifier(t, "7841234", "PY")
	testLookupCountryIdentifier(t, "7851234", "PE")
	testLookupCountryIdentifier(t, "7861234", "EC")
	testLookupCountryIdentifier(t, "7891234", "BR")
	testLookupCountryIdentifier(t, "8001234", "IT")
	testLookupCountryIdentifier(t, "8401234", "ES")
	testLookupCountryIdentifier(t, "8501234", "CU")
	testLookupCountryIdentifier(t, "8581234", "SK")
	testLookupCountryIdentifier(t, "8591234", "CZ")
	testLookupCountryIdentifier(t, "8601234", "YU")
	testLookupCountryIdentifier(t, "8651234", "MN")
	testLookupCountryIdentifier(t, "8671234", "KP")
	testLookupCountryIdentifier(t, "8681234", "TR")
	testLookupCountryIdentifier(t, "8701234", "NL")
	testLookupCountryIdentifier(t, "8801234", "KR")
	testLookupCountryIdentifier(t, "8851234", "TH")
	testLookupCountryIdentifier(t, "8881234", "SG")
	testLookupCountryIdentifier(t, "8901234", "IN")
	testLookupCountryIdentifier(t, "8931234", "VN")
	testLookupCountryIdentifier(t, "8961234", "PK")
	testLookupCountryIdentifier(t, "8991234", "ID")
	testLookupCountryIdentifier(t, "9001234", "AT")
	testLookupCountryIdentifier(t, "9301234", "AU")
	testLookupCountryIdentifier(t, "9401234", "AZ")
	testLookupCountryIdentifier(t, "9551234", "MY")
	testLookupCountryIdentifier(t, "9581234", "MO")
}