package util

import (
	"gifthub/conf"
	"testing"
)

const image = "https://upload.wikimedia.org/wikipedia/commons/c/ca/1x1.png"

var line = []string{
	"123456", "Product", "12.4", "EUR", "1", "online", "Best product", image, "164", "gifts;garden", "1234", "color:blue",
}

var header = []string{"sku", "title", "price", "currency", "quantity", "status", "description", "images", "weight", "tags", "links", "options"}

// TestCsvImportRequiredHeadersMisnumber calls util.CsvImport with csv data
// without enough headers
func TestCsvImportRequiredHeadersMisnumber(t *testing.T) {
	h := make([]string, 3)
	copy(h, header)

	csv := CsvLines{h}

	count, err := CsvImport(csv, conf.DefaultMID)

	if err == nil {
		t.Fatal(`The import should failed because the csv header is not correct`)

	}

	if count != 0 {
		t.Fatal(`The lines processed should be 0`)
	}
}

// TestCsvImportHeaderFirstCellMisvalue calls util.CsvImport with a bad value
// for the first csv data header
func TestCsvImportHeaderFirstCellMisvalue(t *testing.T) {
	h := make([]string, len(header))
	copy(h, header)

	h[0] = "id"

	csv := CsvLines{h}

	count, err := CsvImport(csv, conf.DefaultMID)

	if err == nil {
		t.Fatal(`The import should failed because the first csv header value is not sku`)
	}

	if count != 0 {
		t.Fatal(`The lines processed should be 0`)
	}
}

// TestCsvImportRequiredLineValueMisnumber calls util.CsvImport with a csv data
// line that does not contains the required values
func TestCsvImportRequiredLineValuesMisnumber(t *testing.T) {
	h := make([]string, len(header))
	copy(h, header)

	l := make([]string, 3)
	copy(l, line)

	csv := CsvLines{h, l}

	count, err := CsvImport(csv, conf.DefaultMID)

	if err != nil {
		t.Fatal(`The import should failed because the first line does not contains the required fields`)
	}

	if count != 0 {
		t.Fatal(`The lines processed should be 0`)
	}
}

// TestCsvImportSkuMissing calls util.CsvImport with a csv data
// line without sku
func TestCsvImportSkuMissing(t *testing.T) {
	h := make([]string, len(header))
	copy(h, header)

	l := make([]string, len(line))
	copy(l, line)

	l[0] = ""

	csv := CsvLines{h, l}

	count, err := CsvImport(csv, conf.DefaultMID)

	if err != nil {
		t.Fatal(`The import should not failed`)
	}

	if count != 0 {
		t.Fatal(`The lines processed should be 0`)
	}
}

// TestCsvImportSkuMisvalue calls util.CsvImport with a csv data
// line that contains bad sku value
func TestCsvImportSkuMisvalue(t *testing.T) {
	h := make([]string, len(header))
	copy(h, header)

	l := make([]string, len(line))
	copy(l, line)

	l[0] = "sku!"

	csv := CsvLines{h, l}

	count, err := CsvImport(csv, conf.DefaultMID)

	if err != nil {
		t.Fatal(`The import should not failed`)
	}

	if count != 0 {
		t.Fatal(`The lines processed should be 0`)
	}
}

// TestCsvImportTitleMissing calls util.CsvImport with a csv data
// line without title
func TestCsvImportTitleMissing(t *testing.T) {
	h := make([]string, len(header))
	copy(h, header)

	l := make([]string, len(line))
	copy(l, line)

	l[1] = ""

	csv := CsvLines{h, l}

	count, err := CsvImport(csv, conf.DefaultMID)

	if err != nil {
		t.Fatal(`The import should not failed`)
	}

	if count != 0 {
		t.Fatal(`The lines processed should be 0`)
	}
}

// TestCsvImportPriceMissing calls util.CsvImport with a csv data
// line without price
func TestCsvImportPriceMissing(t *testing.T) {
	h := make([]string, len(header))
	copy(h, header)

	l := make([]string, len(line))
	copy(l, line)

	l[2] = ""

	csv := CsvLines{h, l}

	count, err := CsvImport(csv, conf.DefaultMID)

	if err != nil {
		t.Fatal(`The import should not failed`)
	}

	if count != 0 {
		t.Fatal(`The lines processed should be 0`)
	}
}

// TestCsvImportPriceMisvalue calls util.CsvImport with a csv data
// line that contains a bad price value
func TestCsvImportPriceMisvalue(t *testing.T) {
	h := make([]string, len(header))
	copy(h, header)

	l := make([]string, len(line))
	copy(l, line)

	l[2] = "toto"

	csv := CsvLines{h, l}

	count, err := CsvImport(csv, conf.DefaultMID)

	if err != nil {
		t.Fatal(`The import should not failed`)
	}

	if count != 0 {
		t.Fatal(`The lines processed should be 0`)
	}
}

// TestCsvImportCurrencyMisvalue calls util.CsvImport with a csv data
// line that contains a bad currency value
func TestCsvImportCurrencyMisvalue(t *testing.T) {
	h := make([]string, len(header))
	copy(h, header)

	l := make([]string, len(line))
	copy(l, line)

	l[3] = "toto"

	csv := CsvLines{h, l}

	count, err := CsvImport(csv, conf.DefaultMID)

	if err != nil {
		t.Fatal(`The import should not failed`)
	}

	if count != 0 {
		t.Fatal(`The lines processed should be 0`)
	}
}

// TestCsvImportQuantityMissing calls util.CsvImport with a csv data
// line without quantity
func TestCsvImportQuantityMissing(t *testing.T) {
	h := make([]string, len(header))
	copy(h, header)

	l := make([]string, len(line))
	copy(l, line)

	l[4] = ""

	csv := CsvLines{h, l}

	count, err := CsvImport(csv, conf.DefaultMID)

	if err != nil {
		t.Fatal(`The import should not failed`)
	}

	if count != 0 {
		t.Fatal(`The lines processed should be 0`)
	}
}

// TestCsvImportQuantiyMisvalue calls util.CsvImport with a csv data
// line that contains a bad quantity value
func TestCsvImportQuantiyMisvalue(t *testing.T) {
	h := make([]string, len(header))
	copy(h, header)

	l := make([]string, len(line))
	copy(l, line)

	l[4] = "toto"

	csv := CsvLines{h, l}

	count, err := CsvImport(csv, conf.DefaultMID)

	if err != nil {
		t.Fatal(`The import should not failed`)
	}

	if count != 0 {
		t.Fatal(`The lines processed should be 0`)
	}
}

// TestCsvImportStatusMisvalue calls util.CsvImport with a csv data
// line that contains a bad status value
func TestCsvImportStatusMisvalue(t *testing.T) {
	h := make([]string, len(header))
	copy(h, header)

	l := make([]string, len(line))
	copy(l, line)

	l[5] = "toto"

	csv := CsvLines{h, l}

	count, err := CsvImport(csv, conf.DefaultMID)

	if err != nil {
		t.Fatal(`The import should not failed`)
	}

	if count != 0 {
		t.Fatal(`The lines processed should be 0`)
	}
}

// TestCsvImportDescriptionMissing calls util.CsvImport with a csv data
// line without description
func TestCsvImportDescriptionMissing(t *testing.T) {
	h := make([]string, len(header))
	copy(h, header)

	l := make([]string, len(line))
	copy(l, line)

	l[6] = ""

	csv := CsvLines{h, l}

	count, err := CsvImport(csv, conf.DefaultMID)

	if err != nil {
		t.Fatal(`The import should not failed`)
	}

	if count != 0 {
		t.Fatal(`The lines processed should be 0`)
	}
}

// TestCsvImportImagesMisvalue calls util.CsvImport with a csv data
// line that contains a bad images value
func TestCsvImportImagesMisvalue(t *testing.T) {
	h := make([]string, len(header))
	copy(h, header)

	l := make([]string, len(line))
	copy(l, line)

	l[7] = "toto"

	csv := CsvLines{h, l}

	count, err := CsvImport(csv, conf.DefaultMID)

	if err != nil {
		t.Fatal(`The import should not failed`)
	}

	if count != 0 {
		t.Fatal(`The lines processed should be 0`)
	}
}

// TestCsvImportImagesMissing calls util.CsvImport with a csv data
// line that contains a bad images value
func TestCsvImportImagesMissing(t *testing.T) {
	h := make([]string, len(header))
	copy(h, header)

	l := make([]string, len(line))
	copy(l, line)

	l[7] = ""

	csv := CsvLines{h, l}

	count, err := CsvImport(csv, conf.DefaultMID)

	if err != nil {
		t.Fatal(`The import should not failed`)
	}

	if count != 0 {
		t.Fatal(`The lines processed should be 0`)
	}
}

// TestCsvImportImagesMisvalue calls util.CsvImport with a csv data
// line that contains a not found image
func TestCsvImportImagesNotFound(t *testing.T) {
	h := make([]string, len(header))
	copy(h, header)

	l := make([]string, len(line))
	copy(l, line)

	const image = "https://upload.wikimedia.org/wikipedia/commons/c/ca/1x1_toto.png"
	l[7] = image

	csv := CsvLines{h, l}

	count, err := CsvImport(csv, conf.DefaultMID)

	if err != nil {
		t.Fatal(`The import should not failed`)
	}

	if count != 0 {
		t.Fatal(`The lines processed should be 0`)
	}
}

// TestCsvImportImagesBadExtension calls util.CsvImport with a csv data
// line that contains a file with a bad extension
func TestCsvImportImagesBadExtension(t *testing.T) {
	h := make([]string, len(header))
	copy(h, header)

	l := make([]string, len(line))
	copy(l, line)

	const image = "https://upload.wikimedia.org/wikipedia/commons/c/ca/1x1.svg"
	l[7] = image

	csv := CsvLines{h, l}

	count, err := CsvImport(csv, conf.DefaultMID)

	if err != nil {
		t.Fatal(`The import should not failed`)
	}

	if count != 0 {
		t.Fatal(`The lines processed should be 0`)
	}
}

// TestCsvImportLocalImageBadExtension calls util.CsvImport with a csv data
// line that a local file with a bad extension
func TestCsvImportLocalImageBadExtension(t *testing.T) {
	h := make([]string, len(header))
	copy(h, header)

	l := make([]string, len(line))
	copy(l, line)

	l[7] = "../static/fake/product.svg"

	csv := CsvLines{h, l}

	count, err := CsvImport(csv, conf.DefaultMID)

	if err != nil {
		t.Fatal(`The import should not failed`)
	}

	if count != 0 {
		t.Fatal(`The lines processed should be 0`)
	}
}

// TestCsvImportLocalImageNotFound calls util.CsvImport with a csv data
// line that a local file not found
func TestCsvImportLocalImageNotFound(t *testing.T) {
	h := make([]string, len(header))
	copy(h, header)

	l := make([]string, len(line))
	copy(l, line)

	l[7] = "../static/fake/toto.png"

	csv := CsvLines{h, l}

	count, err := CsvImport(csv, conf.DefaultMID)

	if err != nil {
		t.Fatal(`The import should not failed`)
	}

	if count != 0 {
		t.Fatal(`The lines processed should be 0`)
	}
}

// TestCsvImportWeightMisvalue calls util.CsvImport with a csv data
// line that contains a bad weight value
func TestCsvImportWeightMisvalue(t *testing.T) {
	h := make([]string, len(header))
	copy(h, header)

	l := make([]string, len(line))
	copy(l, line)

	l[8] = "toto"

	csv := CsvLines{h, l}

	count, err := CsvImport(csv, conf.DefaultMID)

	if err != nil {
		t.Fatal(`The import should not failed`)
	}

	if count != 0 {
		t.Fatal(`The lines processed should be 0`)
	}
}

// TestCsvImportOptionsMisvalue calls util.CsvImport with a csv data
// line that contains a bad options value
func TestCsvImportOptionsMisvalue(t *testing.T) {
	h := make([]string, len(header))
	copy(h, header)

	l := make([]string, len(line))
	copy(l, line)

	l[11] = "toto"

	csv := CsvLines{h, l}

	count, err := CsvImport(csv, conf.DefaultMID)

	if err != nil {
		t.Fatal(`The import should not failed`)
	}

	if count != 0 {
		t.Fatal(`The lines processed should be 0`)
	}
}

// TestCsvImportOk calls util.CsvImport with a csv data
// line that contains a bad options value
func TestCsvImportOk(t *testing.T) {
	csv := CsvLines{header, line}

	count, err := CsvImport(csv, conf.DefaultMID)

	if err != nil {
		t.Fatal(`The import should not failed`)
	}

	if count != 1 {
		t.Fatal(`The lines processed should be 1`)
	}
}

// TestCsvImportLocalImageOk calls util.CsvImport with a csv data
// line that contains a local image
func TestCsvImportLocalImageOk(t *testing.T) {
	l := make([]string, len(line))
	copy(l, line)

	l[7] = "../static/fake/product.png"

	csv := CsvLines{header, l}

	count, err := CsvImport(csv, conf.DefaultMID)

	if err != nil {
		t.Fatal(`The import should not failed`)
	}

	if count != 1 {
		t.Fatal(`The lines processed should be 1`)
	}
}

// TestCsvImportWihoutOptionsOk calls util.CsvImport with a csv data
// line that does not contains options
func TestCsvImportWihoutOptionsOk(t *testing.T) {
	l := make([]string, len(line))
	copy(l, line)

	l[11] = ""

	csv := CsvLines{header, l}

	count, err := CsvImport(csv, conf.DefaultMID)

	if err != nil {
		t.Fatal(`The import should not failed`)
	}

	if count != 1 {
		t.Fatal(`The lines processed should be 1`)
	}
}

// TestCsvImportWihoutLinksOk calls util.CsvImport with a csv data
// line that does not contains links
func TestCsvImportWihoutLinksOk(t *testing.T) {
	l := make([]string, len(line))
	copy(l, line)

	l[10] = ""

	csv := CsvLines{header, l}

	count, err := CsvImport(csv, conf.DefaultMID)

	if err != nil {
		t.Fatal(`The import should not failed`)
	}

	if count != 1 {
		t.Fatal(`The lines processed should be 1`)
	}
}

// TestCsvImportWihoutWeightOk calls util.CsvImport with a csv data
// line that does not contains weight
func TestCsvImportWihoutWeightOk(t *testing.T) {
	l := make([]string, len(line))
	copy(l, line)

	l[8] = ""

	csv := CsvLines{header, l}

	count, err := CsvImport(csv, conf.DefaultMID)

	if err != nil {
		t.Fatal(`The import should not failed`)
	}

	if count != 1 {
		t.Fatal(`The lines processed should be 1`)
	}
}
