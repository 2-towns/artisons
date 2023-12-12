package tracking

import (
	"fmt"
	"gifthub/conf"
	"gifthub/http/contexts"
	"gifthub/tests"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func TestSaveTrackingLogWhenSuccess(t *testing.T) {
	c := tests.Context()
	data := map[string]string{"hello": "world"}

	if err := Log(c, "action", data); err != nil {
		t.Fatalf(`Log(c, "action", data) = %s, want nil`, err.Error())
	}

	folder := conf.WorkingSpace + "web/tracking"
	now := time.Now()
	p := fmt.Sprintf("%s/tracking-%s.log", folder, now.Format("20060102"))
	buf, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf(`os.ReadFile(p) = %s, want empty`, err.Error())
	}

	rid := c.Value(middleware.RequestIDKey).(string)
	cid := c.Value(contexts.Cart).(string)
	l := fmt.Sprintf("rid:%s cid:%s lang:en hello:world", rid, cid)
	s := string(buf)
	if !strings.Contains(s, l) {
		t.Fatalf(`strings.Contains(s, '%s') = false, want true`, l)
	}
}
