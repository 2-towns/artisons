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

// TestCsvImportRequiredHeadersMisnumber expects to fail because of missing headers
func TestCsvImportRequiredHeadersMisnumber(t *testing.T) {
	h := make([]string, 3)
	copy(h, header)

	csv := lines{h}

	count, err := Import(csv, conf.DefaultMID)
	if err == nil || count != 0 || err.Error() != "csv_not_valid" {
		t.Fatalf(`Import = %d, %v, want 0, 'csv_not_valid'`, count, err)
	}
}

// TestCsvImportHeaderFirstCellMisvalue expects to fail because of first cell bad value
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

// TestCsvImportRequiredLineValueMisnumber expects to fail because of missing fields
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

// TestCsvImportSkuMissing expects to fail because of missing sku
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

// TestCsvImportSkuMisvalue  expects to fail because of sku misvalue
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

// TestCsvImportTitleMissing expects to fail because of missing title
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

// TestCsvImportPriceMissing expects to fail because of missing price
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

// TestCsvImportPriceMisvalue  expects to fail because of price misvalue
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

// TestCsvImportCurrencyMisvalue  expects to fail because of currency misvalue
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

// TestCsvImportQuantityMissing expects to fail because of missing quantity
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

// TestCsvImportQuantiyMisvalue expects to fail because of quantity misvalue
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

// TestCsvImportStatusMisvalue expects to fail because of status misvalue
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

// TestCsvImportDescriptionMissing expects to fail because of missing description
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

// TestCsvImportImagesMisvalue expects to fail because of images misvalue
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

// TestCsvImportImagesMissing expects to fail because of missing images
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

// TestCsvImportImagesMisvalue expects to fail because of images not found
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

// TestCsvImportImagesBadExtension expects to fail because of extension misvalue
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

// TestCsvImportLocalImageBadExtension expects to fail because of local extension misvalue
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

// TestCsvImportLocalImageNotFound expects to fail because of local file not found
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

// TestCsvImportWeightMisvalue expects to fail because of weight misvalue
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

// TestCsvImportOptionsMisvalue expects to fail because of options misvalue
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

// TestCsvImportOk expects to succeed
func TestCsvImportOk(t *testing.T) {
	csv := lines{header, line}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 1 {
		t.Fatalf(`Import = %d, %v, want 1, nil`, count, err)
	}
}

// TestCsvImportOk expects to succeed with a local image
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

// TestCsvImportWithoutOptionsOk expects to succeed without options
func TestCsvImportWithoutOptionsOk(t *testing.T) {
	l := make([]string, len(line))
	copy(l, line)

	l[11] = ""
	csv := lines{header, l}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 1 {
		t.Fatalf(`Import = %d, %v, want 1, nil`, count, err)
	}
}

// TestCsvImportWihoutLinksOk expects to succeed without links
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

// TestCsvImportWihoutWeightOk expects to succeed without weight
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
