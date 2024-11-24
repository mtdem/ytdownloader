package learning

// Notes of weird, quirky things that go does

/*
https://go.dev/play/
https://google.github.io/styleguide/go/decisions#receiver-type
https://yourbasic.org/golang/default-zero-value/
https://www.geeksforgeeks.org/strings-join-function-in-golang-with-examples/
https://go.dev/ref/spec
*/

// 1 - For a func/type property,
// Capital beginning = public/exposed
// lowercase beginning = private/not exposed
// VideoId => public, videoId => private

// 2 - explicit class declaration is not used
// structs are used instead
// and methods reference the structs that call them
