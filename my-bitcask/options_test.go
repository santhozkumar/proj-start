package bitcask

// func TestNewDefaultConfig(t *testing.T) {
// 	got := newDefaultConfig()
// 	want := &config.Config{
// 		MaxDataFileSize: 1048576,
// 		MaxKeySize:      64,
// 		MaxValueSize:    65536,
// 		Sync:            false,
// 		AutoRecovery:    false,
// 		AutoReadOnly:    false,
// 		DirMode:         448,
// 		FileMode:        384,
// 	}
//
// 	if *got != *want {
// 		t.Errorf("got %v want %v", got, want)
// 	}
//
// }
//
// func TestOption(t *testing.T) {
//
// 	opts := []Option{WithAutoReadOnly(true), WithAutoRecovery(true)}
// 	want := &config.Config{
// 		MaxDataFileSize: 1048576,
// 		MaxKeySize:      64,
// 		MaxValueSize:    65536,
// 		Sync:            true,
// 		AutoRecovery:    true,
// 		AutoReadOnly:    true,
// 		DirMode:         448,
// 		FileMode:        384,
// 	}
//
// 	got := newDefaultConfig()
//
// 	for _, opt := range opts {
// 		opt(got)
// 	}
//
// 	if *got != *want {
// 		t.Errorf("got %v want %v", got, want)
// 	}
//
// }
