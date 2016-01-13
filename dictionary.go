package browscap_go

type dictionary struct {
	browsers map[string]*Browser

	tree *ExpressionTree
}

type section map[string]string

func newDictionary() *dictionary {
	return &dictionary{
		browsers: make(map[string]*Browser),

		tree: NewExpressionTree(),
	}
}

func (dict *dictionary) getBrowser(name string) *Browser {
	if d, ok := dict.browsers[name]; ok {
		d.build(dict.browsers)
		return d
	}
	return nil
}
