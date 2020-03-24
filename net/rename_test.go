package net

import "testing"

func TestFormatURLfileName(t *testing.T) {
	url := "https://pdfs.journals.lww.com/ejanaesthesiology/2018/02000/European_guidelines_on_perioperative_venous.4.pdf?token=method|ExpireAbsolute;source|Journals;ttl|1585024745140;payload|mY8D3u1TCCsNvP5E421JYK6N6XICDamxByyYpaNzk7FKjTaa1Yz22MivkHZqjGP4kdS2v0J76WGAnHACH69s21Csk0OpQi3YbjEMdSoz2UhVybFqQxA7lKwSUlA502zQZr96TQRwhVlocEp/sJ586aVbcBFlltKNKo+tbuMfL73hiPqJliudqs17cHeLcLbV/CqjlP3IO0jGHlHQtJWcICDdAyGJMnpi6RlbEJaRheGeh5z5uvqz3FLHgPKVXJzdGZnEagBFgfcfP0kYnmKqyvZmvve3Z5Pif7IrCfBJwKdlADvekr+x2HZOXmESdfVD;hash|ZNzjBogZCFnBkf1qpBBFCw=="
	FormatURLfileName(url, false, 20, "")
}
