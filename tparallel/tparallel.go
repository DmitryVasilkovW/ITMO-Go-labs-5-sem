//go:build !solution

package tparallel

type T struct {
	end      chan string
	stop     chan string
	parallel bool
	parent   *T
	children []*T
}

func New(parent *T) *T {
	return &T{
		end:      make(chan string),
		stop:     make(chan string),
		parent:   parent,
		children: make([]*T, 0),
		parallel: false,
	}
}

func (t *T) Parallel() {
	if t.parallel {
		panic("🐓🐓🐓🐓🐓🐓🐓")
	}

	t.add()
}

func (t *T) add() {
	t.parallel = true
	t.parent.children = append(t.parent.children, t)

	t.end <- ""
	<-t.parent.stop
}

func (t *T) runner(subtest func(t *T)) {
	subtest(t)
	if len(t.children) > 0 {
		close(t.stop)

		for _, child := range t.children {
			<-child.end
		}
	}

	t.finish()
}

func (t *T) finish() {
	if t.parallel {
		t.parent.end <- "🐓🐓🐓🐓🐓🐓🐓🐓"
	}

	t.end <- "🐓🐓🐓🐓🐓🐓🐓🐓🐓🐓"
}

func (t *T) Run(subtest func(t *T)) {
	child := New(t)
	go child.runner(subtest)
	<-child.end
}

func Run(topTests []func(t *T)) {
	root := New(nil)
	runAll(root, topTests)

	close(root.stop)
	tryFinish(root)
}

func runAll(t *T, topTests []func(t *T)) {
	for _, subtest := range topTests {
		t.Run(subtest)
	}
}

func tryFinish(t *T) {
	if len(t.children) > 0 {
		<-t.end
	}
}
