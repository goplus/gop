package extractor

// -----------------------------------------------------------------------------

const (
	cateLeaf      = 0
	cateNode      = 1
	cateArray     = 2
	cateLeafArray = cateLeaf | cateArray
	cateNodeArray = cateNode | cateArray
)

type Cons struct {
	tag  string
	val  interface{}
	tail *Cons
	cate int
}

type ConsList struct {
	hd *Cons
}

// -----------------------------------------------------------------------------

type Node struct {
	ConsList
}

func (p *Node) clear() {

	p.hd = nil
}

func (p *Node) insertLeaf(tag string, val interface{}, isArray bool) {

	cate := cateLeaf
	if isArray {
		cate = cateLeafArray
	}
	p.hd = &Cons{tag: tag, val: val, tail: p.hd, cate: cate}
}

func (p *Node) insertNode(tag string, isArray bool) *Node {

	node := new(Node)
	cate := cateNode
	if isArray {
		cate = cateNodeArray
	}
	p.hd = &Cons{tag: tag, val: node, tail: p.hd, cate: cate}
	return node
}

func (p *Node) toMap() map[string]interface{} {

	hd := p.hd
	if hd == nil {
		return nil
	}

	data := make(map[string]interface{})
	for hd != nil {
		tag := hd.tag
		var v interface{}
		switch hd.cate {
		case cateLeaf:
			v, hd = hd.val, hd.tail
		case cateNode:
			v, hd = hd.val.(*Node).toMap(), hd.tail
		case cateNodeArray:
			n := arrayCollect(hd)
			items := make([]map[string]interface{}, n)
			for i := n - 1; i >= 0; i-- {
				items[i] = hd.val.(*Node).toMap()
				hd = hd.tail
			}
			v = items
		case cateLeafArray:
			n := arrayCollect(hd)
			items := make([]interface{}, n)
			for i := n - 1; i >= 0; i-- {
				items[i] = hd.val
				hd = hd.tail
			}
			v = items
		default:
			panic("unknown node category")
		}
		data[tag] = v
	}
	return data
}

func arrayCollect(hd *Cons) (n int) {

	n = 1
	cate, tag := hd.cate, hd.tag
	for {
		hd = hd.tail
		if hd == nil || hd.cate != cate || hd.tag != tag {
			return
		}
		n++
	}
}

// -----------------------------------------------------------------------------
