package bitcask

// func TestBatch(t *testing.T) {
//
//     tmpDir, err := os.MkdirTemp("", "bitcask")
//     assert.NoError(t, err)
//
//     var db DB
//
// 	t.Run("Setup", func(t *testing.T) {
// 		t.Run("Open", func(t *testing.T){
//             db, err = Open(tmpDir)
//             assert.NoError(t, err)
//         })
// 	})
//
//     t.Run("Batch", func(t *testing.T){
//         b := db.Batch()
//         b.Put([]byte("Hello"), []byte("World"))
//         b.Put(Key("foo"), Value("Bar"))
//         fmt.Println(b.Entries())
//         assert.NoError(t, err, db.WriteBatch(b))
//
//
//         tests := map[string]Value{
//             "foo": Value("Bar"),
//             "Hello": Value("World"),
//         }
//
//         for key, val := range(tests){
//             actual,err := db.Get(Key(key))
//             assert.NoError(t, err)
//             assert.Equal(t, val, actual)
//         }
//     })
// }
