package tracking

// import (
// 	"artisons/conf"
// 	"artisons/http/contexts"
// 	"artisons/tests"
// 	"errors"
// 	"fmt"
// 	"log"
// 	"os"
// 	"path"
// 	"strings"
// 	"testing"
// 	"time"

// )

// func TestSave(t *testing.T) {
// 	c := tests.Context()
// 	data := map[string]string{"hello": "world"}

// 	if err := Log(c, "action", data); err != nil {
// 		t.Fatalf(`err = %s, want nil`, err.Error())
// 	}

// 	folder := conf.WorkingSpace + "web/tracking"
// 	now := time.Now()
// 	name := fmt.Sprintf("tracking-%s.log", now.Format("20060102"))
// 	p := path.Join(folder, name)
// 	if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
// 		log.Println("fdsfdsfds!!!!")
// 		_, err := os.Create(p)
// 		if err != nil {
// 			t.Fatalf(`err = %s, want empty`, err.Error())
// 		}
// 	}

// 	buf, err := os.ReadFile(p)
// 	if err != nil {
// 		t.Fatalf(`err = %s, want empty`, err.Error())
// 	}

// 	rid := c.Value(middleware.RequestIDKey).(string)
// 	cid := c.Value(contexts.Device).(string)
// 	l := fmt.Sprintf("rid:%s cid:%s lang:en hello:world", rid, cid)
// 	s := string(buf)
// 	if !strings.Contains(s, l) {
// 		t.Fatalf(`s contains %s = false, want true`, l)
// 	}
// }
