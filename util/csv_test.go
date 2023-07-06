package util

import (
	"testing"
)

// TestCsvImportRequiredHeadersMisnumber calls util.CsvImport with csv data
// without enough headers
func TestCsvImportRequiredHeadersMisnumber(t *testing.T) {
	csv := CsvLines{{"sku", "title", "price", "currency", "quantity", "status", "description"}}

	_, err := CsvImport(csv)

	if err == nil {
		t.Fatal(`The csv should failed because the csv header is not correct`)

	}
}

// TestCsvImportHeaderFirstCellMisvalue calls util.CsvImport with a bad value
// for the first csv data header
func TestCsvImportHeaderFirstCellMisvalue(t *testing.T) {
	csv := CsvLines{{"id", "title", "price", "currency", "quantity", "status", "description", "images"}}

	_, err := CsvImport(csv)

	if err == nil {
		t.Fatal(`The csv should failed because the first csv header value is not sku`)
	}
}

// TestCsvImportRequiredLineValueMisnumber calls util.CsvImport with a csv data
// line that does not contains the required values
func TestCsvImportRequiredLineValuesMisnumber(t *testing.T) {

}

// TestCsvImportCurrencyMisvalue calls util.CsvImport with a csv data
// line that contains a bad currency value
func TestCsvImportCurrencyMisvalue(t *testing.T) {

}

// TestCsvImportImagesMisvalue calls util.CsvImport with a csv data
// line that contains a bad images value
func TestCsvImportImagesMisvalue(t *testing.T) {

}

// TestCsvImportImagesMisvalue calls util.CsvImport with a csv data
// line that contains a not found image
func TestCsvImportImagesNotFound(t *testing.T) {

}

// TestCsvImportStatusMisvalue calls util.CsvImport with a csv data
// line that contains a bad status value
func TestCsvImportStatusMisvalue(t *testing.T) {

}

// TestCsvImportQuantityMisvalue calls util.CsvImport with a csv data
// line that contains a bad quantity value
func TestCsvImportQuantityMisvalue(t *testing.T) {

}

// TestCsvImportWeightMisvalue calls util.CsvImport with a csv data
// line that contains a bad weight value
func TestCsvImportWeightMisvalue(t *testing.T) {

}

// TestCsvImportTagsMisvalue calls util.CsvImport with a csv data
// line that contains a bad tags value
func TestCsvImportTagsMisvalue(t *testing.T) {

}

// TestCsvImportLinksMisvalue calls util.CsvImport with a csv data
// line that contains a bad links value
func TestCsvImportLinksMisvalue(t *testing.T) {

}

// TestCsvImportOptionsMisvalue calls util.CsvImport with a csv data
// line that contains a bad options value
func TestCsvImportOptionsMisvalue(t *testing.T) {

}
