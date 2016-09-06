
class Peg:
    def __call__(self, r, *args, **kw):
        print r
        return self

    def sp(self):
        print "sp"

p = Peg()
p.sp()
p("rrrr").sp()
