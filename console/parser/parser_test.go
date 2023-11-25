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

func TestImportReturnsErrorWhenHeadersAreMissing(t *testing.T) {
	h := make([]string, 3)
	copy(h, header)

	csv := lines{h}

	count, err := Import(csv, conf.DefaultMID)
	if err == nil || count != 0 || err.Error() != "error_csv_fileinvalid" {
		t.Fatalf(`Import = %d, %v, want 0, 'error_csv_fileinvalid'`, count, err)
	}
}

func TestImportReturnsErrorWhenFirstCellIsInvalid(t *testing.T) {
	h := make([]string, len(header))
	copy(h, header)

	h[0] = "id"

	csv := lines{h}

	count, err := Import(csv, conf.DefaultMID)
	if err == nil || count != 0 || err.Error() != "error_csv_fileinvalid" {
		t.Fatalf(`Import = %d, %v, want 0, 'error_csv_fileinvalid'`, count, err)
	}
}

func TestImportReturnsCountZeroWhenFieldsAreMissing(t *testing.T) {
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

func TestImportReturnsCountZeroWhenSkuIsMissing(t *testing.T) {
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

func TestImportReturnsCountZeroWhenSkuIsInvalid(t *testing.T) {
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

func TestImportReturnsCountZeroWhenTitleIsMissing(t *testing.T) {
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

func TestImportReturnsCountZeroWhenPriceIsMissing(t *testing.T) {
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

func TestImportReturnsCountZeroWhenPriceIsInvalid(t *testing.T) {
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

func TestImportReturnsCountZeroWhenCurrencyIsInvalid(t *testing.T) {
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
func TestImportReturnsCountZeroWhenQuantityIsMissing(t *testing.T) {
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

func TestImportReturnsCountZeroWhenQuantityIsInvalid(t *testing.T) {
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

func TestImportReturnsCountZeroWhenStatusIsInvalid(t *testing.T) {
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

func TestImportReturnsCountZeroWhenDescriptionIsMissing(t *testing.T) {
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

func TestImportReturnsCountZeroImagesAreInvalid(t *testing.T) {
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

func TestImportReturnsCountZeroWhenSkuImagesAreMissing(t *testing.T) {
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

func TestImportReturnsCountZeroWhenImagesAreNotFound(t *testing.T) {
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

func TestImportReturnsCountZeroWhenImageHaveBadExtension(t *testing.T) {
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

func TestImportReturnsCountZeroWhenLocalImageHaveBadExtension(t *testing.T) {
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

func TestImportReturnsCountZeroWhenLocalImageIsNotFound(t *testing.T) {
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

func TestImportReturnsCountZeroWhenWeightIsInvalid(t *testing.T) {
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

func TestImportReturnsCountZeroWhenOptionsAreInvalid(t *testing.T) {
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

func TestImportReturnsOneCountWhenSuccess(t *testing.T) {
	csv := lines{header, line}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 1 {
		t.Fatalf(`Import = %d, %v, want 1, nil`, count, err)
	}
}

func TestImportReturnsOneCountWhenLocalImageAndSuccess(t *testing.T) {
	l := make([]string, len(line))
	copy(l, line)

	l[7] = "../../web/testdata/product.png"
	csv := lines{header, l}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 1 {
		t.Fatalf(`Import = %d, %v, want 1, nil`, count, err)
	}
}

func TestImportReturnsOneCountWhenNoOptionsAndSuccess(t *testing.T) {
	l := make([]string, len(line))
	copy(l, line)

	l[11] = ""
	csv := lines{header, l}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 1 {
		t.Fatalf(`Import = %d, %v, want 1, nil`, count, err)
	}
}

func TestImportReturnsOneCountWhenNoLinksAndSuccess(t *testing.T) {
	l := make([]string, len(line))
	copy(l, line)

	l[10] = ""
	csv := lines{header, l}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 1 {
		t.Fatalf(`Import = %d, %v, want 1, nil`, count, err)
	}
}

func TestImportReturnsOneCountWhenNoWeightAndSuccess(t *testing.T) {
	l := make([]string, len(line))
	copy(l, line)

	l[8] = ""
	csv := lines{header, l}

	count, err := Import(csv, conf.DefaultMID)
	if err != nil || count != 1 {
		t.Fatalf(`Import = %d, %v, want 1, nil`, count, err)
	}
}
