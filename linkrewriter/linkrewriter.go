package linkrewriter

import(
	"net/url"
	"strings"
)

type LinkRewriter interface {	
	CreateContext(content string) RewriteContext
	GetRewrittenUrl(u *url.URL) (string, bool)
}

type RewriteContext interface {
	Rewrite(u *url.URL, match string)
	String() string
}

func Create() LinkRewriter {
	return &linkRewriterImpl{
		rewrittenUrls: make(map[string]string),
	}
}

type linkRewriterImpl struct {	
	rewrittenUrls map[string]string
}

func (self *linkRewriterImpl) CreateContext(content string) RewriteContext {
	return &rewriteContextImpl {
		content: content,
		linkRewriter: self,
	}
}

func (self *linkRewriterImpl) GetRewrittenUrl(u *url.URL) (string, bool) {
	if rewrite, ok := self.rewrittenUrls[u.String()]; ok {
		return rewrite, true
	} else {
		return "", false 
	}
}

func (self *linkRewriterImpl) addRewrite(u *url.URL, rewrite string) {
	self.rewrittenUrls[u.String()] = rewrite
}

type rewriteContextImpl struct {
	content string
	linkRewriter *linkRewriterImpl
}

func (self *rewriteContextImpl) Rewrite(u *url.URL, match string) {
	rewrittenUrl := "/" + strings.ReplaceAll(u.Host, ".", "_") + u.RequestURI()

	self.content = strings.ReplaceAll(self.content, match, rewrittenUrl)

	self.linkRewriter.addRewrite(u, rewrittenUrl)
}

func (self *rewriteContextImpl) String() string {
	return self.content
}