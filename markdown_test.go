package main

import (
	"fmt"
	"testing"
)

func TestMDToENML(t *testing.T) {
	fmt.Println(MDToENML("### Header\n"))
}

const enml = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE en-note SYSTEM "http://xml.evernote.com/pub/enml2.dtd">
<en-note><div><br/></div><div><br/></div><div><b>that old navy thick sweater</b></div><div><br/></div><div>doesn't look as nice as I thoguht it would.</div><div>it is cheap, it has no elasticity.</div><div>its already stretched out at the cuffs.</div><div><br/></div><div>it does look frumpy, especially around the waste.</div><div><br/></div><div><br/></div><div><br/></div><div><b>trying to dress multiple ways at once</b></div><div><br/></div><div>I have to settle on who I am.</div><div>or at least switch betwee things, but know what I&quot;m switching between.</div><div><br/></div><div><br/></div><div><br/></div><div><br/></div><div><b>layers are bad for me, can't do them.</b></div><div><br/></div><div>can't deal with the heat.</div><div>need breathability.</div><div><br/></div><div>so how do I deal with this?</div><div><br/></div><div><br/></div><div><b>have been getting more eyes though.</b></div><div>when I wear these new cloths I have.</div><div><br/></div><div>? want a good black leather jacket.</div><div><span>    something complicated, not just simple and boring</span><br/></div><div><br/></div><div><br/></div><div><br/></div><div><b>what to get tomorrow</b></div><div><b><br/></b></div><div><ul><li>dark black jeans. not blue</li></ul><div><br/></div><div><ul><li>satchel</li></ul><div><br/></div><div><ul><li>some vneck tees, XL</li></ul><div><br/></div><div><ul><li>dark pattern button downs</li></ul><div><br/></div><div><br/></div><div>maybe, since I can't wear layers,</div><div>I have to get more complicated single layers.</div><div>liek shirts with texture and patterns on them.</div><div><br/></div><div><br/></div></div></div></div><div><br/></div><div><br/></div><b/></div><div><b><br/></b></div><div><br/></div></en-note>
`

func TestENMLToMD(t *testing.T) {
	fmt.Println(ENMLToMD(enml))
}
