package cartio

import (
    "bytes"
    "testing"
)

func TestWriteReadRoundTrip(t *testing.T) {
    m := Manifest{Title:"Hello", Author:"RF", Description:"d", Genre:"Action", Tags:[]string{"a","b"}, Entry:"main.lua"}
    assets := []Asset{{Name:"main.lua", Data:[]byte("print('hi')")}, {Name:"sprites.png", Data:[]byte{1,2,3}}}

    var buf bytes.Buffer
    if err := Write(&buf, m, assets); err != nil { t.Fatalf("write: %v", err) }

    gotM, files, err := Read(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
    if err != nil { t.Fatalf("read: %v", err) }
    if gotM.Entry != "main.lua" || gotM.Title != "Hello" { t.Fatalf("bad manifest: %+v", gotM) }
    if len(files) != 2 { t.Fatalf("expected 2 assets, got %d", len(files)) }
    if _, ok := files["assets/main.lua"]; !ok { t.Fatalf("missing main.lua") }
}


