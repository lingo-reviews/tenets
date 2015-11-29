package log

import "io"

// Logger mimics golang's standard Logger as an interface.
type Logger interface {
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalln(args ...interface{})
	Print(args ...interface{})
	Printf(format string, args ...interface{})
	Println(args ...interface{})
}

type CloserLogger interface {
	Logger
	io.Closer
}

type closer struct {
	Logger
	io.Closer
}

func (c *closer) Close() error {
	if c.Closer != nil {
		return c.Closer.Close()
	}
	return nil
}

type nolog struct{}

func (n nolog) Fatal(args ...interface{})                 {}
func (n nolog) Fatalf(format string, args ...interface{}) {}
func (n nolog) Fatalln(args ...interface{})               {}
func (n nolog) Print(args ...interface{})                 {}
func (n nolog) Printf(format string, args ...interface{}) {}
func (n nolog) Println(args ...interface{})               {}

// logging off by default
var logger = func() *closer { return &closer{Logger: nolog{}} }

// TODO(waigani) write a logger 1. to stream logs back to client 2. to write
// to file. Turn logs on with --debug flag.
// For now, to see logs, uncomment the following line.
// var logger = func() *closer { return &closer{Logger: log.New(os.Stderr, "", log.LstdFlags)} }

// log to file
// var logger = func() *closer {
// 	home := os.Getenv("HOME")
// 	// TODO(waigani) this should not be hardcoded here.
// 	dir := filepath.Join(home, ".lingo_home", "logs")
// 	filename := filepath.Join(dir, "all.log")

// 	if _, err := os.Stat(dir); os.IsNotExist(err) {
// 		err := os.MkdirAll(dir, 0777)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}

// 	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return &closer{
// 		Logger: log.New(f, "", log.LstdFlags),
// 		Closer: f,
// 	}
// }

func GetLogger() Logger {
	return logger()
}

// SetLogger sets the logger that is used in grpc.
func SetLogger(l Logger) {
	logger = func() *closer { return &closer{Logger: l} }
}

// Fatal is equivalent to Print() followed by a call to os.Exit() with a non-zero exit code.
func Fatal(args ...interface{}) {
	l := logger()
	defer l.Close()
	l.Fatal(args...)
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit() with a non-zero exit code.
func Fatalf(format string, args ...interface{}) {
	l := logger()
	defer l.Close()
	l.Fatalf(format, args...)
}

// Fatalln is equivalent to Println() followed by a call to os.Exit()) with a non-zero exit code.
func Fatalln(args ...interface{}) {
	l := logger()
	defer l.Close()
	l.Fatalln(args...)
}

// Print prints to the logger. Arguments are handled in the manner of fmt.Print.
func Print(args ...interface{}) {
	l := logger()
	defer l.Close()
	l.Print(args...)
}

// Printf prints to the logger. Arguments are handled in the manner of fmt.Printf.
func Printf(format string, args ...interface{}) {
	l := logger()
	defer l.Close()
	l.Printf(format, args...)
}

// Println prints to the logger. Arguments are handled in the manner of fmt.Println.
func Println(args ...interface{}) {
	l := logger()
	defer l.Close()
	l.Println(args...)
}
