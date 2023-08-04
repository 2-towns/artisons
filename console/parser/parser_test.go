package parser

import (
	"gifthub/conf"
	"testing"
)

const image = "https://upload.wikimedia.org/wikipedia/commons/c/ca/1x1.png"

var line = []string{
	"123456", "Product", "12.4", "EUR", "1", "online", "Best product", image, "164", "gifts;garden", "1234", "color:blue",
}

var header = []string{"sku", "title", "price", "currency", "quantity", "status", "description", "images", "weight", "tags", "links", "options"}

func init() {
	conf.ImgProxyPath = "../../../" + conf.ImgProxyPath
}

// TestCsvImportRequiredHeadersMisnumber calls util.CsvImport with csv data
// without enough headers
func TestCsvImportRequiredHeadersMisnumber(t *testing.T) {
	h := make([]string, 3)
	copy(h, header)

	csv := lines{h}

	count, err := Import(csv, conf.DefaultMID)
	if err == nil || count != 0 || err.Error() != "csv_not_valid" {
		t.Fatalf(`Import = %d, %v, want 0, 'csv_not_valid'`, count, err)
	}
}

// TestCsvImportHeaderFirstCellMisvalue calls util.CsvImport with a bad value
// for the first csv data header
func TestCsvImportHeaderFirstCellMisvalue(t *testing.T) {
	h := make([]string, len(header))
	copy(h, header)

	h[0] = "id"

	csv := lines{h}

	count, err := Import(csv, conf.DefaultMID)
	if err == nil || count != 0 || err.Error() != "csv_not_valid" {
		t.Fatalf(`Import = %d, %v, want 0, 'csv_not_valid'`, count, err)
	}
}

// TestCsvImportRequiredLineValueMisnumber calls util.CsvImport with a csv data
// line that does not contains the required values
func TestCsvImportRequiredLineValuesMisnumber(t *testing.T) {
	h := make([]string, len(header))
	copy(h, header)

	l := make([]string, 3)
	copy(l, line)

	csv := lines{h, l}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 0 {
		t.Fatalf(`Import = %d, %v, want 0, nil`, count, err)
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
	csv := lines{h, l}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 0 {
		t.Fatalf(`Import = %d, %v, want 0, nil`, count, err)
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
	csv := lines{h, l}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 0 {
		t.Fatalf(`Import = %d, %v, want 0, nil`, count, err)
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
	csv := lines{h, l}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 0 {
		t.Fatalf(`Import = %d, %v, want 0, nil`, count, err)
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
	csv := lines{h, l}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 0 {
		t.Fatalf(`Import = %d, %v, want 0, nil`, count, err)
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
	csv := lines{h, l}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 0 {
		t.Fatalf(`Import = %d, %v, want 0, nil`, count, err)
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
	csv := lines{h, l}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 0 {
		t.Fatalf(`Import = %d, %v, want 0, nil`, count, err)
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
	csv := lines{h, l}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 0 {
		t.Fatalf(`Import = %d, %v, want 0, nil`, count, err)
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
	csv := lines{h, l}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 0 {
		t.Fatalf(`Import = %d, %v, want 0, nil`, count, err)
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
	csv := lines{h, l}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 0 {
		t.Fatalf(`Import = %d, %v, want 0, nil`, count, err)
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
	csv := lines{h, l}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 0 {
		t.Fatalf(`Import = %d, %v, want 0, nil`, count, err)
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
	csv := lines{h, l}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 0 {
		t.Fatalf(`Import = %d, %v, want 0, nil`, count, err)
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
	csv := lines{h, l}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 0 {
		t.Fatalf(`Import = %d, %v, want 0, nil`, count, err)
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
	csv := lines{h, l}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 0 {
		t.Fatalf(`Import = %d, %v, want 0, nil`, count, err)
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
	csv := lines{h, l}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 0 {
		t.Fatalf(`Import = %d, %v, want 0, nil`, count, err)
	}
}

// TestCsvImportLocalImageBadExtension calls util.CsvImport with a csv data
// line that a local file with a bad extension
func TestCsvImportLocalImageBadExtension(t *testing.T) {
	h := make([]string, len(header))
	copy(h, header)

	l := make([]string, len(line))
	copy(l, line)
	l[7] = "../../web/testdata/product.svg"
	csv := lines{h, l}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 0 {
		t.Fatalf(`Import = %d, %v, want 0, nil`, count, err)
	}
}

// TestCsvImportLocalImageNotFound calls util.CsvImport with a csv data
// line that a local file not found
func TestCsvImportLocalImageNotFound(t *testing.T) {
	h := make([]string, len(header))
	copy(h, header)

	l := make([]string, len(line))
	copy(l, line)
	l[7] = "../../web/testdata/toto.png"
	csv := lines{h, l}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 0 {
		t.Fatalf(`Import = %d, %v, want 0, nil`, count, err)
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
	csv := lines{h, l}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 0 {
		t.Fatalf(`Import = %d, %v, want 0, nil`, count, err)
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
	csv := lines{h, l}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 0 {
		t.Fatalf(`Import = %d, %v, want 0, nil`, count, err)
	}
}

// TestCsvImportOk calls util.CsvImport with a csv data
// line that contains a bad options value
func TestCsvImportOk(t *testing.T) {
	csv := lines{header, line}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 1 {
		t.Fatalf(`Import = %d, %v, want 1, nil`, count, err)
	}
}

// TestCsvImportLocalImageOk calls util.CsvImport with a csv data
// line that contains a local image
func TestCsvImportLocalImageOk(t *testing.T) {
	l := make([]string, len(line))
	copy(l, line)

	l[7] = "../../web/testdata/product.png"
	csv := lines{header, l}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 1 {
		t.Fatalf(`Import = %d, %v, want 1, nil`, count, err)
	}
}

// TestCsvImportWihoutOptionsOk calls util.CsvImport with a csv data
// line that does not contains options
func TestCsvImportWihoutOptionsOk(t *testing.T) {
	l := make([]string, len(line))
	copy(l, line)

	l[11] = ""
	csv := lines{header, l}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 1 {
		t.Fatalf(`Import = %d, %v, want 1, nil`, count, err)
	}
}

// TestCsvImportWihoutLinksOk calls util.CsvImport with a csv data
// line that does not contains links
func TestCsvImportWihoutLinksOk(t *testing.T) {
	l := make([]string, len(line))
	copy(l, line)

	l[10] = ""
	csv := lines{header, l}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 1 {
		t.Fatalf(`Import = %d, %v, want 1, nil`, count, err)
	}
}

// TestCsvImportWihoutWeightOk calls util.CsvImport with a csv data
// line that does not contains weight
func TestCsvImportWihoutWeightOk(t *testing.T) {
	l := make([]string, len(line))
	copy(l, line)

	l[8] = ""
	csv := lines{header, l}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 1 {
		t.Fatalf(`Import = %d, %v, want 1, nil`, count, err)
	}
}
