package gossie

type application struct {
	name          string
	description   string
	version       string
	author        string
	commands      map[string]*Command
	defaultAction func(*Context) error
}
