package myCobra

import "testing"

func BenchmarkCommand_Execute(b *testing.B) {
	var test = &Command {
		Use: "test",
		Short: "test",
		Long: "test long",
	}
	for i := 0; i < b.N; i++ {
		test.Execute()
	}
}
