package browscap_go

type dictionary struct {
	sorted		[]*section
	expressions	map[string][]*expression
	mapped		map[string]*section
}

type section struct {
	Name	string
	Data	map[string]string
}

func newDictionary() *dictionary {
	return &dictionary{
		sorted:			[]*section{},
		expressions:	make(map[string][]*expression),
		mapped:			make(map[string]*section),
	}
}

func newSection(name string) *section {
	return &section{
		Name:	name,
		Data:	make(map[string]string),
	}
}

func (self *dictionary) findData(name string) (map[string]string) {
	res := make(map[string]string)

	if item, found := self.mapped[name]; found {
		// Parent's data
		if parentName, hasParent := item.Data["Parent"]; hasParent {
			parentData := self.findData(parentName)
			if len(parentData) > 0 {
				for k, v := range parentData {
					if k == "Parent" {
						continue
					}
					res[k] = v
				}
			}
		}
		// It's item data
		if len(item.Data) > 0 {
			for k, v := range item.Data {
				if k == "Parent" {
					continue
				}
				res[k] = v
			}
		}
	}

	return res
}
